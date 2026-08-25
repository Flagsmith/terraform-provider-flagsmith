package flagsmith

import (
	"context"
	"fmt"

	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &featureDataResource{}

func newFeatureDataResource() datasource.DataSource {
	return &featureDataResource{}
}

type featureDataResource struct {
	client *flagsmithapi.Client
}

func (f *featureDataResource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature"
}

func (f *featureDataResource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*flagsmithapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *flagsmithapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	f.client = client
}

func (f *featureDataResource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Feature/ Remote config. Use this to reference a feature that is managed outside of Terraform, for example to attach an environment specific `flagsmith_feature_state` override to it.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the feature",
			},
			"feature_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the feature",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the feature",
			},
			"project_uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the project the feature belongs to",
			},
			"project_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the project the feature belongs to",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of the feature, either STANDARD or MULTIVARIATE",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the feature",
			},
			"initial_value": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Initial value of the feature",
			},
			"default_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Determines if the feature is enabled by default",
			},
			"is_archived": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Determines if the feature is archived",
			},
			"owners": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of user IDs representing the owners of the feature.",
			},
			"group_owners": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of group IDs representing the group owners of the feature.",
			},
			"tags": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of tag IDs representing the tags attached to the feature.",
			},
		},
	}
}

func (f *featureDataResource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FeatureResourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	feature, err := f.client.GetFeature(data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read feature, got error: %s", err))
		return
	}
	resourceData := MakeFeatureResourceDataFromClientFeature(feature)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}
