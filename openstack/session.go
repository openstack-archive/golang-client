// session - REST client session
// Copyright 2015 Dean Troyer
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openstack

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var Debug = new(bool)

type Response struct {
	Resp *http.Response
	Body []byte
}

type TokenInterface interface {
	GetTokenId() string
}

type Token struct {
	Expires string
	Id      string
	Project struct {
		Id   string
		Name string
	}
}

func (t Token) GetTokenId() string {
	return t.Id
}

// Generic callback to get a token from the auth plugin
type AuthFunc func(s *Session, opts interface{}) (TokenInterface, error)

type Session struct {
	httpClient   *http.Client
	endpoint     string
	authenticate AuthFunc
	Token        TokenInterface
	Headers      http.Header
	//	  ServCat map[string]ServiceEndpoint
}

func NewSession(af AuthFunc, endpoint string, tls *tls.Config) (session *Session, err error) {
	tr := &http.Transport{
		TLSClientConfig:    tls,
		DisableCompression: true,
	}
	session = &Session{
		// TODO(dtroyer): httpClient needs to be able to be passed in, or set externally
		httpClient:   &http.Client{Transport: tr},
		endpoint:     strings.TrimRight(endpoint, "/"),
		authenticate: af,
		Headers:      http.Header{},
	}
	return session, nil
}

func (s *Session) NewRequest(method, url string, headers *http.Header, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// add token, get one if needed
	if s.Token == nil && s.authenticate != nil {
		var tok TokenInterface
		tok, err = s.authenticate(s, nil)
		if err != nil {
			// (re-)auth failure!!
			return nil, err
		}
		s.Token = tok
	}
	if headers != nil {
		req.Header = *headers
	}
	if s.Token != nil {
		req.Header.Add("X-Auth-Token", s.Token.GetTokenId())
	}
	return
}

func (s *Session) Do(req *http.Request) (*Response, error) {
	if *Debug {
		d, _ := httputil.DumpRequestOut(req, true)
		log.Printf(">>>>>>>>>> REQUEST:\n", string(d))
	}

	// Add session headers
	for k := range s.Headers {
		req.Header.Set(k, s.Headers.Get(k))
	}

	hresp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if *Debug {
		dr, _ := httputil.DumpResponse(hresp, true)
		log.Printf("<<<<<<<<<< RESULT:\n", string(dr))
	}

	resp := new(Response)
	resp.Resp = hresp
	return resp, nil
}

// Perform a simple get to an endpoint
func (s *Session) Request(
	method string,
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte,
) (resp *Response, err error) {
	// add params to url here
	if params != nil {
		url = url + "?" + params.Encode()
	}

	// Get the body if one is present
	var buf io.Reader
	if body != nil {
		buf = bytes.NewReader(*body)
	}

	req, err := s.NewRequest(method, url, headers, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err = s.Do(req)
	if err != nil {
		return nil, err
	}
	// do we need to parse this in this func? yes...
	defer resp.Resp.Body.Close()

	resp.Body, err = ioutil.ReadAll(resp.Resp.Body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Session) Get(
	url string,
	params *url.Values,
	headers *http.Header) (resp *Response, err error) {
	return s.Request("GET", url, params, headers, nil)
}

func (s *Session) Post(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte) (resp *Response, err error) {
	return s.Request("POST", url, params, headers, body)
}

func (s *Session) Put(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte) (resp *Response, err error) {
	return s.Request("PUT", url, params, headers, body)
}

// Get sends a GET request.
func Get(
	url string,
	params *url.Values,
	headers *http.Header) (resp *Response, err error) {
	s, _ := NewSession(nil, "", nil)
	return s.Get(url, params, headers)
}

// Post sends a POST request.
func Post(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte) (resp *Response, err error) {
	s, _ := NewSession(nil, "", nil)
	return s.Post(url, params, headers, body)
}

// Put sends a PUT request.
func Put(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte) (resp *Response, err error) {
	s, _ := NewSession(nil, "", nil)
	return s.Put(url, params, headers, body)
}
