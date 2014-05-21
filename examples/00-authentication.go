package main

import (
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/identity"
	"time"
)

// Authentication examples.
func main() {
	config := getConfig()

	auth, err := identity.AuthUserName(config.Host,
		config.Username,
		config.Password)
	if err != nil {
		fmt.Println("There was an error authenticating:", err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
	}
}
