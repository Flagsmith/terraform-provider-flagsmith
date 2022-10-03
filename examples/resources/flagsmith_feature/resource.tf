resource "flagsmith_feature" "new_mv_feature" {
  feature_name = "new_mv_feature"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  description  = "This is a new multivariate feature"
  type         = "MULTIVARIATE"
  multivariate_options = [
    {
      type : "unicode",
      string_value : "option_value_10",
      default_percentage_allocation : 10
    },
    {
      type : "bool",
      boolean_value : true,
      default_percentage_allocation : 10
    }
  ]
}

resource "flagsmith_feature" "new_standard_feature" {
  feature_name = "new_standard_feature"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  description  = "This is a new standard feature"
  type         = "STANDARD"
}

