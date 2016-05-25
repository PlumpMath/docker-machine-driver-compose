// Copyright (C) 2016-2017 ATOS - All rights reserved.
package brooklyn

import (
	"github.com/docker/machine/libmachine/drivers"
)

type deviceConfig struct {
	DiskSize int
	Cpu      int
	Hostname string
}

type Driver struct {
	*drivers.BaseDriver
	deviceConfig deviceConfig
}

const (
	driverName      = "brooklyn"
	defaultMemory   = 1024
	defaultDiskSize = 0
	defaultRegion   = "dal01"
	defaultCpus     = 1
)


// init function
func init() {

}

func GetDriverName() string {
	return driverName
}

/*
// Create a host using the driver's config
func (d *Driver) Create() error {

}
*/
// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return driverName
}
/*
// GetCreateFlags returns the mcnflag.Flag slice representing the flags
// that can be set, their descriptions and defaults.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {

}

// GetIP returns an IP or hostname that this host is available at
// e.g. 1.2.3.4 or docker-host-d60b70a14d3a.cloudapp.net
func (d *Driver) GetIP() (string, error) {

}

// GetMachineName returns the name of the machine
func (d *Driver) GetMachineName() string {

}

// GetSSHHostname returns hostname for use with ssh
func (d *Driver) GetSSHHostname() (string, error) {

}

// GetSSHKeyPath returns key path for use with ssh
func (d *Driver) GetSSHKeyPath() string {

}

// GetSSHPort returns port for use with ssh
func (d *Driver) GetSSHPort() (int, error) {

}

// GetSSHUsername returns username for use with ssh
func (d *Driver) GetSSHUsername() string {

}

// GetURL returns a Docker compatible host URL for connecting to this host
// e.g. tcp://1.2.3.4:2376
func (d *Driver) GetURL() (string, error) {

}

// GetState returns the state that the host is in (running, stopped, etc)
func (d *Driver) GetState() (state.State, error) {

}

// Kill stops a host forcefully
func (d *Driver) Kill() error {

}

// PreCreateCheck allows for pre-create operations to make sure a driver is ready for creation
func (d *Driver) PreCreateCheck() error {

}

// Remove a host
func (d *Driver) Remove() error {

}

// Restart a host. This may just call Stop(); Start() if the provider does not
// have any special restart behaviour.
func (d *Driver) Restart() error {

}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {

}

// Start a host
func (d *Driver) Start() error {

}

// Stop a host gracefully
func (d *Driver) Stop() error {

}
*/