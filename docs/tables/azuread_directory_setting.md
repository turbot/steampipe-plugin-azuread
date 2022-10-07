# Table: azuread_directory_setting

Directory settings define the configurations that can be used to customize the tenant-wide and object-specific restrictions and allowed behavior.

By default, all entities inherit the preset defaults.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  values
from
  azuread_directory_setting;
```

### Check if user admin consent workflow is enabled

```sql
select
  display_name,
  id,
  name,
  value
from
  azuread_directory_setting
where
  display_name = 'Consent Policy Settings'
  and name = 'EnableAdminConsentRequests'
  and value::bool;
```

### Check if banned password protection is enabled

```sql
select
  display_name,
  id,
  name,
  value
from
  azuread_directory_setting
where
  display_name = 'Password Rule Settings'
  and (
    name = 'EnableBannedPasswordCheckOnPremises'
    and value::bool
  ) and (
    name = 'BannedPasswordCheckOnPremisesMode'
    and value = 'Enforced'
  );
```
