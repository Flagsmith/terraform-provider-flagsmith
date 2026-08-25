package flagsmith_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Acceptance tests for custom field (metadata) support. These need custom fields to
// already exist in the Flagsmith project under test; see the accessors in
// provider_test.go for the environment variables that name them.

func metadataPreCheck(t *testing.T, extra ...string) func() {
	return func() {
		testAccPreCheck(t)
		mustHaveEnv(t, "FLAGSMITH_METADATA_FIELD_NAME")
		for _, name := range extra {
			mustHaveEnv(t, name)
		}
	}
}

func TestAccFeatureResourceMetadata(t *testing.T) {
	featureName := "tf_test_" + strings.ToLower(acctest.RandString(12))
	resourceName := "flagsmith_feature.test_feature"

	resource.Test(t, resource.TestCase{
		PreCheck:                 metadataPreCheck(t, "FLAGSMITH_METADATA_FIELD_NAME_WITH_SPACE"),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFeatureResourceDestroy,
		Steps: []resource.TestStep{
			// Create with two custom field values
			{
				Config: testAccFeatureMetadataConfig(featureName, map[string]string{
					metadataFieldName():          "first-value",
					metadataFieldNameWithSpace(): "PROD-123",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata."+metadataFieldName(), "first-value"),
					resource.TestCheckResourceAttr(resourceName, "metadata."+metadataFieldNameWithSpace(), "PROD-123"),
				),
			},
			// Import: metadata must round trip through the read path, which for features
			// needs the extra hydration request.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFeatureImportID(resourceName),
			},
			// Change one value
			{
				Config: testAccFeatureMetadataConfig(featureName, map[string]string{
					metadataFieldName():          "second-value",
					metadataFieldNameWithSpace(): "PROD-123",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata."+metadataFieldName(), "second-value"),
				),
			},
			// Drop one key
			{
				Config: testAccFeatureMetadataConfig(featureName, map[string]string{
					metadataFieldName(): "second-value",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "metadata."+metadataFieldNameWithSpace()),
				),
			},
			// Set an explicit empty map
			{
				Config: testAccFeatureMetadataConfig(featureName, map[string]string{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "0"),
				),
			},
			// Remove the attribute entirely
			{
				Config: testAccFeatureMetadataConfig(featureName, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "metadata.%"),
				),
			},
		},
	})
}

func TestAccFeatureResourceMetadataUnknownField(t *testing.T) {
	featureName := "tf_test_" + strings.ToLower(acctest.RandString(12))

	resource.Test(t, resource.TestCase{
		PreCheck:                 metadataPreCheck(t),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureMetadataConfig(featureName, map[string]string{
					"tf_acc_definitely_not_a_custom_field": "value",
				}),
				ExpectError: regexp.MustCompile("Unknown custom field"),
			},
		},
	})
}

func TestAccFeatureResourceMetadataWrongEntity(t *testing.T) {
	featureName := "tf_test_" + strings.ToLower(acctest.RandString(12))

	resource.Test(t, resource.TestCase{
		PreCheck:                 metadataPreCheck(t, "FLAGSMITH_METADATA_SEGMENT_ONLY_FIELD_NAME"),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureMetadataConfig(featureName, map[string]string{
					metadataSegmentOnlyFieldName(): "value",
				}),
				ExpectError: regexp.MustCompile("not enabled for features"),
			},
		},
	})
}

func TestAccFeatureResourceMetadataInvalidValue(t *testing.T) {
	featureName := "tf_test_" + strings.ToLower(acctest.RandString(12))

	resource.Test(t, resource.TestCase{
		PreCheck:                 metadataPreCheck(t, "FLAGSMITH_METADATA_INT_FIELD_NAME"),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureMetadataConfig(featureName, map[string]string{
					metadataIntFieldName(): "not-a-number",
				}),
				ExpectError: regexp.MustCompile("Invalid custom field value"),
			},
		},
	})
}

func TestAccSegmentResourceMetadata(t *testing.T) {
	segmentName := "tf_test_" + strings.ToLower(acctest.RandString(12))
	resourceName := "flagsmith_segment.test_segment"

	resource.Test(t, resource.TestCase{
		PreCheck:                 metadataPreCheck(t),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSegmentResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentMetadataConfig(segmentName, map[string]string{
					metadataFieldName(): "first-value",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "metadata."+metadataFieldName(), "first-value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getSegmentImportID(resourceName),
			},
			{
				Config: testAccSegmentMetadataConfig(segmentName, map[string]string{
					metadataFieldName(): "second-value",
				}),
				Check: resource.TestCheckResourceAttr(resourceName, "metadata."+metadataFieldName(), "second-value"),
			},
			{
				Config: testAccSegmentMetadataConfig(segmentName, nil),
				Check:  resource.TestCheckNoResourceAttr(resourceName, "metadata.%"),
			},
		},
	})
}

func TestAccEnvironmentResourceMetadata(t *testing.T) {
	environmentName := "tf_test_" + strings.ToLower(acctest.RandString(12))
	resourceName := "flagsmith_environment.test_environment"

	resource.Test(t, resource.TestCase{
		PreCheck:                 metadataPreCheck(t),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentMetadataConfig(environmentName, map[string]string{
					metadataFieldName(): "first-value",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "metadata."+metadataFieldName(), "first-value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getEnvironmentImportID(resourceName),
			},
			{
				Config: testAccEnvironmentMetadataConfig(environmentName, map[string]string{
					metadataFieldName(): "second-value",
				}),
				Check: resource.TestCheckResourceAttr(resourceName, "metadata."+metadataFieldName(), "second-value"),
			},
			{
				Config: testAccEnvironmentMetadataConfig(environmentName, nil),
				Check:  resource.TestCheckNoResourceAttr(resourceName, "metadata.%"),
			},
		},
	})
}

// metadataBlock renders a `metadata` attribute. A nil map omits the attribute entirely,
// which is a different case from an empty map.
func metadataBlock(values map[string]string) string {
	if values == nil {
		return ""
	}

	entries := ""
	for name, value := range values {
		entries += fmt.Sprintf("    %q = %q\n", name, value)
	}

	return fmt.Sprintf("  metadata = {\n%s  }\n", entries)
}

func testAccFeatureMetadataConfig(featureName string, metadata map[string]string) string {
	return fmt.Sprintf(`
%s

resource "flagsmith_feature" "test_feature" {
  feature_name = "%s"
  project_uuid = "%s"
  description  = "created by the custom fields acceptance test"
%s}
`, providerConfig(), featureName, projectUUID(), metadataBlock(metadata))
}

func testAccSegmentMetadataConfig(segmentName string, metadata map[string]string) string {
	return fmt.Sprintf(`
%s

resource "flagsmith_segment" "test_segment" {
  name         = "%s"
  description  = "created by the custom fields acceptance test"
  project_uuid = "%s"
  rules = [
    {
      "rules" : [{
        "conditions" : [{
          "operator" : "EQUAL",
          "property" : "device_type",
          "value" : "mobile"
        }],
        "type" : "ANY"
      }],
      "type" : "ALL"
    }
  ]
%s}
`, providerConfig(), segmentName, projectUUID(), metadataBlock(metadata))
}

func testAccEnvironmentMetadataConfig(environmentName string, metadata map[string]string) string {
	return fmt.Sprintf(`
%s

resource "flagsmith_environment" "test_environment" {
  name        = "%s"
  project_id  = %d
  description = "created by the custom fields acceptance test"
%s}
`, providerConfig(), environmentName, projectID(), metadataBlock(metadata))
}
