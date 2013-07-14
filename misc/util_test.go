package misc_test

import (
	"bytes"
	"errors"
	misc "golang-client/misc"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCallAPI(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			w.WriteHeader(200) //ok
		}))
	zeroByte := &([]byte{})
	if _, err := misc.CallAPI("HEAD", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
	if _, err := misc.CallAPI("DELETE", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
	if _, err := misc.CallAPI("POST", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
}

func TestCallAPIGetContent(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	fContent, err := ioutil.ReadFile("./util.go")
	if err != nil {
		t.Error(err)
	}
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
			w.Write(body)
		}))
	var resp *http.Response
	if resp, err = misc.CallAPI("GET", apiServer.URL, &fContent, "X-Auth-Token", tokn,
		"Etag", "md5hash-blahblah"); err != nil {
		t.Error(err)
	}
	if strconv.Itoa(len(fContent)) != resp.Header.Get("Content-Length") {
		t.Error(errors.New("Failed: Content-Length comparison"))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(fContent, body) {
		t.Error(errors.New("Failed: Content body comparison"))
	}
}

func TestCallAPIPutContent(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	fContent, err := ioutil.ReadFile("./util.go")
	if err != nil {
		t.Error(err)
	}
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if strconv.Itoa(len(fContent)) != r.Header.Get("Content-Length") {
				t.Error(errors.New("Failed: Content-Length comparison"))
			}
			if !bytes.Equal(fContent, body) {
				t.Error(errors.New("Failed: Content body comparison"))
			}
			w.WriteHeader(200)
		}))
	if _, err = misc.CallAPI("PUT", apiServer.URL, &fContent, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
}
