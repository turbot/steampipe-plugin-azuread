---
title: "Steampipe Table: azuread_authorization_policy - Query Azure Active Directory Authorization Policies using SQL"
description: "Allows users to query Azure Active Directory Authorization Policies, specifically the policy settings, providing insights into access management and security configurations."
---

# Table: azuread_authorization_policy - Query Azure Active Directory Authorization Policies using SQL

Azure Active Directory (Azure AD) Authorization Policy is a feature of Microsoft Azure that defines how resources in your organization are accessed. It provides a centralized way to manage authorization settings, control access, and enforce security configurations across your Azure resources. Azure AD Authorization Policy enables you to manage and secure access to your resources effectively.

## Table Usage Guide

The `azuread_authorization_policy` table provides insights into authorization policies within Azure Active Directory. As a security administrator, explore policy-specific details through this table, including policy settings, associated metadata, and security configurations. Utilize it to uncover information about policies, such as those with specific access controls, the enforcement of security configurations, and the verification of authorization settings.

## Examples

### Basic info
Analyze the settings to understand the display name, ID, and invite permissions for a given Azure AD authorization policy. This can be useful for auditing and managing access controls within your Azure environment.

```sql+postgres
select
  display_name,
  id,
  allow_invites_from
from
  azuread_authorization_policy;
```

```sql+sqlite
select
  display_name,
  id,
  allow_invites_from
from
  azuread_authorization_policy;
```

### Check if user consent to apps accessing company data on their behalf is not allowed
Determine the areas in which users have not granted permission for apps to access company data on their behalf. This can be useful to maintain data privacy and prevent unauthorized access.

```sql+postgres
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  default_user_role_permissions ->> 'permissionGrantPoliciesAssigned' = '[]';
```

```sql+sqlite
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  json_extract(default_user_role_permissions, '$.permissionGrantPoliciesAssigned') = '[]';
```

### Check if all members are allowed to invite external users to the organization
Determine if your organization's settings permit all members to invite external users. This is useful for assessing the openness of your organization's communication and collaboration policies.

```sql+postgres
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  allow_invites_from = 'everyone';
```

```sql+sqlite
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  allow_invites_from = 'everyone';
```

### Check if email validation is not required to join the tenant
Determine if your organization's settings allow users to join without verifying their email first. This could be a potential security risk, as it may enable unauthorized individuals to gain access to your system.

```sql+postgres
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  not allowed_email_verified_users_to_join_organization;
```

```sql+sqlite
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  not allowed_email_verified_users_to_join_organization;
```

### Check if non-admin users can create tenants
To enhance security practices, it is highly recommended to disable the feature that allows non-admin users to create Azure AD tenants.

```sql+postgres
select
  id,
  display_name,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  (default_user_role_permissions ->> 'allowedToCreateTenants')::boolean = True;
```

```sql+sqlite
select
  id,
  display_name,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  json_extract(default_user_role_permissions, '$.allowedToCreateTenants')::boolean = True;
```
