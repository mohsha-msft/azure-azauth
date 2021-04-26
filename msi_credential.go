package azauth

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	//"github.com/Azure/go-autorest/autorest/azure"
	"net/url"
)

type MSIAuthConfig struct {
	AccountName 	string
	ApplicationID 	string
	ResourceID    	string
}

type MSIAuth struct {
	BaseAuthConfig
	MSIAuthConfig
}

func (msi *MSIAuth) GetCredentials() (azcore.Credential, error)  {
	panic("Not Implemented")
}

//// getOAuthToken : Generate token and register the refresh call backs
//func (msi *MSIAuth) getOAuthToken(resourceURL string, callbacks ...adal.TokenRefreshCallback) (*azblob.TokenCredential, error) {
//	// Generate the token based on configured inputs
//	spt, err := msi.fetchToken(resourceURL, callbacks...)
//	if err != nil {
//		return nil, err
//	}
//
//	// Refresh obtains a fresh token
//	err = spt.Refresh()
//	if err != nil {
//		log.Fatalf("Failed to refresh token (%s)", err.Error())
//		return nil, err
//	}
//
//	// Using token create the credential object, here also register a call back which refreshes the token
//	callback := func(tc adal.Token) time.Duration {
//		err := spt.Refresh()
//		if err != nil {
//			log.Panicf("Failed to refresh token (%s)", err.Error())
//			return 0
//		}
//
//		// set the new token value
//		tc.Se(spt.Token().AccessToken)
//		log.Printf("MSI Token retreived %s (%d)", spt.Token().AccessToken, spt.Token().Expires())
//
//		// Get the next token slightly before the current one expires
//		return time.Until(spt.Token().Expires()) - 10*time.Second
//	}
//	tc := adal.NewServicePrincipalToken(spt.Token().AccessToken, )
//
//	return &tc, nil
//}
//
//// fetchToken : Based on the input config generate a token
//func (msi *MSIAuth) fetchToken(resourceURL string, callbacks ...adal.TokenRefreshCallback) (*adal.ServicePrincipalToken, error) {
//	msiEndpoint, _ := adal.
//
//	var spt *adal.ServicePrincipalToken
//	var err error
//
//	if msi.ApplicationID == "" && msi.ResourceID == "" {
//		log.Debug("Generating MSI token using resource url (%s)", resourceURL)
//		spt, err = adal.NewServicePrincipalTokenFromMSI(msiEndpoint, resourceURL, callbacks...)
//	} else if msi.ApplicationID != "" {
//		log.Debug("Generating MSI token using application id (%s)", msi.ApplicationID)
//		spt, err = adal.NewServicePrincipalTokenFromMSIWithUserAssignedID(msiEndpoint, resourceURL, msi.ApplicationID, callbacks...)
//	} else if msi.ResourceID != "" {
//		log.Debug("Generating MSI token using resource id (%s)", msi.ResourceID)
//		spt, err = adal.NewServicePrincipalTokenFromMSIWithIdentityResourceID(msiEndpoint, resourceURL, msi.ResourceID, callbacks...)
//	}
//
//	if err != nil {
//		log.Err("Failed to generate MSI token (%s)", err.Error())
//		return nil, err
//	}
//
//	return spt, nil
//}

func (msi *MSIAuth) GetEndpoints() (*url.URL, error)  {
	if msi.BaseAuthConfig.PrivateEndPoint != "" {
		return url.Parse(msi.BaseAuthConfig.PrivateEndPoint)
	}
	return url.Parse(msi.ProtoType() + "//" + msi.AccountName + ".blob.core.windows.net/")
}


