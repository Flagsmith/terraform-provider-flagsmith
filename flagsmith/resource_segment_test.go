package flagsmith_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)



func TestAccSegmentResource(t *testing.T) {
	segmentName :=  acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSegmentResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSegmentResourceConfig(segmentName, "new segment description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "name", segmentName),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "description", "new segment description"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "project_uuid", projectUUID()),

					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.type","ALL"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.type","ANY"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.operator","EQUAL"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.property","device_type"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.value","mobile"),

					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "project_id"),
					resource.TestCheckNoResourceAttr("flagsmith_segment.test_segment", "feature_id"),

				),
			},

			// ImportState testing
			{
				ResourceName:      "flagsmith_segment.test_segment",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getSegmentImportID("flagsmith_segment.test_segment"),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "name", segmentName),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "description", "new segment description"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "project_uuid", projectUUID()),

					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.type","ALL"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.type","ANY"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.operator","EQUAL"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.property","device_type"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.value","mobile"),


					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "project_id"),
					resource.TestCheckNoResourceAttr("flagsmith_segment.test_segment", "feature_id"),

				),
			},

			// Update testing
			{
				Config: testAccSegmentResourceConfig(segmentName, "segment description updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "name", segmentName),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "description", "segment description updated"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "project_uuid", projectUUID()),

					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.type","ALL"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.type","ANY"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.operator","EQUAL"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.property","device_type"),
					resource.TestCheckResourceAttr("flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.value","mobile"),

					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "id"),
					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "uuid"),
					resource.TestCheckResourceAttrSet("flagsmith_segment.test_segment", "project_id"),
					resource.TestCheckNoResourceAttr("flagsmith_segment.test_segment", "feature_id"),


				),
			},
		},
	})
}



func getSegmentImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		return getAttributefromState(s, n, "uuid")
	}
}

func testAccCheckSegmentResourceDestroy(s *terraform.State) error {
	uuid, err := getAttributefromState(s, "flagsmith_segment.test_segment", "uuid")
	if err != nil {
		return err
	}

	_, err = testClient().GetSegment(uuid)
	if err == nil {
		return fmt.Errorf("segment still exists")
	}
	return nil



}


func testAccSegmentResourceConfig(segmentName, description string) string {
	return fmt.Sprintf(`
provider "flagsmith" {

}

resource "flagsmith_segment" "test_segment" {
  name         = "%s"
  description = "%s"
  project_uuid = "%s"
  rules = [
    {
      "rules" : [{
        "conditions" : [{
          "operator" : "EQUAL",
          "property" : "device_type",
          "value" : "mobile"
        }],
        "type" : "ANY"
      }],
      "type" : "ALL"

    }
  ]

}


`, segmentName, description, projectUUID())
}


