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
