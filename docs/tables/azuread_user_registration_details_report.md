---
title: "Steampipe Table: azuread_user_registration_details_report - Query Azure AD Sign-In Reports using SQL"
description: "Allows users to query Azure AD Sign-In Reports, providing detailed information about user sign-in activities, including the location, device, and application used for sign-in."
---

# Table: azuread_user_registration_details_report - Query Azure AD User-Registration-Details Reports using SQL

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service. It helps your employees sign in and access resources in external resources, such as Microsoft Office 365, the Azure portal, and thousands of other SaaS applications. Azure AD User-Registration-Details Reports Represents the state of a user's authentication methods, including which methods are registered and which features the user is registered and capable of

## Table Usage Guide

The `azuread_user_registration_details_report` table provides insights into state of a user's authentication methods within Microsoft's Azure Active Directory. As a security analyst, explore User-Registration-Details through this table, including which methods are registered and which features the user is registered and capable of, such as multifactor authentication, self-service password reset, and passwordless authentication. Utilize it to uncover information about User-Registration-Details.

## Examples

### Basic info
Explore the state of user's authentication methods in your Azure Active Directory. This can be useful for understanding the user authentication methods of your organization, including their display names, principal names, IDs.

```sql+postgres
select
  user_display_name,
  id,
  user_principal_name
from
  azuread_user_registration_details_report;
```

```sql+sqlite
select
  user_display_name,
  id,
  user_principal_name
from
  azuread_user_registration_details_report;
```

### List a user auth methods
Explore the state of a specific user's authentication methods in your Azure Active Directory. This can be useful for understanding the user authentication methods of your organization, including their display names, principal names, IDs.

```sql+postgres
select
  user_display_name,
  id,
  user_principal_name
from
  azuread_user_registration_details_report
where
  user_principal_name = 'abc@myacc.onmicrosoft.com';
```

```sql+sqlite
select
  user_display_name,
  id,
  user_principal_name
from
  azuread_user_registration_details_report
where
  user_principal_name = 'abc@myacc.onmicrosoft.com';
```