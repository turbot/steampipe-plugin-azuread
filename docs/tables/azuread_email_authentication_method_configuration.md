# azuread_email_authentication_method_configuration

Represents the email OTP authentication method policy for the Microsoft Entra tenant.

## Examples

### Basic info

```sql
select
  id,
  state,
  allow_external_id_to_use_email_otp,
  title
from
  azuread_email_authentication_method_configuration;
```

### Get email authentication method configuration details

```sql
select
  id,
  state,
  allow_external_id_to_use_email_otp,
  include_targets,
  exclude_targets
from
  azuread_email_authentication_method_configuration;
```

## Columns

| Name                               | Type  | Description                                                                                                                   |
| ---------------------------------- | ----- | ----------------------------------------------------------------------------------------------------------------------------- |
| tenant_id                          | text  | The Azure Active Directory tenant ID.                                                                                         |
| id                                 | text  | The identifier for the authentication method configuration.                                                                   |
| state                              | text  | The state of the authentication method configuration. Possible values are: enabled, disabled.                                 |
| allow_external_id_to_use_email_otp | text  | Determines whether email OTP is usable by external users for authentication. Possible values are: default, enabled, disabled. |
| include_targets                    | jsonb | A collection of users or groups who are enabled to use the authentication method.                                             |
| exclude_targets                    | jsonb | A collection of users or groups who are excluded from using the authentication method.                                        |
| title                              | text  | Title of the resource.                                                                                                        |
