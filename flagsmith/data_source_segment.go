package flagsmith

import (
	"context"
	"fmt"

	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &segmentDataResource{}

func newSegmentDataResource() datasource.DataSource {
	return &segmentDataResource{}
}

type segmentDataResource struct {
	client *flagsmithapi.Client
}

func (s *segmentDataResource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (s *segmentDataResource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	s.client = client
}

func (s *segmentDataResource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	conditions := schema.ListNestedAttribute{
		Computed:            true,
		MarkdownDescription: "List of Conditions for the Rule",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"property": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Property of the condition",
				},
				"operator": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Operator of the condition",
				},
				"value": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Value of the condition",
				},
			},
		},
	}

	nestedRules := schema.ListNestedAttribute{
		Computed:            true,
		MarkdownDescription: "List of Nested Rules",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Type of the rule, one of: `ALL`, `ANY`, `NONE`",
				},
				"conditions": conditions,
			},
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Segment. Use this to reference a segment that is managed outside of Terraform, for example to attach a segment override to a feature.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the segment",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the segment",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the segment",
			},
			"project_uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the project the segment belongs to",
			},
			"project_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the project the segment belongs to",
			},
			"feature_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the feature the segment is specific to, if any",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the segment",
			},
			"metadata": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom field ([metadata](https://docs.flagsmith.com/administration-and-security/governance-and-compliance/custom-fields)) values for this segment, keyed by custom field name.",
			},
			"rules": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Rules for the segment",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Type of the rule, one of: `ALL`, `ANY`, `NONE`",
						},
						"rules":      nestedRules,
						"conditions": conditions,
					},
				},
			},
		},
	}
}

func (s *segmentDataResource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SegmentResourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	segment, err := s.client.GetSegment(data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read segment, got error: %s", err))
		return
	}
	metadata, diags := metadataFromClient(ctx, s.client, *segment.ProjectID,
		flagsmithapi.MetadataEntitySegment, segment.Metadata, types.MapNull(types.StringType))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceData := MakeSegmentResourceDataFromClientSegment(segment, metadata)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}
