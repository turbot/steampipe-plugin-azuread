---
title: "Steampipe Table: azuread_directory_role - Query Azure Active Directory Directory Roles using SQL"
description: "Allows users to query Directory Roles in Azure Active Directory, providing insights into role-specific details, permissions, and associated metadata."
---

# Table: azuread_directory_role - Query Azure Active Directory Directory Roles using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. It combines core directory services, application access management, and identity protection into a single solution. Directory Roles in Azure AD provide access to various features and capabilities in the Azure portal and Azure AD administrative features.

## Table Usage Guide

The `azuread_directory_role` table provides insights into Directory Roles within Azure Active Directory. As a DevOps engineer or IT professional, you can explore role-specific details through this table, including permissions, and associated metadata. Utilize it to uncover information about roles, such as their assigned permissions, the users associated with each role, and the verification of role-specific settings.

## Examples

### Basic info
Explore the roles within your Azure Active Directory to understand their functions and who has been assigned to them. This can be useful for auditing purposes or to ensure the correct permissions have been granted.

```sql+postgres
select
  id,
  display_name,
  description,
  member_ids
from
  azuread_directory_role;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  member_ids
from
  azuread_directory_role;
```

### List users with access to directory roles
Explore which users have access to specific directory roles. This is useful for managing and reviewing user permissions in a system.

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