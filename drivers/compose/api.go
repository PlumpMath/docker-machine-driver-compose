// Copyright (C) 2016-2017 ATOS - All rights reserved.
package compose

import (
	"encoding/json"
	"github.com/apache/brooklyn-client/api/application"
	"github.com/apache/brooklyn-client/net"
	"github.com/docker/machine/libmachine/log"
	"github.com/apache/brooklyn-client/models"
	"github.com/apache/brooklyn-client/api/entity_effectors"
	"fmt"
	"github.com/apache/brooklyn-client/api/locations"
	"errors"
	"github.com/apache/brooklyn-client/api/entity_sensors"
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

func GetNodeId(network *net.Network, applicationId string) (string, error) {
	sensorInfo, err := application.DescendantsSensor(network, applicationId, HOST_ADDRESS_SENSOR_NAME)
	var nodeId string

	m := map[string]string{}
	if err !=nil {
		return nodeId, err
	}

	err = json.Unmarshal([]byte(sensorInfo), &m)
	if err != nil {
		return nodeId, err
	}
	for key, _ := range m {
		log.Info("Key: ", key)
		nodeId = key
		break
	}
	return nodeId, nil
}

func GetNodeState(network *net.Network, applicationId, entityId string) (string, error) {
	serviceState, err := entity_sensors.SensorValue(network,applicationId,entityId, SERVICE_STATE_SENSOR_NAME)

	if err != nil {
		return "", err
	} else if state, ok := serviceState.(string); ok {
		return state, nil
	}

	return "UNKNOWN", nil
}

func DescendantsSensor(network *net.Network, applicationId string, sensor string) (map[string]int, error) {
	sensor, err := application.DescendantsSensor(network,applicationId,sensor)
	m := map[string]int{}
	if err !=nil {
		return m, err
	}

	err = json.Unmarshal([]byte(sensor), &m)
	if err != nil {
		return m, err
	}
	log.Debug(m)
	return m, nil
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

func TriggerStart(network *net.Network, applicationId string, entityId string) error {
	params := []string{}
	args := []string{}
	_, err := entity_effectors.TriggerEffector(network,applicationId,entityId, "start", params, args)
	return err
}

func TriggerStop(network *net.Network, applicationId string, entityId string) error {
	params := []string{"stopProcessMode","stopMachineMode"}
	args := []string{"ALWAYS","NEVER"}
	_, err := entity_effectors.TriggerEffector(network,applicationId,entityId, "stop", params, args)
	return err
}

func TriggerRestart(network *net.Network, applicationId string, entityId string) error {
	params := []string{"restartChildren","restartMachine"}
	args := []string{"true","false"}
	_, err := entity_effectors.TriggerEffector(network,applicationId,entityId, "restart", params, args)
	return err
}