package flagsmith

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"math/big"
	"strconv"
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
	var data flagResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.CreateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.Number{Value: big.NewFloat(42)}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
func (r flagResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data flagResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.UpdateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r flagResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

	var data flagResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	// Early return if the state is wrong
	if diags.HasError() {
		return
	}
	featureStateID := data.ID.Value

	// TODO: some error handling
	intFeatureStateID, _  := featureStateID.Int64()
	featureState, err := r.provider.client.GetFeatureState(intFeatureStateID)
	if err != nil {
		panic(err)
	}
	fsValue := FeatureStateValueType{
		Type:         types.String{Value: featureState.FeatureStateValue.Type},
		StringValue:  types.String{Value: featureState.FeatureStateValue.StringValue},
		IntegerValue: types.Number{Value: big.NewFloat(float64(featureState.FeatureStateValue.IntegerValue))},
		BooleanValue: types.Bool{Value: featureState.FeatureStateValue.BooleanValue},
	}
	tflog.Trace(ctx, "response: "+string(featureState.ID))
	var result = flagResourceData{
		ID:                types.Number{Value: featureStateID},
		Enabled:           types.Bool{Value: featureState.Enabled},
		FeatureStateValue: &fsValue,
		Feature:           types.Number{Value: big.NewFloat(float64(featureState.Feature))},
		Environment:       types.Number{Value: big.NewFloat(float64(featureState.Environment))},
	}

	diags = resp.State.Set(ctx, &result)
	if diags.HasError() {
		return
	}
	resp.Diagnostics.Append(diags...)
}

func (r flagResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data flagResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.DeleteExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r flagResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	featureStateID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Import ID must be an integer")
		return
	}

	fsID := types.Number{Value: big.NewFloat(float64(featureStateID))}
	// Add ID to the state
	diags := resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &fsID)

	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(diags...)
}
