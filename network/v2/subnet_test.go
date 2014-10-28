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

package network_test

import (
	"encoding/json"
	"errors"
	"git.openstack.org/stackforge/golang-client.git/network/v2"
	"git.openstack.org/stackforge/golang-client.git/testutil"
	"testing"
)

var allocationPool1 = network.AllocationPool{Start: "13.14.14.15", End: "13.14.14.15"}
var allocationPools = []network.AllocationPool{allocationPool1}
var sampleSubNetResponse = network.SubnetResponse{
	ID:              "subnet",
	Name:            "name",
	NetworkID:       "networkid",
	TenantID:        "tenantid",
	EnableDHCP:      true,
	DNSNameserver:   []string{"13.14.14.15"},
	AllocationPools: allocationPools,
	HostRoutes:      []string{"35.15.15.15"},
	IPVersion:       network.IPV4,
	GatewayIP:       "35.15.15.1",
	CIDR:            "35.15.15.3"}

func TestGetSubnets(t *testing.T) {
	mockResponseObject := subnetsResp{Subnets: []network.SubnetResponse{sampleSubNetResponse}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/subnets")

	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	subnets, err := networkService.Subnets()
	if err != nil {
		t.Error(err)
	}

	if len(subnets) != 1 {
		t.Error(errors.New("Error: Expected 1 subnet to be listed"))
	}
	testUtil.Equals(t, sampleSubNetResponse, subnets[0])
}

func TestGetSubnet(t *testing.T) {
	mockResponseObject := subnetResp{Subnet: sampleSubNetResponse}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/subnets/5270u2tg0")
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	subnet, err := networkService.Subnet("5270u2tg0")
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, sampleSubNetResponse, subnet)
}

func TestDeleteSubnet(t *testing.T) {
	name := "subnetName"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/subnets/"+name)
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	err := networkService.DeleteSubnet(name)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateSubnet(t *testing.T) {
	mockResponse, err := json.Marshal(subnetResp{Subnet: sampleSubNetResponse})

	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, string(mockResponse), "/subnets",
		`{"subnet":{"network_id":"subnetid","ip_version":4,"cidr":"12.14.76.87","allocation_pools":[{"start":"13.14.14.15","end":"13.14.14.15"}]}}`)
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	createSubnetParameters := network.CreateSubnetParameters{NetworkID: "subnetid", AllocationPools: allocationPools, CIDR: "12.14.76.87", IPVersion: 4}
	actualPort, err := networkService.CreateSubnet(createSubnetParameters)
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, sampleSubNetResponse, actualPort)
}

type subnetsResp struct {
	Subnets []network.SubnetResponse `json:"subnets"`
}

type subnetResp struct {
	Subnet network.SubnetResponse `json:"subnet"`
}
