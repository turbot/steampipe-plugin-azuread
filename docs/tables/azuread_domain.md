# Table: azuread_domain

Azure Active Directory (Azure AD) lets you use domains to manage access to your cloud-based resources.

## Examples

### Basic info

```sql
select
  id,
  is_admin_managed,
  is_verified,
  supported_services
from
  azuread_domain;
```

### List verified domains

```sql
select
  id
from
  azuread_domain
where
  is_verified;
```
