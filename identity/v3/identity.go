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

// Package v3 provides an OpenStack Identity Service v3 compatible client.
package identity

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Identity provides an identity service v3 session. Creating an Identity object
// typically looks like:
//
// 		ident := identity.Identity{
//     		ID:       "Username/ID",
// 			Password: "Your Password",
// 			ProjectID: "Your Project ID",
// 			Endpoint: "https://idenity.example.com/v3/",
// 			Client:   http.DefaultClient,
// 		}
//
// The identity contains the endpoint to connect to, the creadentials to use,
// and the http client to use. The client is important as it can be configured
// to use a proxy or to handle any other special needs.
type Identity struct {
	// The Endpoint to connect to. For example, https://identity.example.com/v3/
	Endpoint string

	// The ID and password associated with the account.
	ID       string
	Password string

	// (Optionall) set the Project (formerly tenant) to use.
	ProjectID string

	// The HTTP Client. Setting the client allows you to set a proxt or other
	// parameters.
	Client *http.Client

	token string
}

// Token returns an authentication token if one exists. Otherwise it returns an
// empty string.
func (i *Identity) Token() string {
	// Note, Token is implemented as a method so we can regenerate a token if the
	// one we have has expired. This can happen on the fly.
	// TODO: Regenerate the token on the fly.
	return i.token
}

// Authenticate attempts to use the credentials and authenticate against the
// identity service endpoint and retrieve a token. For example,
//
// 		ident := identity.Identity{
//     		ID:       "Username/ID",
// 			Password: "Your Password",
// 			ProjectID: "Your Project ID",
// 			Endpoint: "https://idenity.example.com/v3/",
// 			Client:   http.DefaultClient,
// 		}
// 		token, err := ident.Authenticate()
//
// The Client is a standard http client used in Go. This can be used to configure
// a proxy or any other custom settings.
func (i *Identity) Authenticate() (token string, err error) {

	// Put together the auth request JSON
	auth := authContainer{
		Auth: authRequestContainer{
			Identity: identityContainer{
				Methods: []string{"password"},
				Password: passwordCreds{
					User: userCreds{
						ID:       i.ID,
						Password: i.Password,
					},
				},
			},
		},
	}

	// If there is a ProjectID make sure to include it.
	if i.ProjectID != "" {
		auth.Auth.Scope = scopeContainer{
			Project: userProject{
				ID: i.ProjectID,
			},
		}
	}

	jsonStr, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", i.Endpoint+"/auth/tokens", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := i.Client.Do(req)
	if err != nil {
		return "", err
	}

	i.token = resp.Header.Get("X-Subject-Token")

	defer resp.Body.Close()

	// TODO: unmarshall the response and store that data on the session.

	return i.token, nil
}

// =============================================================================
// authContainer, authRequestContainer, identityContainer, passwordCreds,
// userCreds, and the rest of this section describe the authentication request
// json.
type authContainer struct {
	Auth authRequestContainer `json:"auth"`
}

type authRequestContainer struct {
	Identity identityContainer `json:"identity"`
	Scope    scopeContainer    `json:"scope,omitempty"`
}

type identityContainer struct {
	Methods  []string      `json:"methods"`
	Password passwordCreds `json:"password"`
}

type passwordCreds struct {
	User userCreds `json:"user"`
}

type userCreds struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type scopeContainer struct {
	Project userProject `json:"project,omitempty"`
}

type userProject struct {
	ID string `json:"id,omitempty"`
}

// =============================================================================
