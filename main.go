// Copyright (C) 2016-2017 ATOS - All rights reserved.
package main

import (
	"fmt"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn"
)

func main() {
	fmt.Println(brooklyn.GetDriverName())
	fmt.Println("Brooklyn docker machine driver implementation work in progress")
}
