# Table: azuread_device

Azure Active Directory (Azure AD) detects devices connecting to your Azure accounts. Devices managed with Intune are also tracked in Azure AD.

## Examples

### Basic info

```sql
select
  display_name,
  is_managed,
  is_compliant,
  member_of
from
  azuread_device;
```

### List managed devices

```sql
select
  display_name,
  profile_type,
  id,
  operating_system,
  operating_system_version
from
  azuread_device
where
  is_managed;
```

### List non-compliant devices

```sql
select
  display_name,
  profile_type,
  id,
  operating_system,
  operating_system_version
from
  azuread_device
where
  not is_complaint;
```
