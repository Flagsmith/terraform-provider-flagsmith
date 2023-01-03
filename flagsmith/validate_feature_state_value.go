package flagsmith

import (
	"context"
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

	if (attrs["string_value"].IsNull() || attrs["string_value"].IsUnknown()) && (attrs["integer_value"].IsNull() || attrs["integer_value"].IsUnknown()) && (attrs["boolean_value"].IsNull() || attrs["boolean_value"].IsUnknown()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid feature_state_value",
			"One of string_value, integer_value or boolean_value must be set",
		)

	}
}
