package flagsmith

import (
	flagsmithapi "github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
)

type FeatureStateValue struct {
	Type         types.String `tfsdk:"type"`
	StringValue  types.String `tfsdk:"string_value"`
	IntegerValue types.Int64  `tfsdk:"integer_value"`
	BooleanValue types.Bool   `tfsdk:"boolean_value"`
}

func (f *FeatureStateValue) ToClientFSV() *flagsmithapi.FeatureStateValue {
	switch f.Type.ValueString() {
	case "unicode":
		value := f.StringValue.ValueString()
		return &flagsmithapi.FeatureStateValue{
			Type:        "unicode",
			StringValue: &value,
		}
	case "int":
		intValue := f.IntegerValue.ValueInt64()
		return &flagsmithapi.FeatureStateValue{
			Type:         "int",
			IntegerValue: &intValue,
		}
	case "bool":
		value := f.BooleanValue.ValueBool()
		return &flagsmithapi.FeatureStateValue{
			Type:         "bool",
			BooleanValue: &value,
		}
	}
	return nil
}

func MakeFeatureStateValueFromClientFSV(clientFSV *flagsmithapi.FeatureStateValue) FeatureStateValue {
	fsvType := clientFSV.Type
	fsValue := FeatureStateValue{
		Type:         types.StringValue(fsvType),
		StringValue:  types.StringNull(),
		IntegerValue: types.Int64Null(),
		BooleanValue: types.BoolNull(),
	}
	switch fsvType {
	case "unicode":
		if clientFSV.StringValue == nil {
			fsValue.StringValue = types.StringValue("")
		} else {
			fsValue.StringValue = types.StringValue(*clientFSV.StringValue)
		}
		return fsValue

	case "int":
		fsValue.IntegerValue = types.Int64Value(*clientFSV.IntegerValue)
		return fsValue

	case "bool":
		fsValue.BooleanValue = types.BoolValue(*clientFSV.BooleanValue)
		return fsValue

	}
	return FeatureStateValue{}
}

type FeatureStateResourceData struct {
	ID                types.Int64        `tfsdk:"id"`
	UUID              types.String       `tfsdk:"uuid"`
	Enabled           types.Bool         `tfsdk:"enabled"`
	FeatureStateValue *FeatureStateValue `tfsdk:"feature_state_value"`
	Feature           types.Int64        `tfsdk:"feature_id"`
	Environment       types.Int64        `tfsdk:"environment_id"`
	EnvironmentKey    types.String       `tfsdk:"environment_key"`
	Segment           types.Int64        `tfsdk:"segment_id"`
	SegmentPriority   types.Int64        `tfsdk:"segment_priority"`
	FeatureSegment    types.Int64        `tfsdk:"feature_segment_id"`
}

func (f *FeatureStateResourceData) ToClientFS() *flagsmithapi.FeatureState {
	fs := flagsmithapi.FeatureState{
		ID:                f.ID.ValueInt64(),
		UUID:              f.UUID.ValueString(),
		Enabled:           f.Enabled.ValueBool(),
		FeatureStateValue: f.FeatureStateValue.ToClientFSV(),
		Feature:           f.Feature.ValueInt64(),
		EnvironmentKey:    f.EnvironmentKey.ValueString(),
	}
	featureSegment := f.FeatureSegment.ValueInt64()
	segment := f.Segment.ValueInt64()
	if featureSegment != 0 {
		fs.FeatureSegment = &featureSegment
	}
	if segment != 0 {
		fs.Segment = &segment
	}
	if !f.SegmentPriority.IsNull() && !f.SegmentPriority.IsUnknown() {
		int64SegmentPriority := f.SegmentPriority.ValueInt64()
		fs.SegmentPriority = &int64SegmentPriority
	}
	environment := f.Environment.ValueInt64()
	if environment != 0 {
		fs.Environment = &environment
	}
	return &fs
}

// Generate a new FeatureStateResourceData from client `FeatureState`
func MakeFeatureStateResourceDataFromClientFS(clientFS *flagsmithapi.FeatureState) FeatureStateResourceData {
	fsValue := MakeFeatureStateValueFromClientFSV(clientFS.FeatureStateValue)
	fs := FeatureStateResourceData{
		ID:                types.Int64Value(clientFS.ID),
		UUID:              types.StringValue(clientFS.UUID),
		Enabled:           types.BoolValue(clientFS.Enabled),
		FeatureStateValue: &fsValue,
		Feature:           types.Int64Value(clientFS.Feature),
		Environment:       types.Int64Value(*clientFS.Environment),
		EnvironmentKey:    types.StringValue(clientFS.EnvironmentKey),
		Segment:           types.Int64Null(),
		SegmentPriority:   types.Int64Null(),
		FeatureSegment:    types.Int64Null(),
	}
	if clientFS.FeatureSegment != nil {
		featureSegment := types.Int64Value(*clientFS.FeatureSegment)
		fs.FeatureSegment = featureSegment

		if clientFS.SegmentPriority != nil {
			segmentPriority := types.Int64Value(*clientFS.SegmentPriority)
			fs.SegmentPriority = segmentPriority
		}

		if clientFS.Segment != nil {
			segment := types.Int64Value(*clientFS.Segment)
			fs.Segment = segment
		}

	}
	return fs

}

type MultivariateOptionResourceData struct {
	Type                        types.String `tfsdk:"type"`
	ID                          types.Int64  `tfsdk:"id"`
	UUID                        types.String `tfsdk:"uuid"`
	FeatureID                   types.Int64  `tfsdk:"feature_id"`
	FeatureUUID                 types.String `tfsdk:"feature_uuid"`
	ProjectID                   types.Int64  `tfsdk:"project_id"`
	IntegerValue                types.Int64  `tfsdk:"integer_value"`
	StringValue                 types.String `tfsdk:"string_value"`
	BooleanValue                types.Bool   `tfsdk:"boolean_value"`
	DefaultPercentageAllocation types.Number `tfsdk:"default_percentage_allocation"`
}

func NewMultivariateOptionFromClientOption(clientMvOption *flagsmithapi.FeatureMultivariateOption) MultivariateOptionResourceData {
	mvOption := MultivariateOptionResourceData{
		Type:                        types.StringValue(clientMvOption.Type),
		ID:                          types.Int64Value(clientMvOption.ID),
		UUID:                        types.StringValue(clientMvOption.UUID),
		FeatureID:                   types.Int64Value(*clientMvOption.FeatureID),
		FeatureUUID:                 types.StringValue(clientMvOption.FeatureUUID),
		ProjectID:                   types.Int64Value(*clientMvOption.ProjectID),
		DefaultPercentageAllocation: types.NumberValue(big.NewFloat(clientMvOption.DefaultPercentageAllocation)),
	}
	switch clientMvOption.Type {
	case "unicode":
		mvOption.StringValue = types.StringValue(*clientMvOption.StringValue)
	case "int":
		mvOption.IntegerValue = types.Int64Value(*clientMvOption.IntegerValue)
	case "bool":
		mvOption.BooleanValue = types.BoolValue(*clientMvOption.BooleanValue)
	}
	return mvOption

}

func (m *MultivariateOptionResourceData) ToClientMultivariateOption() *flagsmithapi.FeatureMultivariateOption {
	defaultPercentageAllocation, _ := m.DefaultPercentageAllocation.ValueBigFloat().Float64()
	stringValue := m.StringValue.ValueString()
	booleanValue := m.BooleanValue.ValueBool()

	mo := flagsmithapi.FeatureMultivariateOption{
		Type:                        m.Type.ValueString(),
		UUID:                        m.UUID.ValueString(),
		FeatureUUID:                 m.FeatureUUID.ValueString(),
		DefaultPercentageAllocation: defaultPercentageAllocation,
	}
	if !m.ID.IsNull() && !m.ID.IsUnknown() {
		mvOptionID := m.ID.ValueInt64()
		mo.ID = mvOptionID
	}
	if !m.FeatureID.IsNull() && !m.FeatureID.IsUnknown() {
		featureID := m.FeatureID.ValueInt64()
		mo.FeatureID = &featureID
	}
	if !m.ProjectID.IsNull() && !m.ProjectID.IsUnknown() {
		projectID := m.ProjectID.ValueInt64()
		mo.ProjectID = &projectID
	}
	switch m.Type.ValueString() {
	case "unicode":
		mo.StringValue = &stringValue
	case "int":
		integerValue := m.IntegerValue.ValueInt64()
		mo.IntegerValue = &integerValue
	case "bool":
		mo.BooleanValue = &booleanValue
	}

	return &mo
}

type FeatureResourceData struct {
	UUID           types.String   `tfsdk:"uuid"`
	ID             types.Int64    `tfsdk:"id"`
	Name           types.String   `tfsdk:"feature_name"`
	Type           types.String   `tfsdk:"type"`
	Description    types.String   `tfsdk:"description"`
	InitialValue   types.String   `tfsdk:"initial_value"`
	DefaultEnabled types.Bool     `tfsdk:"default_enabled"`
	IsArchived     types.Bool     `tfsdk:"is_archived"`
	Owners         *[]types.Int64 `tfsdk:"owners"`
	Tags           *[]types.Int64 `tfsdk:"tags"`
	ProjectID      types.Int64    `tfsdk:"project_id"`
	ProjectUUID    types.String   `tfsdk:"project_uuid"`
}

func (f *FeatureResourceData) ToClientFeature() *flagsmithapi.Feature {
	typeValue := f.Type.ValueString()
	if typeValue == "" {
		typeValue = "STANDARD"
	}
	descriptionValue := f.Description.ValueString()
	feature := flagsmithapi.Feature{
		UUID:           f.UUID.ValueString(),
		Name:           f.Name.ValueString(),
		Type:           &typeValue,
		Description:    &descriptionValue,
		InitialValue:   f.InitialValue.ValueString(),
		DefaultEnabled: f.DefaultEnabled.ValueBool(),
		IsArchived:     f.IsArchived.ValueBool(),
		ProjectUUID:    f.ProjectUUID.ValueString(),
		Tags:           []int64{},
		Owners:         &[]int64{},
	}
	if !f.ID.IsNull() && !f.ID.IsUnknown() {
		featureID := f.ID.ValueInt64()
		feature.ID = &featureID
	}
	if !f.ProjectID.IsNull() && !f.ProjectID.IsUnknown() {
		projectID := f.ProjectID.ValueInt64()
		feature.ProjectID = &projectID
	}
	if f.Owners == nil {
		feature.Owners = nil
	}

	if f.Owners != nil && len(*f.Owners) > 0 {
		for _, owner := range *f.Owners {
			ownerID := owner.ValueInt64()
			*feature.Owners = append(*feature.Owners, ownerID)
		}
	}
	if f.Tags != nil {
		for _, tag := range *f.Tags {
			tagID := tag.ValueInt64()
			feature.Tags = append(feature.Tags, tagID)
		}
	}
	return &feature

}

func MakeFeatureResourceDataFromClientFeature(clientFeature *flagsmithapi.Feature) FeatureResourceData {
	resourceData := FeatureResourceData{
		UUID:           types.StringValue(clientFeature.UUID),
		ID:             types.Int64Value(*clientFeature.ID),
		Name:           types.StringValue(clientFeature.Name),
		Type:           types.StringValue(*clientFeature.Type),
		DefaultEnabled: types.BoolValue(clientFeature.DefaultEnabled),
		IsArchived:     types.BoolValue(clientFeature.IsArchived),
		InitialValue:   types.StringValue(clientFeature.InitialValue),
		ProjectID:      types.Int64Value(*clientFeature.ProjectID),
		ProjectUUID:    types.StringValue(clientFeature.ProjectUUID),
		Owners:         &[]types.Int64{},
	}
	if clientFeature.Description != nil {
		resourceData.Description = types.StringValue(*clientFeature.Description)
	}

	if clientFeature.Owners == nil {
		resourceData.Owners = nil
	}

	if clientFeature.Owners != nil && len(*clientFeature.Owners) > 0 {
		for _, owner := range *clientFeature.Owners {
			*resourceData.Owners = append(*resourceData.Owners, types.Int64Value(owner))
		}
	}
	if clientFeature.Tags != nil && len(clientFeature.Tags) > 0 {
		resourceData.Tags = &[]types.Int64{}
		for _, tag := range clientFeature.Tags {
			*resourceData.Tags = append(*resourceData.Tags, types.Int64Value(tag))
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
		Operator: c.Operator.ValueString(),
		Property: c.Property.ValueString(),
		Value:    c.Value.ValueString(),
	}
}

func MakeConditionFromClientCondition(clientCondition *flagsmithapi.Condition) Condition {
	return Condition{
		Operator: types.StringValue(clientCondition.Operator),
		Property: types.StringValue(clientCondition.Property),
		Value:    types.StringValue(clientCondition.Value),
	}

}

type NestedRule struct {
	Conditions []Condition  `tfsdk:"conditions"`
	Type       types.String `tfsdk:"type"`
}

func (r *NestedRule) ToClientRule() *flagsmithapi.Rule {
	conditions := make([]flagsmithapi.Condition, 0)
	for _, condition := range r.Conditions {
		conditions = append(conditions, *condition.ToClientCondition())
	}
	return &flagsmithapi.Rule{
		Conditions: conditions,
		Type:       r.Type.ValueString(),
	}
}
func MakeNestedRuleFromClientRule(clientRule *flagsmithapi.Rule) NestedRule {
	var conditions []Condition
	for _, clientCondition := range clientRule.Conditions {
		conditions = append(conditions, MakeConditionFromClientCondition(&clientCondition))
	}
	return NestedRule{
		Conditions: conditions,
		Type:       types.StringValue(clientRule.Type),
	}

}

type Rule struct {
	Type       types.String `tfsdk:"type"`
	Rules      []NestedRule `tfsdk:"rules"`
	Conditions []Condition  `tfsdk:"conditions"`
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
		Type:       r.Type.ValueString(),
		Rules:      rules,
		Conditions: conditions,
	}
}

func MakeRuleFromClientRule(clientRule *flagsmithapi.Rule) Rule {
	rule := Rule{
		Type: types.StringValue(clientRule.Type),
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
	ID          types.Int64  `tfsdk:"id"`
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ProjectID   types.Int64  `tfsdk:"project_id"`
	ProjectUUID types.String `tfsdk:"project_uuid"`
	FeatureID   types.Int64  `tfsdk:"feature_id"`
	Rules       []Rule       `tfsdk:"rules"`
}

func (s *SegmentResourceData) ToClientSegment() *flagsmithapi.Segment {
	segment := flagsmithapi.Segment{
		UUID:        s.UUID.ValueString(),
		Name:        s.Name.ValueString(),
		ProjectUUID: s.ProjectUUID.ValueString(),
	}
	if s.Description.ValueString() != "" {
		value := s.Description.ValueString()
		segment.Description = &value
	}
	if !s.ID.IsNull() && !s.ID.IsUnknown() {
		segmentID := s.ID.ValueInt64()
		segment.ID = &segmentID
	}
	if !s.FeatureID.IsNull() && !s.FeatureID.IsUnknown() {
		featureID := s.FeatureID.ValueInt64()
		segment.FeatureID = &featureID
	}
	if !s.ProjectID.IsNull() && !s.ProjectID.IsUnknown() {
		projectID := s.ProjectID.ValueInt64()
		segment.ProjectID = &projectID
	}
	for _, rule := range s.Rules {
		segment.Rules = append(segment.Rules, *rule.ToClientRule())
	}
	return &segment
}

func MakeSegmentResourceDataFromClientSegment(clientSegment *flagsmithapi.Segment) SegmentResourceData {
	resourceData := SegmentResourceData{
		ID:          types.Int64Value(*clientSegment.ID),
		UUID:        types.StringValue(clientSegment.UUID),
		Name:        types.StringValue(clientSegment.Name),
		ProjectID:   types.Int64Value(*clientSegment.ProjectID),
		ProjectUUID: types.StringValue(clientSegment.ProjectUUID),
	}
	if clientSegment.Description != nil {
		resourceData.Description = types.StringValue(*clientSegment.Description)
	}
	if clientSegment.FeatureID != nil {
		resourceData.FeatureID = types.Int64Value(*clientSegment.FeatureID)
	}

	for _, clientRule := range clientSegment.Rules {
		resourceData.Rules = append(resourceData.Rules, MakeRuleFromClientRule(&clientRule))
	}
	return resourceData
}

type TagResourceData struct {
	ID          types.Int64  `tfsdk:"id"`
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"tag_name"`
	Description types.String `tfsdk:"description"`
	ProjectID   types.Int64  `tfsdk:"project_id"`
	ProjectUUID types.String `tfsdk:"project_uuid"`
	Colour      types.String `tfsdk:"tag_colour"`
}

func (t *TagResourceData) ToClientTag() *flagsmithapi.Tag {
	tag := flagsmithapi.Tag{
		UUID:        t.UUID.ValueString(),
		Name:        t.Name.ValueString(),
		ProjectUUID: t.ProjectUUID.ValueString(),
		Colour:      t.Colour.ValueString(),
	}
	if t.Description.ValueString() != "" {
		value := t.Description.ValueString()
		tag.Description = &value
	}
	if !t.ID.IsNull() && !t.ID.IsUnknown() {
		tagID := t.ID.ValueInt64()
		tag.ID = &tagID
	}
	if !t.ProjectID.IsNull() && !t.ProjectID.IsUnknown() {
		projectID := t.ProjectID.ValueInt64()
		tag.ProjectID = &projectID
	}
	return &tag
}

func MakeTagResourceDataFromClientTag(clientTag *flagsmithapi.Tag) TagResourceData {
	resourceData := TagResourceData{
		ID:          types.Int64Value(*clientTag.ID),
		UUID:        types.StringValue(clientTag.UUID),
		Name:        types.StringValue(clientTag.Name),
		ProjectID:   types.Int64Value(*clientTag.ProjectID),
		ProjectUUID: types.StringValue(clientTag.ProjectUUID),
		Colour:      types.StringValue(clientTag.Colour),
	}
	if clientTag.Description != nil {
		resourceData.Description = types.StringValue(*clientTag.Description)
	}

	return resourceData
}
