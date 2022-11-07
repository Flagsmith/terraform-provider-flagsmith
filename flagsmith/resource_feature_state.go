package flagsmith

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = featureStateResourceType{}
var _ resource.Resource = featureStateResource{}
var _ resource.ResourceWithImportState = featureStateResource{}

type featureStateResourceType struct{}

func (t featureStateResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Feature state/ Remote config value associated with an environment",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the featurestate",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.NumberType,
			},
			"environment_key": {
				Required:            true,
				MarkdownDescription: "Client side environment key associated with the environment",
				Type:                types.StringType,
			},
			"feature_name": {
				Required:            true,
				MarkdownDescription: "Name of the feature",
				Type:                types.StringType,
			},
			"feature_state_value": {
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
				}),
			},

			"enabled": {
				MarkdownDescription: "Used for enabling/disabling the feature",
				Required:            true,
				Type:                types.BoolType,
			},
			"feature": {
				MarkdownDescription: "ID of the feature",
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.NumberType,
			},
			"environment": {
				MarkdownDescription: "ID of the environment",
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.NumberType,
			},
		},
	}, nil
}

type featureStateResource struct {
	provider fsProvider
}

func (t featureStateResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource , diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return featureStateResource{
		provider: provider,
	}, diags
}

func (r featureStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeatureStateResourceData

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
func (r featureStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureStateResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}
	featureState, err := r.provider.client.GetEnvironmentFeatureState(data.EnvironmentKey.Value, data.FeatureName.Value)
	if err != nil {
		panic(err)
	}
	resourceData := MakeFeatureStateResourceDataFromClientFS(featureState)

	resourceData.EnvironmentKey = data.EnvironmentKey
	resourceData.FeatureName = data.FeatureName

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r featureStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan FeatureStateResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state FeatureStateResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	intFeatureStateID, _ := state.ID.Value.Int64()
	intFeature, _ := state.Feature.Value.Int64()
	intEnvironment, _ := state.Environment.Value.Int64()
	clientFeatureState := plan.ToClientFS(intFeatureStateID, intFeature, intEnvironment)

	updatedClientFS, err := r.provider.client.UpdateFeatureState(clientFeatureState)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature state, got error: %s", err))
		return
	}
	resourceData := MakeFeatureStateResourceDataFromClientFS(updatedClientFS)
	resourceData.EnvironmentKey = plan.EnvironmentKey
	resourceData.FeatureName = plan.FeatureName

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r featureStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Since deleting a feature state does not make sense, we do nothing
	// TODO: maybe reset to the default feature values
	resp.State.RemoveResource(ctx)
	return

}

func (r featureStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importKey := strings.Split(req.ID, ",")
	if len(importKey) != 2 || importKey[0] == "" || importKey[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: environment,feature_name Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_key"), importKey[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("feature_name"), importKey[1])...)

}
