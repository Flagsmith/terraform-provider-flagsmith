package flagsmith

import (
	"context"
	"fmt"
	"regexp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &featureStateResource{}
var _ resource.ResourceWithImportState = &featureStateResource{}

func newFeatureStateResource() resource.Resource {
	return &featureStateResource{}
}

type featureStateResource struct {
	client *flagsmithapi.Client
}

func (r *featureStateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_state"
}

func (r *featureStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (t *featureStateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Feature state/ Remote config value",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the featurestate",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the featurestate",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"environment_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Client side environment key associated with the environment",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"feature_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the feature",
				Required:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"feature_state_value": schema.SingleNestedAttribute{
				Required: true,
				MarkdownDescription: "Value for the feature State. NOTE: One of string_value, integer_value or boolean_value must be set",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Type of the feature state value, can be `unicode`, `int` or `bool`",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"unicode", "int", "bool"}...),
						},
					},
					"string_value": schema.StringAttribute{
						MarkdownDescription: "String value of the feature if the type is `unicode`.",
						Optional:            true,
						Validators: []validator.String{
							// Validate string value satisfies the regular expression for no leading or trailing whitespace
							// but allow empty string
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^\S[\s\S]*\S$|^$`),
								"Leading and trailing whitespace is not allowed",
							),
						},
					},
					"integer_value": schema.Int64Attribute{
						MarkdownDescription: "Integer value of the feature if the type is `int`",
						Optional:            true,
					},
					"boolean_value": schema.BoolAttribute{
						MarkdownDescription: "Boolean value of the feature if the type is `bool`",
						Optional:            true,
					},
				},
			},

			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Used for enabling/disabling the feature",
				Required:            true,
			},
			"environment_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the environment",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"segment_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the segment, used for creating segment overrides",
				Optional:            true,
			},
			"segment_priority": schema.Int64Attribute{
				MarkdownDescription: "Priority of the segment overrides.",
				Optional:            true,
				Computed:            true,
			},
			"feature_segment_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the feature_segment, used internally to bind a feature state to a segment",
				Computed:            true,

				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (f *featureStateResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        resourcevalidator.ExactlyOneOf(
            path.MatchRoot("feature_state_value").AtName("string_value"),
            path.MatchRoot("feature_state_value").AtName("integer_value"),
            path.MatchRoot("feature_state_value").AtName("boolean_value"),
        ),
    }
}

func (r *featureStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeatureStateResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Create segment override if segment is set
	if data.Segment.ValueInt64() != 0 {
		clientFeatureState := data.ToClientFS()
		err := r.client.CreateSegmentOverride(clientFeatureState)
		if err != nil {
			resp.Diagnostics.AddError("Error creating segment override", err.Error())
			return
		}
		// set the state with the new values
		resourceData := MakeFeatureStateResourceDataFromClientFS(clientFeatureState)
		diags = resp.State.Set(ctx, &resourceData)
		resp.Diagnostics.Append(diags...)
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
func (r *featureStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureStateResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	var featureState *flagsmithapi.FeatureState
	var err error

	if data.UUID.ValueString() != "" {
		featureState, err = r.client.GetFeatureState(data.UUID.ValueString())

	} else {
		featureState, err = r.client.GetEnvironmentFeatureState(data.EnvironmentKey.ValueString(), data.Feature.ValueInt64())
	}
	if err != nil {
		if _, ok := err.(flagsmithapi.FeatureStateNotFoundError); ok {
			resp.State.RemoveResource(ctx)
			return
		}
		panic(err)

	}
	resourceData := MakeFeatureStateResourceDataFromClientFS(featureState)

	resourceData.EnvironmentKey = data.EnvironmentKey

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r *featureStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	clientFeatureState := plan.ToClientFS()

	// Load computed data from the state
	clientFeatureState.ID = state.ID.ValueInt64()
	clientFeatureState.Feature = state.Feature.ValueInt64()
	intEnvironment := state.Environment.ValueInt64()
	clientFeatureState.Environment = &intEnvironment

	updateSegmentPriority := state.SegmentPriority.ValueInt64() != plan.SegmentPriority.ValueInt64()
	err := r.client.UpdateFeatureState(clientFeatureState, updateSegmentPriority)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature state, got error: %s", err))
		return
	}
	resourceData := MakeFeatureStateResourceDataFromClientFS(clientFeatureState)
	resourceData.EnvironmentKey = plan.EnvironmentKey

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r *featureStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state FeatureStateResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Delete: Error state data")
		return
	}

	// Delete feature segment if it exists
	if state.FeatureSegment.ValueInt64() != 0 {
		err := r.client.DeleteFeatureSegment(state.FeatureSegment.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete feature segment, got error: %s", err))
			return
		}
	}
	resp.State.RemoveResource(ctx)
	return

}

func (r *featureStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importKey := strings.Split(req.ID, ",")
	if len(importKey) != 2 || importKey[0] == "" || importKey[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: environment,feature_state_uuid Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_key"), importKey[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), importKey[1])...)

}
