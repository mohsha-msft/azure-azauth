package azauth

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
)

type ClientOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transport
	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions
}

func (o *ClientOptions) getConnectionOptions() *connectionOptions {
	if o == nil {
		return nil
	}

	return &connectionOptions{
		HTTPClient: o.HTTPClient,
		Retry:      o.Retry,
		Telemetry:  o.Telemetry,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type client struct {
	connection 	*Connection
	credential 	azcore.Credential
}

type ServiceClient struct {
	client 		*client
	url 		string
}

type Client struct {
	*ServiceClient
	*ContainerClient
	*BlobClient
}

func (o *GetClientOptions) GetClient() (*Client, error) {
	c := Client{}
	switch o.ClientType {
	case EClientType.ServiceClient():
		svcClient, err := NewServiceClient(o.Endpoint, o.Credential, o.ClientOptions)
		if err != nil {
			return nil, err
		}
		c.ServiceClient = svcClient
	case EClientType.ContainerClient():
		containerClient, err := NewContainerClient(o.Endpoint, o.Credential, o.ClientOptions)
		if err != nil {
			return nil, err
		}
		c.ContainerClient = containerClient
	case EClientType.BlobClient():
		blobClient, err := NewBlobClient(o.Endpoint, o.Credential, o.ClientOptions)
		if err != nil {
			return nil, err
		}
		c.BlobClient = blobClient
	default:
		panic("Not Implemented")
	}
	return &c, nil
}

func NewServiceClient(endpoint string, credential azcore.Credential, options *ClientOptions) (*ServiceClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	connection := newConnection(endpoint, credential, options.getConnectionOptions())
	return &ServiceClient{
		client: &client{
			connection: connection,
			credential: credential,
		},
		url: u.String(),
	}, nil
}

type ContainerClient struct {
	client 		*client
	url 		string
}

func (s ServiceClient) NewContainerClient(containerName string) (*ContainerClient, error) {
	containerURL := appendToURLPath(s.url, containerName)
	containerConnection := &Connection{containerURL, s.client.connection.p}
	return &ContainerClient{
		client: &client{
			connection: containerConnection,
			credential: s.client.credential,
		},
		url: containerURL,
	}, nil
}

func NewContainerClient(endpoint string, credential azcore.Credential, options *ClientOptions) (*ContainerClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	connection := newConnection(endpoint, credential, options.getConnectionOptions())
	return &ContainerClient{
		client: &client{
			connection: connection,
			credential: credential,
		},
		url: u.String(),
	}, nil
}

type BlobClient struct {
	client 		*client
	url 		string
}

func NewBlobClient(endpoint string, credential azcore.Credential, options *ClientOptions) (*BlobClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	connection := newConnection(endpoint, credential, options.getConnectionOptions())
	return &BlobClient{
		client: &client{
			connection: connection,
			credential: credential,
		},
		url: u.String(),
	}, nil
}




