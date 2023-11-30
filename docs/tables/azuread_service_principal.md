---
title: "Steampipe Table: azuread_service_principal - Query Azure Active Directory Service Principals using SQL"
description: "Allows users to query Service Principals in Azure Active Directory, specifically the details about the service principals including their roles, permissions, and other related information."
---

# Table: azuread_service_principal - Query Azure Active Directory Service Principals using SQL

Azure Active Directory (Azure AD) is Microsoft’s cloud-based identity and access management service. It helps your employees sign in and access resources in external resources, such as Microsoft Office 365, the Azure portal, and thousands of other SaaS applications. Service Principals in Azure AD are the security identities that user-created apps, services, and automation tools use to access specific Azure resources.

## Table Usage Guide

The `azuread_service_principal` table provides insights into Service Principals within Azure Active Directory. As a security analyst or a DevOps engineer, explore details about the service principals through this table, including their roles, permissions, and other related information. Utilize it to uncover details about the service principals, such as their associated applications, permissions, and the roles they play in your Azure environment.

## Examples

### Basic info
Explore the relationship between display names and application names in your Azure Active Directory. This can be useful in understanding how your applications are connected and organized within the directory.

```sql
select
  display_name,
  id,
  app_display_name
from
  azuread_service_principal;
```

### List disabled service principals
Uncover the details of disabled service principals within your Azure Active Directory. This is useful in ensuring that disabled accounts are not posing a security risk or cluttering your system.

```sql
select
  display_name,
  id
from
  azuread_service_principal
where
  not account_enabled;
```

### List service principals related to applications
Explore which service principals are directly related to applications in Azure Active Directory. This can be useful to determine which applications have active accounts, aiding in both security and account management.

```sql
select
  id,
  app_display_name,
  account_enabled
from
  azuread_service_principal
where
  service_principal_type = 'Application'
  and tenant_id = app_owner_organization_id;
```