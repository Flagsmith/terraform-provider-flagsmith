package flagsmith

import (
	flagsmithapi "github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type FeatureStateResourceData struct {
	ID                types.Number       `tfsdk:"id"`
	Enabled           types.Bool         `tfsdk:"enabled"`
	FeatureStateValue *FeatureStateValue `tfsdk:"feature_state_value"`
	Feature           types.Number       `tfsdk:"feature"`
	Environment       types.Number       `tfsdk:"environment"`
	FeatureName       types.String       `tfsdk:"feature_name"`
	EnvironmentKey    types.String       `tfsdk:"environment_key"`
}

func (f *FeatureStateResourceData) ToClientFS(featureStateID int64, feature int64, environment int64) *flagsmithapi.FeatureState {
	return &flagsmithapi.FeatureState{
		ID:                featureStateID,
		Enabled:           f.Enabled.Value,
		FeatureStateValue: f.FeatureStateValue.ToClientFSV(),
		Feature:           feature,
		Environment:       environment,
	}
}

// Generate a new FeatureStateResourceData from client `FeatureState`
func MakeFeatureStateResourceDataFromClientFS(clientFS *flagsmithapi.FeatureState) FeatureStateResourceData {
	fsValue := MakeFeatureStateValueFromClientFSV(clientFS.FeatureStateValue)
	return FeatureStateResourceData{
		ID:                types.Number{Value: big.NewFloat(float64(clientFS.ID))},
		Enabled:           types.Bool{Value: clientFS.Enabled},
		FeatureStateValue: &fsValue,
		Feature:           types.Number{Value: big.NewFloat(float64(clientFS.Feature))},
		Environment:       types.Number{Value: big.NewFloat(float64(clientFS.Environment))},
	}
}

type MultivariateOptionResourceData struct {
	Type                        types.String `tfsdk:"type"`
	ID                          types.Number `tfsdk:"id"`
	UUID                        types.String `tfsdk:"uuid"`
	FeatureID                   types.Number `tfsdk:"feature_id"`
	FeatureUUID                 types.String `tfsdk:"feature_uuid"`
	ProjectID                   types.Number `tfsdk:"project_id"`
	IntegerValue                types.Number `tfsdk:"integer_value"`
	StringValue                 types.String `tfsdk:"string_value"`
	BooleanValue                types.Bool   `tfsdk:"boolean_value"`
	DefaultPercentageAllocation types.Number `tfsdk:"default_percentage_allocation"`
}

func NewMultivariateOptionFromClientOption(clientMvOption *flagsmithapi.FeatureMultivariateOption) MultivariateOptionResourceData {
	mvOption := MultivariateOptionResourceData{
		Type:                        types.String{Value: clientMvOption.Type},
		ID:                          types.Number{Value: big.NewFloat(float64(clientMvOption.ID))},
		UUID:                        types.String{Value: clientMvOption.UUID},
		FeatureID:                   types.Number{Value: big.NewFloat(float64(*clientMvOption.FeatureID))},
		FeatureUUID:                 types.String{Value: clientMvOption.FeatureUUID},
		ProjectID:                   types.Number{Value: big.NewFloat(float64(*clientMvOption.ProjectID))},
		DefaultPercentageAllocation: types.Number{Value: big.NewFloat(clientMvOption.DefaultPercentageAllocation)},
		StringValue:                 types.String{Null: true},
		IntegerValue:                types.Number{Null: true, Value: nil},
		BooleanValue:                types.Bool{Null: true},
	}
	switch clientMvOption.Type {
	case "unicode":
		mvOption.StringValue = types.String{Value: *clientMvOption.StringValue}
	case "int":
		mvOption.IntegerValue = types.Number{Value: big.NewFloat(float64(*clientMvOption.IntegerValue))}
	case "bool":
		mvOption.BooleanValue = types.Bool{Value: *clientMvOption.BooleanValue}
	}
	return mvOption

}

func (m *MultivariateOptionResourceData) ToClientMultivariateOption() *flagsmithapi.FeatureMultivariateOption {
	defaultPercentageAllocation, _ := m.DefaultPercentageAllocation.Value.Float64()
	stringValue := m.StringValue.Value
	booleanValue := m.BooleanValue.Value

	mo := flagsmithapi.FeatureMultivariateOption{
		Type:                        m.Type.Value,
		UUID:                        m.UUID.Value,
		FeatureUUID:                 m.FeatureUUID.Value,
		DefaultPercentageAllocation: defaultPercentageAllocation,
	}
	if m.ID.Value != nil {
		mvOptionID, _ := m.ID.Value.Int64()
		mo.ID = mvOptionID
	}
	if m.FeatureID.Value != nil {
		featureID, _ := m.FeatureID.Value.Int64()
		mo.FeatureID = &featureID
	}
	if m.ProjectID.Value != nil {
		projectID, _ := m.ProjectID.Value.Int64()
		mo.ProjectID = &projectID
	}
	switch m.Type.Value {
	case "unicode":
		mo.StringValue = &stringValue
	case "int":
		integerValue, _ := m.IntegerValue.Value.Int64()
		mo.IntegerValue = &integerValue
	case "bool":
		mo.BooleanValue = &booleanValue
	}

	return &mo
}

type FeatureResourceData struct {
	UUID           types.String    `tfsdk:"uuid"`
	ID             types.Number    `tfsdk:"id"`
	Name           types.String    `tfsdk:"feature_name"`
	Type           types.String    `tfsdk:"type"`
	Description    types.String    `tfsdk:"description"`
	InitialValue   types.String    `tfsdk:"initial_value"`
	DefaultEnabled types.Bool      `tfsdk:"default_enabled"`
	IsArchived     types.Bool      `tfsdk:"is_archived"`
	Owners         *[]types.Number `tfsdk:"owners"`
	ProjectID      types.Number    `tfsdk:"project_id"`
	ProjectUUID    types.String    `tfsdk:"project_uuid"`
}

func (f *FeatureResourceData) ToClientFeature() *flagsmithapi.Feature {
	feature := flagsmithapi.Feature{
		UUID:           f.UUID.Value,
		Name:           f.Name.Value,
		Type:           &f.Type.Value,
		Description:    &f.Description.Value,
		InitialValue:   f.InitialValue.Value,
		DefaultEnabled: f.DefaultEnabled.Value,
		IsArchived:     f.IsArchived.Value,
		ProjectUUID:    f.ProjectUUID.Value,
	}
	if f.ID.Value != nil {
		featureID, _ := f.ID.Value.Int64()
		feature.ID = &featureID
	}
	if f.ProjectID.Value != nil {
		projectID, _ := f.ProjectID.Value.Int64()
		feature.ProjectID = &projectID
	}

	if f.Owners != nil {
		for _, owner := range *f.Owners {
			ownerID, _ := owner.Value.Int64()
			*feature.Owners = append(*feature.Owners, ownerID)

		}
	}
	return &feature

}

func MakeFeatureResourceDataFromClientFeature(clientFeature *flagsmithapi.Feature) FeatureResourceData {
	resourceData := FeatureResourceData{
		UUID:           types.String{Value: clientFeature.UUID},
		ID:             types.Number{Value: big.NewFloat(float64(*clientFeature.ID))},
		Name:           types.String{Value: clientFeature.Name},
		Type:           types.String{Value: *clientFeature.Type},
		DefaultEnabled: types.Bool{Value: clientFeature.DefaultEnabled},
		IsArchived:     types.Bool{Value: clientFeature.IsArchived},
		InitialValue:   types.String{Value: clientFeature.InitialValue},
		ProjectID:      types.Number{Value: big.NewFloat(float64(*clientFeature.ProjectID))},
		ProjectUUID:    types.String{Value: clientFeature.ProjectUUID},
	}
	if clientFeature.Description != nil {
		resourceData.Description = types.String{Value: *clientFeature.Description}
	}
	if clientFeature.Owners != nil {
		for _, owner := range *clientFeature.Owners {
			*resourceData.Owners = append(*resourceData.Owners, types.Number{Value: big.NewFloat(float64(owner))})
		}
	}
	return resourceData
}

type Condition struct {
	Operator types.String `tfsdk:"operator"`
	Property types.String `tfsdk:"property"`
	Value    types.String `tfsdk:"value"`
}
func (c *Condition) ToClientCondition() *flagsmithapi.Condition {
	return &flagsmithapi.Condition{
		Operator: c.Operator.Value,
		Property: c.Property.Value,
		Value:    c.Value.Value,
	}
}

func MakeConditionFromClientCondition(clientCondition * flagsmithapi.Condition) Condition {
	return Condition{
		Operator: types.String{Value: clientCondition.Operator},
		Property: types.String{Value: clientCondition.Property},
		Value:    types.String{Value: clientCondition.Value},
	}

}
type NestedRule struct {
	Conditions []Condition `tfsdk:"conditions"`
	Type      types.String `tfsdk:"type"`
}
func (t *NestedRule) ToClientRule() *flagsmithapi.Rule {
	conditions := make([]flagsmithapi.Condition, 0)
	for _, condition := range t.Conditions {
		conditions = append(conditions, *condition.ToClientCondition())
	}
	return &flagsmithapi.Rule{
		Conditions: conditions,
		Type:       t.Type.Value,
	}
}
func MakeNestedRuleFromClientRule(clientRule * flagsmithapi.Rule) NestedRule {
	var conditions []Condition
	for _, clientCondition := range clientRule.Conditions {
		conditions = append(conditions, MakeConditionFromClientCondition(&clientCondition))
	}
	return NestedRule{
		Conditions: conditions,
		Type:       types.String{Value: clientRule.Type},
	}

}

type Rule struct {
	Type types.String `tfsdk:"type"`
	Rules []NestedRule `tfsdk:"rules"`
	Conditions []Condition `tfsdk:"conditions"`
}

func (r *Rule) ToClientRule() *flagsmithapi.Rule {
	var conditions []flagsmithapi.Condition
	for _, condition := range r.Conditions {
		conditions = append(conditions, *condition.ToClientCondition())
	}
	var rules []flagsmithapi.Rule
	for _, rule := range r.Rules {
		rules = append(rules, *rule.ToClientRule())
	}
	return &flagsmithapi.Rule{
		Type: r.Type.Value,
		Rules: rules,
		Conditions: conditions,
	}
}

func MakeRuleFromClientRule(clientRule * flagsmithapi.Rule) Rule {
	rule := Rule{
		Type: types.String{Value: clientRule.Type},
	}
	for _, clientCondition := range clientRule.Conditions {
		rule.Conditions = append(rule.Conditions, MakeConditionFromClientCondition(&clientCondition))
	}
	for _, clientSubRule := range clientRule.Rules {
		rule.Rules = append(rule.Rules, MakeNestedRuleFromClientRule(&clientSubRule))
	}
	return rule
}

type SegmentResourceData struct {
	ID types.Number `tfsdk:"id"`
	UUID types.String `tfsdk:"uuid"`
	Name types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ProjectID types.Number `tfsdk:"project_id"`
	ProjectUUID types.String `tfsdk:"project_uuid"`
	FeatureID  types.Number `tfsdk:"feature_id"`
	Rules []Rule `tfsdk:"rules"`
}

func (s *SegmentResourceData) ToClientSegment() *flagsmithapi.Segment {
	segment := flagsmithapi.Segment{
		UUID: s.UUID.Value,
		Name: s.Name.Value,
		ProjectUUID: s.ProjectUUID.Value,
	}
	if s.Description.Value != "" {
		segment.Description = &s.Description.Value
	}
	if s.ID.Value != nil {
		segmentID, _ := s.ID.Value.Int64()
		segment.ID = &segmentID
	}
	if s.FeatureID.Value != nil {
		featureID, _ := s.FeatureID.Value.Int64()
		segment.FeatureID = &featureID
	}
	if s.ProjectID.Value != nil {
		projectID, _ := s.ProjectID.Value.Int64()
		segment.ProjectID = &projectID
	}
	for _, rule := range s.Rules {
		segment.Rules = append(segment.Rules, *rule.ToClientRule())
	}
	return &segment
}


func MakeSegmentResourceDataFromClientSegment(clientSegment *flagsmithapi.Segment) SegmentResourceData {
	resourceData := SegmentResourceData{
		ID: types.Number{Value: big.NewFloat(float64(*clientSegment.ID))},
		UUID: types.String{Value: clientSegment.UUID},
		Name: types.String{Value: clientSegment.Name},
		Description: types.String{Null: true},
		ProjectID: types.Number{Value: big.NewFloat(float64(*clientSegment.ProjectID))},
		ProjectUUID: types.String{Value: clientSegment.ProjectUUID},

	}
	if clientSegment.Description != nil {
		resourceData.Description = types.String{Value: *clientSegment.Description}
	}
	if clientSegment.FeatureID != nil {
		resourceData.FeatureID = types.Number{Value: big.NewFloat(float64(*clientSegment.FeatureID))}
	}

	for _, clientRule := range clientSegment.Rules {
		resourceData.Rules = append(resourceData.Rules, MakeRuleFromClientRule(&clientRule))
	}
	return resourceData
}
