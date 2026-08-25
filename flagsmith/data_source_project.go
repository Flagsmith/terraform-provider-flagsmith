package flagsmith

import (
	"context"
	"fmt"

	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &projectDataResource{}

func newProjectDataResource() datasource.DataSource {
	return &projectDataResource{}
}

type projectDataResource struct {
	client *flagsmithapi.Client
}

func (p *projectDataResource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (p *projectDataResource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	p.client = client
}

func (p *projectDataResource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Project",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the project",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the project",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the project",
			},
			"organisation_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the organisation project belongs to",
			},
			"hide_disabled_flags": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true will exclude flags from SDK which are disabled",
			},
			"prevent_flag_defaults": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Prevent defaults from being set in all environments when creating a feature.",
			},
			"enable_realtime_updates": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, a realtime(sse) event is triggered whenever the value of a flag changes",
			},
			"only_allow_lower_case_feature_names": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Used by UI to validate feature names",
			},
			"feature_name_regex": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Used for validating feature names",
			},
			"stale_flags_limit_days": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of days without modification in any environment before a flag is considered stale.",
			},
			"enforce_feature_owners": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, feature creation requires at least one owner or group owner.",
			},
		},
	}
}

func (p *projectDataResource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectResourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	project, err := p.client.GetProject(data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	}
	resourceData := MakeProjectResourceDataFromClientProject(project)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}
