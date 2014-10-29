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

package compute

import (
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// KeyPair is a structure for all properties of for a tenant
type KeyPair struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	FingerPrint string `json:"fingerprint,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

// KeyPairs will issue a get request to retrieve the all keypairs.
func (computeService Service) KeyPairs() ([]KeyPair, error) {
	var kp = keyPairsResp{}
	err := misc.GetJSON(computeService.URL+"/os-keypairs", computeService.TokenID, computeService.Client, &kp)
	if err != nil {
		return nil, err
	}

	keypairs := make([]KeyPair, 0)
	for _, keyPairContainedItem := range kp.KeyPairs {
		keypairs = append(keypairs, keyPairContainedItem.KeyPair)
	}

	return keypairs, nil
}

// KeyPair will issue a get request to retrieve the specified keypair.
func (computeService Service) KeyPair(name string) (KeyPair, error) {
	var kp = keyPairContainer{}
	err := misc.GetJSON(computeService.URL+"/os-keypairs/"+name, computeService.TokenID, computeService.Client, &kp)
	return kp.KeyPair, err
}

// CreateKeyPair will send a POST request to create a new keypair with the specified parameters.
func (computeService Service) CreateKeyPair(name string, publickey string) (KeyPair, error) {
	reqURL := computeService.URL + "/os-keypairs"
	createKeypairContainerValues := keyPairContainer{KeyPair{Name: name, PublicKey: publickey}}
	outKeyContainer := keyPairContainer{}
	err := misc.PostJSON(reqURL, computeService.TokenID, computeService.Client, createKeypairContainerValues, &outKeyContainer)
	return outKeyContainer.KeyPair, err
}

// DeleteKeyPair will delete the keypair.
func (computeService Service) DeleteKeyPair(name string) (err error) {
	reqURL := computeService.URL + "/os-keypairs/" + name
	return misc.Delete(reqURL, computeService.TokenID, computeService.Client)
}

type keyPairsResp struct {
	KeyPairs []keyPairContainer `json:"keypairs"`
}

type keyPairContainer struct {
	KeyPair KeyPair `json:"keypair"`
}
