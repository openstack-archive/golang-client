//PRE-REQUISITE: Must have valid IdentityService account, either internally
//hosted or with one of the OpenStack providers.  See identitytest/ for the
//JSON specification.
//The JSON file ought to be in .hgignore / .gitignore for security reason.
package identity_test

import (
	"golang-client/identity"
	"golang-client/identity/identitytest"
	"testing"
	"time"
)

var account = identitytest.SetupUser("identitytest/user.json")

func TestAuthKey(t *testing.T) {
	//Not in OpenStack api doc, but in HPCloud api doc.
	auth, err := identity.AuthKey(account.Host,
		account.AccessKey,
		account.SecretKey)
	if err != nil {
		t.Error(err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		t.Error("expiry is wrong")
	}
}

func TestAuthKeyTenantId(t *testing.T) {
	//Not in OpenStack nor HPCloud api doc, but in HPCloud curl example.
	auth, err := identity.AuthKeyTenantId(account.Host,
		account.AccessKey,
		account.SecretKey,
		account.TenantId)
	if err != nil {
		t.Error(err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		t.Error("expiry is wrong")
	}
}

func TestAuthUserName(t *testing.T) {
	//Not in OpenStack api doc, but in HPCloud api doc.
	auth, err := identity.AuthUserName(account.Host,
		account.UserName,
		account.Password)
	if err != nil {
		t.Error(err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		t.Error("expiry is wrong")
	}
}
func TestAuthUserNameTenantName(t *testing.T) {
	//In OpenStack api doc, but not in HPCloud api doc, but tested valid in HPCloud.
	auth, err := identity.AuthUserNameTenantName(account.Host,
		account.UserName,
		account.Password,
		account.TenantName)
	if err != nil {
		t.Error(err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		t.Error("expiry is wrong")
	}
}

func TestAuthUserNameTenantId(t *testing.T) {
	//Not in OpenStack api doc, but in HPCloud api doc.
	auth, err := identity.AuthUserNameTenantId(account.Host,
		account.UserName,
		account.Password,
		account.TenantId)
	if err != nil {
		t.Error(err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		t.Error("expiry is wrong")
	}
}

func TestAuthTenantNameTokenId(t *testing.T) {
	//Not in OpenStack api doc, but in HPCloud api doc.
	auth, err := identity.AuthUserNameTenantId(account.Host,
		account.UserName,
		account.Password,
		account.TenantId)
	if err != nil {
		t.Error(err)
	}
	auth, err = identity.AuthTenantNameTokenId(account.Host,
		account.TenantName,
		auth.Access.Token.Id)
	if err != nil {
		t.Error(err)
	}
	if !auth.Access.Token.Expires.After(time.Now()) {
		t.Error("expiry is wrong")
	}
}
