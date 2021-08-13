# Table: azuread_application

Azure Active Directory (Azure AD) lets you use applications to manage access to your cloud-based resources.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  app_id
from
  azuread_application;
```

### List applications with service authorization disabled

```sql
select
  display_name,
  id
from
  azuread_application
where
  not is_authorization_service_enabled;
```
