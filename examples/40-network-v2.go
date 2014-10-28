// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package main

import (
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/network/v2"
	"net/http"
	"reflect"
	"time"
)

func main() {
	config := getConfig()

	// Before working with object storage we need to authenticate with a project
	// that has active object storage.
	auth, err := identity.AuthUserNameTenantName(config.Host,
		config.Username,
		config.Password,
		config.ProjectName)
	if err != nil {
		panicString := fmt.Sprint("There was an error authenticating:", err)
		panic(panicString)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		panic("There was an error. The auth token has an invalid expiration.")
	}

	// Find the endpoint for networke.
	url := ""
	for _, svc := range auth.Access.ServiceCatalog {
		if svc.Type == "network" {
			for _, ep := range svc.Endpoints {
				if ep.VersionId == "2.0" {
					url = ep.PublicURL + "/"
					break
				}
			}
		}
	}
	if url == "" {
		panic("network url not found during authentication")
	}

	// this is Helion public cloud workaround. This is terrible
	// that the actual full link with version is not given
	// unsure if there is another workaround.
	finalUrl := url + "v2.0"

	if finalUrl == "" {
		panic("v2 network url not found during endpoint lookup")
	}

	networkService := network.Service{TokenID: auth.Access.Token.Id, TenantID: auth.Access.Token.Tenant.Id,
		Client: *http.DefaultClient, URL: finalUrl}

	networkName := "OtherNetwork"
	activeNetwork := CreateNewNetworkVerifyExistsAndActive(networkService, networkName)

	// Create list, get and delete a subnet
	activeSubnet := CreateSubnetAndVerify(networkService, activeNetwork)
	activePort := CreatePortAndVerify(networkService, activeSubnet, activeNetwork)
	DeletePortAndVerify(networkService, activePort)
	DeleteSubnetAndVerify(networkService, activeSubnet)
	DeleteNetworkAndVerify(networkService, activeNetwork)
}

func CreatePortAndVerify(networkService network.Service, subnet network.SubnetResponse, activeNetwork network.Response) network.PortResponse {
	var portToCreate = network.CreatePortParameters{AdminStateUp: false, Name: "testPort", NetworkID: activeNetwork.ID}

	portCreated, err := networkService.CreatePort(portToCreate)
	if err != nil {
		panicString := fmt.Sprint("Error in creating port:", err)
		panic(panicString)
	}

	foundPorts, err := networkService.Ports()
	if err != nil {
		panicString := fmt.Sprint("Error in getting list of subnets:", err)
		panic(panicString)
	}

	foundCreatedPort := false
	for _, portFound := range foundPorts {
		if reflect.DeepEqual(portCreated, portFound) {
			foundCreatedPort = true
		}
	}

	if !foundCreatedPort {
		panic("Cannot find the newly created port.")
	}

	return portCreated
}

func DeletePortAndVerify(networkService network.Service, port network.PortResponse) {
	err := networkService.DeletePort(port.ID)
	if err != nil {
		panicString := fmt.Sprint("Error in deleting port:", err)
		panic(panicString)
	}

	ports, err := networkService.Ports()
	if err != nil {
		panicString := fmt.Sprint("There was an error getting a list of ports to verify the port was deleted:", err)
		panic(panicString)
	}

	portDeleted := true
	for _, portFound := range ports {
		if reflect.DeepEqual(portFound, port) {
			portDeleted = false
		}
	}

	if !portDeleted {
		panic("port was not deleted.")
	}
}

func CreateSubnetAndVerify(networkService network.Service, activeNetwork network.Response) network.SubnetResponse {
	var allocationPools = []network.AllocationPool{network.AllocationPool{Start: "10.1.2.5", End: "10.1.2.15"}}
	subnetToCreate := network.CreateSubnetParameters{NetworkID: activeNetwork.ID, IPVersion: network.IPV4,
		CIDR: "10.1.2.1/25", AllocationPools: allocationPools}

	subnetCreated, err := networkService.CreateSubnet(subnetToCreate)
	if err != nil {
		panicString := fmt.Sprint("Error in creating subnet:", err)
		panic(panicString)
	}

	foundSubnets, err := networkService.Subnets()
	if err != nil {
		panicString := fmt.Sprint("Error in getting list of subnets:", err)
		panic(panicString)
	}

	foundCreatedSubnet := false
	for _, subnetFound := range foundSubnets {
		if reflect.DeepEqual(subnetCreated, subnetFound) {
			foundCreatedSubnet = true
		}
	}

	if !foundCreatedSubnet {
		panic("Cannot find the newly created subnet.")
	}

	return subnetCreated
}

func DeleteSubnetAndVerify(networkService network.Service, subnet network.SubnetResponse) {
	err := networkService.DeleteSubnet(subnet.ID)
	if err != nil {
		panicString := fmt.Sprint("Error in deleting subnet:", err)
		panic(panicString)
	}

	subnets, err := networkService.Subnets()
	if err != nil {
		panicString := fmt.Sprint("There was an error getting a list of networks to verify the network was deleted:", err)
		panic(panicString)
	}

	subnetDeleted := true
	for _, subnetFound := range subnets {
		if reflect.DeepEqual(subnetFound, subnet) {
			subnetDeleted = false
		}
	}

	if !subnetDeleted {
		panic("subnet was not deleted.")
	}
}

func DeleteNetworkAndVerify(networkService network.Service, activeNetwork network.Response) {
	err := networkService.DeleteNetwork(activeNetwork.ID)
	if err != nil {
		panicString := fmt.Sprint("Error in deleting 'OtherNetwork'", err)
		panic(panicString)
	}

	networks, err := networkService.Networks()
	if err != nil {
		panicString := fmt.Sprint("There was an error getting a list of networks to verify the network was deleted:", err)
		panic(panicString)
	}

	networkDeleted := true
	for _, networkFound := range networks {
		if reflect.DeepEqual(activeNetwork, networkFound) {
			networkDeleted = false
		}
	}

	if !networkDeleted {
		panic("network was not deleted.")
	}

}
func CreateNewNetworkVerifyExistsAndActive(networkService network.Service, networkName string) network.Response {
	createdNetwork, err := networkService.CreateNetwork(true, networkName, false)
	if err != nil {
		panicString := fmt.Sprint("There was an error creating a network:", err)
		panic(panicString)
	}

	networks, err := networkService.Networks()
	if err != nil {
		panicString := fmt.Sprint("There was an error getting a list of networks:", err)
		panic(panicString)
	}

	foundCreatedNetwork := false
	for _, networkFound := range networks {
		if reflect.DeepEqual(createdNetwork, networkFound) {
			foundCreatedNetwork = true
		}
	}

	if !foundCreatedNetwork {
		panic("Cannot find network called 'OtherNetwork' when getting a list of networks.")
	}

	// Might be nice to have some sugar api that can do this easily for a developer...
	// Keep iterating until active or until more than 5 tries has been exceeded.
	//numTries := 0
	//activeNetwork := createdNetwork
	//for activeNetwork.Status != "ACTIVE" && numTries < 5 {
	//	activeNetwork, _ = networkService.Network(createdNetwork.ID)
	//	numTries++
	//	fmt.Println("Sleeping 50ms on try:" + string(numTries) + " with status currently " + activeNetwork.Status)
	//	sleepDuration, _ := time.ParseDuration("50ms")
	//	time.Sleep(sleepDuration)
	//}

	//if activeNetwork.Status != "ACTIVE" {
	//	panic("The network is not in the active state and cannot be deleted")
	//}

	return createdNetwork
}
