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
	"git.openstack.org/stackforge/golang-client.git/identity/v3"
	"net/http"
)

func main() {
	// Retrieve the config to an OpenStack environment.
	// See config.json.dist for details.
	config := getConfig()

	// As a user, I can authenticate with a username and password retriving a token.
	ident := identity.Identity{
		ID:       config.Username,
		Password: config.Password,
		Endpoint: config.EndpointV3,
		Client:   http.DefaultClient,
	}
	token, err := ident.Authenticate()

	if err != nil {
		fmt.Println("There was an error authenticating:", err)
		return
	}

	if ident.Token() != token {
		fmt.Println("Returned tokens don't match.")
		return
	}

	// As a user, I can authenticate with a username, password, and projectID
	// where I retrieve a token.
	ident = identity.Identity{
		ID:        config.Username,
		Password:  config.Password,
		Endpoint:  config.EndpointV3,
		Client:    http.DefaultClient,
		ProjectID: config.ProjectID,
	}
	token, err = ident.Authenticate()

	if err != nil {
		fmt.Println("There was an error authenticating:", err)
		return
	}

	if ident.Token() != token {
		fmt.Println("Returned tokens don't match.")
		return
	}
}
