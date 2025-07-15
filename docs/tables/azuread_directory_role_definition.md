---
title: "Steampipe Table: azuread_directory_role_definition - Query Azure Active Directory Role Definitions using SQL"
description: "Allows users to query Role Definitions in Azure Active Directory, providing insights into available roles and their permissions."
---

# Table: azuread_directory_role_definition - Query Azure Active Directory Role Definitions using SQL

Azure Active Directory (Azure AD) Role Definitions represent the available roles within Azure AD that can be assigned to principals (users, groups, or service principals). These definitions describe the permissions, capabilities, and constraints associated with each role, providing the blueprint for role-based access control within the directory.

## Table Usage Guide

The `azuread_directory_role_definition` table provides insights into Role Definitions within Azure Active Directory. As a DevOps engineer or IT professional, you can explore role-specific details through this table, including permissions, display names, descriptions, built-in status, and inheritance relationships. Utilize it to understand available roles, audit role permissions, and design appropriate access control strategies.

## Examples

### Basic info
Explore the available role definitions to understand what roles are available in your Azure AD tenant.

```sql+postgres
select
  id,
  display_name,
  description,
  template_id,
  version,
  is_built_in,
  is_enabled
from
  azuread_directory_role_definition;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  template_id,
  version,
  is_built_in,
  is_enabled
from
  azuread_directory_role_definition;
```

### List built-in role definitions
Find all built-in role definitions provided by Microsoft.

```sql+postgres
select
  id,
  display_name,
  description,
  template_id,
  is_enabled
from
  azuread_directory_role_definition
where
  is_built_in = true
order by
  display_name;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  template_id,
  is_enabled
from
  azuread_directory_role_definition
where
  is_built_in = 1
order by
  display_name;
```

### List custom role definitions
Find all custom role definitions created in your organization.

```sql+postgres
select
  id,
  display_name,
  description,
  template_id,
  version,
  is_enabled
from
  azuread_directory_role_definition
where
  is_built_in = false
order by
  display_name;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  template_id,
  version,
  is_enabled
from
  azuread_directory_role_definition
where
  is_built_in = 0
order by
  display_name;
```

### Find disabled role definitions
Identify role definitions that are currently disabled and cannot be assigned.

```sql+postgres
select
  id,
  display_name,
  description,
  is_built_in,
  template_id
from
  azuread_directory_role_definition
where
  is_enabled = false;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  is_built_in,
  template_id
from
  azuread_directory_role_definition
where
  is_enabled = 0;
```

### Role definitions with specific permissions

Find role definitions that contain specific permissions or actions.

```sql+postgres
select
  id,
  display_name,
  description,
  role_permissions
from
  azuread_directory_role_definition
where
  role_permissions::text like '%microsoft.directory/users%';
```

```sql+sqlite
select
  id,
  display_name,
  description,
  role_permissions
from
  azuread_directory_role_definition
where
  json_extract(role_permissions, '$') like '%microsoft.directory/users%';
```

### Role definitions by resource scope
Analyze role definitions grouped by their resource scopes.

```sql+postgres
select
  display_name,
  description,
  resource_scopes,
  is_built_in
from
  azuread_directory_role_definition
where
  resource_scopes is not null
order by
  display_name;
```

```sql+sqlite
select
  display_name,
  description,
  resource_scopes,
  is_built_in
from
  azuread_directory_role_definition
where
  resource_scopes is not null
order by
  display_name;
```

### Find role definitions with inheritance
Identify role definitions that inherit permissions from other roles.

```sql+postgres
select
  id,
  display_name,
  description,
  inherits_permissions_from
from
  azuread_directory_role_definition
where
  inherits_permissions_from is not null;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  inherits_permissions_from
from
  azuread_directory_role_definition
where
  inherits_permissions_from is not null;
```

### Role definitions by template
Group role definitions by their template to understand role families.

```sql+postgres
select
  template_id,
  count(*) as role_count,
  array_agg(display_name) as role_names
from
  azuread_directory_role_definition
where
  template_id is not null
group by
  template_id
order by
  role_count desc;
```

```sql+sqlite
select
  template_id,
  count(*) as role_count,
  group_concat(display_name) as role_names
from
  azuread_directory_role_definition
where
  template_id is not null
group by
  template_id
order by
  role_count desc;
```

### Find highly privileged roles
Identify role definitions that have broad permissions, which may indicate highly privileged roles.

```sql+postgres
select
  id,
  display_name,
  description,
  is_built_in,
  json_array_length(role_permissions) as permission_count
from
  azuread_directory_role_definition
where
  role_permissions is not null
order by
  json_array_length(role_permissions) desc
limit 10;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  is_built_in,
  json_array_length(role_permissions) as permission_count
from
  azuread_directory_role_definition
where
  role_permissions is not null
order by
  json_array_length(role_permissions) desc
limit 10;
```