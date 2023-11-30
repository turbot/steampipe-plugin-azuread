---
title: "Steampipe Table: azuread_application - Query Azure Active Directory Applications using SQL"
description: "Allows users to query Azure Active Directory Applications, specifically to get information about the applications registered in the Azure AD tenant."
---

# Table: azuread_application - Query Azure Active Directory Applications using SQL

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service. It helps your employees sign in and access resources in external resources, such as Microsoft Office 365, the Azure portal, and thousands of other SaaS applications. Azure AD Applications are the entities that are used to manage and secure app resources within your Azure AD tenant.

## Table Usage Guide

The `azuread_application` table provides insights into applications registered within Azure Active Directory. As a security administrator, explore application-specific details through this table, including the application's ID, display name, and whether it's available to other tenants. Utilize it to uncover information about applications, such as those that are multi-tenanted, the types of permissions they have, and their associated service principals.

## Examples

### Basic info
Explore which applications are registered in your Azure Active Directory by identifying their display names and associated IDs. This can help you manage and monitor your applications effectively.

```sql
select
  display_name,
  id,
  app_id
from
  azuread_application;
```

### List owners of an application
This query helps to identify the owners of a specific application within a system, which is useful for understanding who has control over and responsibility for that application. It's particularly beneficial in scenarios where there is a need to audit access rights or investigate potential security issues.

```sql
select
  app.display_name as application_name,
  app.id as application_id,
  o as owner_id,
  u.display_name as owner_display_name
from
  azuread_application as app,
  jsonb_array_elements_text(owner_ids) as o
  left join azuread_user as u on u.id = o
where
  app.id = 'a6656898-3879-4d35-8a58-b34237095a70';
```