---
title: "Steampipe Table: azuread_access_review_schedule_definition - Query Azure Active Directory Access Review Schedule Definitions using SQL"
description: "Allows users to query Azure Active Directory Access Review Schedule Definitions, specifically to get information about recurring access reviews configured in the Azure AD tenant."
---

# Table: azuread_access_review_schedule_definition - Query Azure Active Directory Access Review Schedule Definitions using SQL

Azure Active Directory (Azure AD) Access Reviews help organizations maintain security and compliance by regularly reviewing and certifying user access to resources. Access Review Schedule Definitions define the settings and scope for recurring access reviews, including who should be reviewed, who should perform the review, and how often the reviews should occur.

## Table Usage Guide

The `azuread_access_review_schedule_definition` table provides insights into access review schedule definitions within Azure Active Directory. As a security administrator, explore access review configurations through this table, including the review scope, reviewers, settings, and recurrence patterns. Utilize it to monitor and manage access review processes, ensure compliance with security policies, and track the status of recurring access reviews.

## Examples

### Basic info
Explore which access review schedule definitions are configured in your Azure Active Directory by identifying their display names, status, and creation dates. This can help you monitor and manage your access review processes effectively.

```sql+postgres
select
  display_name,
  id,
  status,
  created_date_time,
  description_for_admins
from
  azuread_access_review_schedule_definition;
```

```sql+sqlite
select
  display_name,
  id,
  status,
  created_date_time,
  description_for_admins
from
  azuread_access_review_schedule_definition;
```

### List active access review schedule definitions
Identify access review schedule definitions that are currently active or in progress. This helps you understand which reviews are currently running and need attention.

```sql+postgres
select
  display_name,
  id,
  status,
  created_date_time,
  last_modified_date_time
from
  azuread_access_review_schedule_definition
where
  status in ('InProgress', 'NotStarted');
```

```sql+sqlite
select
  display_name,
  id,
  status,
  created_date_time,
  last_modified_date_time
from
  azuread_access_review_schedule_definition
where
  status in ('InProgress', 'NotStarted');
```

### Access review schedule definitions with specific scope
Find access review schedule definitions that target specific resources or groups. This is useful for understanding which reviews are scoped to particular organizational units or resource types.

```sql+postgres
select
  display_name,
  id,
  status,
  scope
from
  azuread_access_review_schedule_definition
where
  scope->>'@odata.type' = '#microsoft.graph.principalResourceMembershipsScope';
```

```sql+sqlite
select
  display_name,
  id,
  status,
  json_extract(scope, '$.@odata.type') as scope_type
from
  azuread_access_review_schedule_definition
where
  json_extract(scope, '$.@odata.type') = '#microsoft.graph.principalResourceMembershipsScope';
```

### Access review schedule definitions with role-based scope
Identify access review schedule definitions that are specifically targeting role assignments, such as administrative roles or custom roles. This helps you track reviews of privileged access.

```sql+postgres
select
  display_name,
  id,
  status,
  scope
from
  azuread_access_review_schedule_definition
where
  scope->'resourceScopes'->0->>'query' like '%roleManagement/directory/roleDefinitions%';
```

```sql+sqlite
select
  display_name,
  id,
  status,
  scope
from
  azuread_access_review_schedule_definition
where
  json_extract(scope, '$.resourceScopes[0].query') like '%roleManagement/directory/roleDefinitions%';
```

### Access review schedule definitions with recurrence patterns
Find access review schedule definitions that have specific recurrence patterns, such as monthly or quarterly reviews. This helps you understand the frequency of access reviews.

```sql+postgres
select
  display_name,
  id,
  status,
  settings->'recurrence'->'pattern'->>'type' as recurrence_type,
  settings->'recurrence'->'pattern'->>'interval' as interval,
  settings->'recurrence'->'range'->>'type' as range_type
from
  azuread_access_review_schedule_definition
where
  settings->'recurrence' is not null;
```

```sql+sqlite
select
  display_name,
  id,
  status,
  json_extract(settings, '$.recurrence.pattern.type') as recurrence_type,
  json_extract(settings, '$.recurrence.pattern.interval') as interval,
  json_extract(settings, '$.recurrence.range.type') as range_type
from
  azuread_access_review_schedule_definition
where
  json_extract(settings, '$.recurrence') is not null;
```

### Access review schedule definitions with notification settings
Identify access review schedule definitions that have specific notification settings enabled. This helps you understand how reviewers are notified about their review tasks.

```sql+postgres
select
  display_name,
  id,
  status,
  settings->>'mailNotificationsEnabled' as mail_notifications,
  settings->>'reminderNotificationsEnabled' as reminder_notifications,
  additional_notification_recipients
from
  azuread_access_review_schedule_definition
where
  settings->>'mailNotificationsEnabled' = 'true'
  or settings->>'reminderNotificationsEnabled' = 'true';
```

```sql+sqlite
select
  display_name,
  id,
  status,
  json_extract(settings, '$.mailNotificationsEnabled') as mail_notifications,
  json_extract(settings, '$.reminderNotificationsEnabled') as reminder_notifications,
  additional_notification_recipients
from
  azuread_access_review_schedule_definition
where
  json_extract(settings, '$.mailNotificationsEnabled') = 'true'
  or json_extract(settings, '$.reminderNotificationsEnabled') = 'true';
```
