// Copyright (C) 2016-2017 ATOS - All rights reserved.
package brooklyn

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"

	"errors"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/api"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/client"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/models"
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
	UBUNTU = "ubuntu"
	SUSE   = "suse"
)

var (
	defaultBrooklynBaseUrl = "http://localhost:8081"
	defaultOperatingSystem = "centos"
	defaultTShirtSize      = MEDIUM

	tShirtSizes      = []string{SMALL, MEDIUM, LARGE, XLARGE, XXLARGE}
	operatingSystems = []string{CENTOS, UBUNTU, SUSE}

	errorMissingUser       = errors.New("Brooklyn user requires the --brooklyn-user option")
	errorMissingPassword   = errors.New("Brooklyn password requires the --brooklyn-password option")
	errorMissingLocation   = errors.New("Brooklyn target location requires the --brooklyn-target-location option")
	errorInvalidTShirtSize = errors.New("Brooklyn t shirt size is invalid, supports only small, medium, large, xlarge, xxlarge")
	errorInvalidOS         = errors.New("Brooklyn requested operating system is not yet supported, currently supported are centos, ubuntu or suse")
)

type Driver struct {
	*drivers.BaseDriver
	Id string

	BrooklynClient *client.BrooklynClient

	Location        string
	OperatingSystem string
	TShirtSize      string

	ApplicationId string
}

type brooklynClient struct {
	Url      string
	User     string
	Password string
}

func GetDriverName() string {
	return driverName
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

	application := models.Application{
		Name:     d.Id,
		Location: d.Location,
		Type:     "com.canopy.compose.centos:1.3",
	}
	taskSummary, err := api.Create(
		d.BrooklynClient.GoRequestWithProxy("http://MC0WBVEC.ww930.my-it-solutions.net:3128"), application)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(taskSummary.Id)
		d.ApplicationId=taskSummary.EntityId
	}
	return err
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
			Name:   "operating-system",
			Usage:  "Operating System",
			Value:  defaultOperatingSystem,
			EnvVar: "OPERATING_SYSTEM",
		},
		mcnflag.StringFlag{
			Name:   "t-shirt-size",
			Usage:  "T Shirt Size",
			Value:  defaultTShirtSize,
			EnvVar: "T_SHIRT_SIZE",
		},
	}
}

// GetIP returns an IP or hostname that this host is available at
// e.g. 1.2.3.4 or docker-host-d60b70a14d3a.cloudapp.net
func (d *Driver) GetIP() (string, error) {
	return "1.2.3.4", nil
}

// GetMachineName returns the name of the machine
func (d *Driver) GetMachineName() string {
	return d.MachineName
}

// GetSSHHostname returns hostname for use with ssh
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetSSHKeyPath returns key path for use with ssh
func (d *Driver) GetSSHKeyPath() string {
	return "DummySSHKey"
}

// GetSSHPort returns port for use with ssh
func (d *Driver) GetSSHPort() (int, error) {
	return 2376, nil
}

// GetSSHUsername returns username for use with ssh
func (d *Driver) GetSSHUsername() string {
	return defaultSSHUser
}

// GetURL returns a Docker compatible host URL for connecting to this host
// e.g. tcp://1.2.3.4:2376
func (d *Driver) GetURL() (string, error) {
	return "tcp://1.2.3.4:2376", nil
}

// GetState returns the state that the host is in (running, stopped, etc)
func (d *Driver) GetState() (state.State, error) {
	return state.Running, nil
}

// Kill stops a host forcefully
func (d *Driver) Kill() error {
	_, err :=api.Delete(d.BrooklynClient.GoRequestWithProxy("http://MC0WBVEC.ww930.my-it-solutions.net:3128"),d.ApplicationId)
	return err
}

// PreCreateCheck allows for pre-create operations to make sure a driver is ready for creation
func (d *Driver) PreCreateCheck() error {
	return nil
}

// Remove a host
func (d *Driver) Remove() error {
	return nil
}

// Restart a host. This may just call Stop(); Start() if the provider does not
// have any special restart behaviour.
func (d *Driver) Restart() error {
	return nil
}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.BrooklynClient = &client.BrooklynClient{}
	d.BrooklynClient.BaseUrl = opts.String("brooklyn-base-url")
	d.BrooklynClient.User = opts.String("brooklyn-user")                // mandatory
	d.BrooklynClient.Password = opts.String("brooklyn-password")        // mandatory
	d.Location = opts.String("brooklyn-target-location") // mandatory
	d.OperatingSystem = opts.String("operating-system")
	d.TShirtSize = opts.String("t-shirt-size")

	if d.BrooklynClient.User == "" {
		return errorMissingUser
	}

	if d.BrooklynClient.Password == "" {
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
	return nil
}

func contains(element string, elements []string) bool {
	for _, s := range elements {
		if element == s {
			return true
		}
	}
	return false
}

// Start a host
func (d *Driver) Start() error {
	return nil
}

// Stop a host gracefully
func (d *Driver) Stop() error {
	return nil
}
