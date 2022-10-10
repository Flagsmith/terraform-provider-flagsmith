package flagsmith

import (
	"context"
	"strings"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = multivariateResourceType{}
var _ resource.Resource = multivariateResource{}
var _ resource.ResourceWithImportState = multivariateResource{}

type multivariateResourceType struct{}

func (t multivariateResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Feature Multivariate Option",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the multivariate option",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.NumberType,
			},
			"uuid": {
				Computed:            true,
				MarkdownDescription: "UUID of the multivariate option",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},

			"type": {
				Type:                types.StringType,
				MarkdownDescription: "Type of the multivariate option can be `unicode`, `int` or `bool`",
				Required:            true,
			},
			"string_value": {
				Type:                types.StringType,
				MarkdownDescription: "String value of the multivariate option if the type is `unicode`",
				Optional:            true,
			},
			"integer_value": {
				Type:                types.NumberType,
				MarkdownDescription: "Integer value of the multivariate option if the type is `int`",
				Optional:            true,
			},
			"boolean_value": {
				Type:                types.BoolType,
				MarkdownDescription: "Boolean value of the multivariate option if the type is `bool`",
				Optional:            true,
			},
			"default_percentage_allocation": {
				Type:                types.NumberType,
				MarkdownDescription: "Percentage allocation of the current multivariate option",
				Required:            true,
			},
			"feature_id": {
				Computed:            true,
				MarkdownDescription: "ID of the feature to which the multivariate option belongs",
				Type:                types.NumberType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"feature_uuid": {
				MarkdownDescription: "UUID of the feature to which the multivariate option belongs",
				Required:            true,
				Type:                types.StringType,
			},
			"project_id": {
				Computed: 	  true,
				MarkdownDescription: "Project ID of the feature to which the multivariate option belongs",
				Type:                types.NumberType,
			},
		},
	}, nil
}

type multivariateResource struct {
	provider fsProvider
}
func (t multivariateResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return multivariateResource{
		provider: provider,
	}, diags
}



func (r multivariateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MultivariateOptionResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	mvOption := data.ToClientMultivariateOption()

	err := r.provider.client.CreateFeatureMVOption(mvOption)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create feature multivariate option, got error: %s", err))
		return
	}

	resourceData := NewMultivariateOptionFromClientOption(mvOption)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}

func (r multivariateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MultivariateOptionResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	mvOption, err := r.provider.client.GetFeatureMVOption(data.FeatureUUID.Value, data.UUID.Value)
	if err != nil {
		panic(err)
	}
	resourceData := NewMultivariateOptionFromClientOption(mvOption)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}


func (r multivariateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan MultivariateOptionResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading plan data")
		return
	}

	// Get current state
	var state MultivariateOptionResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading state data")
		return
	}

	// Update state with plan values
	state.Type = plan.Type
	state.StringValue = plan.StringValue
	state.IntegerValue = plan.IntegerValue
	state.BooleanValue = plan.BooleanValue
	state.DefaultPercentageAllocation = plan.DefaultPercentageAllocation

	// Generate API request body from plan
	mvOption := state.ToClientMultivariateOption()

	err := r.provider.client.UpdateFeatureMVOption(mvOption)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature multivariate option, got error: %s", err))
		return
	}

	resourceData := NewMultivariateOptionFromClientOption(mvOption)

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resourceData)

	resp.Diagnostics.Append(diags...)

}


func (r multivariateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//Get current state
	var state MultivariateOptionResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Delete: Error reading state data")
		return
	}
	// Generate API request body from plan
	mvOption:= state.ToClientMultivariateOption()

	err := r.provider.client.DeleteFeatureMVOption(*mvOption.ProjectID, *mvOption.FeatureID, mvOption.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete feature multivariate option, got error: %s", err))
		return
	}
	resp.State.RemoveResource(ctx)

}


func (r multivariateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importKey := strings.Split(req.ID, ",")
	if len(importKey) != 2 || importKey[0] == "" || importKey[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: feature_uuid,mv_option_uuid Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("feature_uuid"), importKey[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), importKey[1])...)


}
