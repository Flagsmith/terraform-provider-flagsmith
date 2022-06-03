package flagsmith

import (
	"context"
	"fmt"

	"github.com/Flagsmith/flagsmith-go-api-client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const BaseAPIURL = "https://api.flagsmith.com/api/v1"

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.Provider = &provider{}
// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//

	client *flagsmithapi.Client
	// TODO: If appropriate, implement upstream provider SDK or HTTP client.
	// client vendorsdk.ExampleClient

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	MasterAPIKey types.String `tfsdk:"api_key"`
	BaseAPIURL types.String `tfsdk:"base_api_url"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {

	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: do we need this check on master api key?
	if data.MasterAPIKey.Null  || data.MasterAPIKey.Unknown || data.MasterAPIKey.Value == " " {
		resp.Diagnostics.AddError("Unable to find master_api_key", "master_api_key cannot be an empty string")
	}

	baseAPIURL := BaseAPIURL;
	if data.BaseAPIURL.Value != "" {
		baseAPIURL = data.BaseAPIURL.Value

	}

	client := flagsmithapi.NewClient(data.MasterAPIKey.Value, baseAPIURL)
	p.client = client


	// Configuration values are now available.
	// if data.Example.Null { /* ... */ }

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.

	p.configured = true

	return
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"flagsmith_flag": flagResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	// Does not define any data source
	return map[string]tfsdk.DataSourceType{}, nil
	// return map[string]tfsdk.DataSourceType{
	// 	"flag": flagResourceType{},
	// }, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"master_api_key": {
				MarkdownDescription: "Master API key used by flagsmith api client",
				Required:            true,
				Type:                types.StringType,
				Sensitive: true,
			},
			"base_api_url": {
				MarkdownDescription: "Used by api client to connect to flagsmith instance. NOTE: update this if you are running a self hosted version",
				Required: false,
				Optional: true,
				Type: types.StringType,
			},

		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
