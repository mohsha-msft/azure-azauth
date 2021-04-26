package azauth

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
)

type SASAuthConfig struct {

}

type SASAuth struct {
	BaseAuthConfig
	SASAuthConfig
}

func (sas *SASAuth) GetCredentials() (azcore.Credential, error)  {
	panic("Not Implemented")
}

func (sas *SASAuth) GetEndpoints() (*url.URL, error)  {
	panic("Not Implemented")
}
