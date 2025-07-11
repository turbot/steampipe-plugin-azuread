---
title: "Steampipe Table: azuread_directory_role_assignment - Query Azure Active Directory Role Assignments using SQL"
description: "Allows users to query Role Assignments in Azure Active Directory, providing insights into role assignments and their properties."
---

# Table: azuread_directory_role_assignment - Query Azure Active Directory Role Assignments using SQL

Azure Active Directory (Azure AD) Role Assignments represent the assignment of roles to principals (users, groups, or service principals) within Azure AD. These assignments define which principals have specific permissions and access rights within the directory, enabling fine-grained access control and governance.

## Table Usage Guide

The `azuread_directory_role_assignment` table provides insights into Role Assignments within Azure Active Directory. As a DevOps engineer or IT professional, you can explore assignment-specific details through this table, including principal assignments, role definitions, scopes, conditions, and creation details. Utilize it to monitor access permissions, audit role assignments, and ensure compliance with access governance policies.

## Examples

### Basic info
Explore the role assignments to understand which principals have been assigned which roles and when.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  condition,
  principal,
  role_definition
from
  azuread_directory_role_assignment;
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  condition,
  principal,
  role_definition
from
  azuread_directory_role_assignment;
```

### List role assignments with conditions
Find role assignments that have specific conditions applied to them.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  condition
from
  azuread_directory_role_assignment
where
  condition is not null;
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  condition
from
  azuread_directory_role_assignment
where
  condition is not null;
```

### Role assignments by scope type
Analyze role assignments grouped by their scope type to understand assignment patterns.

```sql+postgres
select
  case
    when directory_scope_id is not null then 'Directory Scope'
    when app_scope_id is not null then 'App Scope'
    else 'Unknown'
  end as scope_type,
  count(*) as assignment_count
from
  azuread_directory_role_assignment
group by
  scope_type
order by
  assignment_count desc;
```

```sql+sqlite
select
  case
    when directory_scope_id is not null then 'Directory Scope'
    when app_scope_id is not null then 'App Scope'
    else 'Unknown'
  end as scope_type,
  count(*) as assignment_count
from
  azuread_directory_role_assignment
group by
  scope_type
order by
  assignment_count desc;
```

### Find principals with multiple role assignments
Identify principals who have multiple role assignments, which may indicate high-privilege accounts.

```sql+postgres
select
  principal_id,
  count(*) as role_assignment_count
from
  azuread_directory_role_assignment
group by
  principal_id
having
  count(*) > 1
order by
  role_assignment_count desc;
```

```sql+sqlite
select
  principal_id,
  count(*) as role_assignment_count
from
  azuread_directory_role_assignment
group by
  principal_id
having
  count(*) > 1
order by
  role_assignment_count desc;
```

### Role assignments with application scopes
Find role assignments that are scoped to specific applications rather than directory-wide.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  app_scope_id,
  condition,
  app_scope
from
  azuread_directory_role_assignment
where
  app_scope_id is not null;
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  app_scope_id,
  condition,
  app_scope
from
  azuread_directory_role_assignment
where
  app_scope_id is not null;
```

### Role assignments with detailed role and principal information
Get detailed information about the role definitions and principals assigned to roles.

```sql+postgres
select
  id,
  principal_id,
  role_definition ->> 'display_name' as role_name,
  role_definition ->> 'description' as role_description,
  principal ->> '@odata.type' as principal_type,
  directory_scope_id,
  condition
from
  azuread_directory_role_assignment
where
  role_definition is not null
  and principal is not null;
```

```sql+sqlite
select
  id,
  principal_id,
  json_extract(role_definition, '$.display_name') as role_name,
  json_extract(role_definition, '$.description') as role_description,
  json_extract(principal, '$."@odata.type"') as principal_type,
  directory_scope_id,
  condition
from
  azuread_directory_role_assignment
where
  role_definition is not null
  and principal is not null;
```

### Role assignments with specific directory scope
Find role assignments that are scoped to a specific directory scope.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  condition
from
  azuread_directory_role_assignment
where
  directory_scope_id = '/';
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  condition
from
  azuread_directory_role_assignment
where
  directory_scope_id = '/';
```
