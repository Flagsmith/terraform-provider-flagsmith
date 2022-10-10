package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccMultivariateFeatureOptionResource(t *testing.T) {
	featureName := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFeatureMVOptionDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFeatureMVOptionResourceConfig(featureName, "option_value_43.13", 43.13),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "type", "unicode"),
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "string_value", "option_value_43.13"),
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "default_percentage_allocation", "43.13"),

					resource.TestCheckNoResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "integer_value"),
					resource.TestCheckNoResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "boolean_value"),

					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "feature_id"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "feature_uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "project_id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_mv_feature_option.feature_1_mv_option",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getMvFeatureOptionImportID("flagsmith_mv_feature_option.feature_1_mv_option"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "type", "unicode"),
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "string_value", "option_value_43.13"),
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "default_percentage_allocation", "43.13"),

					resource.TestCheckNoResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "integer_value"),
					resource.TestCheckNoResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "boolean_value"),

					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "feature_id"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "feature_uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "project_id"),
				),
			},

			// Update testing
			{
				Config: testAccFeatureMVOptionResourceConfig(featureName, "updated_option_value", 99.99),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "type", "unicode"),
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "string_value", "updated_option_value"),
					resource.TestCheckResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "default_percentage_allocation", "99.99"),

					resource.TestCheckNoResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "integer_value"),
					resource.TestCheckNoResourceAttr("flagsmith_mv_feature_option.feature_1_mv_option", "boolean_value"),

					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "feature_id"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "feature_uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_mv_feature_option.feature_1_mv_option", "project_id"),
				),
			},
		},
	})
}

func getMvFeatureOptionImportID(resource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resource)
		}

		featureUUID := rs.Primary.Attributes["feature_uuid"]
		uuid := rs.Primary.Attributes["uuid"]

		return fmt.Sprintf("%s,%s", featureUUID, uuid), nil
	}
}

func testAccCheckFeatureMVOptionDestroy(s *terraform.State) error {
	resource := "flagsmith_mv_feature_option.feature_1_mv_option"
	rs, ok := s.RootModule().Resources[resource]
	if !ok {
		return fmt.Errorf("Not found: %s", resource)
	}

	uuid := rs.Primary.Attributes["uuid"]
	if uuid == "" {
		return fmt.Errorf("No UUID is set")
	}
	featureUUID := rs.Primary.Attributes["feature_uuid"]
	if featureUUID == "" {
		return fmt.Errorf("No feature UUID is set")
	}
	_, err := testClient().GetFeatureMVOption(featureUUID, uuid)
	if err == nil {
		return fmt.Errorf("Feature MV Option still exists")
	}
	return nil

}

func testAccFeatureMVOptionResourceConfig(featureName, optionValue string, defaultPercentageAllocation float64) string {
	return fmt.Sprintf(`
provider "flagsmith" {

  base_api_url   = "http://localhost:8000/api/v1"
}

resource "flagsmith_feature" "test_feature" {
  feature_name = "%s"
  description = "mv_feature_option_test"
  project_uuid = "%s"
  type = "MULTIVARIATE"
}
resource "flagsmith_mv_feature_option" "feature_1_mv_option" {
  type                          = "unicode"
  feature_uuid                  = flagsmith_feature.test_feature.uuid
  string_value                  = "%s"
  default_percentage_allocation = %.2f

}
`, featureName, projectUUID(), optionValue, defaultPercentageAllocation)
}
