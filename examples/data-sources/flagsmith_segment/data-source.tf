# Look up a feature and a segment that were created outside of Terraform, and manage
# only the segment override for them.
data "flagsmith_feature" "my_feature" {
  uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
}

data "flagsmith_segment" "mobile_users" {
  uuid = "f6c714d3-94e7-4b14-9117-8dd9db91bc19"
}

resource "flagsmith_feature_state" "mobile_override" {
  feature_id       = data.flagsmith_feature.my_feature.id
  environment_key  = "some-environment-key"
  segment_id       = data.flagsmith_segment.mobile_users.id
  segment_priority = 0
  enabled          = true
  feature_state_value = {
    type         = "unicode"
    string_value = "mobile-value"
  }
}
