package connector

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

func Run() {
	azidentity.NewChainedTokenCredential()
	cred, _ := azidentity.NewInteractiveBrowserCredential(azidentity.InteractiveBrowserCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			APIVersion:                      "",
			Cloud:                           cloud.Configuration{},
			InsecureAllowCredentialWithHTTP: false,
			Logging:                         policy.LogOptions{},
			Retry:                           policy.RetryOptions{},
			Telemetry:                       policy.TelemetryOptions{},
			TracingProvider:                 tracing.Provider{},
			Transport:                       nil,
			PerCallPolicies:                 nil,
			PerRetryPolicies:                nil,
		},
		AdditionallyAllowedTenants: nil,
		ClientID:                   "",
		DisableInstanceDiscovery:   false,
		LoginHint:                  "",
		RedirectURL:                "",
		TenantID:                   "",
	})
	cred.
		client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"Files.Read"})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
}
