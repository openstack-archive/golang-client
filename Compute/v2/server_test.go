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

package compute_test

import (
	"encoding/json"
	"errors"
	"git.openstack.org/stackforge/golang-client.git/Compute/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCreateServer(t *testing.T) {

	now := time.Now()

	var serverInfoDetail = compute.ServerInfoDetail{
		"my_server_id_1",
		"my_server_name_1",
		"my_server_status_1",
		&now, // created time
		&now, // updated time
		"my_server_hostId_1",
		make(map[string][]compute.Address),
		[]compute.Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}},
		compute.Image{"image_id1", []compute.Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}}},
		compute.Flavor{"image_id1", []compute.Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}}},
		"my_OS-EXT-STS_task_state_1",
		"my_OS-EXT-STS_vm_state_1",
		1, // PowerState
		"my_zone_a",
		"my_user_id_1",
		"my_tenant_id_1",
		"192.168.0.12",
		"my_accessIPv6_1",
		"my_config_drive_1",
		2, // Progress
		make(map[string]string),
		"my_adminPass_1",
	}

	marshaledserverInfoDetail, err := json.Marshal(serverInfoDetail)
	if err != nil {
		t.Error(err)
	}

	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				w.Header().Set("Version", "my_server_version.1.0.0")
				w.WriteHeader(201)
				w.Write([]byte(marshaledserverInfoDetail))
				return
			}
			t.Error(errors.New("Failed: r.Method == POST"))
		}))
	defer apiServer.Close()

	serverService := compute.ServerService{
		apiServer.URL,
		"my_token"}

	var serverServiceParameters = compute.ServerServiceParameters{
		"my_server",
		"8c3cd338-1282-4fbb-bbaf-2256ff97c7b7", // imageRef
		"my_key_name",
		101, // flavorRef
		1,   // maxcount
		1,   // mincount
		"my_user_data",
		[]compute.ServerNetworkParameters{
			compute.ServerNetworkParameters{"1111d337-0282-4fbb-bbaf-2256ff97c7b7", "881"},
			compute.ServerNetworkParameters{"2222d337-0282-4fbb-bbaf-2256ff97c7b7", "882"},
			compute.ServerNetworkParameters{"3333d337-0282-4fbb-bbaf-2256ff97c7b7", "883"}},
		[]compute.SecurityGroup{
			compute.SecurityGroup{"my_security_group_123"},
			compute.SecurityGroup{"my_security_group_456"}}}

	var serverRequestInfo = make(compute.ServerRequestInfo)
	serverRequestInfo["server"] = serverServiceParameters

	result, err := compute.CreateServer(new(http.Client), serverService, serverRequestInfo)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result, serverInfoDetail) {
		t.Error(errors.New("Failed: result != expected serverInfoDetail"))
	}

	err = nil
	return

}
