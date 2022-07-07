package flagsmith

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"

	flagsmithapi "github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/stretchr/testify/assert"
)

func TestIntFeatureStateValueToClientFSV(t *testing.T) {
	// Given
	intFSV := FeatureStateValue{Type: types.String{Value: "int"},
		StringValue:  types.String{Null: true},
		IntegerValue: types.Number{Value: big.NewFloat(1)},
		BooleanValue: types.Bool{Null: true},
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
	stringFSV := FeatureStateValue{Type: types.String{Value: "unicode"},
		StringValue:  types.String{Value: "string"},
		IntegerValue: types.Number{Null: true},
		BooleanValue: types.Bool{Null: true},
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
	boolFSV := FeatureStateValue{Type: types.String{Value: "bool"},
		StringValue:  types.String{Null: true},
		IntegerValue: types.Number{Null: true},
		BooleanValue: types.Bool{Value: true},
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
	assert.Equal(t, "int", fsv.Type.Value)
	assert.Equal(t, big.NewFloat(float64(intValue)), fsv.IntegerValue.Value)
	assert.Equal(t, true, fsv.StringValue.Null)
	assert.Equal(t, true, fsv.BooleanValue.Null)
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
	assert.Equal(t, "unicode", fsv.Type.Value)
	assert.Equal(t, stringValue, fsv.StringValue.Value)
	assert.Equal(t, true, fsv.IntegerValue.Null)
	assert.Equal(t, true, fsv.BooleanValue.Null)

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
	assert.Equal(t, "bool", fsv.Type.Value)
	assert.Equal(t, boolValue, fsv.BooleanValue.Value)
	assert.Equal(t, true, fsv.StringValue.Null)
	assert.Equal(t, true, fsv.IntegerValue.Null)

}

func TestMakeFlagResourceDataFromClientFS(t *testing.T) {
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
		Environment:       int64(1),
	}
	// When
	flagResourceData := MakeFlagResourceDataFromClientFS(&clientFS)

	// Then
	assert.Equal(t, big.NewFloat(1), flagResourceData.ID.Value)
	assert.Equal(t, isEnabled, flagResourceData.Enabled.Value)
	assert.Equal(t, big.NewFloat(1), flagResourceData.Feature.Value)
	assert.Equal(t, big.NewFloat(1), flagResourceData.Environment.Value)
	assert.Equal(t, "int", flagResourceData.FeatureStateValue.Type.Value)
	assert.Equal(t, big.NewFloat(float64(intValue)), flagResourceData.FeatureStateValue.IntegerValue.Value)
	assert.Equal(t, true, flagResourceData.FeatureStateValue.StringValue.Null)
	assert.Equal(t, true, flagResourceData.FeatureStateValue.BooleanValue.Null)

}

func TestFlagResourceDataToClientFS(t *testing.T) {
	//Given
	flagResourceData := FlagResourceData{
		Enabled: types.Bool{Value: true},
		FeatureStateValue: &FeatureStateValue{
			Type:         types.String{Value: "int"},
			StringValue:  types.String{Null: true},
			IntegerValue: types.Number{Value: big.NewFloat(1)},
			BooleanValue: types.Bool{Null: true},
		},
	}

	// When
	featureStateID := int64(1)
	environment := int64(1)
	feture := int64(1)
	clientFS := flagResourceData.ToClientFS(featureStateID, environment, feture)

	// Then
	assert.Equal(t, featureStateID, clientFS.ID)
	assert.Equal(t, true, clientFS.Enabled)
	assert.Equal(t, int64(1), clientFS.Feature)
	assert.Equal(t, int64(1), clientFS.Environment)
	assert.Equal(t, "int", clientFS.FeatureStateValue.Type)
	assert.Equal(t, int64(1), *clientFS.FeatureStateValue.IntegerValue)

	var nilString *string
	var nilBool *bool

	assert.Equal(t, nilString, clientFS.FeatureStateValue.StringValue)
	assert.Equal(t, nilBool, clientFS.FeatureStateValue.BooleanValue)

}
