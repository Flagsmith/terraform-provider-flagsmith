package flagsmith

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FeatureStateValueType struct {
	Type types.String `tfsdk:"type"`
	StringValue  types.String `tfsdk:"string_value"`
	IntegerValue types.Number `tfsdk:"integer_value"`
	BooleanValue types.Bool   `tfsdk:"boolean_value"`

}

type June struct {
	ID types.Number `tfsdk:"id"`
}

type flagResourceData struct {
	ID                types.Number          `tfsdk:"id"`
	Enabled           types.Bool            `tfsdk:"enabled"`
	FeatureStateValue *FeatureStateValueType `tfsdk:"feature_state_value"`
	Feature           types.Number          `tfsdk:"feature"`
	Environment       types.Number          `tfsdk:"environment"`
}
