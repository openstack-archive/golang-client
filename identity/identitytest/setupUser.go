package identitytest

import (
	"encoding/json"
	"io/ioutil"
)

//SetupUser() is used to retrieve externally stored testing credentials.
//The testing credentials are stored outside
//the source code so they do not get checked in, assuming the user.json is
//in .gitignore / .hgignore. "user.json" should contain the following where
//... is the actual value from the test user account credentials.
//{
// "TenantId":"...",
// "TenantName": "...",
// "AccessKey": "...",
// "SecretKey": "...",
// "UserName": "...",
// "Password": "...",
// "Host": "https://.../v2.0/tokens"
//}
func SetupUser(jsonFile string) (acct struct {
	TenantId, TenantName, AccessKey, SecretKey, UserName, Password, Host string
},) {
	usrJson, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic("ReadFile json failed")
	}
	if err = json.Unmarshal(usrJson, &acct); err != nil {
		panic("Unmarshal json failed")
	}
	return acct
}
