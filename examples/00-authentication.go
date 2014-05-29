package main

import (
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/identity"
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
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
	}

	// Authenticate with a username, password, tenant name.
	auth, err = identity.AuthUserNameTenantName(config.Host,
		config.Username,
		config.Password,
		config.ProjectName)
	if err != nil {
		fmt.Println("There was an error authenticating:", err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
	}

	// Authenticate with a username, password, tenant id.
	auth, err = identity.AuthUserNameTenantId(config.Host,
		config.Username,
		config.Password,
		config.ProjectID)
	if err != nil {
		fmt.Println("There was an error authenticating:", err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
	}
}
