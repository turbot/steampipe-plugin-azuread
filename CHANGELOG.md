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
