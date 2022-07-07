package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFlagResourceConfig("one", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature_name", featureName()),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "environment", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature", strconv.Itoa(featureID())),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature_state_value.string_value", "one"),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "enabled", "true"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_flag.test_feature",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s,%s", environmentKey(), featureName()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature_name", featureName()),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "environment", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature", strconv.Itoa(featureID())),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature_state_value.string_value", "one"),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "enabled", "true"),
				),
			},
			// Update testing
			{
				Config: testAccFlagResourceConfig("two", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature_name", featureName()),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "environment_key", environmentKey()),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "environment", strconv.Itoa(environmentID())),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature", strconv.Itoa(featureID())),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "feature_state_value.string_value", "two"),
					resource.TestCheckResourceAttr("flagsmith_flag.test_feature", "enabled", "false"),
				),
			},
		},
	})
}

func testAccFlagResourceConfig(featureStateValue string, isEnabled bool) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_flag" "test_feature" {
  enabled         = %t
  environment_key = "%s"
  feature_name    = "%s"
  feature_state_value = {
    type         = "unicode"
    string_value = "%s"
  }

}

`, isEnabled, environmentKey(), featureName(), featureStateValue)
}
