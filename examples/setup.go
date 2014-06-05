// The acceptance package is a set of acceptance tests showcasing how the
// contents of the package are meant to be used. This is setup in a similar
// manner to a consuming application.
package main

import (
	"encoding/json"
	"io/ioutil"
)

// testconfig contains the user information needed by the acceptance and
// integration tests.
type testconfig struct {
	Host        string
	Username    string
	Password    string
	ProjectID   string
	ProjectName string
	Container   string
}

// getConfig provides access to credentials in other tests and examples.
func getConfig() *testconfig {
	config := &testconfig{}
	userJSON, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("ReadFile json failed")
	}
	if err = json.Unmarshal(userJSON, &config); err != nil {
		panic("Unmarshal json failed")
	}
	return config
}
