package flagsmith

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = flagResourceType{}
var _ tfsdk.Resource = flagResource{}
var _ tfsdk.ResourceWithImportState = flagResource{}

type flagResourceType struct{}

func (t flagResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Feature State resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the featurestate",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.NumberType,
			},
			"environment_key": {
				Required:            true,
				MarkdownDescription: "API key for the environment",
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
						MarkdownDescription: "Type of the feature state value",
						Required:            true,
					},
					"string_value": {
						Type:     types.StringType,
						Optional: true,
					},
					"integer_value": {
						Type:     types.NumberType,
						Optional: true,
					},
					"boolean_value": {
						Type:     types.BoolType,
						Optional: true,
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
				Required:            true,
				Type:                types.NumberType,
			},
			"environment": {
				MarkdownDescription: "ID of the environment",
				Required:            true,
				Type:                types.NumberType,
			},
		},
	}, nil
}

type flagResource struct {
	provider provider
}

func (t flagResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return flagResource{
		provider: provider,
	}, diags
}

func (r flagResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data FlagResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Read and load the state of the object
	readResponse := tfsdk.ReadResourceResponse{State: resp.State}
	r.Read(ctx, tfsdk.ReadResourceRequest{
		State: tfsdk.State{
			Raw:    req.Plan.Raw,
			Schema: req.Plan.Schema,
		},
		ProviderMeta: req.ProviderMeta,
	}, &readResponse)
	if readResponse.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResponse.Diagnostics...)
		tflog.Error(ctx, "Error reading resource state")
		return
	}
	// Log the state
	elog := fmt.Sprintf("%+v", readResponse.State.Get(ctx, &FlagResourceData{}))

	tflog.Debug(ctx, elog)

	//Now call update to update the state
	updateResponse := tfsdk.UpdateResourceResponse{State: resp.State}
	r.Update(ctx, tfsdk.UpdateResourceRequest{
		Config:       req.Config,
		Plan:         req.Plan,
		State:        readResponse.State,
		ProviderMeta: req.ProviderMeta,
	}, &updateResponse)
	if updateResponse.Diagnostics.HasError() {
		resp.Diagnostics.Append(updateResponse.Diagnostics...)
		tflog.Error(ctx, "Error updating resource state")
		return
	}
	resp.State = updateResponse.State
	resp.Diagnostics.Append(updateResponse.Diagnostics...)

}
func (r flagResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data FlagResourceData
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
	resoureData := MakeFlagResourceDataFromClientFS(featureState)
	resoureData.EnvironmentKey = data.EnvironmentKey
	resoureData.FeatureName = data.FeatureName
	elog := fmt.Sprintf("%+v", resoureData.FeatureStateValue)
	tflog.Debug(ctx, elog)
	diags = resp.State.Set(ctx, &resoureData)
	if diags.HasError() {
		// Log error from diags
		for _, diag := range diags {
			tflog.Error(ctx, diag.Detail())
		}
		resp.Diagnostics.Append(diags...)
		tflog.Error(ctx, "Read: Error setting resource state")
		return
	}
	resp.Diagnostics.Append(diags...)
}

func (r flagResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan FlagResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state FlagResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	intFeatureStateID, _ := state.ID.Value.Int64()
	clientFeatureState := plan.ToClientFS(intFeatureStateID)

	updatedClientFS, err := r.provider.client.UpdateFeatureState(clientFeatureState)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature state, got error: %s", err))
		return
	}
	resoureData := MakeFlagResourceDataFromClientFS(updatedClientFS)
	resoureData.EnvironmentKey = plan.EnvironmentKey
	resoureData.FeatureName = plan.FeatureName

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resoureData)
	resp.Diagnostics.Append(diags...)
}

func (r flagResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	resp.Diagnostics.AddError("Not implemented", "Delete is not implemented; Please use `terraform state rm <address>` to remove the resource")
	return

}

func (r flagResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: attr_one,attr_two. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("environment_key"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("feature_name"), idParts[1])...)

}
