package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
	"strconv"
)

func TestAccProjectResource(t *testing.T) {
	projectName := acctest.RandString(16)
	newProjectName := acctest.RandString(16)


	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		//CheckDestroy:             testAccCheckProjectResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig(projectName),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("flagsmith_project.test_project", "name", projectName),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "organisation_id", strconv.Itoa( organisationID())),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "hide_disabled_flags", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "prevent_flag_defaults", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "enable_realtime_updates", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "only_allow_lower_case_feature_names", "true"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "feature_name_regex", ""),

					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "stale_flags_limit_days"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_project.test_project",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getProjectImportID("flagsmith_project.test_project"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "name", projectName),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "organisation_id", strconv.Itoa( organisationID())),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "hide_disabled_flags", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "prevent_flag_defaults", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "enable_realtime_updates", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "only_allow_lower_case_feature_names", "true"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "feature_name_regex", ""),

					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "stale_flags_limit_days"),

				),
			},

			// Update testing
			{
				Config: testAccProjectResourceConfig(newProjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "name", newProjectName),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "organisation_id", strconv.Itoa( organisationID())),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "hide_disabled_flags", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "prevent_flag_defaults", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "enable_realtime_updates", "false"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "only_allow_lower_case_feature_names", "true"),
					resource.TestCheckResourceAttr("flagsmith_project.test_project", "feature_name_regex", ""),

					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_project.test_project", "stale_flags_limit_days"),

				),
			},
		},
	})
}


func getProjectImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		return getAttributefromState(s, n, "uuid")
	}
}

func testAccCheckProjectResourceDestroy(s *terraform.State) error {
	uuid, err := getAttributefromState(s, "flagsmith_project.test_project", "uuid")
	if err != nil {
		return err
	}

	_, err = testClient().GetProject(uuid)
	if err == nil {
		return fmt.Errorf("project still exists")
	}
	return nil

}

func testAccProjectResourceConfig(projectName string) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}
data "flagsmith_organisation" "test_org" {
  uuid = "%s"
}

resource "flagsmith_project" "test_project" {
  name = "%s"
  organisation_id = data.flagsmith_organisation.test_org.id

}

`,organisationUUID(), projectName)
}
