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
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// SubnetResponse returns a set of values of the a Subnet response
type SubnetResponse struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	NetworkID       string           `json:"network_id"`
	TenantID        string           `json:"tenant_id"`
	EnableDHCP      bool             `json:"enable_dhcp"`
	DNSNameserver   []string         `json:"dns_nameservers"`
	AllocationPools []AllocationPool `json:"allocation_pools"`
	HostRoutes      []string         `json:"host_routes"`
	IPVersion       IPVersion        `json:"ip_version"`
	GatewayIP       string           `json:"gateway_ip"`
	CIDR            string           `json:"cidr"`
}

// CreateSubnetParameters is a set of values to create a new subnet.
type CreateSubnetParameters struct {
	NetworkID       string           `json:"network_id"`
	IPVersion       IPVersion        `json:"ip_version"`
	CIDR            string           `json:"cidr"`
	AllocationPools []AllocationPool `json:"allocation_pools"`
}

// AllocationPool is a set of values for an allocation pool of ip addresses.
type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// IPVersion type indicates whether an ip address is IPV4 or IPV6.
type IPVersion int

const (
	// IPV4 indicates its an ip address version 4.
	IPV4 IPVersion = 4
	// IPV6 indicates its an ip address version 6
	IPV6 IPVersion = 6
)

// Subnets issues a GET request to return all subnets.
func (networkService Service) Subnets() ([]SubnetResponse, error) {
	reqURL := networkService.URL + "/subnets"
	var sn = subnetsResp{}
	err := misc.GetJSON(reqURL, networkService.TokenID, networkService.Client, &sn)
	return sn.Subnets, err
}

// Subnet issues a GET request to a specific url of a subnet and returns a subnet response.
func (networkService Service) Subnet(id string) (SubnetResponse, error) {
	reqURL := networkService.URL + "/subnets/" + id

	subnetContainer := subnetResp{}
	err := misc.GetJSON(reqURL, networkService.TokenID, networkService.Client, &subnetContainer)
	return subnetContainer.Subnet, err
}

// DeleteSubnet issues a DELETE request to remove the subnet.
func (networkService Service) DeleteSubnet(id string) error {

	reqURL := networkService.URL + "/subnets/" + id

	return misc.Delete(reqURL, networkService.TokenID, networkService.Client)
}

// CreateSubnet issues a GET request to add a Subnet with the specified parameters
// and returns the Subnet created.
func (networkService Service) CreateSubnet(parameters CreateSubnetParameters) (SubnetResponse, error) {

	reqURL := networkService.URL + "/subnets"
	parametersContainer := createSubnetContainer{Subnet: parameters}
	response := subnetResp{}

	err := misc.PostJSON(reqURL, networkService.TokenID, networkService.Client, parametersContainer, &response)
	return response.Subnet, err
}

type subnetResp struct {
	Subnet SubnetResponse `json:"subnet"`
}

type subnetsResp struct {
	Subnets []SubnetResponse `json:"subnets"`
}

type createSubnetContainer struct {
	Subnet CreateSubnetParameters `json:"subnet"`
}
