package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccFeatureResource(t *testing.T) {
	featureName := acctest.RandString(16)
	firstUserID := 3936
	secondUserID := 12662
	thirdUserID := 11871

	initialOwners := []int{firstUserID, secondUserID}
	updatedOwners := []int{firstUserID, thirdUserID}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFeatureResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFeatureResourceConfig(featureName, "new feature description", initialOwners),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "description", "new feature description"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "default_enabled", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "is_archived", "false"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.0", fmt.Sprintf("%d", firstUserID)),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.1", fmt.Sprintf("%d", secondUserID)),

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
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.0", fmt.Sprintf("%d", firstUserID)),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.1", fmt.Sprintf("%d", secondUserID)),

					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_feature.test_feature", "project_id"),
				),
			},

			// Update testing
			{
				Config: testAccFeatureResourceConfig(featureName, "feature description updated", updatedOwners),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "feature_name", featureName),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "description", "feature description updated"),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "project_uuid", projectUUID()),

					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.0", fmt.Sprintf("%d", firstUserID)),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.1", fmt.Sprintf("%d", thirdUserID)),
				),
			},
		},
	})
}

func TestAccFeatureResourceOwners(t *testing.T) {
	featureName := acctest.RandString(16)
	firstUserID := 3936
	secondUserID := 12662
	thirdUserID := 11871

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFeatureResourceDestroy,

		Steps: []resource.TestStep{
			// Create without owners
			{
				Config: testAccFeatureResourceConfig(featureName, "new feature description", []int{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.#", "0"),
				),
			},

			// ImportState without owners
			{
				ResourceName:      "flagsmith_feature.test_feature",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFeatureImportID("flagsmith_feature.test_feature"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.#", "0"),
				),
			},

			// UpdateState - Add owners
			{
				Config: testAccFeatureResourceConfig(featureName, "feature description updated", []int{firstUserID, secondUserID}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.0", fmt.Sprintf("%d", firstUserID)),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.1", fmt.Sprintf("%d", secondUserID)),
				),
			},

			// UpdateState - Update owners
			{
				Config: testAccFeatureResourceConfig(featureName, "feature description updated", []int{firstUserID, thirdUserID}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.0", fmt.Sprintf("%d", firstUserID)),
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.1", fmt.Sprintf("%d", thirdUserID)),
				),
			},
			// UpdateState - Remove all owners
			{
				Config: testAccFeatureResourceConfig(featureName, "feature description updated", []int{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_feature.test_feature", "owners.#", "0"),
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
	uuid, err := getUUIDfromState(s, "flagsmith_feature.test_feature")
	if err != nil {
		return err
	}

	_, err = testClient().GetFeature(uuid)
	if err == nil {
		return fmt.Errorf("feature still exists")
	}
	return nil

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

func testAccFeatureResourceConfig(featureName, description string, owners []int) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_feature" "test_feature" {
  feature_name = "%s"
  description = "%s"
  project_uuid = "%s"
  type = "STANDARD"
  owners = %s
}

`, featureName, description, projectUUID(), strings.Join(strings.Fields(fmt.Sprint(owners)), ","))
}
