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

package network

import (
	// "git.openstack.org/openstack/golang-client.git/openstack"
)

// PortResponse returns a set of values of the a port response.
type PortResponse struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Status              string    `json:"status"`
	AdminStateUp        bool      `json:"admin_state_up"`
	PortSecurityEnabled bool      `json:"port_security_enabled"`
	DeviceID            string    `json:"device_id"`
	DeviceOwner         string    `json:"device_owner"`
	NetworkID           string    `json:"network_id"`
	TenantID            string    `json:"tenant_id"`
	MacAddress          string    `json:"mac_address"`
	FixedIPs            []FixedIP `json:"fixed_ips"`
	SecurityGroups      []string  `json:"security_groups"`
}

// CreatePortParameters holds a set of values that specify how
// to create a new port.
type CreatePortParameters struct {
	AdminStateUp   bool      `json:"admin_state_up,omitempty"`
	Name           string    `json:"name,omitempty"`
	NetworkID      string    `json:"network_id"`
	DeviceID       string    `json:"device_id,omitempty"`
	MacAddress     string    `json:"mac_address,omitempty"`
	FixedIPs       []FixedIP `json:"fixed_ips,omitempty"`
	SecurityGroups []string  `json:"security_groups,omitempty"`
}

// PortResponses is a type for a slice of PortResponses.
type PortResponses []PortResponse

func (a PortResponses) Len() int           { return len(a) }
func (a PortResponses) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PortResponses) Less(i, j int) bool { return a[i].Name < a[j].Name }

// FixedIP is holds data that specifies a fixed IP.
type FixedIP struct {
	SubnetID  string `json:"subnet_id"`
	IPAddress string `json:"ip_address,omitempty"`
}

// Ports issues a GET request that returns the found port responses
func (networkService Service) Ports() ([]PortResponse, error) {
	reqURL := networkService.URL + "/ports"
	var portResponse = portsResp{}
	_, err := networkService.Session.GetJSON(reqURL, nil, nil, &portResponse)
	if err != nil {
		return nil, err
	}

	return portResponse.Ports, nil
}

// Port issues a GET request that returns a specific port response.
func (networkService Service) Port(id string) (PortResponse, error) {
	reqURL := networkService.URL + "/ports/" + id
	portResponse := portResp{}
	_, err := networkService.Session.GetJSON(reqURL, nil, nil, &portResponse)
	return portResponse.Port, err
}

// DeletePort issues a DELETE to the specified port url to delete it.
func (networkService Service) DeletePort(id string) error {
	reqURL := networkService.URL + "/ports/" + id
	_, err := networkService.Session.Delete(reqURL, nil, nil)
	return err
}

// CreatePort issues a POST to create the specified port and return a PortResponse.
func (networkService Service) CreatePort(parameters CreatePortParameters) (PortResponse, error) {
	reqURL := networkService.URL + "/ports"
	parametersContainer := createPortContainer{Port: parameters}
	portResponse := portResp{}

	_, err := networkService.Session.PostJSON(reqURL, nil, nil, &parametersContainer, &portResponse)
	return portResponse.Port, err
}

type portsResp struct {
	Ports []PortResponse `json:"ports"`
}

type portResp struct {
	Port PortResponse `json:"port"`
}

type createPortContainer struct {
	Port CreatePortParameters `json:"port"`
}
