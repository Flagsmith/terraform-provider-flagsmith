package flagsmith_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFeatureDataResource(t *testing.T) {
	featureName := "tf_test_" + strings.ToLower(acctest.RandString(12))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the feature on its own first, so that the data source read in the
			// next step cannot race the create.
			{
				Config: testAccFeatureDataResourceFeatureConfig(featureName),
			},
			{
				Config: testAccFeatureDataResourceConfig(featureName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("data.flagsmith_feature.test_feature", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("data.flagsmith_feature.test_feature", "project_id", strconv.Itoa(projectID())),
					resource.TestCheckResourceAttr("data.flagsmith_feature.test_feature", "description", "created by the feature data source acceptance test"),
					resource.TestCheckResourceAttr("data.flagsmith_feature.test_feature", "type", "STANDARD"),
					resource.TestCheckResourceAttr("data.flagsmith_feature.test_feature", "default_enabled", "false"),
					resource.TestCheckResourceAttr("data.flagsmith_feature.test_feature", "is_archived", "false"),
					resource.TestCheckResourceAttrSet("data.flagsmith_feature.test_feature", "id"),

					// The lookup must resolve to the feature we created.
					resource.TestCheckResourceAttrPair(
						"data.flagsmith_feature.test_feature", "uuid",
						"flagsmith_feature.test_feature", "uuid"),
					resource.TestCheckResourceAttrPair(
						"data.flagsmith_feature.test_feature", "id",
						"flagsmith_feature.test_feature", "id"),
				),
			},
		},
	})
}

// The motivating use case from #153: reference an externally managed feature and attach
// an environment specific override to it.
func TestAccFeatureDataResourceDrivesFeatureState(t *testing.T) {
	featureName := "tf_test_" + strings.ToLower(acctest.RandString(12))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureDataResourceFeatureConfig(featureName),
			},
			{
				Config: testAccFeatureDataResourceFeatureStateConfig(featureName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"flagsmith_feature_state.test_feature_state", "feature_id",
						"data.flagsmith_feature.test_feature", "id"),
					resource.TestCheckResourceAttr("flagsmith_feature_state.test_feature_state", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccFeatureDataResourceNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "flagsmith_feature" "test_feature" {
  uuid = "8f1b0e1a-0000-4000-8000-000000000000"
}
`, providerConfig()),
				ExpectError: regexp.MustCompile("Unable to read feature"),
			},
		},
	})
}

func testAccFeatureDataResourceFeatureConfig(featureName string) string {
	return fmt.Sprintf(`
%s

resource "flagsmith_feature" "test_feature" {
  feature_name = "%s"
  project_uuid = "%s"
  description  = "created by the feature data source acceptance test"
}
`, providerConfig(), featureName, projectUUID())
}

func testAccFeatureDataResourceConfig(featureName string) string {
	return testAccFeatureDataResourceFeatureConfig(featureName) + `
data "flagsmith_feature" "test_feature" {
  uuid = flagsmith_feature.test_feature.uuid
}
`
}

func testAccFeatureDataResourceFeatureStateConfig(featureName string) string {
	return testAccFeatureDataResourceConfig(featureName) + fmt.Sprintf(`
resource "flagsmith_feature_state" "test_feature_state" {
  enabled         = true
  environment_key = "%s"
  feature_id      = data.flagsmith_feature.test_feature.id
  feature_state_value = {
    type         = "unicode"
    string_value = "some-value"
  }
}
`, environmentKey())
}
