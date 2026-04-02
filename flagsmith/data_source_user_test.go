package flagsmith_test

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccUserDataResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			mustHaveEnv(t, "FLAGSMITH_USER_EMAIL")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataResourceConfig(organisationUUID(), userEmail()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flagsmith_user.test_user", "email", userEmail()),
					resource.TestCheckResourceAttr("data.flagsmith_user.test_user", "organisation_id", strconv.Itoa(organisationID())),
					resource.TestCheckResourceAttrSet("data.flagsmith_user.test_user", "id"),
					resource.TestCheckResourceAttrSet("data.flagsmith_user.test_user", "first_name"),
					resource.TestCheckResourceAttrSet("data.flagsmith_user.test_user", "last_name"),
					resource.TestCheckResourceAttrSet("data.flagsmith_user.test_user", "role"),
				),
			},
		},
	})
}

func testAccUserDataResourceConfig(organisationUUID, email string) string {
	return fmt.Sprintf(`
%s

data "flagsmith_organisation" "test_org" {
  uuid = "%s"
}

data "flagsmith_user" "test_user" {
  organisation_id = data.flagsmith_organisation.test_org.id
  email           = "%s"
}
`, providerConfig(), organisationUUID, email)
}
