// Copyright (C) 2016-2017 ATOS - All rights reserved.
package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn"

	/*"github.com/docker/machine/libmachine/log"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/api"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/client"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/models"*/
)

func main() {
	plugin.RegisterDriver(brooklyn.NewDriver("", ""))

	// Sample code to demonstrate the REST API calls of brooklyn
	/*client := client.BrooklynClient{
		BaseUrl:  "https://test.compose.canopy-cloud.com",
		User:     "compose.test@canopy-cloud.com", // While running provide user
		Password: "Canopy1!",                      // While running provide password
	}

	request := client.GoRequestWithProxy("http://MC0WBVEC.ww930.my-it-solutions.net:3128")

	versionSummary, err := api.GetServerVersion(request)

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Server Version: ", versionSummary.Version)
	}

	application := models.Application{
		Name:     "SampleApplication",
		Location: "AWS Frankfurt",
		Type:     "com.canopy.compose.centos:1.3",
	}
	taskSummary, err := api.CreateApplication(request, application)

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Application Id:", taskSummary.EntityId)
	}*/
}
