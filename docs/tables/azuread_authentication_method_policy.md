# azuread_authentication_method_policy

Represents the authentication methods policy for the Microsoft Entra tenant.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  description,
  title
from
  azuread_authentication_method_policy;
```

### Get authentication method policy details

```sql
select
  id,
  display_name,
  description,
  last_modified_date_time,
  policy_migration_state,
  policy_version,
  reconfirmation_in_days,
  registration_enforcement,
  authentication_method_configurations
from
  azuread_authentication_method_policy;
```

### Check registration enforcement settings

```sql
select
  id,
  display_name,
  registration_enforcement
from
  azuread_authentication_method_policy;
```

### Review authentication method configurations

```sql
select
  id,
  display_name,
  authentication_method_configurations
from
  azuread_authentication_method_policy;
```

### Policy compliance check

```sql
select
  id,
  display_name,
  description,
  policy_migration_state,
  policy_version,
  registration_enforcement
from
  azuread_authentication_method_policy;
```

## Columns

| Name                                   | Type        | Description                                                                                                                                             |
| -------------------------------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `tenant_id`                            | `text`      | The Azure Active Directory tenant ID.                                                                                                                   |
| `id`                                   | `text`      | The identifier for the authentication methods policy.                                                                                                   |
| `display_name`                         | `text`      | The display name of the authentication methods policy.                                                                                                  |
| `description`                          | `text`      | The description of the authentication methods policy.                                                                                                   |
| `last_modified_date_time`              | `timestamp` | The date and time of the last update to the policy.                                                                                                     |
| `policy_migration_state`               | `text`      | The state of migration of the authentication methods policy from the legacy multifactor authentication and self-service password reset (SSPR) policies. |
| `policy_version`                       | `text`      | The version of the policy in use.                                                                                                                       |
| `reconfirmation_in_days`               | `integer`   | The reconfirmation in days property.                                                                                                                    |
| `registration_enforcement`             | `jsonb`     | Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods.                              |
| `authentication_method_configurations` | `jsonb`     | Represents the settings for each authentication method.                                                                                                 |
| `title`                                | `text`      | Title of the resource.                                                                                                                                  |
| `sp_connection_name`                   | `text`      | Steampipe connection name.                                                                                                                              |
| `sp_ctx`                               | `jsonb`     | Steampipe context in JSON form.                                                                                                                         |
