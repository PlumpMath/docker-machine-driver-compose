// Copyright (C) 2016-2017 ATOS - All rights reserved.
package compose

import (
	"testing"

	"github.com/apache/brooklyn-client/cli/net"
	"github.com/docker/machine/libmachine/log"
)

var (
	network = net.NewNetwork("http://217.115.71.184:5550", "compose", "Canopy1!", false)
)

func TestDelete(t *testing.T) {
	sshHostAddress, err := DescendantsSSHHostAndPortSensor(network, "s0ZNhmV9")

	if err != nil {
		t.Fail()
	}

	log.Info(sshHostAddress)
}

func TestSensor(t *testing.T) {
	sshHostAddress, err := DescendantsSensor(network, "sdpxTJF2", MappedPortSensorName)

	if err != nil {
		t.Fail()
	}

	log.Info(sshHostAddress)
}

func TestCatalogByRegex(t *testing.T) {
	catalogs, err := CatalogByRegex(network, "com.canopy.compose.ubuntu")
	log.Info(catalogs)
	if err != nil || len(catalogs) <= 0 {
		t.Fail()
	}
}
