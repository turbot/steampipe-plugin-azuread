---
title: "Steampipe Table: azuread_directory_role_template - Query Azure Active Directory Directory Role Templates using SQL"
description: "Allows users to query Directory Role Templates in Azure Active Directory, providing insights into role template definitions and their properties."
---

# Table: azuread_directory_role_template - Query Azure Active Directory Directory Role Templates using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. Directory Role Templates in Azure AD define the property values that are used to create directory roles. Each template serves as a blueprint for creating directory roles with predefined permissions and properties.

## Table Usage Guide

The `azuread_directory_role_template` table provides insights into Directory Role Templates within Azure Active Directory. As a DevOps engineer or IT professional, you can explore template-specific details through this table, including descriptions, display names, and template identifiers. Utilize it to understand available role templates, their purposes, and how they can be used to create directory roles with specific permissions.

## Examples

### Basic info
Explore the available directory role templates in Azure Active Directory to understand their purpose and properties.

```sql+postgres
select
  id,
  display_name,
  description
from
  azuread_directory_role_template;
```

```sql+sqlite
select
  id,
  display_name,
  description
from
  azuread_directory_role_template;
```

### List built-in administrative role templates
Identify administrative role templates that are available for creating directory roles with elevated permissions.

```sql+postgres
select
  id,
  display_name,
  description
from
  azuread_directory_role_template
where
  display_name like '%Admin%'
  or display_name like '%Administrator%';
```

```sql+sqlite
select
  id,
  display_name,
  description
from
  azuread_directory_role_template
where
  display_name like '%Admin%'
  or display_name like '%Administrator%';
```

### Find specific role template by display name
Retrieve details of a specific role template to understand its purpose and capabilities.

```sql+postgres
select
  id,
  display_name,
  description
from
  azuread_directory_role_template
where
  display_name = 'Global Administrator';
```

```sql+sqlite
select
  id,
  display_name,
  description
from
  azuread_directory_role_template
where
  display_name = 'Global Administrator';
```

### Compare role templates with active directory roles
Compare available role templates with currently active directory roles to identify unused templates.

```sql+postgres
select
  t.display_name as template_name,
  r.display_name as active_role_name,
  case
    when r.id is null then 'Template not activated'
    else 'Template activated'
  end as status
from
  azuread_directory_role_template as t
  left join azuread_directory_role as r on t.id = r.role_template_id
order by
  t.display_name;
```

```sql+sqlite
select
  t.display_name as template_name,
  r.display_name as active_role_name,
  case
    when r.id is null then 'Template not activated'
    else 'Template activated'
  end as status
from
  azuread_directory_role_template as t
  left join azuread_directory_role as r on t.id = r.role_template_id
order by
  t.display_name;
```
