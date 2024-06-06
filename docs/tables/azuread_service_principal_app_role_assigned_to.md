---
title: "Steampipe Table: azuread_service_principal_app_role_assigned_to - Query Application Role Assignments granted for Azure Active Directory Service Principals using SQL"
description: "Allows users to query Application Role Assignments granted for an Azure Active Directory Service Principal, providing comprehensive details about each app role assignment including its principal id, name, type, and more."
---

# Table: azuread_service_principal_app_role_assigned_to - Query Application Role Assignments granted for Azure Active Directory Service Principal using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. It combines core directory services, application access management, and identity protection into a single solution. Azure AD also offers a rich, standards-based platform that enables developers to deliver access control to their applications, based on centralized policy and rules.

## Table Usage Guide

The `azuread_service_principal_app_role_assigned_to` table provides insights application roles assigned for service principals within Microsoft's Azure Active Directory. As an IT administrator, you can explore app role assignment-specific details through this table, including the application ID, display name, role, and more. Utilize it to uncover granted app permissions, aiding in the management and security of your organization's resources.

## Examples

### Basic info
Explore which Application Role Assignments are granted for an Azure Active Directory service principal. This can be useful for understanding which principles can access the application represented by the service principal.

```sql+postgres
select
  resource_id,
  resource_display_name,
  app_role_id,
  principal_id,
  principal_type,
  principal_display_name,
  created_date_time,
  deleted_date_time
from
  azuread_service_principal_app_role_assigned_to
where
  service_principal_id = '<service_principal_id>';
```

```sql+sqlite
select
  resource_id,
  resource_display_name,
  app_role_id,
  principal_id,
  principal_type,
  principal_display_name,
  created_date_time,
  deleted_date_time
from
  azuread_service_principal_app_role_assigned_to
where
  service_principal_id = '<service_principal_id>';
```

### List all application role assignments granted for service principals
Explore which principals in your Azure Active Directory have Application Role Assignments for an Azure Active Directory service principal. This information is useful for auditing purposes and ensuring adherence to security protocols.

```sql+postgres
select
  azuread_service_principal_app_role_assigned_to.service_principal_id,
  azuread_service_principal_app_role_assigned_to.resource_id,
  azuread_service_principal_app_role_assigned_to.resource_display_name,
  azuread_service_principal_app_role_assigned_to.app_role_id,
  azuread_service_principal_app_role_assigned_to.created_date_time,
  azuread_service_principal_app_role_assigned_to.deleted_date_time
from
  azuread_service_principal
join azuread_service_principal_app_role_assigned_to
  on azuread_service_principal_app_role_assigned_to.service_principal_id = azuread_service_principal.id;
```

```sql+sqlite
select
  azuread_service_principal_app_role_assigned_to.service_principal_id,
  azuread_service_principal_app_role_assigned_to.resource_id,
  azuread_service_principal_app_role_assigned_to.resource_display_name,
  azuread_service_principal_app_role_assigned_to.app_role_id,
  azuread_service_principal_app_role_assigned_to.created_date_time,
  azuread_service_principal_app_role_assigned_to.deleted_date_time
from
  azuread_service_principal
join azuread_service_principal_app_role_assigned_to
  on azuread_service_principal_app_role_assigned_to.service_principal_id = azuread_service_principal.id;
```
