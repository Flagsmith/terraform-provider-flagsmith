---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "flagsmith_organisation Data Source - terraform-provider-flagsmith"
subcategory: ""
description: |-
  Flagsmith Organisation/ Remote config
---

# flagsmith_organisation (Data Source)

Flagsmith Organisation/ Remote config



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `uuid` (String) UUID of the organisation

### Read-Only

- `force_2fa` (Boolean) If true, signup will require 2FA
- `id` (Number) ID of the organisation
- `name` (String) Name of the organisation
- `persist_trait_data` (Boolean) If false, trait data for this organisation identities will not stored
- `restrict_project_create_to_admin` (Boolean) If true, only organisation admin can create projects