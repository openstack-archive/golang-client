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

var fixedIps = []network.FixedIP{network.FixedIP{SubnetID: "125071206726"}, network.FixedIP{SubnetID: "15615526", IPAddress: "132.145.15.11"}}
var samplePortResponse = network.PortResponse{
	ID:                  "2490640zbg",
	Name:                "Port14",
	Status:              "active",
	AdminStateUp:        true,
	PortSecurityEnabled: true,
	DeviceID:            "deviceid",
	DeviceOwner:         "deviceowner",
	NetworkID:           "networkid",
	TenantID:            "tenantid",
	MacAddress:          "macAddress",
	FixedIPs:            fixedIps,
	SecurityGroups:      []string{"sec1", "sec2"}}

func TestGetPorts(t *testing.T) {
	mockResponseObject := portsContainer{Ports: []network.PortResponse{samplePortResponse}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/ports")
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	ports, err := networkService.Ports()
	if err != nil {
		t.Error(err)
	}

	if len(ports) != 1 {
		t.Error(errors.New("Error: Expected 2 networks to be listed"))
	}
	testUtil.Equals(t, samplePortResponse, ports[0])
}

func TestGetPort(t *testing.T) {
	mockResponseObject := portContainer{Port: samplePortResponse}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/ports/23507256")
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	port, err := networkService.Port("23507256")
	if err != nil {
		t.Error(err)
	}
	testUtil.Equals(t, samplePortResponse, port)
}

func TestDeletePort(t *testing.T) {
	portName := "portName"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/ports/"+portName)
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	err := networkService.DeletePort(portName)
	if err != nil {
		t.Error(err)
	}
}

func TestCreatePort(t *testing.T) {
	mockResponse, _ := json.Marshal(portContainer{Port: samplePortResponse})
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, string(mockResponse), "/ports",
		`{"port":{"admin_state_up":true,"name":"name","network_id":"networkid","fixed_ips":[{"subnet_id":"125071206726"},{"subnet_id":"15615526","ip_address":"132.145.15.11"}]}}`)
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	createPortParameters := network.CreatePortParameters{AdminStateUp: true, Name: "name", NetworkID: "networkid", FixedIPs: fixedIps}
	actualPort, err := networkService.CreatePort(createPortParameters)
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, samplePortResponse, actualPort)
}

type portsContainer struct {
	Ports []network.PortResponse `json:"ports"`
}

type portContainer struct {
	Port network.PortResponse `json:"port"`
}
