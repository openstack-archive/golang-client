// keypair.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/misc"

	"github.com/parnurzeal/gorequest"
)

type keyPairsResp struct {
	KeyPairs []Keypair `josn:"keypairs"`
}

type Keypair struct {
	KeyPair KeyPairEntry `json:"keypair"`
}

type KeyPairEntry struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	FingerPrint string `json:"fingerprint"`
}

func GetKeypairs(url string, token identity.Token) (keypairs []Keypair, err error) {

	req := gorequest.New()

	resp, body, errs := req.Get(url+"/os-keypairs").
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Auth-Token", token.Id).
		End()

	if err = misc.CheckHttpResponseStatusCode(resp); err != nil {
		return
	}

	if errs != nil {
		err = errs[len(errs)-1]
		return
	}

	var kp = keyPairsResp{}
	if err = json.Unmarshal([]byte(body), &kp); err != nil {
		return
	}

	keypairs = kp.KeyPairs
	err = nil
	return
}

func GetKeypair(url string, token identity.Token, name string) (keypair KeyPairEntry, err error) {

	keypairs, err := GetKeypairs(url, token)
	if err != nil {
		return
	}

	for _, v := range keypairs {
		if v.KeyPair.Name == name {
			keypair = v.KeyPair
			err = nil
			return
		}
	}

	err = errors.New(fmt.Sprintf("keypair %s not found", name))
	return
}
