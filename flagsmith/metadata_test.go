package flagsmith

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func metadataMap(t *testing.T, values map[string]string) types.Map {
	t.Helper()

	elements := map[string]attr.Value{}
	for name, value := range values {
		elements[name] = types.StringValue(value)
	}

	metadata, diags := types.MapValue(types.StringType, elements)
	require.False(t, diags.HasError(), "%s", diags)

	return metadata
}

func TestMetadataMapToStringMap(t *testing.T) {
	ctx := context.Background()

	t.Run("null", func(t *testing.T) {
		values, diags := metadataMapToStringMap(ctx, types.MapNull(types.StringType))
		assert.False(t, diags.HasError())
		assert.Nil(t, values)
	})

	t.Run("unknown", func(t *testing.T) {
		values, diags := metadataMapToStringMap(ctx, types.MapUnknown(types.StringType))
		assert.False(t, diags.HasError())
		assert.Nil(t, values)
	})

	t.Run("empty", func(t *testing.T) {
		values, diags := metadataMapToStringMap(ctx, metadataMap(t, nil))
		assert.False(t, diags.HasError())
		require.NotNil(t, values)
		assert.Equal(t, 0, len(values))
	})

	t.Run("populated", func(t *testing.T) {
		values, diags := metadataMapToStringMap(ctx, metadataMap(t, map[string]string{"Jira Ticket": "PROD-123"}))
		assert.False(t, diags.HasError())
		assert.Equal(t, map[string]string{"Jira Ticket": "PROD-123"}, values)
	})
}

// The null vs empty distinction is what keeps a configuration that omits `metadata` and
// one that sets `metadata = {}` both converging.
func TestMetadataToState(t *testing.T) {
	t.Run("no values, null config stays null", func(t *testing.T) {
		state := metadataToState(nil, types.MapNull(types.StringType))
		assert.True(t, state.IsNull())
	})

	t.Run("no values, unknown config becomes null", func(t *testing.T) {
		state := metadataToState(nil, types.MapUnknown(types.StringType))
		assert.True(t, state.IsNull())
	})

	t.Run("no values, empty config stays empty", func(t *testing.T) {
		state := metadataToState(map[string]string{}, metadataMap(t, nil))
		assert.False(t, state.IsNull())
		assert.Equal(t, 0, len(state.Elements()))
	})

	t.Run("values with null config are surfaced as drift", func(t *testing.T) {
		state := metadataToState(map[string]string{"Jira Ticket": "PROD-123"}, types.MapNull(types.StringType))
		assert.False(t, state.IsNull())
		assert.Equal(t, map[string]attr.Value{"Jira Ticket": types.StringValue("PROD-123")}, state.Elements())
	})

	t.Run("values round trip", func(t *testing.T) {
		configured := metadataMap(t, map[string]string{"Jira Ticket": "PROD-123"})
		state := metadataToState(map[string]string{"Jira Ticket": "PROD-123"}, configured)
		assert.Equal(t, configured, state)
	})
}

// A zero types.Map has a nil element type and cannot be written to state, so every
// branch must produce a properly typed map.
func TestMetadataToStateAlwaysHasStringElementType(t *testing.T) {
	ctx := context.Background()

	states := []types.Map{
		metadataToState(nil, types.MapNull(types.StringType)),
		metadataToState(map[string]string{}, metadataMap(t, nil)),
		metadataToState(map[string]string{"Jira Ticket": "PROD-123"}, types.MapNull(types.StringType)),
	}

	for _, state := range states {
		assert.Equal(t, types.StringType, state.ElementType(ctx))
	}
}

// When the endpoint did not report metadata at all there is nothing to reconcile, so the
// configured value must be carried through untouched.
func TestMetadataFromClientWithNilClientMetadata(t *testing.T) {
	configured := metadataMap(t, map[string]string{"Jira Ticket": "PROD-123"})

	state, diags := metadataFromClient(context.Background(), nil, 1, "feature", nil, configured)

	assert.False(t, diags.HasError())
	assert.Equal(t, configured, state)
}

func TestMetadataAttributeSchemaIsOptionalNotComputed(t *testing.T) {
	attribute := metadataAttributeSchema("feature")

	// Optional+Computed would hide drift and make removal impossible.
	assert.True(t, attribute.IsOptional())
	assert.False(t, attribute.IsComputed())
	assert.False(t, attribute.IsRequired())
	assert.Equal(t, types.StringType, attribute.ElementType)
}
