// AlsTool project main.go
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/Compute/v2"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/network/v2"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var (
	configFile  string
	authUrl     string
	tenantId    string
	tenantName  string
	username    string
	password    string
	uninstall   bool
	cloudConfig bool
	install     bool
	status      bool
)

type Config struct {
	Nodes     map[string]Node `yaml:"hosts"`
	ImageRef  string          `yaml:"imageRef"`
	SSHKey    string          `yaml:"sshkey"`
	FlavorRef string          `json:"flavorref"`
	MaxCount  string          `json:"maxcount"`
	MinCount  string          `json:"mincount"`
	UserData  string          `json:"userdata"`
}

type Node struct {
	IP       string `yaml:"ip"`
	IsMaster bool   `yaml:"ismaster"`
	ServerId string
}

func init() {
	flag.StringVar(&configFile, "c", "config.yml", "ALS tool configuration file")
	flag.StringVar(&authUrl, "a", "", "OpenStack authentication URL (OS_AUTH_URL)")
	flag.StringVar(&tenantId, "i", "", "OpenStack tenant id (OS_TENANT_ID)")
	flag.StringVar(&tenantName, "n", "", "OpenStack tenant name (OS_TENANT_NAME)")
	flag.StringVar(&username, "u", "", "OpenStack user name (OS_USERNAME)")
	flag.StringVar(&password, "p", "", "OpenStack passsword (OS_PASSWORD)")
	flag.BoolVar(&uninstall, "U", false, "Uninstall cluster (defined in config file -c)")
	flag.BoolVar(&cloudConfig, "C", false, "Create CloudConfig files for cluster")
	flag.BoolVar(&install, "I", false, "Install cluster (defined in config file -c)")
	flag.BoolVar(&status, "S", false, "Cluster status (defined in config file -c)")
}

func main() {
	flag.Parse()

	config, err := readConfigFile(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	config.Log()

	openStackConfig, err := InitializeFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}
	updateConfigFromCommandLine(&openStackConfig)
	openStackConfig.Log()

	auth, err := Authenticate(openStackConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	token := auth.Access.Token
	log.Printf("%-20s - %s\n", "token:", token.Id)

	endpointList := auth.EndpointList()

	subnets, err := network.GetSubnets(endpointList["network"], token.Id)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	ports, err := network.GetPorts(endpointList["network"], token.Id)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}
	sort.Sort(network.ByName(ports))

	servers, err := compute.GetServerInfos(endpointList["compute"], token.Id)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}
	sort.Sort(compute.ByName(servers))

	// uninstall/cleanup
	for _, v := range servers {

		if _, ok := config.Nodes[v.Name]; ok {

			log.Printf("%-20s - %s\n", "delete server", v.Name)

			err := compute.DeleteServer(endpointList["compute"], token.Id, v.Id)
			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}

			log.Printf("%-20s - %s %s\n", "delete server", v.Name, "COMPLETED")
		}
	}

	for _, v := range ports {

		if _, ok := config.Nodes[v.Name]; ok {

			log.Printf("%-20s - %s\n", "delete port", v.Name)

			err := network.DeletePort(endpointList["network"], token.Id, v.Id)
			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}

			log.Printf("%-20s - %s %s\n", "delete port", v.Name, "COMPLETED")
		}
	}

	flavorRef, err := strconv.Atoi(config.FlavorRef)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	maxCount, err := strconv.Atoi(config.MaxCount)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	minCount, err := strconv.Atoi(config.MinCount)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	//install
	for hostname, hostnode := range config.Nodes {

		log.Printf("%-20s - %s %s\n", "create port", hostname, hostnode.IP)
		port, err := network.CreatePort(endpointList["network"], token.Id, hostname, hostnode.IP, subnets[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}
		log.Printf("%-20s - %s %s\n", "create port", port.Id, "COMPLETED")

		log.Printf("%-20s - %s %s\n", "create server", hostname, hostnode.IP)

		serverService := compute.ServerService{
			endpointList["compute"],
			token.Id}

		var serverServiceParameters = compute.ServerServiceParameters{
			hostname,
			config.ImageRef,
			config.SSHKey,
			int32(flavorRef),
			int32(maxCount),
			int32(minCount),
			config.UserData,
			[]compute.ServerNetworkParameters{
				compute.ServerNetworkParameters{port.NetworkId, port.Id}},
			[]compute.SecurityGroup{
				compute.SecurityGroup{"default"}}}

		var serverRequestInfo = make(compute.ServerRequestInfo)
		serverRequestInfo["server"] = serverServiceParameters

		result, err := compute.CreateServer(new(http.Client), serverService, serverRequestInfo)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		marshaledResult, err := json.Marshal(result)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}
		misc.LogDebug("Result = " + string(marshaledResult))

		log.Printf("%-20s - %s %s %s\n", "create server", hostname, hostnode.IP, "COMPLETED")
	}

	err = nil
	os.Exit(1)
}

func updateConfigFromCommandLine(config *OpenStackConfig) {

	if len(authUrl) > 0 {
		config.AuthUrl = authUrl
	}
	if len(tenantId) > 0 {
		config.TenantId = tenantId
	}
	if len(tenantName) > 0 {
		config.TenantName = tenantName
	}
	if len(username) > 0 {
		config.Username = username
	}
	if len(password) > 0 {
		config.Password = password
	}
}

func readConfigFile(filename string) (config Config, err error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return
	}

	err = nil
	return
}

func (config Config) Log() {
	for k, v := range config.Nodes {
		log.Printf("%-20s - %s %v\n", "config file", k, v)
	}
	log.Printf("%-20s - %s %s\n", "config file", "ImageRef", config.ImageRef)
	log.Printf("%-20s - %s %s\n", "config file", "SSHKey", config.SSHKey)
	log.Printf("%-20s - %s %s\n", "config file", "FlavorRef", config.FlavorRef)
	log.Printf("%-20s - %s %s\n", "config file", "MaxCount", config.MaxCount)
	log.Printf("%-20s - %s %s\n", "config file", "MinCount", config.MinCount)
	log.Printf("%-20s - %s %s\n", "config file", "UserData", config.UserData)
}

func getUserData(filename string) (encodedStr string, err error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	encodedStr = base64.StdEncoding.EncodeToString(b)
	err = nil
	return
}
