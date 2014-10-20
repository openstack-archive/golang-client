// AlsTool project openstackconfig.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/identity/v2"
	"github.com/parnurzeal/gorequest"
	"log"
	"os"
	"time"
)

type OpenStackConfig struct {
	AuthUrl        string
	TenantId       string
	TenantName     string
	Username       string
	Password       string
	DebugOpenStack bool
}

type Links struct {
	Links []Link `json:"links"`
}

type Link struct {
	href string `json:"href"`
	rel  string `json:"rel"`
}

func (config OpenStackConfig) Log() {
	log.Printf("%-20s - %s\n", "OS_AUTH_URL", config.AuthUrl)
	log.Printf("%-20s - %s\n", "OS_TENANT_ID", config.TenantId)
	log.Printf("%-20s - %s\n", "OS_TENANT_NAME", config.TenantName)
	log.Printf("%-20s - %s\n", "OS_USERNAME", config.Username)
}

func InitializeFromEnv() (config OpenStackConfig, err error) {

	var c = OpenStackConfig{}

	c.AuthUrl = os.Getenv("OS_AUTH_URL")
	c.TenantId = os.Getenv("OS_TENANT_ID")
	c.TenantName = os.Getenv("OS_TENANT_NAME")
	c.Username = os.Getenv("OS_USERNAME")
	c.Password = os.Getenv("OS_PASSWORD")

	if len(c.AuthUrl) == 0 {
		err = errors.New("Error: no authentication URL specified")
		return
	}
	if len(c.Username) == 0 {
		err = errors.New("Error: no username specified")
		return
	}
	if len(c.Password) == 0 {
		err = errors.New("Error: no password specified")
		return
	}
	if len(c.TenantName) == 0 {
		err = errors.New("Error: no tenant name specified")
		return
	}
	if len(c.TenantId) == 0 {
		err = errors.New("Error: no tenant ID specified")
		return
	}

	c.DebugOpenStack = false
	if len(os.Getenv("DebugOpenStack")) != 0 {
		c.DebugOpenStack = true
	}

	config = c
	err = nil
	return
}

func Authenticate(openStackConfig OpenStackConfig) (auth identity.Auth, err error) {

	req := gorequest.New()

	reqUrl := fmt.Sprintf("%s/tokens", openStackConfig.AuthUrl)

	reqBody := fmt.Sprintf(`
	{"auth":
		{
			"passwordCredentials":
			{
				"username":"%s",
				"password":"%s"
			},
			"tenantName":"%s"
		}
	}`, openStackConfig.Username, openStackConfig.Password, openStackConfig.TenantName)

	_, body, errs := req.Post(reqUrl).
		Set(`Accept-Encoding`, `gzip,deflate`).
		Set(`Accept`, `application/json`).
		Set(`Content-Type`, `application/json`).
		Send(reqBody).
		End()

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	if err = json.Unmarshal([]byte(body), &auth); err != nil {
		return
	}

	if !auth.Access.Token.Expires.After(time.Now()) {
		err = errors.New("Error: The AuthN token is expired")
		return
	}

	err = nil
	return
}
