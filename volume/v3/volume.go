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
	"io/ioutil"
	"log"
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

type RequestBody struct {
	Name         string `json:"name"`
	Size         int    `json:"size"`
	HostName     string `json:"host_name"`
	Mountpoint   string `json:"mountpoint"`
	AttachmentID string `json:"attachment_id"`
}

type CreateBody struct {
	VolumeBody RequestBody `json:"volume"`
}

type MountBody struct {
	VolumeBody RequestBody `json:"os-attach"`
}

type UnmountBody struct {
	VolumeBody RequestBody `json:"os-detach"`
}

// Response is a structure for all properties of
// an volume for a non detailed query
type Response struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Consistencygroup_id string `json:"consistencygroup_id"`
}

// DetailResponse is a structure for all properties of
// an volume for a detailed query
type DetailResponse struct {
	ID              string               `json:"id"`
	Attachments     []map[string]string  `json:"attachments"`
	Links           []map[string]string  `json:"links"`
	Metadata        map[string]string    `json:"metadata"`
	Protected       bool                 `json:"protected"`
	Status          string               `json:"status"`
	MigrationStatus string               `json:"migration_status"`
	UserID          string               `json:"user_id"`
	Encrypted       bool                 `json:"encrypted"`
	Multiattach     bool                 `json:"multiattach"`
	CreatedAt       util.RFC8601DateTime `json:"created_at"`
	Description     string               `json:"description"`
	Volume_type     string               `json:"volume_type"`
	Name            string               `json:"name"`
	Source_volid    string               `json:"source_volid"`
	Snapshot_id     string               `json:"snapshot_id"`
	Size            int64                `json:"size"`

	Aavailability_zone  string `json:"availability_zone"`
	Rreplication_status string `json:"replication_status"`
	Consistencygroup_id string `json:"consistencygroup_id"`
}

type VolumeResponse struct {
	Volume Response `json:"volume"`
}

type VolumesResponse struct {
	Volumes []Response `json:"volumes"`
}

type DetailVolumeResponse struct {
	DetailVolume DetailResponse `json:"volume"`
}

type DetailVolumesResponse struct {
	DetailVolumes []DetailResponse `json:"volumes"`
}

func (volumeService Service) Create(reqBody *CreateBody) (Response, error) {
	return volumeService.createVolume(reqBody)
}

func (volumeService Service) createVolume(reqBody *CreateBody) (Response, error) {
	nullResponse := Response{}

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponse, err
	}
	urlPostFix := "/volumes"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start POST request to create volume, body = %s\n", body)
	resp, err := volumeService.Session.Post(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("POST response error:", err)
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponse, err
	}

	volumeResponse := new(VolumeResponse)
	if err = json.Unmarshal(rbody, volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse.Volume, nil
}

func (volumeService Service) Show(id string) (DetailResponse, error) {
	return volumeService.getVolume(id)
}

func (volumeService Service) getVolume(id string) (DetailResponse, error) {
	nullResponse := DetailResponse{}

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponse, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to get volume!")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("GET response error:", err)
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponse, err
	}

	detailVolumeResponse := new(DetailVolumeResponse)
	if err = json.Unmarshal(rbody, detailVolumeResponse); err != nil {
		return nullResponse, err
	}
	return detailVolumeResponse.DetailVolume, nil
}

func (volumeService Service) List() ([]Response, error) {
	return volumeService.getAllVolumes()
}

func (volumeService Service) getAllVolumes() ([]Response, error) {
	nullResponses := make([]Response, 0)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponses, err
	}
	urlPostFix := "/volumes"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to get all volumes!")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("GET response error:", err)
		return nullResponses, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponses, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponses, err
	}

	volumesResponse := new(VolumesResponse)
	if err = json.Unmarshal(rbody, volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse.Volumes, nil
}

func (volumeService Service) Detail() ([]DetailResponse, error) {
	return volumeService.detailAllVolumes()
}

func (volumeService Service) detailAllVolumes() ([]DetailResponse, error) {
	nullResponses := make([]DetailResponse, 0)

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponses, err
	}
	urlPostFix := "/volumes/detail"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start GET request to detail all volumes!")
	resp, err := volumeService.Session.Get(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("GET response error:", err)
		return nullResponses, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponses, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponses, err
	}

	detailVolumesResponse := new(DetailVolumesResponse)
	if err = json.Unmarshal(rbody, detailVolumesResponse); err != nil {
		return nullResponses, err
	}
	return detailVolumesResponse.DetailVolumes, nil
}

func (volumeService Service) Update(id string, reqBody *CreateBody) (Response, error) {
	return volumeService.updateVolume(id, reqBody)
}

func (volumeService Service) updateVolume(id string, reqBody *CreateBody) (Response, error) {
	nullResponse := Response{}

	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return nullResponse, err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start PUT request to update volume, body = %s\n", body)
	resp, err := volumeService.Session.Put(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("PUT response error:", err)
		return nullResponse, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nullResponse, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nullResponse, err
	}

	volumeResponse := new(VolumeResponse)
	if err = json.Unmarshal(rbody, volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse.Volume, nil
}

func (volumeService Service) Delete(id string) error {
	return volumeService.deleteVolume(id)
}

func (volumeService Service) deleteVolume(id string) error {
	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/volumes" + "/" + id
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	log.Println("Start DELETE request to delete volume!")
	resp, err := volumeService.Session.Delete(reqURL.String(), nil, &headers)
	if err != nil {
		log.Println("DELETE response error:", err)
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}

func (volumeService Service) Mount(id string, reqBody *MountBody) error {
	return volumeService.mountVolume(id, reqBody)
}

func (volumeService Service) mountVolume(id string, reqBody *MountBody) error {
	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/volumes" + "/" + id + "/action"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start POST request to mount volume, body = %s\n", body)
	resp, err := volumeService.Session.Post(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("POST response error:", err)
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}

func (volumeService Service) Unmount(id string, reqBody *UnmountBody) error {
	return volumeService.unmountVolume(id, reqBody)
}

func (volumeService Service) unmountVolume(id string, reqBody *UnmountBody) error {
	reqURL, err := url.Parse(volumeService.URL)
	if err != nil {
		log.Println("Parse URL error:", err)
		return err
	}
	urlPostFix := "/volumes" + "/" + id + "/action"
	reqURL.Path += urlPostFix

	var headers http.Header = http.Header{}
	headers.Set("Content-Type", "application/json")
	body, _ := json.Marshal(reqBody)
	log.Printf("Start PUT request to unmount volume, body = %s\n", body)
	resp, err := volumeService.Session.Put(reqURL.String(), nil, &headers, &body)
	if err != nil {
		log.Println("PUT response error:", err)
		return err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}
