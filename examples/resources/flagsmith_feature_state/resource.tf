resource "flagsmith_feature" "new_standard_feature" {
  feature_name = "new_standard_feature"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  description  = "This is a new standard feature"
  type         = "STANDARD"
}


resource "flagsmith_feature_state" "feature_1_dev" {
  enabled         = true
  environment_key = "<environment_key>"
  feature         = flagsmith_feature.new_standard_feature.id
  feature_state_value = {
    type         = "unicode"
    string_value = "some_flag_value"
  }

}

resource "flagsmith_segment" "device_type_segment" {
  name         = "device_type"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
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

resource "flagsmith_feature_state" "feature_1_dev_segment_override" {
  enabled          = true
  environment_key  = "<environment_key>"
  feature          = flagsmith_feature.new_standard_feature.id
  segment          = flagsmith_segment.device_type_segment.id
  segment_priority = 0
  feature_state_value = {
    type         = "unicode"
    string_value = "segment_override_value"
  }
}
