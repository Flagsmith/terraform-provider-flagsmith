package flagsmith

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Flagsmith/flagsmith-go-api-client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &segmentResource{}
var _ resource.ResourceWithImportState = &segmentResource{}

func newSegmentResource() resource.Resource {
	return &segmentResource{}
}
type segmentResource struct{
	client *flagsmithapi.Client
}

func (r *segmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (r *segmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (t *segmentResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	conditions := tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
		"property": {
			Optional:            true,
			MarkdownDescription: "Property of the condition",
			Type:                types.StringType,
		},
		"operator": {
			Required:            true,
			MarkdownDescription: "Operator of the condition, can be one of `EQUAL`, `GREATER_THAN`, `LESS_THAN`, `LESS_THAN_INCLUSIVE` `CONTAINS`, `GREATER_THAN_INCLUSIVE`, `NOT_CONTAINS`, `NOT_EQUAL`,  `REGEX`, `PERCENTAGE_SPLIT`,  `MODULO`, `IS_SET`, `IS_NOT_SET`, `IN` ",
			Type:                types.StringType,
		},
		"value": {
			Optional:            true,
			MarkdownDescription: "Value of the condition",
			Type:                types.StringType,
		},
	})

	nestedRules := tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
		"type": {
			Required:            true,
			MarkdownDescription: "Type of the rule",
			Type:                types.StringType,
		},
		"conditions": {
			Optional:            true,
			MarkdownDescription: "List of conditions for the nested rule",
			Attributes:          conditions,
		},
	})

	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Segment",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the segment",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.Int64Type,
			},
			"uuid": {
				Computed:            true,
				MarkdownDescription: "UUID of the segment",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"project_id": {
				Computed:            true,
				MarkdownDescription: "ID of the project",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.Int64Type,
			},
			"feature_id": {
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Set this to create a feature specific segment",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.Int64Type,
			},
			"name": {
				Required:            true,
				MarkdownDescription: "Name of the segment",
				Type:                types.StringType,
			},
			"description": {
				Optional:            true,
				MarkdownDescription: "Description of the segment",
				Type:                types.StringType,
			},
			"project_uuid": {
				MarkdownDescription: "UUID of project the segment belongs to",
				Required:            true,
				Type:                types.StringType,
			},
			"rules": {
				MarkdownDescription: "Rules for the segment",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"type": {
						Required:            true,
						MarkdownDescription: "Type of the rule, can be of: `ALL`, `ANY`, `NONE`",
						Type:                types.StringType,
					},
					"rules": {
						Optional:            true,
						MarkdownDescription: "List of Nested Rules",
						Attributes:          nestedRules,
					},

					"conditions": {
						Optional:            true,
						MarkdownDescription: "Conditions for the rule",
						Attributes:          conditions,
					},
				}),
			},
		},
	}, nil

}



func (r *segmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SegmentResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	clientSegment := data.ToClientSegment()

	// Create the segment
	err := r.client.CreateSegment(clientSegment)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create segment, got error: %s", err))
		return
	}
	resourceData := MakeSegmentResourceDataFromClientSegment(clientSegment)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}

func (r *segmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SegmentResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	segment, err := r.client.GetSegment(data.UUID.ValueString())
	if err != nil {
		panic(err)
	}
	resourceData := MakeSegmentResourceDataFromClientSegment(segment)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}

func (r *segmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get plan values
	var plan SegmentResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading plan data")
		return
	}

	// Get current state
	var state SegmentResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading state data")
		return
	}
	// Generate API request body from plan
	clientSegment := plan.ToClientSegment()

	err := r.client.UpdateSegment(clientSegment)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update segment, got error: %s", err))
		return
	}

	resourceData := MakeSegmentResourceDataFromClientSegment(clientSegment)

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resourceData)

	resp.Diagnostics.Append(diags...)

}

func (r *segmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state SegmentResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Delete: Error reading state data")
		return
	}
	//Generate API request body from plan
	clientSegment := state.ToClientSegment()

	err := r.client.DeleteSegment(*clientSegment.ProjectID, *clientSegment.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete segment, got error: %s", err))
		return
	}
	resp.State.RemoveResource(ctx)

}
func (r *segmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
