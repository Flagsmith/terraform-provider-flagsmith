package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccFeatureResource(t *testing.T) {
	featureName :=  acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFeatureResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFeatureResourceConfig(featureName, "new feature description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "description", "new feature description"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "default_enabled", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "is_archived", "false"),

					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "project_id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_feature.test_feature",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFeatureImportID("flagsmith_feature.test_feature"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "default_enabled", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "is_archived", "false"),

					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "project_id"),
				),
			},

			// Update testing
			{
				Config: testAccFeatureResourceConfig(featureName, "feature description updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "description", "feature description updated"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),
				),
			},
		},
	})
}

func TestAccMVFeatureResouce(t *testing.T) {
	featureName :=  acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFeatureResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMVFeatureResourceConfig(featureName, "new feature description", 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "description", "new feature description"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "default_enabled", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "is_archived", "false"),

					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.#", "3"),

					// Mv value 0
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.string_value", "option_value_10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.type", "unicode"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.0.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.integer_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.boolean_value"),

					// Mv value 1
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.integer_value", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.type", "int"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.1.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.string_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.boolean_value"),

					// Mv value 2
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.boolean_value", "true"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.type", "bool"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.2.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.string_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.integer_value"),

					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "project_id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_feature.test_feature",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFeatureImportID("flagsmith_feature.test_feature"),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "description", "new feature description"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "default_enabled", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "is_archived", "false"),

					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.#", "3"),

					// Mv value 0
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.string_value", "option_value_10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.type", "unicode"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.0.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.integer_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.boolean_value"),

					// Mv value 1
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.integer_value", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.type", "int"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.1.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.string_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.boolean_value"),

					// Mv value 2
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.boolean_value", "true"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.type", "bool"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.2.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.string_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.integer_value"),

					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "project_id"),
				),
			},

			// Update testing
			{
				Config: testAccMVFeatureResourceConfig(featureName, "feature description updated", 60),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "description", "feature description updated"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "default_enabled", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "is_archived", "false"),

					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.#", "3"),

					// Mv value 0
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.string_value", "option_value_10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.type", "unicode"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.0.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.integer_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.0.boolean_value"),

					// Mv value 1
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.default_percentage_allocation", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.integer_value", "10"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.type", "int"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.1.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.string_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.1.boolean_value"),

					// Mv value 2
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.default_percentage_allocation", "60"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.boolean_value", "true"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.type", "bool"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "multivariate_options.2.id"),

					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.string_value"),
					resource.TestCheckNoResourceAttr("flagsmith_feature.test_feature", "multivariate_options.2.integer_value"),

					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "project_id"),
				),
			},
		},
	})
}

func getFeatureImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		return getUUIDfromState(s, n)
	}
}

func testAccCheckFeatureResourceDestroy(s *terraform.State) error {
	return testAccFeatureDestroy("flagsmith_feature.test_feature")(s)

}

func getUUIDfromState(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("not found: %s", resourceName)
	}

	uuid := rs.Primary.Attributes["uuid"]

	if uuid == "" {
		return "", fmt.Errorf("no uuid is set")
	}
	return uuid, nil
}

func testAccFeatureDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		uuid, err := getUUIDfromState(s, n)
		if err != nil {
			return err
		}

		_, err = testClient().GetFeature(uuid)
		if err == nil {
			return fmt.Errorf("feature still exists")
		}
		return nil

	}
}

func testAccFeatureResourceConfig(featureName, description string) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_feature" "test_feature" {
  feature_name = "%s"
  description = "%s"
  project_uuid = "%s"
  type = "STANDARD"
}

`, featureName, description, projectUUID())
}

func testAccMVFeatureResourceConfig(featureName, description string, boolPercentageAllocation int) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_feature" "test_feature" {
  feature_name = "%s"
  description = "%s"
  project_uuid = "%s"
  type = "MULTIVARIATE"
  multivariate_options = [
    {
      type : "unicode",
      string_value : "option_value_10",
      default_percentage_allocation : 10
    },
    {
      type : "int",
      integer_value : 10,
      default_percentage_allocation : 10
    },
    {
      type : "bool",
      boolean_value : true,
      default_percentage_allocation : %d
    }
  ]
}

`, featureName, description, projectUUID(), boolPercentageAllocation)
}

