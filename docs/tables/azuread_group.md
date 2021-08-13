# Table: azuread_group

Azure Active Directory (Azure AD) lets you use groups to manage access to your cloud-based apps, on-premises apps, and your resources.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  description,
  mail
from
  azuread_group;
```

### List groups with public visibility

```sql
select
  display_name,
  id,
  description,
  mail
from
  azuread_group
where
  visibility = 'Public';
```

### List security enabled groups

```sql
select
  display_name,
  id,
  description,
  mail
from
  azuread_group
where
  security_enabled;
```

### List groups that can be assigned to roles

```sql
select
  display_name,
  id,
  description,
  mail
from
  azuread_group
where
  is_assignable_to_role;
```

### Get owner details of an specific group

```sql
select
  gr.display_name as group_name,
  u.display_name as user_name,
  owner_id
from
  azuread_user u,
  azuread_group gr,
  jsonb_array_elements_text(gr.owner_ids) as owner_id
where
  owner_id = u.id and
  gr.display_name = 'turbot'
order by
  user_name;
```
