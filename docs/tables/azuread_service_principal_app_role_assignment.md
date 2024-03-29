---
title: "Steampipe Table: azuread_service_principal_app_role_assignment - Query Application Role Assignments granted to Azure Active Directory Service Principals using SQL"
description: "Allows users to query Application Role Assignments granted to an Azure Active Directory Service Principal, providing comprehensive details about each app role assignment including its application name, role, and more."
---

# Table: azuread_service_principal_app_role_assignment - Query Application Role Assignments granted to Azure Active Directory Service Principal using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. It combines core directory services, application access management, and identity protection into a single solution. Azure AD also offers a rich, standards-based platform that enables developers to deliver access control to their applications, based on centralized policy and rules.

## Table Usage Guide

The `azuread_service_principal_app_role_assignment` table provides insights into application roles assigned to service principals within Microsoft's Azure Active Directory. As an IT administrator, you can explore app role assignment-specific details through this table, including the application ID, display name, role, and more. Utilize it to uncover service principal app permissions, aiding in the management and security of your organization's resources.

## Examples

### Basic info
Explore which Application Role Assignments are granted to an Azure Active Directory service principal. This can be useful for understanding what applications are available to service principals.

```sql+postgres
select
  resource_id,
  resource_display_name,
  app_role_id,
  created_date_time,
  deleted_date_time
from
  azuread_service_principal_app_role_assignment
where
  service_principal_id = '<service_principal_id>';
```

```sql+sqlite
select
  resource_id,
  resource_display_name,
  app_role_id,
  created_date_time,
  deleted_date_time
from
  azuread_service_principal_app_role_assignment
where
  service_principal_id = '<service_principal_id>';
```

### List all application role assignments granted to service principals
Explore which service principals in your Azure Active Directory have Application Role Assignments. This information is useful for auditing purposes and ensuring adherence to security protocols.

```sql+postgres
select
  azuread_service_principal_app_role_assignment.service_principal_id,
  azuread_service_principal_app_role_assignment.resource_id,
  azuread_service_principal_app_role_assignment.resource_display_name,
  azuread_service_principal_app_role_assignment.app_role_id,
  azuread_service_principal_app_role_assignment.created_date_time,
  azuread_service_principal_app_role_assignment.deleted_date_time
from
  azuread_service_principal
join azuread_service_principal_app_role_assignment
  on azuread_service_principal_app_role_assignment.service_principal_id = azuread_service_principal.id;
```

```sql+sqlite
select
  azuread_service_principal_app_role_assignment.service_principal_id,
  azuread_service_principal_app_role_assignment.resource_id,
  azuread_service_principal_app_role_assignment.resource_display_name,
  azuread_service_principal_app_role_assignment.app_role_id,
  azuread_service_principal_app_role_assignment.created_date_time,
  azuread_service_principal_app_role_assignment.deleted_date_time
from
  azuread_service_principal
join azuread_service_principal_app_role_assignment
  on azuread_service_principal_app_role_assignment.service_principal_id = azuread_service_principal.id;
```
