// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
Package volume implements a client library for accessing OpenStack Volume service

The CRUD operation of volumes can be retrieved using the api.

*/

package v3

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"git.openstack.org/openstack/golang-client.git/openstack"
	"git.openstack.org/openstack/golang-client.git/util"
)

type Service struct {
	Session openstack.Session
	Client  http.Client
	URL     string
}

// Response is a structure for all properties of
// an volume for a non detailed query
type Response struct {
	ID    string              `json:"id"`
	Links []map[string]string `json:"links"`
	Name  string              `json:"name"`
}

type RequestBody struct {
	// The size of the volume, in gibibytes (GiB) [REQUIRED]
	Size int
	// The volume name [OPTIONAL]
	Name string
}

func (volumeService Service) Create(reqBody *RequestBody) (Response, error) {
	return volumeService.createVolume(reqBody)
}

func (volumeService Service) createVolume(reqBody *RequestBody) (Response, error) {
	nullResponse := new(Response)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nil, err
	}
	urlPostFix := "/volumes"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Accept", "application/json")
	body := buildBody(volumeService, reqBody)
	resp, err := volumeService.Session.Post(reqURL.String(), nil, &headers, json.Marshal(body))
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}
	result, _ := json.Unmarshal(rbody)
	return result, nil
}

func (volumeService Service) Show(id string) (Response, error) {
	return volumeService.getVolume(id)
}

func (volumeService Service) getVolume(id string) (Response, error) {
	nullResponse := new(Response)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nil, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Accept", "application/json")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}
	result, _ := json.Unmarshal(rbody)
	return result, nil
}

func (volumeService Service) List() ([]Response, error) {
	return volumeService.getAllVolumes()
}

func (volumeService Service) getAllVolumes() (Response, error) {
	nullResponse := new(Response)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nil, err
	}
	urlPostFix := "/volumes"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Accept", "application/json")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}
	result, _ := json.Unmarshal(rbody)
	return result, nil
}

func (volumeService Service) Update(params Parameters) (Response, error) {
	return volumeService.updateVolume(params)
}

func (volumeService Service) updateVolume(id string, reqBody *RequestBody) (Response, error) {
	nullResponse := new(Response)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nil, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Accept", "application/json")
	body := buildBody(volumeService, reqBody)
	resp, err := volumeService.Session.Put(reqURL.String(), nil, &headers, json.Marshal(body))
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}
	result, _ := json.Unmarshal(rbody)
	return result, nil
}

func (volumeService Service) Delete(id int) (Response, error) {
	return volumeService.deleteVolume(params)
}

func (volumeService Service) deleteVolume(id string) (Response, error) {
	nullResponse := new(Response)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		return nil, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Accept", "application/json")
	resp, err := volumeService.Session.Delete(reqURL.String(), nil, &headers)
	if err != nil {
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nullResponse, errors.New("aaa")
	}
	result, _ := json.Unmarshal(rbody)
	return result, nil
}

func buildBody(volumeService Service, rbody *RequestBody) map[string]interface{} {
	body := make(map[string]interface{})
	if rbody != nil {
		r := make(map[string][]string)
		if rbody.Size != 0 {
			r.Set("size", fmt.Sprintf("%d", rbody.Size))
		}
		if rbody.Name != "" {
			r.Set("name", fmt.Sprintf("%s", rbody.Name))
		}
		body["volume"] = r
	}
	return body
}
