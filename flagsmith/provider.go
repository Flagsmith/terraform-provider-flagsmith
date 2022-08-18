package flagsmith

import (
	"context"
	"fmt"
	"os"

	"github.com/Flagsmith/flagsmith-go-api-client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const BaseAPIURL = "https://api.flagsmith.com/api/v1"

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &fsProvider{}

type fsProvider struct {
	// client contains the upstream provider SDK used to
	// communicate with the flagsmith api
	client *flagsmithapi.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData is used to store data from the Terraform configuration.
type providerData struct {
	MasterAPIKey types.String `tfsdk:"master_api_key"`
	BaseAPIURL   types.String `tfsdk:"base_api_url"`
}

func (p *fsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	var masterAPIKey string
	if data.MasterAPIKey.Unknown {
		resp.Diagnostics.AddError("Unable to find master_api_key", "Cannot use unknown value for master_api_key")
		return
	}
	if data.MasterAPIKey.Null {
		masterAPIKey = os.Getenv("FLAGSMITH_MASTER_API_KEY")

	} else {
		masterAPIKey = data.MasterAPIKey.Value
	}
	if masterAPIKey == "" {
		resp.Diagnostics.AddError("Unable to find master_api_key", "master_api_key cannot be an empty string")
	}

	baseAPIURL := BaseAPIURL
	if data.BaseAPIURL.Value != "" {
		baseAPIURL = data.BaseAPIURL.Value

	}

	client := flagsmithapi.NewClient(masterAPIKey, baseAPIURL)
	p.client = client

	p.configured = true

	return
}

func (p *fsProvider) GetResources(ctx context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"flagsmith_flag": flagResourceType{},
	}, nil
}

func (p *fsProvider) GetDataSources(ctx context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	// Does not define any data source
	return map[string]provider.DataSourceType{}, nil
}

func (p *fsProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `The flagsmith provider is used to enable/disable and/or update flag values.`,
		Attributes: map[string]tfsdk.Attribute{
			"master_api_key": {
				MarkdownDescription: "Master API key used by flagsmith api client. Can also be set using the environment variable `FLAGSMITH_MASTER_API_KEY`",
				Optional:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"base_api_url": {
				MarkdownDescription: "Used by api client to connect to flagsmith instance. NOTE: update this if you are running a self hosted version. e.g: https://your.flagsmith.com/api/v1",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &fsProvider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in provider.Provider) (fsProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*fsProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return fsProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return fsProvider{}, diags
	}

	return *p, diags
}
