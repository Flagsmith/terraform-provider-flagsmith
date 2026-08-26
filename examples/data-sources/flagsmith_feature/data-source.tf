# Look up a feature that was created outside of Terraform, so that environment
# specific overrides can be managed independently of the feature itself.
data "flagsmith_feature" "my_feature" {
  uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
}

resource "flagsmith_feature_state" "staging_override" {
  feature_id      = data.flagsmith_feature.my_feature.id
  environment_key = "some-environment-key"
  enabled         = true
  feature_state_value = {
    type         = "unicode"
    string_value = "staging-value"
  }
}
