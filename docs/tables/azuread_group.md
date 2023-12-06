---
title: "Steampipe Table: azuread_group - Query Azure Active Directory Groups using SQL"
description: "Allows users to query Azure Active Directory Groups, providing comprehensive details about each group including its ID, display name, security identifier, and more."
---

# Table: azuread_group - Query Azure Active Directory Groups using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. It combines core directory services, application access management, and identity protection into a single solution. Azure AD also offers a rich, standards-based platform that enables developers to deliver access control to their applications, based on centralized policy and rules.

## Table Usage Guide

The `azuread_group` table provides insights into groups within Microsoft's Azure Active Directory. As an IT administrator, you can explore group-specific details through this table, including the group's ID, display name, security identifier, and more. Utilize it to uncover information about groups, such as their membership and associated metadata, aiding in the management and security of your organization's resources.

## Examples

### Basic info
Explore which Azure Active Directory groups are present in your system, along with their associated email addresses. This can be useful for understanding your group structure and managing group communication.

```sql+postgres
select
  display_name,
  id,
  description,
  mail
from
  azuread_group;
```

```sql+sqlite
select
  display_name,
  id,
  description,
  mail
from
  azuread_group;
```

### List groups with public visibility
Explore which user groups within your Azure Active Directory have been set to public visibility. This can help in managing data security and privacy by identifying potential areas of risk.

```sql+postgres
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

```sql+sqlite
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
Explore which groups in your Azure Active Directory have security features enabled. This information is useful for auditing purposes and ensuring adherence to security protocols.

```sql+postgres
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

```sql+sqlite
select
  display_name,
  id,
  description,
  mail
from
  azuread_group
where
  security_enabled = 1;
```

### List groups that can be assigned to roles
Explore which groups within your Azure Active Directory can be assigned to roles. This enables better management of access permissions, ensuring appropriate role assignments within your organization.

```sql+postgres
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

```sql+sqlite
select
  display_name,
  id,
  description,
  mail
from
  azuread_group
where
  is_assignable_to_role = 1;
```

### Get owner details of an specific group
Discover the segments that identify the owner of a specific group in the Azure Active Directory. This could be useful in scenarios where you need to understand access control and ownership structures within your organization.

```sql+postgres
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

```sql+sqlite
select
  gr.display_name as group_name,
  u.display_name as user_name,
  owner_id.value as owner_id
from
  azuread_user u,
  azuread_group gr,
  json_each(gr.owner_ids) as owner_id
where
  owner_id.value = u.id and
  gr.display_name = 'turbot'
order by
  user_name;
```