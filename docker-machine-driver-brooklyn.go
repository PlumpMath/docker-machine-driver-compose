// Copyright (C) 2016-2017 ATOS - All rights reserved.
package main

import (
	"fmt"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/api"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/client"
)

func main() {
	fmt.Println(brooklyn.GetDriverName())
	fmt.Println("Brooklyn docker machine driver implementation work in progress")

	// Sample code to demonstrate the REST API calls of brooklyn
	client := client.BrooklynClient{
		BaseUrl:  "https://test.compose.canopy-cloud.com",
		User:     "", // While running provide user
		Password: "", // While running provide password
	}
	err := api.GetServerVersion(client)

	if err != nil {
		fmt.Println(err)
	}
}
