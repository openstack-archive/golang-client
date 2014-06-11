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
	"time"
)

// Authentication examples.
func main() {
	config := getConfig()

	// Authenticate with just a username and password. The returned token is
	// unscoped to a tenant.
	auth, err := identity.AuthUserName(config.Host,
		config.Username,
		config.Password)
	if err != nil {
		fmt.Println("There was an error authenticating:", err)
		return
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
		return
	}

	// Authenticate with a username, password, tenant name.
	auth, err = identity.AuthUserNameTenantName(config.Host,
		config.Username,
		config.Password,
		config.ProjectName)
	if err != nil {
		fmt.Println("There was an error authenticating:", err)
		return
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
		return
	}

	// Authenticate with a username, password, tenant id.
	auth, err = identity.AuthUserNameTenantId(config.Host,
		config.Username,
		config.Password,
		config.ProjectID)
	if err != nil {
		fmt.Println("There was an error authenticating:", err)
		return
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
		return
	}
}
