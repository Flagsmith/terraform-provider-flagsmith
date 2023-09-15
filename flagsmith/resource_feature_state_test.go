package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"strconv"
	"testing"
	"regexp"
)

func TestAccEnvironmentFeatureStateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test feature State value validator
			{
				Config: testAccInvalidFeatureStateValueConfig(),
				ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured:\n\[feature_state_value.string_value,feature_state_value.integer_value,feature_state_value.boolean_value\]`),

			},
			// Test feature State string value validator
			{
				Config: testAccEnvironmentFeatureStateResourceConfig(" some_value ", true),
				ExpectError: regexp.MustCompile(`Attribute feature_state_value.string_value Leading and trailing whitespace is\n.*not allowed`),

			},

			// Create and Read testing
			{
				Config: testAccEnvironmentFeatureStateResourceConfig("one", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "environment_id", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "feature_id", strconv.Itoa(featureID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "feature_state_value.string_value", "one"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "enabled", "true"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_feature_state.dummy_environment_feature_x",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFeatureStateImportID("flagsmith_feature_state.dummy_environment_feature_x"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "environment_id", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "feature_id", strconv.Itoa(featureID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "feature_state_value.string_value", "one"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "enabled", "true"),
				),
			},

			//Update testing
			{
				Config: testAccEnvironmentFeatureStateResourceConfig("two", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "environment_id", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "feature_id", strconv.Itoa(featureID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "feature_state_value.string_value", "two"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccSegmentFeatureStateResource(t *testing.T) {
	featureName := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSegmentFeatureStateResourceConfig("one", featureName, true, 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "environment_id", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_state_value.string_value", "one"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "enabled", "true"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "segment_priority", "0"),

					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "segment_id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_segment_id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_feature_state.dummy_environment_feature_x_segment_override",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFeatureStateImportID("flagsmith_feature_state.dummy_environment_feature_x_segment_override"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "environment_id", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_state_value.string_value", "one"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "enabled", "true"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "segment_priority", "0"),

					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "segment_id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_segment_id"),
				),
			},

			// Update testing
			{
				Config: testAccSegmentFeatureStateResourceConfig("two", featureName, false, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "environment_id", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_state_value.string_value", "two"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "enabled", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "segment_priority", "2"),

					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "segment_id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature_state.dummy_environment_feature_x_segment_override", "feature_segment_id"),
				),
			},
		},
	})
}

func getFeatureStateImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		uuid, err := getUUIDfromState(s, n)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s,%s", environmentKey(), uuid), nil

	}
}

func testAccCheckSegmentFeatureStateDestroy(s *terraform.State) error {
	uuid, err := getUUIDfromState(s, "flagsmith_feature_state.dummy_environment_feature_x_segment_override")
	if err != nil {
		return err
	}

	_, err = testClient().GetFeature(uuid)
	if err == nil {
		return fmt.Errorf("feature still exists")
	}
	return nil

}
func testAccSegmentFeatureStateResourceConfig(featureStateValue string, featureName string, isEnabled bool, segmentPriority int) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_segment" "test_segment" {
  name         = "test_segment"
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
}

resource "flagsmith_feature" "test_feature" {
  feature_name = "%s"
  project_uuid = "%s"
  description = "feature created for terraform segment override test"
  type = "STANDARD"
}

resource "flagsmith_feature_state" "dummy_environment_feature_x_segment_override" {
  enabled         = %t
  environment_key = "%s"
  feature_id = flagsmith_feature.test_feature.id
  segment_id = flagsmith_segment.test_segment.id
  segment_priority = %d
  feature_state_value = {
    type         = "unicode"
    string_value = "%s"
  }

}

`, projectUUID(), featureName, projectUUID(), isEnabled, environmentKey(), segmentPriority, featureStateValue)
}

func testAccEnvironmentFeatureStateResourceConfig(featureStateValue string, isEnabled bool) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_feature_state" "dummy_environment_feature_x" {
  enabled         = %t
  environment_key = "%s"
  feature_id = %d
  feature_state_value = {
    type         = "unicode"
    string_value = "%s"
  }

}

`, isEnabled, environmentKey(), featureID(), featureStateValue)
}
func testAccInvalidFeatureStateValueConfig() string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_feature_state" "dummy_environment_feature_x" {
  enabled         = true
  environment_key = "%s"
  feature_id = %d
  feature_state_value = {
    type         = "unicode"
  }
}

`,  environmentKey(), featureID())
}
