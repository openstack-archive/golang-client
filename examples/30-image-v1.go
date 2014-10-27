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
	"git.openstack.org/stackforge/golang-client.git/image/v1"
	"net/http"
	"time"
)

// Image examples.
func main() {
	config := getConfig()

	// Authenticate with a username, password, tenant id.
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

	// Find the endpoint for the image service.
	url := ""
	for _, svc := range auth.Access.ServiceCatalog {
		if svc.Type == "image" {
			for _, ep := range svc.Endpoints {
				if ep.VersionId == "1.0" && ep.Region == config.ImageRegion {
					url = ep.PublicURL
					break
				}
			}
		}
	}

	if url == "" {
		panic("v1 image service url not found during authentication")
	}

	imageService := image.Service{TokenID: auth.Access.Token.Id, Client: *http.DefaultClient, URL: url}
	imagesDetails, err := imageService.ImagesDetail()
	if err != nil {
		panicString := fmt.Sprint("Cannot access images:", err)
		panic(panicString)
	}

	var imageIDs = make([]string, 0)
	for _, element := range imagesDetails {
		imageIDs = append(imageIDs, element.ID)
	}

	if len(imageIDs) == 0 {
		panicString := fmt.Sprint("No images found, check to make sure access is correct")
		panic(panicString)
	}
}
