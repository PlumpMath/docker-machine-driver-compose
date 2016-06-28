// Copyright (C) 2016-2017 ATOS - All rights reserved.
package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/compose"
)

func main() {
	plugin.RegisterDriver(compose.NewDriver("", ""))
}
