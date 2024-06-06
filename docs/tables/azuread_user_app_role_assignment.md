---
title: "Steampipe Table: azuread_user_app_role_assignment - Query Application Role Assignments granted to Azure Active Directory User using SQL"
description: "Allows users to query Application Role Assignments granted to an Azure Active Directory User, providing comprehensive details about each app role assignment including its application name, role, and more."
---

# Table: azuread_user_app_role_assignment - Query Application Role Assignments granted to Azure Active Directory User using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. It combines core directory services, application access management, and identity protection into a single solution. Azure AD also offers a rich, standards-based platform that enables developers to deliver access control to their applications, based on centralized policy and rules.

## Table Usage Guide

The `azuread_user_app_role_assignment` table provides insights into application roles assigned to users within Microsoft's Azure Active Directory. As an IT administrator, you can explore app role assignment-specific details through this table, including the application ID, display name, role, and more. Utilize it to uncover user app permissions, aiding in the management and security of your organization's resources.

> This table also includes application role assignments granted to groups that the user is a direct member of.

## Examples

### Basic info
Explore which Application Role Assignments are granted to an Azure Active Directory user. This can be useful for understanding what applications are available to a user.

```sql+postgres
select
  resource_id,
  resource_display_name,
  app_role_id,
  created_date_time,
  deleted_date_time
from
  azuread_user_app_role_assignment
where
  user_id = '<user_id>';
```

```sql+sqlite
select
  resource_id,
  resource_display_name,
  app_role_id,
  created_date_time,
  deleted_date_time
from
  azuread_user_app_role_assignment
where
  user_id = '<user_id>';
```

### List all application role assignments granted to users
Explore which users in your Azure Active Directory have Application Role Assignments. This information is useful for auditing purposes and ensuring adherence to security protocols.

```sql+postgres
select
  azuread_user_app_role_assignment.user_id,
  azuread_user_app_role_assignment.resource_id,
  azuread_user_app_role_assignment.resource_display_name,
  azuread_user_app_role_assignment.app_role_id,
  azuread_user_app_role_assignment.principal_id,
  azuread_user_app_role_assignment.principal_display_name,
  azuread_user_app_role_assignment.created_date_time,
  azuread_user_app_role_assignment.deleted_date_time
from
  azuread_user
join azuread_user_app_role_assignment
  on azuread_user_app_role_assignment.user_id = azuread_user.id;
```

```sql+sqlite
select
  azuread_user_app_role_assignment.group_id,
  azuread_user_app_role_assignment.resource_id,
  azuread_user_app_role_assignment.resource_display_name,
  azuread_user_app_role_assignment.app_role_id,
  azuread_user_app_role_assignment.principal_id,
  azuread_user_app_role_assignment.principal_display_name,
  azuread_user_app_role_assignment.created_date_time,
  azuread_user_app_role_assignment.deleted_date_time
from
  azuread_user
join azuread_user_app_role_assignment
  on azuread_user_app_role_assignment.user_id = azuread_user.id;
```
