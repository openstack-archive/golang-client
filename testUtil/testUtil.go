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

// image.go
package testUtil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// equals fails the test if exp is not equal to act.
// Code was copied from https://github.com/benbjohnson/testing
// MIT license
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func HeaderValuesEqual(t *testing.T, req *http.Request, name string, expectedValue string) {
	actualValue := req.Header.Get(name)
	if actualValue != expectedValue {
		t.Error(errors.New(fmt.Sprintf("Expected Header {Name:'%s', Value:'%s', actual value '%s'", name, expectedValue, actualValue)))
	}
}

func CreateGetJsonTestRequestServer(t *testing.T, expectedAuthTokenValue string, jsonPayload string, verifyRequest func(*http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			HeaderValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
			HeaderValuesEqual(t, r, "Accept", "application/json")
			verifyRequest(r)
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(jsonPayload))
				w.WriteHeader(200)
				return
			}

			t.Error(errors.New("Failed: r.Method == GET"))
		}))
}
