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
