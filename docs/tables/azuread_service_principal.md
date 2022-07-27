# Table: azuread_service_principal

Azure Active Directory (Azure AD) lets you use service principal to manage access to your cloud-based resources.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  app_display_name
from
  azuread_service_principal;
```

### List disabled service principals

```sql
select
  display_name,
  id
from
  azuread_service_principal
where
  not account_enabled;
```

### List service principals related to applications

```sql
select
  id,
  app_display_name,
  account_enabled
from
  azuread_service_principal
where
  service_principal_type = 'Application'
  and tenant_id = app_owner_organization_id;
```
