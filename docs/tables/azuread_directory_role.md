# Table: azuread_directory_role

Azure Active Directory (Azure AD) lets you use directory roles to manage access to your cloud-based resources. Azure AD directory roles are also known as administrator roles.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  description,
  member_ids
from
  azuread_directory_role;
```

### List users with access to directory roles

```sql
select
  u.display_name as username,
  role.display_name as directory_role
from
  azuread_directory_role as role,
  jsonb_array_elements_text(member_ids) as m_id,
  azuread_user as u
where
  u.id = m_id;
```
