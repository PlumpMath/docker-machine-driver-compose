// Copyright (C) 2016-2017 ATOS - All rights reserved.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/docker/machine/libmachine/log"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/client"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/models"
	"text/template"
)

func CreateApplication(request *client.BrooklynAgent, application models.Application) (models.TaskSummary, error) {
	// Define the template
	const applicationTmpl = `name: {{.Name}}
location: {{.Location}}
services:
  - type: {{.Type}}
`

	// Create a new template and parse the application into it.
	t := template.Must(template.New("application").Parse(applicationTmpl))

	var appYml bytes.Buffer
	err := t.Execute(&appYml, application)
	log.Info(appYml.String())
	var taskSummary models.TaskSummary
	if err != nil {
		return taskSummary, err
	}

	resp, body, errs := request.Post(request.BaseUrl+"/v1/applications").
		Set("Content-Type", "text/plain").
		Set("Accept", "application/json").
		SetDebug(true).
		Send(appYml.String()).
		End()

	if errs != nil {
		return taskSummary, errs[0]
	}

	if resp.StatusCode != 201 {
		return taskSummary, errors.New(resp.Status)
	}

	log.Info(resp)
	err = json.Unmarshal([]byte(body), &taskSummary)
	return taskSummary, err
}
