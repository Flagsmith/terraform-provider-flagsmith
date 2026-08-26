package flagsmith_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSegmentDataResource(t *testing.T) {
	segmentName := "tf_test_" + strings.ToLower(acctest.RandString(12))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the segment on its own first, so that the data source read in the
			// next step cannot race the create.
			{
				Config: testAccSegmentDataResourceSegmentConfig(segmentName, "null"),
			},
			{
				Config: testAccSegmentDataResourceConfig(segmentName, "null"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "name", segmentName),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "project_uuid", projectUUID()),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "project_id", strconv.Itoa(projectID())),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "description", "created by the segment data source acceptance test"),
					resource.TestCheckResourceAttrSet("data.flagsmith_segment.test_segment", "id"),

					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "rules.0.type", "ALL"),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "rules.0.rules.0.type", "ANY"),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.operator", "EQUAL"),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.property", "device_type"),
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "rules.0.rules.0.conditions.0.value", "mobile"),

					resource.TestCheckResourceAttrPair(
						"data.flagsmith_segment.test_segment", "uuid",
						"flagsmith_segment.test_segment", "uuid"),
					resource.TestCheckResourceAttrPair(
						"data.flagsmith_segment.test_segment", "id",
						"flagsmith_segment.test_segment", "id"),
				),
			},
		},
	})
}

func TestAccFeatureSpecificSegmentDataResource(t *testing.T) {
	segmentName := "tf_test_" + strings.ToLower(acctest.RandString(12))
	featureSpecific := strconv.Itoa(featureID())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentDataResourceSegmentConfig(segmentName, featureSpecific),
			},
			{
				Config: testAccSegmentDataResourceConfig(segmentName, featureSpecific),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.flagsmith_segment.test_segment", "feature_id", featureSpecific),
					resource.TestCheckResourceAttrPair(
						"data.flagsmith_segment.test_segment", "uuid",
						"flagsmith_segment.test_segment", "uuid"),
				),
			},
		},
	})
}

func TestAccSegmentDataResourceNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "flagsmith_segment" "test_segment" {
  uuid = "8f1b0e1a-0000-4000-8000-000000000000"
}
`, providerConfig()),
				ExpectError: regexp.MustCompile("Unable to read segment"),
			},
		},
	})
}

// featureID is either "null" or a feature ID, to create a feature specific segment.
func testAccSegmentDataResourceSegmentConfig(segmentName, featureID string) string {
	return fmt.Sprintf(`
%s

resource "flagsmith_segment" "test_segment" {
  name         = "%s"
  description  = "created by the segment data source acceptance test"
  project_uuid = "%s"
  feature_id   = %s
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
`, providerConfig(), segmentName, projectUUID(), featureID)
}

func testAccSegmentDataResourceConfig(segmentName, featureID string) string {
	return testAccSegmentDataResourceSegmentConfig(segmentName, featureID) + `
data "flagsmith_segment" "test_segment" {
  uuid = flagsmith_segment.test_segment.uuid
}
`
}
