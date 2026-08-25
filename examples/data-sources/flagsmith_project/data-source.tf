data "flagsmith_project" "my_project" {
  uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
}

resource "flagsmith_environment" "staging" {
  name       = "Staging"
  project_id = data.flagsmith_project.my_project.id
}
