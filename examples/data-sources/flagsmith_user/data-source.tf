data "flagsmith_organisation" "my_org" {
  uuid = "0ee0578e-f2b8-467d-ba49-4cda324cea91"
}

data "flagsmith_user" "john" {
  organisation_id = data.flagsmith_organisation.my_org.id
  email           = "john@example.com"
}

resource "flagsmith_feature" "my_feature" {
  feature_name = "my_feature"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  owners       = [data.flagsmith_user.john.id]
}
