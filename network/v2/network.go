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

// Package network is used to create, delete, and query, networks, ports and subnets
package network

import (
	"git.openstack.org/stackforge/golang-client.git/misc"
	"net/http"
)

// Service holds state that is use to make requests and get responses for networks,
// ports and subnets
type Service struct {
	Client   http.Client
	TokenID  string
	TenantID string
	URL      string
}

// Response returns a set of values of the a network response.
type Response struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	Subnets             []string `json:"subnets"`
	TenantID            string   `json:"tenant_id"`
	RouterExternal      bool     `json:"router:external"`
	AdminStateUp        bool     `json:"admin_state_up"`
	Shared              bool     `json:"shared"`
	PortSecurityEnabled bool     `json:"port_security_enabled"`
}

// Networks will issue a get query that returns a list of networks
func (networkService Service) Networks() ([]Response, error) {
	reqURL := networkService.URL + "/networks"
	nwsContainer := networksResp{}
	err := misc.GetJSON(reqURL, networkService.TokenID, networkService.Client, &nwsContainer)
	return nwsContainer.Networks, err
}

// Network will issue a get request for a specific network.
func (networkService Service) Network(id string) (Response, error) {
	reqURL := networkService.URL + "/networks/" + id
	nwContainer := networkResp{}
	err := misc.GetJSON(reqURL, networkService.TokenID, networkService.Client, &nwContainer)
	return nwContainer.Network, err
}

// CreateNetwork will send a POST request to create a new network with the specified parameters.
func (networkService Service) CreateNetwork(adminStateUp bool, name string, shared bool) (Response, error) {
	createParameters := createNetworkValuesContainer{createNetworkValues{Name: name, AdminStateUp: adminStateUp, Shared: shared, TenantID: networkService.TenantID}}
	reqURL := networkService.URL + "/networks"
	nContainer := networkResp{}
	err := misc.PostJSON(reqURL, networkService.TokenID, networkService.Client, createParameters, &nContainer)
	return nContainer.Network, err
}

// DeleteNetwork will delete the specified network.
func (networkService Service) DeleteNetwork(name string) (err error) {
	reqURL := networkService.URL + "/networks/" + name
	return misc.Delete(reqURL, networkService.TokenID, networkService.Client)
}

type createNetworkValues struct {
	Name         string `json:"name"`
	AdminStateUp bool   `json:"admin_state_up"`
	Shared       bool   `json:"shared"`
	TenantID     string `json:"tenant_id"`
}

type createNetworkValuesContainer struct {
	Network createNetworkValues `json:"network"`
}

type networksResp struct {
	Networks []Response `json:"networks"`
}

type networkResp struct {
	Network Response `json:"network"`
}
