package flagsmith

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func validateFeatureStateValue() validatorFeatureStateValue {
	return validatorFeatureStateValue{}
}

type validatorFeatureStateValue struct {
}

func (v validatorFeatureStateValue) Description(ctx context.Context) string {
	return "One of string_value, integer_value or boolean_value must be set"
}

func (v validatorFeatureStateValue) MarkdownDescription(ctx context.Context) string {
	return "One of string_value, integer_value or boolean_value must be set"
}

func (v validatorFeatureStateValue) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	attrs := req.ConfigValue.Attributes()

	if !hasValue(attrs) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid feature_state_value",
			"One of string_value, integer_value or boolean_value must be set",
		)
	}
}

func hasValue(attrs map[string]attr.Value) bool {
	values := []string{"string_value", "integer_value", "boolean_value"}
	for _, value := range values {
		if !(attrs[value].IsNull() || attrs[value].IsUnknown()) {
			return true
		}
	}
	return false
}
