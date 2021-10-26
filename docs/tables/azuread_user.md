# Table: azuread_user

An Azure AD user account, which helps employees sign in and access resources.

## Examples

### Basic info

```sql
select
  display_name,
  user_principal_name,
  id,
  given_name,
  mail
from
  azuread_user;
```

### List guest users

```sql
select
  display_name,
  user_principal_name,
  id,
  mail
from
  azuread_user
where
  user_type = 'Guest';
```

### List disabled users

```sql
select
  display_name,
  user_principal_name,
  id,
  mail
from
  azuread_user
where
  not account_enabled;
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

### List users with information of groups they are attached

```sql
select
  grp.display_name as group_name,
  grp.id as group_id,
  u.display_name as username,
  u.user_principal_name as user_principal_name,
  u.id as user_id
from
  azuread_group as grp,
  jsonb_array_elements_text(member_ids) as m_id,
  azuread_user as u
where
  u.id = m_id
order by
  group_id,
  username;
```
