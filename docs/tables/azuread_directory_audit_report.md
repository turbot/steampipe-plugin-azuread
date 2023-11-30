---
title: "Steampipe Table: azuread_directory_audit_report - Query Azure Active Directory Audit Reports using SQL"
description: "Allows users to query Azure Active Directory Audit Reports, providing insights into audit logs and activity data."
---

# Table: azuread_directory_audit_report - Query Azure Active Directory Audit Reports using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. It combines core directory services, application access management, and identity protection into a single solution. Azure AD also offers a rich, standards-based platform that enables developers to deliver access control to their applications, based on centralized policy and rules.

## Table Usage Guide

The `azuread_directory_audit_report` table provides insights into the audit reports within Azure Active Directory. As a security analyst, explore audit-specific details through this table, including activity data, changes made, and the entities affected. Utilize it to uncover information about user activities, such as login attempts, password changes, and the creation of new entities, aiding in the detection of unusual or potentially harmful behavior.

## Examples

### Basic info
Analyze the settings to understand the activities within your Azure Active Directory. This query allows you to identify who initiated specific operations, what those operations were, and when they occurred, helping you maintain security and compliance.

```sql
select
  activity_display_name,
  activity_date_time,
  category,
  operation_type,
  initiated_by -> 'user' ->> 'userPrincipalName' as initiated_user,
  result
from
  azuread_directory_audit_report;
```

### List all activities related to policy
Determine the areas in which policy-related activities have occurred within your Azure Active Directory. This can help you gain insights into the operations and users involved, as well as the results of these activities, enhancing your understanding and management of policy-related actions.

```sql
select
  activity_display_name,
  activity_date_time,
  category,
  operation_type,
  initiated_by -> 'user' ->> 'userPrincipalName' as initiated_user,
  result
from
  azuread_directory_audit_report
where
  category = 'Policy';
```

### List all activities initiated by a specific user
Explore the types of activities initiated by a particular user within an organization. This is useful for auditing purposes, allowing you to monitor user actions and identify any unusual or suspicious activities.

```sql
select
  activity_display_name,
  activity_date_time,
  category,
  operation_type,
  initiated_by -> 'user' ->> 'userPrincipalName' as initiated_user,
  result
from
  azuread_directory_audit_report
where
  filter = 'initiatedBy/user/userPrincipalName eq ''test@org.onmicrosoft.com''';
```

### List activities related to user creation in last 7 days
Explore recent user creation activities within the past week. This allows you to identify who initiated the creation and the username of the new user, providing insights into your user management activities.

```sql
select
  activity_date_time,
  category,
  operation_type,
  initiated_by -> 'user' ->> 'userPrincipalName' as initiated_user,
  t ->> 'userPrincipalName' as new_user_name
from
  azuread_directory_audit_report,
  jsonb_array_elements(target_resources) as t
where
  activity_display_name = 'Add user'
  and activity_date_time >= (current_date - interval '7 days')
order by activity_date_time;
```

### List users who have reset their passwords in last 7 days
This query lets you track recent password resets in your organization, helping you monitor account security. It's useful for identifying any unusual activity, such as an unexpected surge in password resets, that may indicate a security issue.

```sql
select
  activity_date_time,
  category,
  operation_type,
  initiated_by -> 'user' ->> 'userPrincipalName' as initiated_user,
  t ->> 'userPrincipalName' as target_user
from
  azuread_directory_audit_report,
  jsonb_array_elements(target_resources) as t
where
  t ->> 'displayName' = 'Microsoft password reset service'
  and activity_date_time >= (current_date - interval '7 days')
order by activity_date_time;
```