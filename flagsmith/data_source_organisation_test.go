package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
	"strconv"
)

func TestAccOrganisationDataResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganisationDataResourceConfig(organisationUUID()),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.flagsmith_organisation.test_org", "id", strconv.Itoa( organisationID())),
					resource.TestCheckResourceAttr("data.flagsmith_organisation.test_org", "uuid", organisationUUID()),

					resource.TestCheckResourceAttrSet("data.flagsmith_organisation.test_org", "name"),

				),
			},
		},
	})
}


func testAccOrganisationDataResourceConfig(organisationUUID string) string {
	return fmt.Sprintf(`
provider "flagsmith" {
}

data "flagsmith_organisation" "test_org" {
  uuid = "%s"
}

`,organisationUUID)
}
