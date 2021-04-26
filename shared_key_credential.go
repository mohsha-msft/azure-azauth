package azauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

const (
	AccountName  = "ACCOUNT_NAME"
	AccountKey   = "ACCOUNT_KEY"
)
type SharedKeyAuthConfig struct {
	AccountName string
	AccountKey  atomic.Value
}

type SharedKeyCredentials struct {
	BaseAuthConfig
	SharedKeyAuthConfig
}

func (ska *SharedKeyCredentials) GetClient(_ ClientType, options *ClientOptions) (*ServiceClient, error) {
	endpoint, err := ska.GetEndpoints()
	if err != nil {
		return nil, err
	}
	credentials, err := ska.GetCredentials()
	if err != nil {
		return nil, err
	}
	return NewServiceClient(endpoint.String(), credentials, options)
}

func (ska *SharedKeyCredentials) GetCredentials() (azcore.Credential, error)  {
	return ska.newSharedKeyCredential()
}

func (ska *SharedKeyCredentials) GetEndpoints() (*url.URL, error)  {
	return url.Parse(ska.ProtoType() + "//" + ska.AccountName + ".blob.core.windows.net/")
}

func (ska *SharedKeyCredentials) getAccountInfo() (string, string) {
	return os.Getenv(AccountName), os.Getenv(AccountKey)
}

// newSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func (ska *SharedKeyCredentials) newSharedKeyCredential() (*SharedKeyCredentials, error) {
	if err := ska.setAccountInfo(); err != nil {
		return nil, err
	}
	return ska, nil
}

// setAccountInfo sets account name and account key.
func (ska *SharedKeyCredentials) setAccountInfo() error {
	accountName, accountKey := ska.getAccountInfo()
	ska.AccountName = accountName
	byts, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return fmt.Errorf("decode account key: %w", err)
	}
	ska.AccountKey.Store(byts)
	return nil
}

// computeHMACSHA256 generates a hash signature for an HTTP request or for a SAS.
func (ska *SharedKeyCredentials) computeHMACSHA256(message string) (base64String string) {
	h := hmac.New(sha256.New, ska.AccountKey.Load().([]byte))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (ska *SharedKeyCredentials) buildStringToSign(req *http.Request) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/authentication-for-the-azure-storage-services
	headers := req.Header
	contentLength := headers.Get(azcore.HeaderContentLength)
	if contentLength == "0" {
		contentLength = ""
	}

	canonicalizeResource, err := ska.buildCanonicalizeResource(req.URL)
	if err != nil {
		return "", err
	}

	stringToSign := strings.Join([]string{
		req.Method,
		headers.Get(azcore.HeaderContentEncoding),
		headers.Get(azcore.HeaderContentLanguage),
		contentLength,
		headers.Get(azcore.HeaderContentMD5),
		headers.Get(azcore.HeaderContentType),
		"", // Empty date because x-ms-date is expected (as per web page above)
		headers.Get(azcore.HeaderIfModifiedSince),
		headers.Get(azcore.HeaderIfMatch),
		headers.Get(azcore.HeaderIfNoneMatch),
		headers.Get(azcore.HeaderIfUnmodifiedSince),
		headers.Get(azcore.HeaderRange),
		ska.buildCanonicalizeHeader(headers),
		canonicalizeResource,
	}, "\n")
	return stringToSign, nil
}

func (ska *SharedKeyCredentials) buildCanonicalizeHeader(headers http.Header) string {
	cm := map[string][]string{}
	for k, v := range headers {
		headerName := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(headerName, "x-ms-") {
			cm[headerName] = v // NOTE: the value must not have any whitespace around it.
		}
	}
	if len(cm) == 0 {
		return ""
	}

	keys := make([]string, 0, len(cm))
	for key := range cm {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	ch := bytes.NewBufferString("")
	for i, key := range keys {
		if i > 0 {
			ch.WriteRune('\n')
		}
		ch.WriteString(key)
		ch.WriteRune(':')
		ch.WriteString(strings.Join(cm[key], ","))
	}
	return string(ch.Bytes())
}

func (ska *SharedKeyCredentials) buildCanonicalizeResource(u *url.URL) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/authentication-for-the-azure-storage-services
	cr := bytes.NewBufferString("/")
	cr.WriteString(ska.AccountName)

	if len(u.Path) > 0 {
		// Any portion of the CanonicalizeResource string that is derived from
		// the resource's URI should be encoded exactly as it is in the URI.
		// -- https://msdn.microsoft.com/en-gb/library/azure/dd179428.aspx
		cr.WriteString(u.EscapedPath())
	} else {
		// a slash is required to indicate the root path
		cr.WriteString("/")
	}

	// params is a map[string][]string; param name is key; params values is []string
	params, err := url.ParseQuery(u.RawQuery) // Returns URL decoded values
	if err != nil {
		return "", fmt.Errorf("failed to parse query params: %w", err)
	}

	if len(params) > 0 { // There is at least 1 query parameter
		var paramNames []string // We use this to sort the parameter key names
		for paramName := range params {
			paramNames = append(paramNames, paramName) // paramNames must be lowercase
		}
		sort.Strings(paramNames)

		for _, paramName := range paramNames {
			paramValues := params[paramName]
			sort.Strings(paramValues)

			// Join the sorted key values separated by ','
			// Then prepend "keyName:"; then add this string to the buffer
			cr.WriteString("\n" + paramName + ":" + strings.Join(paramValues, ","))
		}
	}
	return string(cr.Bytes()), nil
}

// AuthenticationPolicy implements the Credential interface on SharedKeyCredential.
func (ska *SharedKeyCredentials) AuthenticationPolicy(azcore.AuthenticationPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		// Add a x-ms-date header if it doesn't already exist
		if d := req.Request.Header.Get(azcore.HeaderXmsDate); d == "" {
			req.Request.Header.Set(azcore.HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))
		}
		stringToSign, err := ska.buildStringToSign(req.Request)
		if err != nil {
			return nil, err
		}
		signature := ska.computeHMACSHA256(stringToSign)
		authHeader := strings.Join([]string{"SharedKey ", ska.AccountName, ":", signature}, "")
		req.Request.Header.Set(azcore.HeaderAuthorization, authHeader)

		response, err := req.Next()
		if err != nil && response != nil && response.StatusCode == http.StatusForbidden {
			// Service failed to authenticate request, log it
			azcore.Log().Write(azcore.LogResponse, "===== HTTP Forbidden status, String-to-Sign:\n"+stringToSign+"\n===============================\n")
		}
		return response, err
	})
}
