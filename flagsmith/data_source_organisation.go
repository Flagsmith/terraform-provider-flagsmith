package flagsmith

import (
	"context"
	"fmt"
	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &organisationDataResource{}

func newOrganisationDataResource() datasource.DataSource {
	return &organisationDataResource{}
}

type organisationDataResource struct {
	client *flagsmithapi.Client
}

func (o *organisationDataResource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organisation"
}

func (o *organisationDataResource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	o.client = client
}
func (o *organisationDataResource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Organisation",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "UUID of the organisation",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the organisation",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the organisation",
			},
			"force_2fa": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, signup will require 2FA",
			},
			"persist_trait_data": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If false, trait data for this organisation identities will not stored",
			},
			"restrict_project_create_to_admin": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, only organisation admin can create projects",
			},
		},
	}
}
func (o *organisationDataResource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganisationResourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	organisation, err := o.client.GetOrganisationByUUID(data.UUID.ValueString())
	if err != nil {
		panic(err)

	}
	resourceData := MakeOrganisationResourceDataFromClientOrganisation(organisation)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}
