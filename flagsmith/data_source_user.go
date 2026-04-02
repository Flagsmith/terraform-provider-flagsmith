package flagsmith

import (
	"context"
	"fmt"
	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &userDataResource{}

func newUserDataResource() datasource.DataSource {
	return &userDataResource{}
}

type userDataResource struct {
	client *flagsmithapi.Client
}

func (o *userDataResource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (o *userDataResource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (o *userDataResource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Use this data source to look up a Flagsmith user by email within an organisation.",

		Attributes: map[string]schema.Attribute{
			"organisation_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the organisation the user belongs to",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email address of the user",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the user",
			},
			"first_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "First name of the user",
			},
			"last_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last name of the user",
			},
			"role": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Role of the user in the organisation",
			},
		},
	}
}
func (o *userDataResource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserResourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	user, err := o.client.GetOrganisationUserByEmail(data.OrganisationID.ValueInt64(), data.Email.ValueString())
	if err != nil {
		panic(err)

	}
	resourceData := MakeUserResourceDataFromClientUser(user, data.OrganisationID.ValueInt64())

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}
