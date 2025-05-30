---
title: "Steampipe Table: azuread_user - Query Azure AD Users using SQL"
description: "Allows users to query Azure AD Users, specifically user profiles, providing insights into user information and behavior."
---

# Table: azuread_user - Query Azure AD Users using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. It combines core directory services, application access management, and identity protection into a single solution. Azure AD is the backbone of the Office 365 system, and it can sync with on-premise Active Directory.

## Table Usage Guide

The `azuread_user` table provides insights into user profiles within Azure Active Directory. As a system administrator, explore user-specific details through this table, including user identities, user principal names, and associated metadata. Utilize it to uncover information about users, such as their display names, job titles, and the verification of user identities.

## Examples

### Basic info
Explore the basic information of users in your Azure Active Directory. This can be useful for understanding the user composition of your organization, including their display names, principal names, IDs, given names, and emails.

```sql+postgres
select
  display_name,
  user_principal_name,
  id,
  given_name,
  mail
from
  azuread_user;
```

```sql+sqlite
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
Discover the segments that consist of guest users within your Azure Active Directory, allowing you to better manage and monitor these specific user accounts. This is particularly useful in maintaining security protocols and ensuring guest users have appropriate access permissions.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that consist of disabled user accounts within the Azure Active Directory. This can be useful in monitoring and managing user accessibility for security and compliance purposes.

```sql+postgres
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

```sql+sqlite
select
  display_name,
  user_principal_name,
  id,
  mail
from
  azuread_user
where
  account_enabled = 0;
```

### List users with access to directory roles
Discover the segments that have access to directory roles to better manage permissions and security protocols. This is particularly useful for administrators seeking to optimize access control and understand user-role relationships.

```sql+postgres
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

```sql+sqlite
select
  u.display_name as username,
  role.display_name as directory_role
from
  azuread_directory_role as role,
  json_each(role.member_ids) as m_id,
  azuread_user as u
where
  u.id = m_id.value;
```

### List users with information of groups they are attached
Discover the segments that outline the association between users and groups in your Azure Active Directory. This query is useful for assessing user-group relationships, aiding in the management of access and permissions.

```sql+postgres
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

```sql+sqlite
select
  grp.display_name as group_name,
  grp.id as group_id,
  u.display_name as username,
  u.user_principal_name as user_principal_name,
  u.id as user_id
from
  azuread_group as grp,
  json_each(grp.member_ids) as m_id,
  azuread_user as u
where
  u.id = m_id.value
order by
  group_id,
  username;
```

### List users sign_in_activity
Discover the sign_in_activity of users within your Azure Active Directory, allowing you to query for inactive users. This is particularly useful in maintaining security posture.

```sql+postgres
select
  display_name,
  user_principal_name,
  id,
  mail,
  external_user_state,
  sign_in_activity ->> 'LastSignInDateTime' as last_sign_in,
  sign_in_activity ->> 'LastNonInteractiveSignInDateTime' as last_non_interactive_sign_in
from
  azuread_user
```

```sql+sqlite
select
  display_name,
  user_principal_name,
  id,
  mail,
  external_user_state,
  sign_in_activity ->> 'LastSignInDateTime' as last_sign_in,
  sign_in_activity ->> 'LastNonInteractiveSignInDateTime' as last_non_interactive_sign_in
from
  azuread_user
```