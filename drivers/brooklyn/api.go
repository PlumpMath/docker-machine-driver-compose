// Copyright (C) 2016-2017 ATOS - All rights reserved.
package brooklyn

import (
	"encoding/json"
	"github.com/apache/brooklyn-client/api/application"
	"github.com/apache/brooklyn-client/net"
	"github.com/docker/machine/libmachine/log"
	"github.com/apache/brooklyn-client/models"
	"fmt"
	"github.com/apache/brooklyn-client/api/locations"
	"errors"
)

const (
	HOST_SSH_ADDRESS_SENSOR = "host.sshAddress"
)

func CatalogByRegex(network *net.Network, regex string) ([]models.CatalogItemSummary, error) {
	url := fmt.Sprintf("/v1/catalog/applications/?regex=%s&allVersions=false", regex)
	var response []models.CatalogItemSummary
	body, err := network.SendGetRequest(url)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	log.Info(len(response))
	return response, err
}

func Delete(network *net.Network, application string) (models.TaskSummary, error) {
	url := fmt.Sprintf("/v1/applications/%s/entities/%s/expunge?release=true",application,application)
	var response models.TaskSummary
	body, err := network.SendEmptyPostRequest(url)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	return response, err
}

func DescendantsSshHostAndPortSensor(network *net.Network, applicationId string) (SshHostAddress, error) {
	sensor, err := application.DescendantsSensor(network,applicationId,HOST_SSH_ADDRESS_SENSOR)
	m := map[string]SshHostAddress{}
	var sshHostAddress SshHostAddress
	if err !=nil {
		return sshHostAddress, err
	}

	err = json.Unmarshal([]byte(sensor), &m)
	if err != nil {
		return sshHostAddress, err
	}
	log.Debug(m)

	for key, _ := range m {
		sshHostAddress = m[key]
		break;
	}
	return sshHostAddress, nil
}

func LocationExists(network *net.Network, locationName string) (string, error) {
	locations, err := locations.LocationList(network)

	var locationId string
	if err != nil {
		return locationId, err
	}

	for _, location := range locations {
		if location.Name == locationName {
			return locationId, nil
		}
	}
	return locationId, errors.New("Location with specified name does not exists.")
}
