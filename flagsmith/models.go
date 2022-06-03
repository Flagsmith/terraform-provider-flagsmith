package flagsmith

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FeatureStateValue struct {
	Type types.String `tfsdk:"type"`
	StringValue types.String `tfsdk:"string_value"`
	IntegerValue types.Int64 `tfsdk:"integer_value"`
	BooleanValue types.Bool `tfsdk:"boolean_value"`


}
type flagResourceData struct {
	Id      types.Int64 `tfsdk:"id"`
	Enabled types.Bool  `tfsdk:"enabled"`
	FeatureStateValue types.Bool `tfsdk:"feature_state_value"`
	Feature types.Int64 `tfsdk:"feature"`
	Environment types.Int64 `tfsdk:"environment"`
}
