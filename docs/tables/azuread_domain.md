# Table: azuread_application

Azure Active Directory (Azure AD) lets you use domains to manage access to your cloud-based resources.

## Examples

### Basic info

```sql
select
  id,
  is_admin_managed,
  is_verified,
  supportedServices
from
  azuread_domain;
```

### List domains which are verified

```sql
select
  id
from
  azuread_domain
where
  is_verified;
```
