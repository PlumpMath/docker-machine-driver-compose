package api

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/client"
	"stash.fsc.atos-services.net/scm/cet/bdmd.git/drivers/brooklyn/models"
)

func GetServerVersion(client client.BrooklynClient) error {
	request := gorequest.New(). //Proxy("http://MC0WBVEC.ww930.my-it-solutions.net:3128").
					SetBasicAuth(client.User, client.Password)
	_, body, errs := request.Get(client.BaseUrl + "/v1/server/version").End()

	if errs != nil {
		return errs[0]
	}

	var versionSummary models.VersionSummary
	err := json.Unmarshal([]byte(body), &versionSummary)
	return err
}
