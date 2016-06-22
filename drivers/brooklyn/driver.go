// Copyright (C) 2016-2017 ATOS - All rights reserved.
package brooklyn

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"

	"errors"

	"bytes"
	"text/template"

	"io/ioutil"

	"github.com/apache/brooklyn-client/api/application"
	"github.com/apache/brooklyn-client/net"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/apache/brooklyn-client/api/server"
)

const (
	driverName     = "brooklyn"
	defaultSSHUser = "compose"
	defaultSSHPort = 22
)

const (
	SMALL   = "small"
	MEDIUM  = "medium"
	LARGE   = "large"
	XLARGE  = "xlarge"
	XXLARGE = "xxlarge"

	CENTOS = "centos"
	CENTOS1 = "centos1"
	UBUNTU = "ubuntu"
	SUSE   = "suse"

	COMPOSE_CATALOG_ID_STARTS_WITH = "com.canopy.compose"
)

var (
	dockerPort = 2376
	swarmPort  = 3376

	defaultBrooklynBaseUrl = "http://localhost:8081"
	defaultOperatingSystem = "centos"
	defaultTShirtSize      = MEDIUM

	tShirtSizes      = []string{SMALL, MEDIUM, LARGE, XLARGE, XXLARGE}
	operatingSystems = []string{CENTOS, CENTOS1, UBUNTU, SUSE}

	errorMissingUser       = errors.New("Brooklyn user requires the --brooklyn-user option")
	errorMissingPassword   = errors.New("Brooklyn password requires the --brooklyn-password option")
	errorMissingLocation   = errors.New("Brooklyn target location requires the --brooklyn-target-location option")
	errorInvalidTShirtSize = errors.New("Brooklyn t shirt size is invalid, supports only small, medium, large, xlarge, xxlarge")
	errorInvalidOS         = errors.New("Brooklyn requested operating system is not yet supported, currently supported are centos, ubuntu or suse")
)

type Driver struct {
	*drivers.BaseDriver
	Id string

	BrooklynClient *net.Network

	Location        string
	OperatingSystem string
	TShirtSize      string

	ApplicationId string

	SshHostAddress SshHostAddress
}

func NewDriver(hostName, storePath string) *Driver {
	id := generateId()

	driver := &Driver{
		Id: id,
		BaseDriver: &drivers.BaseDriver{
			SSHUser:     defaultSSHUser,
			SSHPort:     defaultSSHPort,
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
	return driver
}

// Template for DockerHost
const dockerHostAppTmpl = `name: {{.Name}}
location: {{.Location}}
services:
  - type: {{.Type}}
    brooklyn.config:
      dockerhost.port: 2376
      compose.sshUserKey: {{.SshUserKey}}`
// Template for DockerSwarmHost
const swarmHostAppTmpl = `name: {{.Name}}
location: {{.Location}}
services:
  - type: {{.Type}}
    brooklyn.config:
      dockerhost.port: 2376
      swarmhost.port: 3376
      compose.sshUserKey: {{.SshUserKey}}`

func applicationYaml(swarmMaster bool, application Application) ([]byte, error) {
	// Create a new template and parse the application into it.
	var t *template.Template
	if swarmMaster {
		t = template.Must(template.New("application").Parse(swarmHostAppTmpl))
	} else {
		t = template.Must(template.New("application").Parse(dockerHostAppTmpl))
	}
	var appYml bytes.Buffer
	err := t.Execute(&appYml, application)
	log.Infof(appYml.String())
	var app []byte
	if err != nil {
		return app, err
	}
	return appYml.Bytes(), nil
}

// contains validate element exists in slice or not.
func contains(element string, elements []string) bool {
	for _, s := range elements {
		if element == s {
			return true
		}
	}
	return false
}

// generateId generates random id.
func generateId() string {
	rb := make([]byte, 10)
	_, err := rand.Read(rb)
	if err != nil {
		log.Warnf("Unable to generate id: %s", err)
	}

	h := md5.New()
	io.WriteString(h, string(rb))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Create a host using the driver's config
func (d *Driver) Create() error {
	// Create SSH Key and Pair
	publicKey, err := d.createKeyPair()
	if err != nil {
		return fmt.Errorf("unable to create key pair: %s", err)
	}

	regex := fmt.Sprintf("%s.%s",COMPOSE_CATALOG_ID_STARTS_WITH,d.OperatingSystem)
	catalogs, err := CatalogByRegex(d.BrooklynClient, regex);
	catalogId := catalogs[0].Id
	log.Infof(catalogId)

	app := Application{
		Name:       fmt.Sprintf("%s-%s", d.BaseDriver.MachineName, d.Id),
		Location:   d.Location,
		Type:       catalogId,
		SshUserKey: publicKey,
	}

	appYaml, err := applicationYaml(d.SwarmMaster, app)
	if err != nil {
		return err
	}

	taskSummary, err := application.CreateFromBytes(d.BrooklynClient, appYaml)
	d.ApplicationId = taskSummary.EntityId

	if err != nil {
		log.Error(err)
		return err
	} else {
		log.Infof(taskSummary.Id)
	}

	//Wait for Instance to Running.
	d.waitForInstance()

	return nil
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return driverName
}

// GetCreateFlags returns the mcnflag.Flag slice representing the flags
// that can be set, their descriptions and defaults.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			Name:   "brooklyn-base-url",
			Usage:  "Brooklyn Base URL",
			Value:  defaultBrooklynBaseUrl,
			EnvVar: "BROOKLYN_BASE_URL",
		},
		mcnflag.StringFlag{
			Name:   "brooklyn-user",
			Usage:  "Brooklyn User",
			EnvVar: "BROOKLYN_USER",
		},
		mcnflag.StringFlag{
			Name:   "brooklyn-password",
			Usage:  "Brooklyn Password",
			EnvVar: "BROOKLYN_PASSWORD",
		},
		mcnflag.StringFlag{
			Name:   "brooklyn-target-location",
			Usage:  "Brooklyn Target Location",
			EnvVar: "BROOKLYN_TARGET_LOCATION",
		},
		mcnflag.StringFlag{
			Name:   "brooklyn-target-os",
			Usage:  "Brooklyn Target OS",
			Value:  defaultOperatingSystem,
			EnvVar: "BROOKLYN_TARGET_OS",
		},
		mcnflag.StringFlag{
			Name:   "brooklyn-template-size",
			Usage:  "Brooklyn Template Size",
			Value:  defaultTShirtSize,
			EnvVar: "BROOKLYN_TEMPLATE_SIZE",
		},
	}
}

// GetIP returns an IP or hostname that this host is available at
// e.g. 1.2.3.4 or docker-host-d60b70a14d3a.cloudapp.net
func (d *Driver) GetIP() (string, error) {
	log.Infof("GetIP()")
	sshHostAddress, err := DescendantsSshHostAndPortSensor(d.BrooklynClient, d.ApplicationId)
	if err != nil {
		return "", nil
	}
	d.SshHostAddress = sshHostAddress
	log.Infof(d.SshHostAddress.HostAndPort.Host)
	return sshHostAddress.HostAndPort.Host, nil
}

// GetMachineName returns the name of the machine
func (d *Driver) GetMachineName() string {
	return d.MachineName
}

// GetSSHHostname returns hostname for use with ssh
func (d *Driver) GetSSHHostname() (string, error) {
	log.Infof("GetSSHHostname()")
	return d.GetIP()
}

// GetSSHPort returns port for use with ssh
func (d *Driver) GetSSHPort() (int, error) {
	log.Infof("GetSSHPort()")
	log.Info(d.SshHostAddress.HostAndPort.Port)
	return d.SshHostAddress.HostAndPort.Port, nil
}

// GetSSHUsername returns username for use with ssh
func (d *Driver) GetSSHUsername() string {
	/*
	if d.SshHostAddress.User != "" {
		return d.SshHostAddress.User
	}
	*/
	return defaultSSHUser
}

// GetURL returns a Docker compatible host URL for connecting to this host
// e.g. tcp://1.2.3.4:2376
func (d *Driver) GetURL() (string, error) {
	log.Infof("GetURL()")
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}

	url := fmt.Sprintf("tcp://%s:%d", ip, dockerPort)
	return url, nil
}

// GetState returns the state that the host is in (running, stopped, etc)
func (d *Driver) GetState() (state.State, error) {
	log.Infof("GetState()")

	if d.ApplicationId == "" {
		log.Warnf("Application id is nil.")
		return state.Stopped, errors.New("Application id is nil.")
	}

	applicationSummary, err := application.Application(d.BrooklynClient, d.ApplicationId)
	if err != nil {
		return state.Error, err
	}

	log.Info(applicationSummary.Status)
	switch applicationSummary.Status {
	case "RUNNING":
		return state.Running, nil
	case "STARTING":
		return state.Starting, nil
	case "STOPPING":
		return state.Stopping, nil
	case "ERROR":
		return state.Error, nil
	default:
		return state.None, nil
	}
}

// Kill stops a host forcefully
func (d *Driver) Kill() error {

	if d.ApplicationId == "" {
		log.Warnf("Empty ApplicationId")
		return nil
	}

	_, err := application.Application(d.BrooklynClient, d.ApplicationId)

	if err != nil {
		log.Warnf("Application having id [%s] does not exists", d.ApplicationId)
		return nil
	}

	_, err = application.Delete(d.BrooklynClient, d.ApplicationId)
	if err != nil {
		log.Errorf("Error while killing application [%s]", d.ApplicationId)
	}
	return err
}

// PreCreateCheck allows for pre-create operations to make sure a driver is ready for creation
func (d *Driver) PreCreateCheck() error {

	// Validate specified server exists and reachable.
	state,err := server.Healthy(d.BrooklynClient)
	if err != nil {
		return err
 	} else if state != "true" {
		return errors.New("Brooklyn Server not healthy.")
	}

	// Validate specified location exists.
	if _, err = LocationExists(d.BrooklynClient,d.Location); err != nil {
		return err
	}

	// Validate specified operating system catalog exists.
	regex := fmt.Sprintf("%s.%s",COMPOSE_CATALOG_ID_STARTS_WITH,d.OperatingSystem)
	catalogs, err := CatalogByRegex(d.BrooklynClient, regex);
	if  err != nil  {
		return err
	} else if len(catalogs) <= 0 {
		return errors.New("Catalog does not exists.")
	}

	return nil
}

// Remove a host
func (d *Driver) Remove() error {
	if d.ApplicationId == "" {
		log.Warnf("Empty ApplicationId")
		return nil
	}
	_, err := Delete(d.BrooklynClient, d.ApplicationId)

	if err != nil {
		log.Errorf("Error while removing application [%s]", d.ApplicationId)
	}
	return err
}

// Restart a host. This may just call Stop(); Start() if the provider does not
// have any special restart behaviour.
func (d *Driver) Restart() error {
	log.Infof("TODO: Restart not yet implemented")
	return nil
}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	baseUrl := opts.String("brooklyn-base-url")
	user := opts.String("brooklyn-user")         // mandatory
	password := opts.String("brooklyn-password") // mandatory

	d.Location = opts.String("brooklyn-target-location") // mandatory
	d.OperatingSystem = opts.String("brooklyn-target-os")
	d.TShirtSize = opts.String("brooklyn-template-size")
	d.SetSwarmConfigFromFlags(opts)

	if user == "" {
		return errorMissingUser
	}

	if password == "" {
		return errorMissingPassword
	}

	if d.Location == "" {
		return errorMissingLocation
	}

	if !contains(d.TShirtSize, tShirtSizes) {
		return errorInvalidTShirtSize
	}

	if !contains(d.OperatingSystem, operatingSystems) {
		return errorInvalidOS
	}

	d.BrooklynClient = net.NewNetwork(baseUrl, user, password, false)

	return nil
}

// Start a host
func (d *Driver) Start() error {
	log.Infof("TODO: Start not yet implemented.")
	return nil
}

// Stop a host gracefully
func (d *Driver) Stop() error {
	log.Infof("TODO: Stop not yet implemented.")
	return nil
}

func (d *Driver) createKeyPair() (string, error) {
	keyPath := ""

	log.Debugf("Creating New SSH Key")
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return "", err
	}
	keyPath = d.GetSSHKeyPath()

	_, err := ioutil.ReadFile(keyPath + ".pub")
	if err != nil {
		return "", err
	}

	log.Debugf(keyPath)
	//keyName := d.MachineName

	publicKey, err := ioutil.ReadFile(keyPath + ".pub")
	if err != nil {
		return "", err
	}

	return string(publicKey), nil
}

func (d *Driver) waitForInstance() error {
	if err := mcnutils.WaitFor(d.instanceIsRunning); err != nil {
		return err
	}
	return nil
}

func (d *Driver) instanceIsRunning() bool {
	st, err := d.GetState()
	if err != nil {
		log.Debug(err)
	}
	if st == state.Running {
		return true
	}
	return false
}

func (d *Driver) isSwarmMaster() bool {
	return d.SwarmMaster
}
