package flagsmith

import (
	"context"
	"fmt"
	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &featureResource{}
var _ resource.ResourceWithImportState = &featureResource{}

func newFeatureResource() resource.Resource {
	return &featureResource{}
}

type featureResource struct {
	client *flagsmithapi.Client
}

func (r *featureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature"
}

func (r *featureResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (t *featureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Feature/ Remote config",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the feature",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the feature",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the project",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"feature_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the feature",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of the feature, can be STANDARD, or MULTIVARIATE",
			},
			"default_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Determines if the feature is enabled by default. If unspecified, it will default to false",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"initial_value": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Determines the initial value of the feature.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the feature",
			},
			"is_archived": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Can be used to archive/unarchive a feature. If unspecified, it will default to false",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"owners": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of user IDs representing the owners of the feature.",
			},
			"project_uuid": schema.StringAttribute{
				MarkdownDescription: "UUID of project the feature belongs to",
				Required:            true,
			},
		},
	}
}

func (r *featureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeatureResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	clientFeature := data.ToClientFeature()
	owners := clientFeature.Owners

	// Create the feature
	err := r.client.CreateFeature(clientFeature)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create feature, got error: %s", err))
		return
	}
	if owners != nil && len(*owners) > 0 {
		err := r.client.AddFeatureOwners(clientFeature, *owners)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add owners to feature, got error: %s", err))
			return
		}

	}
	clientFeature.Owners = owners
	resourceData := MakeFeatureResourceDataFromClientFeature(clientFeature)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r *featureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	feature, err := r.client.GetFeature(data.UUID.ValueString())
	if err != nil {
		if _, ok := err.(flagsmithapi.FeatureNotFoundError); ok {
			resp.State.RemoveResource(ctx)
			return
		}
		panic(err)

	}
	// This prevents creating unnecessary plan change(from [] -> nil)
	// when owners is not part of the plan
	if data.Owners == nil && feature.Owners != nil && len(*feature.Owners) == 0 {
		feature.Owners = nil
	}
	resourceData := MakeFeatureResourceDataFromClientFeature(feature)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}

func (r *featureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get plan values
	var plan FeatureResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading plan data")
		return
	}

	// Get current state
	var state FeatureResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading state data")
		return
	}
	// Generate API request body from plan
	clientFeature := plan.ToClientFeature()
	planOwners := clientFeature.Owners

	err := r.client.UpdateFeature(clientFeature)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature, got error: %s", err))
		return
	}

	// Update feature owners
	if planOwners != clientFeature.Owners {
		ownerIDsToRemove := Difference(clientFeature.Owners, planOwners)
		if len(ownerIDsToRemove) > 0 {
		  err := r.client.RemoveFeatureOwners(clientFeature, ownerIDsToRemove)
		  if err != nil {
			  resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove feature owners, got error: %s", err))
			  return
		  }
		}
		ownerIDsToAdd := Difference(planOwners, clientFeature.Owners)

		if len(ownerIDsToAdd) > 0 {
			err := r.client.AddFeatureOwners(clientFeature, ownerIDsToAdd)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add feature owners, got error: %s", err))
				return
			}
		}
	}
	clientFeature.Owners = planOwners
	resourceData := MakeFeatureResourceDataFromClientFeature(clientFeature)

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resourceData)

	resp.Diagnostics.Append(diags...)

}


func (r *featureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state FeatureResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Delete: Error reading state data")
		return
	}
	// Generate API request body from plan
	clientFeature := state.ToClientFeature()

	err := r.client.DeleteFeature(*clientFeature.ProjectID, *clientFeature.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete feature, got error: %s", err))
		return
	}
	resp.State.RemoveResource(ctx)

}
func (r *featureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
