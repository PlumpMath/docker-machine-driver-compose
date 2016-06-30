// Copyright (C) 2016-2017 ATOS - All rights reserved.
package compose

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"

	"errors"

	"bytes"
	"text/template"

	"io/ioutil"

	"regexp"
	"strings"

	"github.com/apache/brooklyn-client/api/application"
	"github.com/apache/brooklyn-client/api/server"
	"github.com/apache/brooklyn-client/net"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

const (
	driverName     = "compose"
	defaultSSHUser = "compose"
	defaultSSHPort = 22
)

const (
	SMALL   = "small"
	MEDIUM  = "medium"
	LARGE   = "large"
	XLARGE  = "xlarge"
	XXLARGE = "xxlarge"

	CENTOS7  = "centos:7"
	UBUNTU14 = "ubuntu:14"

	COMPOSE_DOCKERHOST_CATALOG = "com.canopy.compose.dockerhost"
	MAPPED_PORT_SENSOR_NAME    = "mapped.portPart.dockerhost.port"

)

var (
	dockerPort = 2376
	swarmPort  = 3376

	defaultComposeBaseUrl  = "http://localhost:8081"
	defaultOperatingSystem = UBUNTU14
	defaultTemplateSize    = SMALL

	defaultOpenPorts = "tomcat.port: 8080,web.port: 80,ssl.port: 443"

	templateSizes    = []string{SMALL, MEDIUM, LARGE, XLARGE, XXLARGE}
	operatingSystems = []string{CENTOS7, UBUNTU14}

	openPortsRegx = regexp.MustCompile(`((([A-Za-z])\w+[.]port[:][\ ][0-9]{1,5})(?:\,)?)+/g`)

	errorMissingUser         = errors.New("Compose user requires use the --compose-user option")
	errorMissingPassword     = errors.New("Compose password requires use the --compose-password option")
	errorMissingLocation     = errors.New("Compose target location requires use the --compose-target-location option")
	errorInvalidTemplateSize = errors.New("Specified template size not supported, available options are small, medium, large, xlarge, xxlarge")
	errorInvalidOS           = errors.New("Specified operating system not supported, available options are centos:7, ubuntu:14")
	errorInvalidOpenPorts    = errors.New("Invalid open port request, format is > web.port: 2345,tomcat.port: 8080 < etc")
)

type Driver struct {
	*drivers.BaseDriver
	Id string

	ComposeClient *net.Network

	Application *Application

	ApplicationId  string
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
      compose.os.name: {{.OsName}}
      compose.os.version: {{.OsVersion}}
      compose.template.size: {{.TemplateSize}}{{range $_, $val := .OpenPorts}}
      {{$val}}{{end}}
      compose.sshUserKey: {{.SshUserKey}}`

// Template for DockerSwarmHost
const swarmHostAppTmpl = `name: {{.Name}}
location: {{.Location}}
services:
  - type: {{.Type}}
    brooklyn.config:
      dockerhost.port: 2376
      swarmhost.port: 3376
      compose.os.name: {{.OsName}}
      compose.os.version: {{.OsVersion}}
      compose.template.size: {{.TemplateSize}}{{range $_, $val := .OpenPorts}}
      {{$val}}{{end}}
      compose.sshUserKey: {{.SshUserKey}}`

func applicationYaml(swarmMaster bool, application *Application) ([]byte, error) {
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

	catalogs, err := CatalogByRegex(d.ComposeClient, COMPOSE_DOCKERHOST_CATALOG)
	catalogId := catalogs[0].Id
	log.Infof(catalogId)

	d.Application.Type = catalogId
	d.Application.SshUserKey = publicKey

	appYaml, err := applicationYaml(d.SwarmMaster, d.Application)
	if err != nil {
		return err
	}

	taskSummary, err := application.CreateFromBytes(d.ComposeClient, appYaml)
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
			Name:   "compose-base-url",
			Usage:  "Compose Base URL",
			Value:  defaultComposeBaseUrl,
			EnvVar: "COMPOSE_BASE_URL",
		},
		mcnflag.StringFlag{
			Name:   "compose-user",
			Usage:  "Compose User",
			EnvVar: "COMPOSE_USER",
		},
		mcnflag.StringFlag{
			Name:   "compose-password",
			Usage:  "Compose Password",
			EnvVar: "COMPOSE_PASSWORD",
		},
		mcnflag.StringFlag{
			Name:   "compose-target-location",
			Usage:  "Compose Target Location",
			EnvVar: "COMPOSE_TARGET_LOCATION",
		},
		mcnflag.StringFlag{
			Name:   "compose-target-os",
			Usage:  "Compose Target OS",
			Value:  defaultOperatingSystem,
			EnvVar: "COMPOSE_TARGET_OS",
		},
		mcnflag.StringFlag{
			Name:   "compose-template-size",
			Usage:  "Compose Template Size",
			Value:  defaultTemplateSize,
			EnvVar: "COMPOSE_TEMPLATE_SIZE",
		},
		mcnflag.StringFlag{
			Name:   "compose-open-ports",
			Usage:  "Compose Open Ports",
			Value:  defaultOpenPorts,
			EnvVar: "COMPOSE_OPEN_PORTS",
		},
	}
}

// GetIP returns an IP or hostname that this host is available at
// e.g. 1.2.3.4 or docker-host-d60b70a14d3a.cloudapp.net
func (d *Driver) GetIP() (string, error) {
	log.Debugf("Calling .GetIP()")
	sshHostAddress, err := DescendantsSshHostAndPortSensor(d.ComposeClient, d.ApplicationId)
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
	log.Debugf("Calling .GetSSHHostname()")
	return d.GetIP()
}

// GetSSHPort returns port for use with ssh
func (d *Driver) GetSSHPort() (int, error) {
	log.Debugf("Calling .GetSSHPort()")
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
	log.Debugf("Calling .GetURL()")
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

	sensorInfo, err := DescendantsSensor(d.ComposeClient, d.ApplicationId, MAPPED_PORT_SENSOR_NAME)
	if err != nil {
		for key, _ := range sensorInfo {
			dockerPort = sensorInfo[key]
			break
		}
	}

	url := fmt.Sprintf("tcp://%s:%d", ip, dockerPort)
	return url, nil
}

// GetState returns the state that the host is in (running, stopped, etc)
func (d *Driver) GetState() (state.State, error) {
	log.Debugf("Calling .GetState()")

	if d.ApplicationId == "" {
		log.Warnf("Application id is nil.")
		return state.Stopped, nil
	}

	applicationSummary, err := application.Application(d.ComposeClient, d.ApplicationId)
	if err != nil {
		log.Warnf("Application does not exists.")
		return state.Stopped, nil
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

	_, err := application.Application(d.ComposeClient, d.ApplicationId)

	if err != nil {
		log.Warnf("Application having id [%s] does not exists", d.ApplicationId)
		return nil
	}

	_, err = application.Delete(d.ComposeClient, d.ApplicationId)
	if err != nil {
		log.Warnf("Error while killing application [%s]", d.ApplicationId)
		return nil
	}
	return err
}

// PreCreateCheck allows for pre-create operations to make sure a driver is ready for creation
func (d *Driver) PreCreateCheck() error {

	// Validate specified server exists and reachable.
	state, err := server.Healthy(d.ComposeClient)
	if err != nil {
		return err
	} else if state != "true" {
		return errors.New("Compose Server not healthy.")
	}

	// Validate specified location exists.
	if _, err = LocationExists(d.ComposeClient, d.Application.Location); err != nil {
		return err
	}

	// Validate specified operating system catalog exists.
	catalogs, err := CatalogByRegex(d.ComposeClient, COMPOSE_DOCKERHOST_CATALOG)
	if err != nil {
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
	_, err := Delete(d.ComposeClient, d.ApplicationId)

	if err != nil {
		log.Warnf("Error while removing application [%s]", d.ApplicationId)
		return nil
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
	d.SetSwarmConfigFromFlags(opts)
	baseUrl := opts.String("compose-base-url")
	user := opts.String("compose-user")                // mandatory
	password := opts.String("compose-password")        // mandatory
	location := opts.String("compose-target-location") // mandatory
	operatingSystem := opts.String("compose-target-os")
	templateSize := opts.String("compose-template-size")
	openPortsStr := opts.String("compose-open-ports")


	if user == "" {
		return errorMissingUser
	}

	if password == "" {
		return errorMissingPassword
	}

	if location == "" {
		return errorMissingLocation
	}

	if !contains(templateSize, templateSizes) {
		return errorInvalidTemplateSize
	}

	if !contains(operatingSystem, operatingSystems) {
		return errorInvalidOS
	}

	/*
	if openPortsStr != "" && strings.Trim(openPortsStr," ") != "" &&  !openPortsRegx.MatchString(openPortsStr) {
		return errorInvalidOpenPorts
	}*/

	tokens := strings.Split(operatingSystem, ":")
	if len(tokens) != 2 {
		return errorInvalidOS
	}

	d.Application = NewApplication()
	d.Application.Name = d.Id
	d.Application.Location = location
	d.Application.OsName = tokens[0]
	d.Application.OsVersion = tokens[1]
	d.Application.TemplateSize = templateSize
	d.Application.OpenPorts = strings.Split(openPortsStr, ",")

	d.ComposeClient = net.NewNetwork(baseUrl, user, password, false)
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
