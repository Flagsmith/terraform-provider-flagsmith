## 0.12.0
FEATURES:
* Add data source `flagsmith_feature`
* Add data source `flagsmith_segment`
* Add `metadata` to resources `flagsmith_feature`, `flagsmith_segment` and `flagsmith_environment`, to set custom field values

NOTES:
* Custom fields are managed authoritatively: values set outside of Terraform are removed on the next write. The provider previously never sent `metadata` at all, so updating a feature, segment or environment already cleared its UI-set custom field values silently. Those values now appear in the plan as an explicit removal instead — add them to `metadata` to keep them.

ENHANCEMENTS:
* Update dependencies

## 0.11.0
FEATURES:
* Add data source `flagsmith_project`

ENHANCEMENTS:
* Update dependencies

## 0.10.0
FEATURES:
* Add data source `flagsmith_user`
* resource(flagsmith_feature): Add `group_owners`
* resource(flagsmith_project): Add `enforce_feature_owners`

BUG FIXES
* Update `.goreleaser.yml` to the v2 format, so that releases publish again

ENHANCEMENTS:
* Add import block examples for Terraform v1.5.0+
* Update dependencies

## 0.9.1
ENHANCEMENTS:
* Update dependencies

## 0.9.0
FEATURES:
* Add resource `flagsmith_project`
* Add resource `flagsmith_environment`
* Add data resource `flagsmith_organisation`


## 0.8.2
ENHANCEMENTS:
* resource(flagsmith_tag): Make `tag_colour` optional

## 0.8.1
BUG FIXES
fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/155


## 0.8.0
NOTES:
* This Go module(and related dependencies) has been updated to GO 1.21 as per the Go Support policy

## 0.7.0
FEATURES:
* Add resource `flagsmith_tag`
* Update resource `flagsmith_feature` to add support for tags

BUG FIXES
fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/67

## 0.6.0
NOTES:
* This Go module(and related dependencies) has been updated to GO 1.20 as per the Go Support policy

## 0.5.1
BUG FIXES
fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/81


## 0.5.0
BUG_FIXES
* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/76
* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/68

ENHANCEMENTS:
* resource(feature): Make `type` optional and use `STANDARD` as default value
* resource(feature): Add default value(false) for `default_enabled`
* resource(feature): Add RequiresReplace plan modifier to `project_uuid` field

* resource(feature_state): Add RequiresReplace plan modifier to `environment_key` field
* resource(feature_state): Add RequiresReplace plan modifier to `feature_id` field

* resource(multivariate_options): Add RequiresReplace plan modifier to `feature_uuid`

* resource(segment): Add RequiresReplace plan modifier to `project_uuid` field

* Update Go module to GO 1.19
* Update Terraform-plugin-framework
* Update testify
* Update terraform plugin sdk
* Update terraform-plugin-go
* Update terraform-plugin-docs

## 0.4.1
BUG_FIXES
* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/60

## 0.4.0
ENHANCEMENTS:
* Update terraform-plugin-framework

BUG_FIXES
* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/56

BREAKING CHANGES:
* resource(feature_state): make feature_state_value required
* resource(feature_state): make feature_state_value.type required


## 0.3.0
BREAKING CHANGES:
* resource: update import of `flagsmith_feature_state` from `<enviroment_client_key>,<feature_name>` to `<enviroment_client_key>,<feature_state_uuid>`
* resource: replace `feature_name` field with `feature_id` on `flagsmith_feature_state`

FEATURES:

* Add resource `flagsmith_segment`
* Update resource `flagsmith_feature_state` to add support for segment override

ENHANCEMENTS:
* Update testify
* Update terraform-plugin-go
* Update terraform-plugin-sdk/v2
* Update terraform-plugin-framework

## 0.2.0
BREAKING CHANGES:

* resource: rename `flagsmith_flag` to `flagsmith_feature_state`

FEATURES:

* Add resource `flagsmith_feature`
* Add resource `flagsmith_mv_feature_option`

## 0.1.3

ENHANCEMENTS:

* update tf plugin sdk


## 0.1.2

ENHANCEMENTS:

* Remove deprecated release action

## 0.1.1

ENHANCEMENTS:

* Update dependencies
* Update registry url

## 0.1.0

FEATURES:

* Update Feature state for a given environment
