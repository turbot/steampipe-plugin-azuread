## v1.3.0 [2025-05-30]

_Enhancements_

- Added `external_user_state` and `sign_in_activity` columns to `azuread_user` table. ([#250](https://github.com/turbot/steampipe-plugin-azuread/pull/250)) (Thanks [@MarkusGnigler](https://github.com/MarkusGnigler) for the contribution!)
- Added `disable_resilience_defaults` column to `azuread_conditional_access_policy` table. ([#251](https://github.com/turbot/steampipe-plugin-azuread/pull/251)) (Thanks [@MarkusGnigler](https://github.com/MarkusGnigler) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.11.6](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5116-2025-05-22) which improves how errors are handled during query execution. ([#249](https://github.com/turbot/steampipe-plugin-azuread/pull/249))

## v1.2.0 [2025-05-20]

_Enhancements_

- Added `allowedToCreateTenants` field under `default_user_role_permissions` column of `azuread_authorization_policy` table. ([#243](https://github.com/turbot/steampipe-plugin-azuread/pull/243)) (Thanks [@MarkusGnigler](https://github.com/MarkusGnigler) for the contribution!)

## v1.1.0 [2025-04-18]

_Enhancements_

- Added `employee_*` and `on_premises_*` columns to the `azuread_user` table. ([#210](https://github.com/turbot/steampipe-plugin-azuread/pull/210))

_Dependencies_

- Recompiled plugin with Go version `1.23.1`. ([#231](https://github.com/turbot/steampipe-plugin-azuread/pull/231))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.11.5/CHANGELOG.md#v5115-2025-03-31) that addresses critical and high vulnerabilities in dependent packages. ([#231](https://github.com/turbot/steampipe-plugin-azuread/pull/231))

## v0.16.0 [2024-05-14]

_Enhancements_

- The `tenant_id` column has now been assigned as a connection key column across all the tables which facilitates more precise and efficient querying across multiple Azure subscriptions. ([#175](https://github.com/turbot/steampipe-plugin-azuread/pull/175))
- The Plugin and the Steampipe Anywhere binaries are now built with the `netgo` package. ([#180](https://github.com/turbot/steampipe-plugin-azuread/pull/180))
- Added support for `China cloud` endpoint and scope based on the environment. ([#174](https://github.com/turbot/steampipe-plugin-azuread/pull/174))
- Added the `version` flag to the plugin's Export tool. ([#65](https://github.com/turbot/steampipe-export/pull/65))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.10.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v5101-2024-05-09) which ensures that `QueryData` passed to `ConnectionKeyColumns` value callback is populated with `ConnectionManager`. ([#175](https://github.com/turbot/steampipe-plugin-azuread/pull/175))

## v0.15.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/install/steampipe.sh), as a [Postgres FDW](https://steampipe.io/install/postgres.sh), as a [SQLite extension](https://steampipe.io/install/sqlite.sh) and as a standalone [exporter](https://steampipe.io/install/export.sh).
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension.
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-azuread/blob/main/docs/LICENSE).

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server enacapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#158](https://github.com/turbot/steampipe-plugin-azuread/pull/158))

## v0.14.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#142](https://github.com/turbot/steampipe-plugin-azuread/pull/142))

## v0.14.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#139](https://github.com/turbot/steampipe-plugin-azuread/pull/139))
- Recompiled plugin with Go version `1.21`. ([#139](https://github.com/turbot/steampipe-plugin-azuread/pull/139))

## v0.13.0 [2023-09-27]

_Enhancements_

- Added the `department` column to `azuread_user` table. ([#132](https://github.com/turbot/steampipe-plugin-azuread/pull/132))

_Bug fixes_

- Fixed the `title` column in `azuread_device` and `azuread_user` tables to correctly return data instead of null. ([#134](https://github.com/turbot/steampipe-plugin-azuread/pull/134))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v551-2023-07-26). ([#127](https://github.com/turbot/steampipe-plugin-azuread/pull/127))

## v0.12.0 [2023-07-18]

_Enhancements_

- Updated the `docs/index.md` file to include multi-tenant configuration examples. ([#122](https://github.com/turbot/steampipe-plugin-azuread/pull/122))

## v0.11.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#119](https://github.com/turbot/steampipe-plugin-azuread/pull/119))

## v0.10.1 [2023-06-16]

_Bug fixes_

- Fixed the `ListConfig` of `azuread_directory_audit_report` and `azuread_sign_in_report` tables to prevent errors during type conversion caused by inconsistent API responses. ([#116](https://github.com/turbot/steampipe-plugin-azuread/pull/116))

## v0.10.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#112](https://github.com/turbot/steampipe-plugin-azuread/pull/112))

_Bug fixes_

- Fixed the `tags` column in `azuread_group` table to be of JSON type instead of string. ([#111](https://github.com/turbot/steampipe-plugin-azuread/pull/111))

## v0.9.0 [2023-04-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#103](https://github.com/turbot/steampipe-plugin-azuread/pull/103))

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
