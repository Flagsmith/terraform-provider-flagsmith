resource "flagsmith_flag" "feature_1_dev" {
  enabled         = true
  environment     = 2
  feature         = 14
  environment_key = "<environment_key>"
  feature_name    = "feature_1"
  feature_state_value = {
    type         = "unicode"
    string_value = "some_flag_value"
  }

}
