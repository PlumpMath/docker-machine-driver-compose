package client

import "github.com/parnurzeal/gorequest"

// method to send request and response

type BrooklynClient struct {
	BaseUrl  string
	User     string
	Password string
}

type BrooklynAgent struct {
	*gorequest.SuperAgent
	BaseUrl string
}

func (client *BrooklynClient) GoRequest() *BrooklynAgent {
	request := gorequest.New().
	SetBasicAuth(client.User, client.Password)
	return &BrooklynAgent{
		request,
		client.BaseUrl,
	}
}

func (client *BrooklynClient) GoRequestWithProxy(proxy string) *BrooklynAgent {
	request := gorequest.New().Proxy(proxy).
	SetBasicAuth(client.User, client.Password)
	return &BrooklynAgent{
		request,
		client.BaseUrl,
	}
}
