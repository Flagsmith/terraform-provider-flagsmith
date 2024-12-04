package flagsmith

import (
	"context"
	"fmt"
	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &environmentResource{}
var _ resource.ResourceWithImportState = &environmentResource{}

func newEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

type environmentResource struct {
	client *flagsmithapi.Client
}

func (r *environmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *environmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}
func (t *environmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Environment",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the environment",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the environment",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the project",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the environment",
			},
			"api_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Client side API Key",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the environment",
			},
			"banner_text": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Banner text to display in the UI",
			},
			"banner_colour": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "hex code for the UI banner colour",
			},
			"minimum_change_request_approvals": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Minimum number of approvals required for a change request",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
				Default:             int64default.StaticInt64(0),
			},
			"hide_disabled_flags": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "If true will exclude flags from SDK which are disabled",
				Default:             booldefault.StaticBool(false),
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"hide_sensitive_data": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "If true, will hide sensitive data(e.g: traits, description etc) from the SDK endpoints",
				Default:             booldefault.StaticBool(false),
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"allow_client_traits": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Allows clients using the client API key to set traits.",
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"use_identity_composite_key_for_hashing": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enable this to have consistent multivariate and percentage split evaluations across all SDKs (in local and server side mode) ",
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnvironmentResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientEnvironment := data.ToClientEnvironment()

	// Create the environment
	err := r.client.CreateEnvironment(clientEnvironment)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create environment, got error: %s", err))
		return
	}
	resourceData := MakeEnvironmentResourceDataFromClientEnvironment(clientEnvironment)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EnvironmentResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	environment, err := r.client.GetEnvironmentByUUID(data.UUID.ValueString())
	if err != nil {
		panic(err)

	}
	resourceData := MakeEnvironmentResourceDataFromClientEnvironment(environment)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}

func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get plan values
	var plan EnvironmentResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading plan data")
		return
	}

	// Generate API request body from plan
	clientEnvironment := plan.ToClientEnvironment()

	err := r.client.UpdateEnvironment(clientEnvironment)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update environment, got error: %s", err))
		return
	}

	resourceData := MakeEnvironmentResourceDataFromClientEnvironment(clientEnvironment)

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resourceData)

	resp.Diagnostics.Append(diags...)

}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state EnvironmentResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Delete: Error reading state data")
		return
	}
	apiKey := state.APIKey.ValueString()
	if apiKey != "" {
		err := r.client.DeleteEnvironment(apiKey)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete environment, got error: %s", err))
			return
		}
	}
	resp.State.RemoveResource(ctx)

}
func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
