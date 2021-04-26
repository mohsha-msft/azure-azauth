package azauth

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
)

type AuthConfig struct {
	SharedKeyAuthConfig

	MSIAuthConfig

	SPNAuthConfig
}

type BaseAuthConfig struct {
	AuthType        AuthType
	UseHTTP         bool
	PrivateEndPoint string
}

func (baseAuthCfg *BaseAuthConfig) ProtoType() string {
	protocol := "https"
	if baseAuthCfg.UseHTTP { protocol = "http"}
	return protocol
}

// ---------------------------------------------------------------------------------------------------------------------

// AzAuth is an interface common to
type AzAuth interface {
	GetCredentials() (azcore.Credential, error)
	GetEndpoints() (*url.URL, error)
}

type GetClientOptions struct {
	ClientType 		ClientType
	ClientOptions 	*ClientOptions
	Endpoint 		string
	Credential 		azcore.Credential
}

type AzClient interface {
	GetClient() (*Client, error)
}


