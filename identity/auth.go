//Package identity provides functions for client-side access to OpenStack
//IdentityService.
package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"io/ioutil"
	"strings"
	"time"
)

type Auth struct {
	Access Access
}

type Access struct {
	Token          Token
	User           User
	ServiceCatalog []Service
}

type Token struct {
	Id      string
	Expires time.Time
	Tenant  Tenant
}

type Tenant struct {
	Id   string
	Name string
}

type User struct {
	Id          string
	Name        string
	Roles       []Role
	Roles_links []string
}

type Role struct {
	Id       string
	Name     string
	TenantId string
}

type Service struct {
	Name            string
	Type            string
	Endpoints       []Endpoint
	Endpoints_links []string
}

type Endpoint struct {
	TenantId    string
	PublicURL   string
	InternalURL string
	Region      string
	VersionId   string
	VersionInfo string
	VersionList string
}

func AuthKey(url, accessKey, secretKey string) (Auth, error) {
	jsonStr := (fmt.Sprintf(`{"auth":{
		"apiAccessKeyCredentials":{"accessKey":"%s","secretKey":"%s"}}
		}`,
		accessKey, secretKey))
	return auth(&url, &jsonStr)
}

func AuthKeyTenantId(url, accessKey, secretKey, tenantId string) (Auth, error) {
	jsonStr := (fmt.Sprintf(`{"auth":{
		"apiAccessKeyCredentials":{"accessKey":"%s","secretKey":"%s"},"tenantId":"%s"}
		}`,
		accessKey, secretKey, tenantId))
	return auth(&url, &jsonStr)
}

func AuthUserName(url, username, password string) (Auth, error) {
	jsonStr := (fmt.Sprintf(`{"auth":{
		"passwordCredentials":{"username":"%s","password":"%s"}}
		}`,
		username, password))
	return auth(&url, &jsonStr)
}

func AuthUserNameTenantName(url, username, password, tenantName string) (Auth, error) {
	jsonStr := (fmt.Sprintf(`{"auth":{
		"passwordCredentials":{"username":"%s","password":"%s"},"tenantName":"%s"}
		}`,
		username, password, tenantName))
	return auth(&url, &jsonStr)
}

func AuthUserNameTenantId(url, username, password, tenantId string) (Auth, error) {
	jsonStr := (fmt.Sprintf(`{"auth":{
		"passwordCredentials":{"username":"%s","password":"%s"},"tenantId":"%s"}
		}`,
		username, password, tenantId))
	return auth(&url, &jsonStr)
}

func AuthTenantNameTokenId(url, tenantName, tokenId string) (Auth, error) {
	jsonStr := (fmt.Sprintf(`{"auth":{
		"tenantName":"%s","token":{"id":"%s"}}
		}`,
		tenantName, tokenId))
	return auth(&url, &jsonStr)
}

func auth(url, jsonStr *string) (Auth, error) {
	var s []byte = []byte(*jsonStr)
	resp, err := misc.CallAPI("POST", *url, &s,
		"Accept-Encoding", "gzip,deflate",
		"Accept", "application/json",
		"Content-Type", "application/json",
		"Content-Length", string(len(*jsonStr)))
	if err != nil {
		return Auth{}, err
	}
	if err = misc.CheckHttpResponseStatusCode(resp); err != nil {
		return Auth{}, err
	}
	var contentType string = strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(contentType, "json") != true {
		return Auth{}, errors.New("err: header Content-Type is not JSON")
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Auth{}, err
	}
	var auth = Auth{}
	if err = json.Unmarshal(body, &auth); err != nil {
		return Auth{}, err
	}
	return auth, nil
}
