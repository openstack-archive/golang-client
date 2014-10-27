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

package objectstorage

import (
	"git.openstack.org/stackforge/golang-client.git/misc"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var zeroByte = &([]byte{}) //pointer to empty []byte

//ListContainers calls the OpenStack list containers API using
//previously obtained token.
//"limit" and "marker" corresponds to the API's "limit" and "marker".
//"url" can be regular storage or cdn-enabled storage URL.
//It returns []byte which then needs to be unmarshalled to decode the JSON.
func ListContainers(limit int64, marker, url, token string) ([]byte, error) {
	return ListObjects(limit, marker, "", "", "", url, token)
}

//GetAccountMeta calls the OpenStack retrieve account metadata API using
//previously obtained token.
func GetAccountMeta(url, token string) (http.Header, error) {
	return GetObjectMeta(url, token)
}

//DeleteContainer calls the OpenStack delete container API using
//previously obtained token.
func DeleteContainer(url, token string) error {
	return DeleteObject(url, token)
}

//GetContainerMeta calls the OpenStack retrieve object metadata API
//using previously obtained token.
//url can be regular storage or CDN-enabled storage URL.
func GetContainerMeta(url, token string) (http.Header, error) {
	return GetObjectMeta(url, token)
}

//SetContainerMeta calls the OpenStack API to create / update meta data
//for container using previously obtained token.
//url can be regular storage or CDN-enabled storage URL.
func SetContainerMeta(url string, token string, s ...string) (err error) {
	return SetObjectMeta(url, token, s...)
}

//PutContainer calls the OpenStack API to create / update
//container using previously obtained token.
func PutContainer(url, token string, s ...string) error {
	return PutObject(zeroByte, url, token, s...)
}

//ListObjects calls the OpenStack list object API using previously
//obtained token. "Limit", "marker", "prefix", "path", "delim" corresponds
//to the API's "limit", "marker", "prefix", "path", and "delimiter".
func ListObjects(limit int64,
	marker, prefix, path, delim, conURL, token string) ([]byte, error) {
	var query = "?format=json"
	if limit > 0 {
		query += "&limit=" + strconv.FormatInt(limit, 10)
	}
	if marker != "" {
		query += "&marker=" + url.QueryEscape(marker)
	}
	if prefix != "" {
		query += "&prefix=" + url.QueryEscape(prefix)
	}
	if path != "" {
		query += "&path=" + url.QueryEscape(path)
	}
	if delim != "" {
		query += "&delimiter=" + url.QueryEscape(delim)
	}
	resp, err := misc.CallAPI("GET", conURL+query, zeroByte,
		"X-Auth-Token", token)
	if err != nil {
		return nil, err
	}
	if err = misc.CheckHTTPResponseStatusCode(resp); err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

//PutObject calls the OpenStack create object API using previously
//obtained token.
//url can be regular storage or CDN-enabled storage URL.
func PutObject(fContent *[]byte, url, token string, s ...string) (err error) {
	s = append(s, "X-Auth-Token")
	s = append(s, token)
	resp, err := misc.CallAPI("PUT", url, fContent, s...)
	if err != nil {
		return err
	}
	return misc.CheckHTTPResponseStatusCode(resp)
}

//CopyObject calls the OpenStack copy object API using previously obtained
//token.  Note from API doc: "The destination container must exist before
//attempting the copy."
func CopyObject(srcURL, destURL, token string) (err error) {
	resp, err := misc.CallAPI("COPY", srcURL, zeroByte,
		"X-Auth-Token", token,
		"Destination", destURL)
	if err != nil {
		return err
	}
	return misc.CheckHTTPResponseStatusCode(resp)
}

//DeleteObject calls the OpenStack delete object API using
//previously obtained token.
//
//Note from API doc: "A DELETE to a versioned object removes the current version
//of the object and replaces it with the next-most current version, moving it
//from the non-current container to the current." .. "If you want to completely
//remove an object and you have five total versions of it, you must DELETE it
//five times."
func DeleteObject(url, token string) (err error) {
	resp, err := misc.CallAPI("DELETE", url, zeroByte, "X-Auth-Token", token)
	if err != nil {
		return err
	}
	return misc.CheckHTTPResponseStatusCode(resp)
}

//SetObjectMeta calls the OpenStack API to create/update meta data for
//object using previously obtained token.
func SetObjectMeta(url string, token string, s ...string) (err error) {
	s = append(s, "X-Auth-Token")
	s = append(s, token)
	resp, err := misc.CallAPI("POST", url, zeroByte, s...)
	if err != nil {
		return err
	}
	return misc.CheckHTTPResponseStatusCode(resp)
}

//GetObjectMeta calls the OpenStack retrieve object metadata API using
//previously obtained token.
func GetObjectMeta(url, token string) (http.Header, error) {
	resp, err := misc.CallAPI("HEAD", url, zeroByte, "X-Auth-Token", token)
	if err != nil {
		return nil, err
	}
	return resp.Header, misc.CheckHTTPResponseStatusCode(resp)
}

//GetObject calls the OpenStack retrieve object API using previously
//obtained token. It returns http.Header, object / file content downloaded
//from the server, and err.
//
//Since this implementation of GetObject retrieves header info, it
//effectively executes GetObjectMeta also in addition to getting the
//object content.
func GetObject(url, token string) (http.Header, []byte, error) {
	resp, err := misc.CallAPI("GET", url, zeroByte, "X-Auth-Token", token)
	if err != nil {
		return nil, nil, err
	}
	if err = misc.CheckHTTPResponseStatusCode(resp); err != nil {
		return nil, nil, err
	}
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	}
	resp.Body.Close()
	return resp.Header, body, nil
}
