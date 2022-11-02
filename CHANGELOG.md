## v0.8.3 [2022-11-02]

_Bug fixes_

- Updated `ip_address` column type to string in `azuread_sign_in_report` table to handle partial IP addresses. ([#95](https://github.com/turbot/steampipe-plugin-azuread/pull/95))

## v0.8.2 [2022-10-28]

_Bug fixes_

- Updated the `docs/index.md` file to add the missing permission `Policy.Read.All` which is required to query `azuread_admin_consent_request_policy`, `azuread_authorization_policy` and `azuread_conditional_access_policy` tables. ([#92](https://github.com/turbot/steampipe-plugin-azuread/pull/92))

## v0.8.1 [2022-10-21]

_Bug fixes_

- Disabled caching of Graph client to avoid errors when running consecutive queries for some tables. ([#90](https://github.com/turbot/steampipe-plugin-azuread/pull/90))

## v0.8.0 [2022-10-13]

_What's new?_

- New tables added
  - [azuread_admin_consent_request_policy](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_admin_consent_request_policy) ([#73](https://github.com/turbot/steampipe-plugin-azuread/pull/73))
  - [azuread_directory_audit_report](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_directory_audit_report) ([#72](https://github.com/turbot/steampipe-plugin-azuread/pull/72))
  - [azuread_directory_setting](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_directory_setting) ([#75](https://github.com/turbot/steampipe-plugin-azuread/pull/75))
  - [azuread_security_defaults_policy](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_security_defaults_policy) ([#77](https://github.com/turbot/steampipe-plugin-azuread/pull/77))

_Bug fixes_

- Fixed default page size in all tables. ([#85](https://github.com/turbot/steampipe-plugin-azuread/pull/85))

## v0.7.0 [2022-09-29]

_What's new?_

- New tables added
  - [azuread_device](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_device) ([#68](https://github.com/turbot/steampipe-plugin-azuread/pull/68)) (Thanks [@chrichts](https://github.com/chrichts) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#79](https://github.com/turbot/steampipe-plugin-azuread/pull/79))
- Recompiled plugin with Go version `1.19`. ([#79](https://github.com/turbot/steampipe-plugin-azuread/pull/79))

## v0.6.0 [2022-08-09]

_What's new?_

- New tables added
  - [azuread_authorization_policy](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_authorization_policy) ([#67](https://github.com/turbot/steampipe-plugin-azuread/pull/67))
- Updated plugin to use the [Microsoft Graph SDK for Go](https://github.com/microsoftgraph/msgraph-sdk-go). There should be no large changes in functionality, but if you notice any, please create a [new issue](https://github.com/turbot/steampipe-plugin-azuread/issues/new/choose). ([#62](https://github.com/turbot/steampipe-plugin-azuread/pull/62))

_Enhancements_

- Added `description` column to `azuread_application` table.
- Added `filter` column to `azuread_identity_provider` table.
- Added `app_id`, `app_description`, `description`, `login_url`, `logout_url`, and `oauth2_permission_scopes` columns to `azuread_service_principal` table.

_Bug fixes_

- Fixed column name `keyCredentials` to `key_credent4ials` in `azuread_service_principal` table.
- Fixed column type of `ip_address` column in `azuread_sign_in_report` from string to IP address.

_Breaking changes_

- Removed `verified_publisher` column from `azuread_service_principal` table due to lack of API support.
- Removed `published_permission_scopes` column from `azuread_service_principal` table (replaced by `oauth2_permission_scopes` column).
- Removed `refresh_tokens_valid_from_date_time` column from `azuread_user` table (replaced by `sign_in_sessions_valid_from_date_time` column).

## v0.5.1 [2022-07-25]

_Bug fixes_

- Added the missing `</li>` list item element tag in the credentials section of `docs/index.md` page which would cause the plugin build process to fail. ([#60](https://github.com/turbot/steampipe-plugin-azuread/pull/60))

## v0.5.0 [2022-07-22]

_Enhancements_

- Improved `docs/index.md` file information about what permissions are required to query resources. ([#56](https://github.com/turbot/steampipe-plugin-azuread/pull/56))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11) which includes several caching fixes. ([#58](https://github.com/turbot/steampipe-plugin-azuread/pull/58))

## v0.4.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#54](https://github.com/turbot/steampipe-plugin-azuread/pull/54))

## v0.4.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#52](https://github.com/turbot/steampipe-plugin-azuread/pull/52))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#48](https://github.com/turbot/steampipe-plugin-azuread/pull/48))

## v0.3.1 [2022-04-22]

_Bug fixes_

- Fixed the following tables' results being limited to 999 rows unintentionally: ([#50](https://github.com/turbot/steampipe-plugin-azuread/pull/50))
  - azuread_application
  - azuread_conditional_access_policy
  - azuread_domain
  - azuread_group
  - azuread_service_principal
  - azuread_sign_in_report
  - azuread_user

## v0.3.0 [2022-03-25]

_Enhancements_

- Added additional optional key quals, filter support, page limit and context cancellation handling to `azuread_application`, `azuread_conditional_access_policy`, `azuread_directory_role`, `azuread_domain`, `azuread_group`, `azuread_identity_provider`, `azuread_service_principal`, `azuread_sign_in_report`, and `azuread_user` tables ([#43](https://github.com/turbot/steampipe-plugin-azuread/pull/43))

## v0.2.0 [2022-03-23]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v211--2022-03-10) ([#44](https://github.com/turbot/steampipe-plugin-azuread/pull/44))

## v0.1.0 [2021-12-08]

_Enhancements_

- Recompiled plugin with Go version 1.17 ([#40](https://github.com/turbot/steampipe-plugin-azuread/pull/40))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#39](https://github.com/turbot/steampipe-plugin-azuread/pull/39))
- Added `on_premises_immutable_id` column to the `azuread_user` table ([#37](https://github.com/turbot/steampipe-plugin-azuread/pull/37))

## v0.0.3 [2021-11-03]

_What's new?_

- New tables added
  - [azuread_sign_in_report](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_sign_in_report) ([#26](https://github.com/turbot/steampipe-plugin-azuread/pull/26))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.7.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v171--2021-11-01) ([#33](https://github.com/turbot/steampipe-plugin-azuread/pull/33))

_Bug fixes_

- Fixed the title of the `azuread_user` table in the documentation file ([#32](https://github.com/turbot/steampipe-plugin-azuread/pull/32))

## v0.0.2 [2021-09-28]

_What's new?_

- New tables added
  - [azuread_conditional_access_policy](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_conditional_access_policy) ([#21](https://github.com/turbot/steampipe-plugin-azuread/pull/21))

_Bug fixes_

- Updated `on_premises_net_bios_name` column type from `timestamp` to `string` in `azuread_group` table ([#25](https://github.com/turbot/steampipe-plugin-azuread/pull/25))
- Added missing config options for managed identity in credential setup ([#16](https://github.com/turbot/steampipe-plugin-azuread/pull/16))
- Fixed example queries in the documentation of `azuread_user` and `azuread_service_principal` tables ([#20](https://github.com/turbot/steampipe-plugin-azuread/pull/20)) ([#29](https://github.com/turbot/steampipe-plugin-azuread/pull/29))

## v0.0.1 [2021-08-19]

_What's new?_

- New tables added
  - [azuread_application](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_application) ([#8](https://github.com/turbot/steampipe-plugin-azuread/pull/8))
  - [azuread_directory_role](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_directory_role) ([#8](https://github.com/turbot/steampipe-plugin-azuread/pull/8))
  - [azuread_domain](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_domain) ([#8](https://github.com/turbot/steampipe-plugin-azuread/pull/8))
  - [azuread_group](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_group) ([#5](https://github.com/turbot/steampipe-plugin-azuread/pull/5))
  - [azuread_identity_provider](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_identity_provider) ([#8](https://github.com/turbot/steampipe-plugin-azuread/pull/8))
  - [azuread_service_principal](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_service_principal) ([#8](https://github.com/turbot/steampipe-plugin-azuread/pull/8))
  - [azuread_user](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_user) ([#3](https://github.com/turbot/steampipe-plugin-azuread/pull/3))
