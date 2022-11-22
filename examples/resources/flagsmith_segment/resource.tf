resource "flagsmith_feature" "new_standard_feature" {
  feature_name = "new_standard_feature"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  description  = "This is a new standard feature"
  type         = "STANDARD"
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


resource "flagsmith_segment" "new_standard_feature_specific_segment" {
  name         = "device_type"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  feature_id   = flagsmith_feature.new_standard_feature.id
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
