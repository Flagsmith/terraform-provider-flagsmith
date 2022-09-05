package flagsmith

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = featureResourceType{}
var _ resource.Resource = featureResource{}
var _ resource.ResourceWithImportState = featureResource{}

type featureResourceType struct{}



func (t featureResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith feature/ Remote config",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the feature",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.NumberType,
			},
			"project_id": {
				Computed:            true,
				MarkdownDescription: "ID of the project",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.NumberType,
			},
			"feature_name": {
				Required:            true,
				MarkdownDescription: "Name of the feature",
				Type:                types.StringType,
			},
			"default_enabled": {
				Optional:            true,
				MarkdownDescription: "Determines if the feature is enabled by default",
				Type:                types.BoolType,
			},
			"initial_value": {
				Optional:            true,
				MarkdownDescription: "Determines the initial value of the feature",
				Type:                types.StringType,
			},
			"description": {
				Optional:            true,
				MarkdownDescription: "Description of the feature",
				Type:                types.StringType,
			},
			"is_archived": {
				Optional:            true,
				MarkdownDescription: "Can be used to archive/unarchive a feature",
				Type:                types.BoolType,
			},
			"multivariate_options": {
				Optional: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"type": {
						Type:                types.StringType,
						MarkdownDescription: "Type of the feature state value, can be `unicode`, `int` or `bool`",
						Optional:            true,
					},
					"string_value": {
						Type:                types.StringType,
						MarkdownDescription: "String value of the feature if the type is `unicode`",
						Optional:            true,
					},
					"integer_value": {
						Type:                types.NumberType,
						MarkdownDescription: "Integer value of the feature if the type is `int`",
						Optional:            true,
					},
					"boolean_value": {
						Type:                types.BoolType,
						MarkdownDescription: "Boolean value of the feature if the type is `bool`",
						Optional:            true,
					},
					"default_percentage_allocation": {
						Type:                types.NumberType,
						MarkdownDescription: "Percentage allocation of the current multivariate option",
						Optional:            true,
					},
				}),
			},

			"project_uuid": {
				MarkdownDescription: "UUID of project the feature belongs to",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

type featureResource struct {
	provider fsProvider
}

func (t featureResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource , diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return featureResource{
		provider: provider,
	}, diags
}

func (r featureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FlagResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Read and load the state of the object
	readResponse := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{
		State: tfsdk.State{
			Raw:    req.Plan.Raw,
			Schema: req.Plan.Schema,
		},
		ProviderMeta: req.ProviderMeta,
	}, &readResponse)

	if readResponse.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResponse.Diagnostics...)
		tflog.Error(ctx, "Create: Error reading resource state")
		return
	}

	//Now call update to update the state
	updateResponse := resource.UpdateResponse{State: resp.State}
	r.Update(ctx, resource.UpdateRequest{
		Config:       req.Config,
		Plan:         req.Plan,
		State:        readResponse.State,
		ProviderMeta: req.ProviderMeta,
	}, &updateResponse)

	if updateResponse.Diagnostics.HasError() {
		resp.Diagnostics.Append(updateResponse.Diagnostics...)
		tflog.Error(ctx, "Create: Error updating resource state")
		return
	}

	resp.State = updateResponse.State
	resp.Diagnostics.Append(updateResponse.Diagnostics...)

}

func (r featureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	feature, err := r.provider.client.GetFeature(data.EnvironmentKey.Value, data.FeatureName.Value)
	if err != nil {
		panic(err)
	}
	resourceData := MakeFlagResourceDataFromClientFS(featureState)

	resourceData.EnvironmentKey = data.EnvironmentKey
	resourceData.FeatureName = data.FeatureName

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}


func (r featureResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	//	var plan FeatureResourceData
	// diags := req.Plan.Get(ctx, &plan)
	// resp.Diagnostics.Append(diags...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Get current state
	// var state FlagResourceData
	// diags = req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Generate API request body from plan
	// intFeatureStateID, _ := state.ID.Value.Int64()
	// intFeature, _ := state.Feature.Value.Int64()
	// intEnvironment, _ := state.Environment.Value.Int64()
	// clientFeatureState := plan.ToClientFS(intFeatureStateID, intFeature, intEnvironment)

	// updatedClientFS, err := r.provider.client.UpdateFeatureState(clientFeatureState)

	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature state, got error: %s", err))
	// 	return
	// }
	// resourceData := MakeFlagResourceDataFromClientFS(updatedClientFS)
	// resourceData.EnvironmentKey = plan.EnvironmentKey
	// resourceData.FeatureName = plan.FeatureName

	// Update the state with the new values
	// diags = resp.State.Set(ctx, &resourceData)
	// resp.Diagnostics.Append(diags...)

}

func (r featureResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Since deleting a feature state does not make sense, we do nothing
	// TODO: maybe reset to the default feature values
	// resp.State.RemoveResource(ctx)
	// return


}

func (r featureResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// importKey := strings.Split(req.ID, ",")
	// if len(importKey) != 2 || importKey[0] == "" || importKey[1] == "" {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Import Identifier",
	// 		fmt.Sprintf("Expected import identifier with format: environment,feature_name Got: %q", req.ID),
	// 	)
	// 	return
	// }
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_key"), importKey[0])...)
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("feature_name"), importKey[1])...)

}
