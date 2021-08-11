# Table: azuread_application

Azure Active Directory (Azure AD) lets you use directory roles to manage access to your cloud-based resources.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  description
from
  azuread_directory_role;
```

### List domains which are verified

```sql
select
  id,
  display_name,
  member_ids
from
  azuread_directory_role;
```
