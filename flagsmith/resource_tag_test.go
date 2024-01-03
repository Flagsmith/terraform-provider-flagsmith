package flagsmith_test
import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccTagResource(t *testing.T) {
	tagName := acctest.RandString(16)
	tagColour := "#f1d502"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTagResourceConfig(tagName, tagColour, "tag description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "tag_name", tagName),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "description", "tag description"),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "tag_colour", tagColour),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "project_uuid", projectUUID()),

					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "project_id"),

				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_tag.test_tag",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getTagImportID("flagsmith_tag.test_tag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "tag_name", tagName),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "tag_colour", tagColour),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "description", "new tag descriptionnnnn"),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "project_uuid", projectUUID()),

					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "project_id"),
				),
			},

			// Update testing
			{
				Config: testAccTagResourceConfig(tagName, tagColour, "updated tag description"),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "tag_name", tagName),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "description", "updated tag description"),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "tag_colour", tagColour),
					resource.TestCheckResourceAttr("flagsmith_tag.test_tag", "project_uuid", projectUUID()),

					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_tag.test_tag", "project_id"),


				),
			},

		},
	})
}

func testAccTagResourceConfig(tagName, tagColour,  description string) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_tag" "test_tag" {
  tag_name = "%s"
  tag_colour = "%s"
  description = "%s"
  project_uuid = "%s"
}

`, tagName, tagColour, description, projectUUID() )
}

func testAccCheckTagResourceDestroy(s *terraform.State) error {
	uuid, err := getAttributefromState(s, "flagsmith_tag.test_tag","uuid")
	projectUUID, err := getAttributefromState(s, "flagsmith_tag.test_tag","project_uuid")

	if err != nil {
		return err
	}

	_, err = testClient().GetTag(projectUUID, uuid)
	if err == nil {
		return fmt.Errorf("tag still exists")
	}
	return nil

}
func getTagImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		// return a string in the format of projectUUID,tagUUID
		projectUUID, err := getAttributefromState(s, n, "project_uuid")
		if err != nil {
			return "", err
		}
		tagUUID, err := getAttributefromState(s, n, "uuid")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s,%s", projectUUID, tagUUID), nil

	}
}
