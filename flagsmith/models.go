package flagsmith

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	flagsmithapi "github.com/Flagsmith/flagsmith-go-api-client"
	"math/big"
)

type FeatureStateValue struct {
	Type         types.String `tfsdk:"type"`
	StringValue  types.String `tfsdk:"string_value"`
	IntegerValue types.Number `tfsdk:"integer_value"`
	BooleanValue types.Bool   `tfsdk:"boolean_value"`
}

func (f *FeatureStateValue) ToClientFSV() *flagsmithapi.FeatureStateValue {
	switch f.Type.Value {
	case "unicode":
		return &flagsmithapi.FeatureStateValue{
			Type:        "unicode",
			StringValue: &f.StringValue.Value,
		}
	case "int":
		intValue, _ := f.IntegerValue.Value.Int64()
		return &flagsmithapi.FeatureStateValue{
			Type:         "int",
			IntegerValue: &intValue,
		}
	case "bool":
		return &flagsmithapi.FeatureStateValue{
			Type:         "bool",
			BooleanValue: &f.BooleanValue.Value,
		}
	}
	return nil
}

func MakeFeatureStateValueFromClientFSV(clientFSV *flagsmithapi.FeatureStateValue) FeatureStateValue {
	fsvType := clientFSV.Type
	switch fsvType {
	case "unicode":
		return FeatureStateValue{
			Type:         types.String{Value: fsvType},
			StringValue:  types.String{Value: *clientFSV.StringValue},
			IntegerValue: types.Number{Null: true, Value: nil},
			BooleanValue: types.Bool{Null: true},
		}
	case "int":
		return FeatureStateValue{
			Type:         types.String{Value: fsvType},
			StringValue:  types.String{Null: true},
			IntegerValue: types.Number{Value: big.NewFloat(float64(*clientFSV.IntegerValue))},
			BooleanValue: types.Bool{Null: true},
		}
	case "bool":
		return FeatureStateValue{
			Type:         types.String{Value: fsvType},
			StringValue:  types.String{Null: true},
			IntegerValue: types.Number{Null: true},
			BooleanValue: types.Bool{Value: *clientFSV.BooleanValue},
		}

	}
	return FeatureStateValue{}
}

type FlagResourceData struct {
	ID                types.Number       `tfsdk:"id"`
	Enabled           types.Bool         `tfsdk:"enabled"`
	FeatureStateValue *FeatureStateValue `tfsdk:"feature_state_value"`
	Feature           types.Number       `tfsdk:"feature"`
	Environment       types.Number       `tfsdk:"environment"`
	FeatureName       types.String       `tfsdk:"feature_name"`
	EnvironmentKey    types.String       `tfsdk:"environment_key"`
}

func (f *FlagResourceData) ToClientFS(featureStateID int64, feature int64, environment int64) *flagsmithapi.FeatureState {
	return &flagsmithapi.FeatureState{
		ID:                featureStateID,
		Enabled:           f.Enabled.Value,
		FeatureStateValue: f.FeatureStateValue.ToClientFSV(),
		Feature:           feature,
		Environment:       environment,
	}
}

// Generate a new FlagResourceData from client `FeatureState`
func MakeFlagResourceDataFromClientFS(clientFS *flagsmithapi.FeatureState) FlagResourceData {
	fsValue := MakeFeatureStateValueFromClientFSV(clientFS.FeatureStateValue)
	return FlagResourceData{
		ID:                types.Number{Value: big.NewFloat(float64(clientFS.ID))},
		Enabled:           types.Bool{Value: clientFS.Enabled},
		FeatureStateValue: &fsValue,
		Feature:           types.Number{Value: big.NewFloat(float64(clientFS.Feature))},
		Environment:       types.Number{Value: big.NewFloat(float64(clientFS.Environment))},
	}
}
