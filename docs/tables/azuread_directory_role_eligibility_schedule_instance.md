---
title: "Steampipe Table: azuread_directory_role_eligibility_schedule_instance - Query Azure Active Directory Role Eligibility Schedule Instances using SQL"
description: "Allows users to query Role Eligibility Schedule Instances in Azure Active Directory, providing insights into role eligibility schedules and their properties."
---

# Table: azuread_directory_role_eligibility_schedule_instance - Query Azure Active Directory Role Eligibility Schedule Instances using SQL

Azure Active Directory (Azure AD) Privileged Identity Management (PIM) allows organizations to manage, control, and monitor access to important resources. Role Eligibility Schedule Instances represent the schedule instances for role eligibility operations, showing when principals are eligible to activate roles within specified time windows.

## Table Usage Guide

The `azuread_directory_role_eligibility_schedule_instance` table provides insights into Role Eligibility Schedule Instances within Azure Active Directory PIM. As a DevOps engineer or IT professional, you can explore eligibility-specific details through this table, including principal assignments, role definitions, time windows, and membership types. Utilize it to monitor privileged access eligibility, audit role eligibility assignments, and ensure compliance with access governance policies.

## Examples

### Basic info
Explore the role eligibility schedule instances to understand which principals are eligible for which roles and when.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  start_date_time,
  end_date_time,
  member_type,
  role_eligibility_schedule_id
from
  azuread_directory_role_eligibility_schedule_instance;
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  start_date_time,
  end_date_time,
  member_type,
  role_eligibility_schedule_id
from
  azuread_directory_role_eligibility_schedule_instance;
```

### List active role eligibility instances
Find role eligibility instances that are currently active (within their validity period).

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  start_date_time,
  end_date_time,
  member_type
from
  azuread_directory_role_eligibility_schedule_instance
where
  start_date_time <= now()
  and (end_date_time is null or end_date_time > now());
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  directory_scope_id,
  app_scope_id,
  start_date_time,
  end_date_time,
  member_type
from
  azuread_directory_role_eligibility_schedule_instance
where
  start_date_time <= datetime('now')
  and (end_date_time is null or end_date_time > datetime('now'));
```

### Get role eligibility details with complex objects
Query role eligibility instances with detailed information about the principal, role definition, and scopes.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  principal ->> 'displayName' as principal_name,
  principal ->> 'userPrincipalName' as principal_upn,
  role_definition ->> 'displayName' as role_name,
  role_definition ->> 'description' as role_description,
  directory_scope ->> 'displayName' as directory_scope_name,
  app_scope ->> 'displayName' as app_scope_name,
  member_type,
  start_date_time,
  end_date_time
from
  azuread_directory_role_eligibility_schedule_instance
where
  start_date_time <= now()
  and (end_date_time is null or end_date_time > now());
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  json_extract(principal, '$.displayName') as principal_name,
  json_extract(principal, '$.userPrincipalName') as principal_upn,
  json_extract(role_definition, '$.displayName') as role_name,
  json_extract(role_definition, '$.description') as role_description,
  json_extract(directory_scope, '$.displayName') as directory_scope_name,
  json_extract(app_scope, '$.displayName') as app_scope_name,
  member_type,
  start_date_time,
  end_date_time
from
  azuread_directory_role_eligibility_schedule_instance
where
  start_date_time <= datetime('now')
  and (end_date_time is null or end_date_time > datetime('now'));
```

### Filter by member type
Find role eligibility instances by how the eligibility is inherited.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  member_type,
  start_date_time,
  end_date_time
from
  azuread_directory_role_eligibility_schedule_instance
where
  member_type = 'Direct';
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  member_type,
  start_date_time,
  end_date_time
from
  azuread_directory_role_eligibility_schedule_instance
where
  member_type = 'Direct';
```

### Get eligibility instances for a specific principal
Find all role eligibility instances for a specific user or service principal.

```sql+postgres
select
  id,
  principal_id,
  role_definition_id,
  principal ->> 'displayName' as principal_name,
  role_definition ->> 'displayName' as role_name,
  member_type,
  start_date_time,
  end_date_time
from
  azuread_directory_role_eligibility_schedule_instance
where
  principal_id = 'specific-principal-id';
```

```sql+sqlite
select
  id,
  principal_id,
  role_definition_id,
  json_extract(principal, '$.displayName') as principal_name,
  json_extract(role_definition, '$.displayName') as role_name,
  member_type,
  start_date_time,
  end_date_time
from
  azuread_directory_role_eligibility_schedule_instance
where
  principal_id = 'specific-principal-id';
```
