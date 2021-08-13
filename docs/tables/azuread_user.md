# Table: azure_ad_user

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
  memberof ->> 'displayName' as directory_role,
  memberof ->> 'id' as directory_role_id,
  display_name as user_display_name,
  user_principal_name,
  id as user_id
from
  azuread_user,
  jsonb_array_elements(member_of) as memberof
where
  split_part(memberof ->> '@odata.type', '.', 3) = 'directoryRole'
order by
  directory_role_id,
  user_display_name;
```

### List users with information of groups they are attached

```sql
select
  memberof ->> 'displayName' as group_name,
  memberof ->> 'id' as group_id,
  display_name as user_display_name,
  user_principal_name,
  id as user_id
from
  azuread_user,
  jsonb_array_elements(member_of) as memberof
where
  split_part(memberof ->> '@odata.type', '.', 3) = 'group'
order by
  group_id,
  user_display_name;
```
