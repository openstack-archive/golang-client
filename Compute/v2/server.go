// server.go
package compute

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type serversResp struct {
	ServerInfos []ServerInfo `json:"servers"`
}

type serversDetailResp struct {
	ServerInfoDetails []ServerInfoDetail `json:"servers"`
}

type serverDetailResp struct {
	ServerInfoDetail ServerInfoDetail `json:"server"`
}

type ServerInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Links []Link `json:"links"`
}

type ServerInfoDetail struct {
	Id               string               `json:"id"`
	Name             string               `json:"name"`
	Status           string               `json:"status"`
	Created          *time.Time           `json:"created"`
	Updated          *time.Time           `json:"updated"`
	HostId           string               `json:"hostId"`
	Addresses        map[string][]Address `json:"addresses"`
	Links            []Link               `json:"links"`
	Image            Image                `json:"image"`
	Flavor           Flavor               `json:"flavor"`
	TaskState        string               `json:"OS-EXT-STS:task_state"`
	VMState          string               `json:"OS-EXT-STS:vm_state"`
	PowerState       int                  `json:"OS-EXT-STS:power_state"`
	AvailabilityZone string               `json:"OS-EXT-AZ:availability_zone:"`
	UserId           string               `json:"user_id"`
	TenantId         string               `json:"tenant_id"`
	AccessIPv4       string               `json:"accessIPv4"`
	AccessIPv6       string               `json:"accessIPv6"`
	ConfigDrive      string               `json:"config_drive"`
	Progress         int                  `json:"progress"`
	MetaData         map[string]string    `json:"metadata"`
	AdminPass        string               `json:"adminPass"`
}

type Link struct {
	HRef string `json:"href"`
	Rel  string `json:"rel"`
}

type Image struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}

type Flavor struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}

type SecurityGroups struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

type SecurityGroup struct {
	Name string `json:"name"`
}

type Address struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
	Type    string `json:"OS-EXT-IPS:type"`
	MacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
}

type ByName []ServerInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func GetServerInfos(url string, token string) (serverInfos []ServerInfo, err error) {

	_, body, err := misc.CallGetAPI(url+"/servers", "X-Auth-Token", token)
	if err != nil {
		return nil, err
	}

	var sr = serversResp{}
	if err = json.Unmarshal([]byte(body), &sr); err != nil {
		return
	}

	serverInfos = sr.ServerInfos
	err = nil
	return
}

func GetServerInfoDetails(url string, token string) (serverInfoDetails []ServerInfoDetail, err error) {

	_, body, err := misc.CallGetAPI(url+"/servers/detail", "X-Auth-Token", token)
	if err != nil {
		return nil, err
	}

	var sr = serversDetailResp{}
	if err = json.Unmarshal(body, &sr); err != nil {
		return
	}

	serverInfoDetails = sr.ServerInfoDetails
	err = nil
	return
}

func GetServerInfoDetail(url string, token string, id string) (serverInfoDetail ServerInfoDetail, err error) {
	reqUrl := fmt.Sprintf("%s/servers/%s", url, id)
	_, body, err := misc.CallGetAPI(reqUrl, "X-Auth-Token", token)
	if err != nil {
		return serverInfoDetail, err
	}

	serverDetailResp := serverDetailResp{}
	if err = json.Unmarshal(body, &serverDetailResp); err != nil {
		return
	}

	serverInfoDetail = serverDetailResp.ServerInfoDetail
	err = nil
	return
}

func DeleteServer(url string, token string, id string) (err error) {
	reqUrl := fmt.Sprintf("%s/servers/%s", url, id)
	err = misc.CallDeleteAPI(reqUrl, "X-Auth-Token", token)

	if err != nil {
		return err
	}

	err = nil
	return
}

type ServerNetworkParameters struct {
	Uuid string `json:"uuid"`
	Port string `json:"port"`
}

type ServerServiceParameters struct {
	Name          string                    `json:"name"`
	ImageRef      string                    `json:"imageRef"`
	SSHKey        string                    `json:"sshkey"`
	FlavorRef     int32                     `json:"flavorRef"`
	MaxCount      int32                     `json:"maxcount"`
	MinCount      int32                     `json:"mincount"`
	UserData      string                    `json:"userdata"`
	Networks      []ServerNetworkParameters `json:"networks"`
	SecurityGroup []SecurityGroup           `json:"securitygroups"`
}

type ServerRequestInfo map[string]ServerServiceParameters

type ServerService struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

func CreateServer(client *http.Client, ss ServerService, serverRequestInfo ServerRequestInfo) (serverInfoDetail ServerInfoDetail, err error) {

	reqBody, err := json.Marshal(serverRequestInfo)
	if err != nil {
		return serverInfoDetail, err
	}

	reqUrl := ss.Url + "/servers"

	req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(reqBody))
	if err != nil {
		return serverInfoDetail, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", ss.Token)

	misc.LogDebug("CreateServer-----------------------------------httputil.DumpRequestOut-------BEGIN")
	dumpReqByte, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Printf(err.Error())
	}
	misc.LogDebug(string(dumpReqByte))
	misc.LogDebug("CreateServer-----------------------------------httputil.DumpRequestOut-------END")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf(err.Error())
	}

	misc.LogDebug("CreateServer-----------------------------------httputil.DumpResponse-------BEGIN")
	dumpRspByte, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf(err.Error())
	}
	misc.LogDebug(string(dumpRspByte))
	misc.LogDebug("CreateServer-----------------------------------httputil.DumpResponse-------END")

	if err != nil {
		return serverInfoDetail, err
	}

	err = json.NewDecoder(resp.Body).Decode(&serverInfoDetail)
	defer resp.Body.Close()

	if !(resp.StatusCode == 201 || resp.StatusCode == 202) {
		err = errors.New(fmt.Sprintf("Error: status code != 201 or 202, object not created. Status: (%s), reqUrl: %s, reqBody: %s, resp: %s, respBody: %s",
			resp.Status, reqUrl, reqBody, resp, resp.Body))

		log.Printf(err.Error())
		return
	}

	err = nil
	return
}

func ServerAction(url string, token string, id string, action string, key string, value string) (err error) {
	var reqBody = []byte(fmt.Sprintf(`
	{
		"%s": 
		{
			"%s": "%s"
		}
	}`, action, key, value))

	resp, err := misc.CallAPI("POST", url+"/servers/"+id+"/action", &reqBody, "X-Auth-Token", token)

	if err = misc.CheckHttpResponseStatusCode(resp); err != nil {
		return
	}

	if err != nil {
		return err
	}

	err = nil
	return
}
