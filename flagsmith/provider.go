package flagsmith

import (
	"context"
	"fmt"
	"os"

	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const BaseAPIURL = "https://api.flagsmith.com/api/v1"

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &fsProvider{}

type fsProvider struct {
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

func (p *fsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "flagsmith"
	resp.Version = p.version
}

func (p *fsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	var masterAPIKey string
	if data.MasterAPIKey.IsUnknown() {
		resp.Diagnostics.AddError("Unable to find master_api_key", "Cannot use unknown value for master_api_key")
		return
	}
	if data.MasterAPIKey.IsNull() {
		masterAPIKey = os.Getenv("FLAGSMITH_MASTER_API_KEY")

	} else {
		masterAPIKey = data.MasterAPIKey.ValueString()
	}
	if masterAPIKey == "" {
		resp.Diagnostics.AddError("Unable to find master_api_key", "master_api_key cannot be an empty string")
	}

	baseAPIURL := BaseAPIURL
	if data.BaseAPIURL.ValueString() != "" {
		baseAPIURL = data.BaseAPIURL.ValueString()

	}

	client := flagsmithapi.NewClient(masterAPIKey, baseAPIURL)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *fsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newFeatureResource,
		newFeatureStateResource,
		newSegmentResource,
		newMultivariateResource,
		newTagResource,
	}

}

func (p *fsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// Does not define any data source
	return []func() datasource.DataSource{}
}

func (p *fsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The flagsmith provider is used  to interact with the resource supported by Flagsmith.
				      The provider needs to be configured with the proper credentials before it can be used.`,
		Attributes: map[string]schema.Attribute{
			"master_api_key": schema.StringAttribute{
				MarkdownDescription: "Master API key used by flagsmith api client. Can also be set using the environment variable `FLAGSMITH_MASTER_API_KEY`",
				Optional:            true,
				Sensitive:           true,
			},
			"base_api_url": schema.StringAttribute{
				MarkdownDescription: "Used by api client to connect to flagsmith instance. NOTE: update this if you are running a self hosted version. e.g: https://your.flagsmith.com/api/v1",
				Optional:            true,
			},
		},
	}
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
