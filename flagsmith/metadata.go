package flagsmith

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// metadataAttributeSchema returns the schema for the `metadata` attribute, which holds
// custom field values keyed by custom field name.
//
// The attribute is Optional but not Computed. That is deliberate: the Flagsmith API
// replaces the whole set of custom field values on every write, so a value set outside
// Terraform on a managed resource is going to be removed either way. Optional+Computed
// would hide that in the plan and make it impossible to remove a value by deleting it
// from the configuration.
func metadataAttributeSchema(entity string) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:    true,
		ElementType: types.StringType,
		MarkdownDescription: fmt.Sprintf(
			"Custom field ([metadata](https://docs.flagsmith.com/administration-and-security/governance-and-compliance/custom-fields)) "+
				"values for this %[1]s, keyed by custom field name. The field must already exist in Flagsmith "+
				"and be enabled for %[1]ss. Terraform is authoritative: values set outside of Terraform are "+
				"removed on the next write.",
			entity),
		Validators: []validator.Map{
			mapvalidator.ValueStringsAre(
				stringvalidator.LengthBetween(1, flagsmithapi.MetadataFieldValueMaxLength),
			),
		},
	}
}

// resolveProjectID returns the project ID of a resource, looking it up from the project
// UUID only when it is not already known. project_id is Computed on the feature and
// segment resources, so it is unknown while creating one.
func resolveProjectID(client *flagsmithapi.Client, projectID *int64, projectUUID string) (int64, error) {
	if projectID != nil {
		return *projectID, nil
	}
	project, err := client.GetProject(projectUUID)
	if err != nil {
		return 0, err
	}
	return project.ID, nil
}

// metadataMapToStringMap converts the `metadata` attribute to a plain Go map. A null or
// unknown map yields nil.
func metadataMapToStringMap(ctx context.Context, metadata types.Map) (map[string]string, diag.Diagnostics) {
	if metadata.IsNull() || metadata.IsUnknown() {
		return nil, nil
	}

	values := map[string]string{}
	diags := metadata.ElementsAs(ctx, &values, false)

	return values, diags
}

// metadataToState decides what to store in state for the `metadata` attribute.
//
// Keeping null and empty distinct matters: writing an empty map where the configuration
// had null produces "Provider produced inconsistent result after apply", and writing
// null where the configuration had `{}` produces a perpetual diff.
func metadataToState(values map[string]string, configured types.Map) types.Map {
	if len(values) == 0 {
		if configured.IsNull() || configured.IsUnknown() {
			return types.MapNull(types.StringType)
		}
		return types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	elements := make(map[string]attr.Value, len(values))
	for name, value := range values {
		elements[name] = types.StringValue(value)
	}

	return types.MapValueMust(types.StringType, elements)
}

// resolveMetadataForWrite converts the `metadata` attribute into the wire format,
// resolving custom field names to the model field IDs the API requires.
//
// A null or unknown map yields a non-nil empty slice, i.e. an explicit "no custom field
// values", matching the API's behaviour of replacing the whole set on every write.
func resolveMetadataForWrite(
	ctx context.Context,
	client *flagsmithapi.Client,
	projectID int64,
	entity flagsmithapi.MetadataEntity,
	metadata types.Map,
) (*[]flagsmithapi.Metadata, diag.Diagnostics) {
	var diags diag.Diagnostics

	values, valueDiags := metadataMapToStringMap(ctx, metadata)
	diags.Append(valueDiags...)
	if diags.HasError() {
		return nil, diags
	}

	if len(values) == 0 {
		return &[]flagsmithapi.Metadata{}, diags
	}

	resolver, err := client.GetMetadataFieldResolver(projectID, entity)
	if err != nil {
		diags.AddAttributeError(
			path.Root("metadata"),
			"Unable to read custom fields",
			fmt.Sprintf("Could not list the custom fields available to this project: %s", err),
		)
		return nil, diags
	}

	// Sorted so that diagnostics and request bodies are deterministic.
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)

	clientMetadata := make([]flagsmithapi.Metadata, 0, len(names))
	for _, name := range names {
		modelFieldID, err := resolver.ModelFieldID(name)
		if err != nil {
			addMetadataFieldDiagnostic(&diags, name, string(entity), resolver, err)
			continue
		}

		if field, ok := resolver.Field(name); ok {
			if err := field.ValidateValue(values[name]); err != nil {
				diags.AddAttributeError(
					path.Root("metadata").AtMapKey(name),
					"Invalid custom field value",
					fmt.Sprintf("The custom field %q is of type %q, and %q is not a valid value for it.",
						name, field.Type, values[name]),
				)
				continue
			}
		}

		clientMetadata = append(clientMetadata, flagsmithapi.Metadata{
			ModelField: modelFieldID,
			FieldValue: values[name],
		})
	}

	if diags.HasError() {
		return nil, diags
	}

	return &clientMetadata, diags
}

// addMetadataFieldDiagnostic turns a resolution failure into an error scoped to the
// offending map key, so the user is pointed at the custom field that is wrong.
func addMetadataFieldDiagnostic(diags *diag.Diagnostics, name, entity string,
	resolver *flagsmithapi.MetadataFieldResolver, err error) {
	attributePath := path.Root("metadata").AtMapKey(name)

	switch err.(type) {
	case flagsmithapi.MetadataFieldNotBoundError:
		diags.AddAttributeError(
			attributePath,
			"Custom field not enabled for this entity",
			fmt.Sprintf("The custom field %q exists in this project but is not enabled for %ss. "+
				"Enable it for the %s entity in Flagsmith, or remove it from `metadata`.",
				name, entity, entity),
		)
	case flagsmithapi.MetadataFieldNotFoundError:
		available := "none"
		if names := resolver.BoundNames(); len(names) > 0 {
			available = strings.Join(names, ", ")
		}
		diags.AddAttributeError(
			attributePath,
			"Unknown custom field",
			fmt.Sprintf("No custom field named %q is available for %ss in this project.\n\n"+
				"Available custom fields: %s\n\n"+
				"Custom fields are defined per organisation or project in Flagsmith. Create the field, "+
				"enable it for %ss, then re-run Terraform.",
				name, entity, available, entity),
		)
	default:
		diags.AddAttributeError(attributePath, "Unable to resolve custom field", err.Error())
	}
}

// metadataFromClient converts wire format metadata back into the `metadata` attribute.
//
// configured is the value from the configuration, plan or prior state, and is used to
// decide between a null and an empty map, and to carry the value through unchanged when
// the API response did not report metadata at all.
func metadataFromClient(
	ctx context.Context,
	client *flagsmithapi.Client,
	projectID int64,
	entity flagsmithapi.MetadataEntity,
	clientMetadata *[]flagsmithapi.Metadata,
	configured types.Map,
) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	if clientMetadata == nil {
		// The endpoint did not report metadata, so there is nothing to reconcile.
		return configured, diags
	}

	values, unresolved, err := client.ResolveMetadataNames(projectID, entity, *clientMetadata)
	if err != nil {
		diags.AddAttributeError(
			path.Root("metadata"),
			"Unable to read custom fields",
			fmt.Sprintf("Could not list the custom fields available to this project: %s", err),
		)
		return configured, diags
	}

	if len(unresolved) > 0 {
		diags.AddAttributeWarning(
			path.Root("metadata"),
			"Custom field values could not be matched to a custom field",
			fmt.Sprintf("This %s has values for custom fields that no longer exist in the project "+
				"(model field IDs %v). They have been left out of the state.", entity, unresolved),
		)
	}

	return metadataToState(values, configured), diags
}
