package flagsmith

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	flagsmithapi "github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/stretchr/testify/assert"
)

func TestIntFeatureStateValueToClientFSV(t *testing.T) {
	// Given
	intFSV := FeatureStateValue{Type: types.StringValue("int"),
		StringValue:  types.StringNull(),
		IntegerValue: types.Int64Value(1),
		BooleanValue: types.BoolValue(true),
	}
	// When
	clientFSV := intFSV.ToClientFSV()

	// Then
	var nilString *string
	var nilBool *bool

	assert.Equal(t, "int", clientFSV.Type)
	assert.Equal(t, nilString, clientFSV.StringValue)
	assert.Equal(t, nilBool, clientFSV.BooleanValue)

	assert.Equal(t, int64(1), *clientFSV.IntegerValue)
}

func TestStringFeatureStateValueToClientFSV(t *testing.T) {
	// Given
	stringFSV := FeatureStateValue{Type: types.StringValue("unicode"),
		StringValue:  types.StringValue("string"),
		IntegerValue: types.Int64Null(),
		BooleanValue: types.BoolNull(),
	}
	// When
	clientFSV := stringFSV.ToClientFSV()

	// Then
	var nilInt *int64
	var nilBool *bool

	assert.Equal(t, "unicode", clientFSV.Type)
	assert.Equal(t, "string", *clientFSV.StringValue)
	assert.Equal(t, nilInt, clientFSV.IntegerValue)
	assert.Equal(t, nilBool, clientFSV.BooleanValue)
}

func TestBoolFeatureStateValueToClientFSV(t *testing.T) {
	// Given
	boolFSV := FeatureStateValue{Type: types.StringValue( "bool"),
		StringValue:  types.StringNull(),
		IntegerValue: types.Int64Null(),
		BooleanValue: types.BoolValue(true),
	}
	// When
	clientFSV := boolFSV.ToClientFSV()

	// Then
	assert.Equal(t, "bool", clientFSV.Type)
	assert.Equal(t, true, *clientFSV.BooleanValue)

	var nilString *string
	var nilInt *int64

	assert.Equal(t, nilString, clientFSV.StringValue)
	assert.Equal(t, nilInt, clientFSV.IntegerValue)
}

func TestMakeIntFeatureStateValueFromClientFSV(t *testing.T) {
	// Given
	intValue := int64(1)
	clientFSV := flagsmithapi.FeatureStateValue{
		Type:         "int",
		IntegerValue: &intValue,
	}
	// When
	fsv := MakeFeatureStateValueFromClientFSV(&clientFSV)

	// Then
	assert.Equal(t, "int", fsv.Type.ValueString())
	assert.Equal(t, intValue, fsv.IntegerValue.ValueInt64())
	assert.Equal(t, true, fsv.StringValue.IsNull())
	assert.Equal(t, true, fsv.BooleanValue.IsNull())
}

func TestMakeStringFeatureStateValueFromClientFSV(t *testing.T) {
	// Given
	stringValue := "string"
	clientFSV := flagsmithapi.FeatureStateValue{
		Type:        "unicode",
		StringValue: &stringValue,
	}
	// When
	fsv := MakeFeatureStateValueFromClientFSV(&clientFSV)

	// Then
	assert.Equal(t, "unicode", fsv.Type.ValueString())
	assert.Equal(t, stringValue, fsv.StringValue.ValueString())
	assert.Equal(t, true, fsv.IntegerValue.IsNull())
	assert.Equal(t, true, fsv.BooleanValue.IsNull())

}

func TestMakeBooleanFeatureStateValueFromClientFSV(t *testing.T) {
	// Given
	boolValue := true
	clientFSV := flagsmithapi.FeatureStateValue{
		Type:         "bool",
		BooleanValue: &boolValue,
	}

	// When
	fsv := MakeFeatureStateValueFromClientFSV(&clientFSV)

	// Then
	assert.Equal(t, "bool", fsv.Type.ValueString())
	assert.Equal(t, boolValue, fsv.BooleanValue.ValueBool())
	assert.Equal(t, true, fsv.StringValue.IsNull())
	assert.Equal(t, true, fsv.IntegerValue.IsNull())

}

func TestMakeFeatureStateResourceDataFromClientFS(t *testing.T) {
	// Given
	intValue := int64(1)
	isEnabled := true
	clientFSV := flagsmithapi.FeatureStateValue{
		Type:         "int",
		StringValue:  nil,
		IntegerValue: &intValue,
		BooleanValue: nil,
	}
	clientFS := flagsmithapi.FeatureState{
		ID:                1,
		FeatureStateValue: &clientFSV,
		Enabled:           isEnabled,
		Feature:           int64(1),
		Environment:       &intValue,
	}
	// When
	featureStateResourceData := MakeFeatureStateResourceDataFromClientFS(&clientFS)

	// Then
	assert.Equal(t, int64(1), featureStateResourceData.ID.ValueInt64())
	assert.Equal(t, isEnabled, featureStateResourceData.Enabled.ValueBool())
	assert.Equal(t, int64(1), featureStateResourceData.Feature.ValueInt64())
	assert.Equal(t, int64(1), featureStateResourceData.Environment.ValueInt64())
	assert.Equal(t, "int", featureStateResourceData.FeatureStateValue.Type.ValueString())
	assert.Equal(t, intValue, featureStateResourceData.FeatureStateValue.IntegerValue.ValueInt64())
	assert.Equal(t, true, featureStateResourceData.FeatureStateValue.StringValue.IsNull())
	assert.Equal(t, true, featureStateResourceData.FeatureStateValue.BooleanValue.IsNull())

}

func TestFeatureStateResourceDataToClientFS(t *testing.T) {
	//Given
	featureStateID := int64(1)
	environmentID := int64(1)
	featureID := int64(1)

	featureStateResourceData := FeatureStateResourceData{
		Enabled: types.BoolValue(true),
		ID:     types.Int64Value(featureStateID),
		Environment: types.Int64Value(environmentID),
		Feature: types.Int64Value(featureID),
		FeatureStateValue: &FeatureStateValue{
			Type:         types.StringValue("int"),
			StringValue:  types.StringNull(),
			IntegerValue: types.Int64Value(1),
			BooleanValue: types.BoolNull(),
		},
	}

	// When
	clientFS := featureStateResourceData.ToClientFS()

	// Then
	assert.Equal(t, featureStateID, clientFS.ID)
	assert.Equal(t, true, clientFS.Enabled)
	assert.Equal(t, featureID, clientFS.Feature)
	assert.Equal(t, environmentID, *clientFS.Environment)
	assert.Equal(t, "int", clientFS.FeatureStateValue.Type)
	assert.Equal(t, int64(1), *clientFS.FeatureStateValue.IntegerValue)

	var nilString *string
	var nilBool *bool

	assert.Equal(t, nilString, clientFS.FeatureStateValue.StringValue)
	assert.Equal(t, nilBool, clientFS.FeatureStateValue.BooleanValue)

}
