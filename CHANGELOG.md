# Changelog

## [0.12.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.11.0...v0.12.0) (2026-08-26)


### Features

* Add data source `flagsmith_feature`
* Add data source `flagsmith_segment`
* Add `metadata` to resources `flagsmith_feature`, `flagsmith_segment` and `flagsmith_environment`, to set custom field values

### Notes

* Custom fields are managed authoritatively: values set outside of Terraform are removed on the next write. The provider previously never sent `metadata` at all, so updating a feature, segment or environment already cleared its UI-set custom field values silently. Those values now appear in the plan as an explicit removal instead — add them to `metadata` to keep them.

### Enhancements

* Update dependencies

## [0.11.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.10.0...v0.11.0) (2026-08-25)


### Features

* Add data source `flagsmith_project`

### Enhancements

* Update dependencies

## [0.10.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.9.1...v0.10.0) (2026-04-06)


### Features

* Add data source `flagsmith_user`
* resource(flagsmith_feature): Add `group_owners`
* resource(flagsmith_project): Add `enforce_feature_owners`

### Bug Fixes

* Update `.goreleaser.yml` to the v2 format, so that releases publish again

### Enhancements

* Add import block examples for Terraform v1.5.0+
* Update dependencies

## [0.9.1](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.9.0...v0.9.1) (2025-01-03)


### Enhancements

* Update dependencies

## [0.9.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.8.2...v0.9.0) (2024-12-06)


### Features

* Add resource `flagsmith_project`
* Add resource `flagsmith_environment`
* Add data resource `flagsmith_organisation`


## [0.8.2](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.8.1...v0.8.2) (2024-11-08)


### Enhancements

* resource(flagsmith_tag): Make `tag_colour` optional


## [0.8.1](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.8.0...v0.8.1) (2024-09-16)


### Bug Fixes

* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/155


## [0.8.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.7.0...v0.8.0) (2024-05-16)


### Notes

* This Go module(and related dependencies) has been updated to GO 1.21 as per the Go Support policy


## [0.7.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.6.0...v0.7.0) (2024-01-30)


### Features

* Add resource `flagsmith_tag`
* Update resource `flagsmith_feature` to add support for tags


### Bug Fixes

* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/67


## [0.6.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.5.1...v0.6.0) (2023-09-19)


### Notes

* This Go module(and related dependencies) has been updated to GO 1.20 as per the Go Support policy


## [0.5.1](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.5.0...v0.5.1) (2023-04-25)


### Bug Fixes

* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/81


## [0.5.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.4.1...v0.5.0) (2023-04-17)


### Bug Fixes

* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/76
* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/68


### Enhancements

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


## [0.4.1](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.4.0...v0.4.1) (2023-02-24)


### Bug Fixes

* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/60


## [0.4.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.3.0...v0.4.0) (2023-02-06)


### ⚠ BREAKING CHANGES

* resource(feature_state): make feature_state_value required
* resource(feature_state): make feature_state_value.type required


### Bug Fixes

* fix https://github.com/Flagsmith/terraform-provider-flagsmith/issues/56


### Enhancements

* Update terraform-plugin-framework


## [0.3.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.2.0...v0.3.0) (2022-12-08)


### ⚠ BREAKING CHANGES

* resource: update import of `flagsmith_feature_state` from `<enviroment_client_key>,<feature_name>` to `<enviroment_client_key>,<feature_state_uuid>`
* resource: replace `feature_name` field with `feature_id` on `flagsmith_feature_state`


### Features

* Add resource `flagsmith_segment`
* Update resource `flagsmith_feature_state` to add support for segment override


### Enhancements

* Update testify
* Update terraform-plugin-go
* Update terraform-plugin-sdk/v2
* Update terraform-plugin-framework


## [0.2.0](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.1.3...v0.2.0) (2022-11-07)


### ⚠ BREAKING CHANGES

* resource: rename `flagsmith_flag` to `flagsmith_feature_state`


### Features

* Add resource `flagsmith_feature`
* Add resource `flagsmith_mv_feature_option`


## [0.1.3](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.1.2...v0.1.3) (2022-08-18)


### Enhancements

* update tf plugin sdk


## [0.1.2](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.1.1...v0.1.2) (2022-07-26)


### Enhancements

* Remove deprecated release action


## [0.1.1](https://github.com/Flagsmith/terraform-provider-flagsmith/compare/v0.1.0...v0.1.1) (2022-07-26)


### Enhancements

* Update dependencies
* Update registry url


## 0.1.0 (2022-07-07)


### Features

* Update Feature state for a given environment
