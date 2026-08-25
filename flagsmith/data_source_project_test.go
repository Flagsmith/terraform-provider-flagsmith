package flagsmith_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectDataResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataResourceConfig(projectUUID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flagsmith_project.test_project", "id", strconv.Itoa(projectID())),
					resource.TestCheckResourceAttr("data.flagsmith_project.test_project", "uuid", projectUUID()),
					resource.TestCheckResourceAttr("data.flagsmith_project.test_project", "organisation_id", strconv.Itoa(organisationID())),
					resource.TestCheckResourceAttrSet("data.flagsmith_project.test_project", "name"),
				),
			},
		},
	})
}

func testAccProjectDataResourceConfig(projectUUID string) string {
	return fmt.Sprintf(`
%s

data "flagsmith_project" "test_project" {
  uuid = "%s"
}
`, providerConfig(), projectUUID)
}
