package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	//"strings"
	"testing"
	"strconv"
)

func TestAccEnvironmentResource(t *testing.T) {
	environmentName := acctest.RandString(16)


	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEnvironmentResourceConfig(environmentName, projectID(), "new environment"),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "name", environmentName),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "project_id", strconv.Itoa(projectID())),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "description", "new environment"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "minimum_change_request_approvals", "0"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "hide_disabled_flags", "false"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "hide_sensitive_data", "false"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "allow_client_traits", "true"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "use_identity_composite_key_for_hashing", "true"),

					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "api_key"),

				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_environment.test_environment",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getEnvironmentImportID("flagsmith_environment.test_environment"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "name", environmentName),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "project_id", strconv.Itoa(projectID())),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "description", "new environment"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "minimum_change_request_approvals", "0"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "hide_disabled_flags", "false"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "hide_sensitive_data", "false"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "allow_client_traits", "true"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "use_identity_composite_key_for_hashing", "true"),

					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "api_key"),

				),
			},

			// Update testing
			{
				Config: testAccEnvironmentResourceConfig(environmentName, projectID(), "updated environment"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "name", environmentName),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "project_id", strconv.Itoa(projectID())),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "description", "updated environment"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "minimum_change_request_approvals", "0"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "hide_disabled_flags", "false"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "hide_sensitive_data", "false"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "allow_client_traits", "true"),
					resource.TestCheckResourceAttr("flagsmith_environment.test_environment", "use_identity_composite_key_for_hashing", "true"),

					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_environment.test_environment", "api_key"),

				),
			},
		},
	})
}


func getEnvironmentImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		return getAttributefromState(s, n, "uuid")
	}
}

func testAccCheckEnvironmentResourceDestroy(s *terraform.State) error {
	id , err := getAttributefromState(s, "flagsmith_environment.test_environment", "id")
	if err != nil {
		return err
	}

	_, err = testClient().GetEnvironment(id)
	if err == nil {
		return fmt.Errorf("environment still exists")
	}
	return nil

}

func testAccEnvironmentResourceConfig(environmentName string, projectID int,  description string ) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_environment" "test_environment" {
  name = "%s"
  project_id = %d
  description = "%s"
}

`,environmentName, projectID, description)
}
