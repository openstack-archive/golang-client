package objectstorage_test

import (
	"bytes"
	"encoding/json"
	"golang-client/identity"
	"golang-client/identity/identitytest"
	"golang-client/objectstorage"
	"io/ioutil"
	"testing"
)

//PRE-REQUISITE: Must have valid ObjectStorage account, either internally
//hosted or with one of the OpenStack providers.  Identity is assumed to
//use IdentityService mechanism, instead of legacy Swift mechanism.
func TestEndToEnd(t *testing.T) {
	//user.json holds the user account info needed to authenticate
	account := identitytest.SetupUser("../identity/identitytest/user.json")
	auth, err := identity.AuthUserNameTenantId(account.Host,
		account.UserName,
		account.Password,
		account.TenantId)
	if err != nil {
		t.Fatal(err)
	}

	url := ""
	for _, svc := range auth.Access.ServiceCatalog {
		if svc.Type == "object-store" {
			url = svc.Endpoints[0].PublicURL + "/"
			break
		}
	}
	if url == "" {
		t.Fatal("object-store url not found during authentication")
	}

	hdr, err := objectstorage.GetAccountMeta(url, auth.Access.Token.Id)
	if err != nil {
		t.Error("\nGetAccountMeta error\n", err)
	}

	container := "testContainer1"
	if err = objectstorage.PutContainer(url+container, auth.Access.Token.Id,
		"X-Log-Retention", "true"); err != nil {
		t.Fatal("\nPutContainer\n", err)
	}

	containersJson, err := objectstorage.ListContainers(0, "",
		url, auth.Access.Token.Id)
	if err != nil {
		t.Fatal(err)
	}

	type containerType struct {
		Name         string
		Bytes, Count int
	}
	containersList := []containerType{}

	if err = json.Unmarshal(containersJson, &containersList); err != nil {
		t.Error(err)
	}

	found := false
	for i := 0; i < len(containersList); i++ {
		if containersList[i].Name == container {
			found = true
		}
	}
	if !found {
		t.Fatal("created container is missing from downloaded containersList")
	}

	if err = objectstorage.SetContainerMeta(url+container, auth.Access.Token.Id,
		"X-Container-Meta-fubar", "false"); err != nil {
		t.Error(err)
	}
	hdr, err = objectstorage.GetContainerMeta(url+container, auth.Access.Token.Id)
	if err != nil {
		t.Error("\nGetContainerMeta error\n", err)
	}
	if hdr.Get("X-Container-Meta-fubar") != "false" {
		t.Error("container meta does not match")
	}

	var fContent []byte
	srcFile := "objectstorage_integration_test.go"
	fContent, err = ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	object := container + "/" + srcFile
	if err = objectstorage.PutObject(&fContent, url+object, auth.Access.Token.Id,
		"X-Object-Meta-fubar", "false"); err != nil {
		t.Fatal(err)
	}
	objectsJson, err := objectstorage.ListObjects(0, "", "", "", "",
		url+container, auth.Access.Token.Id)

	type objectType struct {
		Name, Hash, Content_type, Last_modified string
		Bytes                                   int
	}
	objectsList := []objectType{}

	if err = json.Unmarshal(objectsJson, &objectsList); err != nil {
		t.Error(err)
	}
	found = false
	for i := 0; i < len(objectsList); i++ {
		if objectsList[i].Name == srcFile {
			found = true
		}
	}
	if !found {
		t.Fatal("created object is missing from the objectsList")
	}

	if err = objectstorage.SetObjectMeta(url+object, auth.Access.Token.Id,
		"X-Object-Meta-fubar", "true"); err != nil {
		t.Error("\nSetObjectMeta error\n", err)
	}
	hdr, err = objectstorage.GetObjectMeta(url+object, auth.Access.Token.Id)
	if err != nil {
		t.Error("\nGetObjectMeta error\n", err)
	}
	if hdr.Get("X-Object-Meta-fubar") != "true" {
		t.Error("\nSetObjectMeta error\n", "object meta does not match")
	}

	_, body, err := objectstorage.GetObject(url+object, auth.Access.Token.Id)
	if err != nil {
		t.Error("\nGetObject error\n", err)
	}
	if !bytes.Equal(fContent, body) {
		t.Error("\nGetObject error\n", "byte comparison of uploaded != downloaded")
	}

	if err = objectstorage.CopyObject(url+object, "/"+object+".dup",
		auth.Access.Token.Id); err != nil {
		t.Fatal("\nCopyObject error\n", err)
	}

	if err = objectstorage.DeleteObject(url+object,
		auth.Access.Token.Id); err != nil {
		t.Fatal("\nDeleteObject error\n", err)
	}
	if err = objectstorage.DeleteObject(url+object+".dup",
		auth.Access.Token.Id); err != nil {
		t.Fatal("\nDeleteObject error\n", err)
	}

	if err = objectstorage.DeleteContainer(url+container,
		auth.Access.Token.Id); err != nil {
		t.Error("\nDeleteContainer error\n", err)
	}
}
