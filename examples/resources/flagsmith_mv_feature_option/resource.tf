resource "flagsmith_feature" "feature_1" {
  feature_name = "feature_1"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  description  = "This is a new multivariate feature"
  type         = "MULTIVARIATE"
}

resource "flagsmith_mv_feature_option" "feature_1_mv_option" {
  type                          = "unicode"
  feature_uuid                  = flagsmith_feature.feature_1.uuid
  string_value                  = "option_value_60_percent_of_the_times"
  default_percentage_allocation = 60
}
