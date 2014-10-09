// Package v3 provides an OpenStack Identity Service v3 compatible client.
package v3

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
)

// Identity provides an identity service v3 session.
type Identity struct {
	// The Endpoint to connect to. For example, https://identity.example.com/v3/
	Endpoint string

	// The username and password associated with the account.
	Username string
	Password string

	// (Optionall) set the Project (formerly tenant) to use.
	ProjectId string

	// The HTTP Client. Setting the client allows you to set a proxt or other
	// parameters.
	Client *http.Client

	token string
}

// Token returns an authentication token.
func (i *Identity) Token() string {
	// Note, Token is implemented as a method so we can regenerate a token if the
	// one we have has expired. This can happen on the fly.
	// TODO: Regenerate the token on the fly.
	return i.token
}

// Authenticate attempts to use the credentials and authenticate against the
// identity service endpoint and retrieve a token.
func (i *Identity) Authenticate() (token string, err error) {

	// Put togethet the auth request JSON
	auth := authContainer{
		Auth: authRequestContainer{
			Identity: identityContainer{
				Methods: []string{"password"},
				Password: passwordCreds{
					User: userCreds{
						Username: i.Username,
						Password: i.Password,
					},
				},
			},
		},
	}

	// If there is a ProjectID make sure to include it.
	if i.ProjectId != "" {
		auth.Auth.Scope = scopeContainer{
			Project: userProject{
				Id: i.ProjectId,
			},
		}
	}

	jsonStr, _ := json.Marshal(auth)

	req, _ := http.NewRequest("POST", i.Endpoint+"/auth/tokens", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err := i.Client.Do(req)
	if err != nil {
		return "", err
	}

	i.token = resp.Header.Get("X-Subject-Token")

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var m struct{}
	err = json.Unmarshal(body, &m)
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
	Username string `json:"name"`
	Password string `json:"password"`
}

type scopeContainer struct {
	Project userProject `json:"project,omitempty"`
}

type userProject struct {
	Id string `json:"id,omitempty"`
}

// =============================================================================
