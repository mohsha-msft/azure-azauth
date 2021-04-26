package azauth

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const scope = "https://storage.azure.com/.default"
const telemetryInfo = "azsdk-go-generated/<version>"

// connectionOptions contains configuration settings for the Connection's pipeline.
// All zero-value fields will be initialized with their default values.
type connectionOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transport
	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions
	// Logging configures the built-in logging policy behavior.
	Logging azcore.LogOptions
	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []azcore.Policy
	// PerRetryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry request.
	PerRetryPolicies []azcore.Policy
}

func (c *connectionOptions) telemetryOptions() *azcore.TelemetryOptions {
	to := c.Telemetry
	if to.Value == "" {
		to.Value = telemetryInfo
	} else {
		to.Value = fmt.Sprintf("%s %s", telemetryInfo, to.Value)
	}
	return &to
}

type Connection struct {
	u string
	p azcore.Pipeline
}

// newConnection creates an instance of the Connection type with the specified endpoint.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func newConnection(endpoint string, cred azcore.Credential, options *connectionOptions) *Connection {
	if options == nil {
		options = &connectionOptions{}
	}
	policies := []azcore.Policy{
		azcore.NewTelemetryPolicy(options.telemetryOptions()),
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, azcore.NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
	policies = append(policies, cred.AuthenticationPolicy(
		azcore.AuthenticationPolicyOptions{
			Options: azcore.TokenRequestOptions{Scopes: []string{scope}},
		}),
	)
	policies = append(policies, azcore.NewLogPolicy(&options.Logging))
	return &Connection{u: endpoint, p: azcore.NewPipeline(options.HTTPClient, policies...)}
}

// Endpoint returns the Connection's endpoint.
func (c *Connection) Endpoint() string {
	return c.u
}

// Pipeline returns the Connection's pipeline.
func (c *Connection) Pipeline() azcore.Pipeline {
	return c.p
}

