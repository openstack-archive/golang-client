// session - REST client session

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

package session

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
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

// Geenric callback to get a token from the auth plugin
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
		httpClient:   &http.Client{Transport: tr},
		endpoint:     strings.TrimRight(endpoint, "/"),
		authenticate: af,
		Headers:      http.Header{},
	}
	return session, nil
}

func (self *Session) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	// add token, get one if needed
	if self.Token == nil && self.authenticate != nil {
		var tok TokenInterface
		tok, err = self.authenticate(self, nil)
		if err != nil {
			// (re-)auth failure!!
			return
		}
		self.Token = tok
	}
	if self.Token != nil {
		req.Header.Add("X-Auth-Token", self.Token.GetTokenId())
	}
	return
}

func (self *Session) Do(req *http.Request) (*Response, error) {
	if *Debug {
		d, _ := httputil.DumpRequestOut(req, true)
		print("----------\nREQUEST:\n", string(d), "\n----------\n")
	}

	// Add session headers
	for k := range self.Headers {
		req.Header.Set(k, self.Headers.Get(k))
	}

	hresp, err := self.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if *Debug {
		dr, _ := httputil.DumpResponse(hresp, true)
		print("\nRESULT:\n", string(dr), "\n----------\n")
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

	req, err := s.NewRequest(method, url, buf)
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

func (self *Session) Get(
	url string,
	params *url.Values) (resp *Response, err error) {
	return self.Request("GET", url, params, nil, nil)
}

func (self *Session) Post(
	url string,
	params *url.Values,
	body *[]byte) (resp *Response, err error) {
	return self.Request("POST", url, params, nil, body)
}

// Post sends a POST request.
func Post(
	url string,
	params *url.Values,
	body *[]byte) (resp *Response, err error) {
	s, _ := NewSession(nil, "", nil)
	return s.Post(url, params, body)
}
