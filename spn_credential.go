package azauth

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
)

type SPNAuthConfig struct {
	AccountName 	string
	TenantID     	string
	ClientID     	string
	ClientSecret 	string
}

type SPNAuth struct {
	BaseAuthConfig
	SPNAuthConfig
}

func (spn *SPNAuth) GetCredentials() (azcore.Credential, error)  {
	panic("Not Implemented")
}

func (spn *SPNAuth) GetEndpoints() (*url.URL, error)  {
	if spn.BaseAuthConfig.PrivateEndPoint != "" {
		return url.Parse(spn.BaseAuthConfig.PrivateEndPoint)
	}
	return url.Parse(spn.ProtoType() + "//" + spn.AccountName + ".blob.core.windows.net/")
}
