// Copyright (C) 2016-2017 ATOS - All rights reserved.
package api

import (
	"encoding/json"

	"errors"

	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/client"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/models"
)

func GetServerVersion(request *client.BrooklynAgent) (models.VersionSummary, error) {
	resp, body, errs := request.Get(request.BaseUrl + "/v1/server/version").End()
	var versionSummary models.VersionSummary
	if errs != nil {
		return versionSummary, errs[0]
	}
	if resp.StatusCode != 200 {
		return versionSummary, errors.New(resp.Status)
	}
	err := json.Unmarshal([]byte(body), &versionSummary)
	return versionSummary, err
}
