---
title: "Steampipe Table: azuread_conditional_access_policy - Query Azure AD Conditional Access Policies using SQL"
description: "Allows users to query Azure AD Conditional Access Policies, providing detailed information about the policies that control access based on conditions."
---

# Table: azuread_conditional_access_policy - Query Azure AD Conditional Access Policies using SQL

Azure AD Conditional Access is a feature in Azure Active Directory that allows administrators to define policies that control access to applications and resources based on conditions. These conditions can include user role, location, device status, and risk level. This feature is crucial for managing security and compliance in organizations.

## Table Usage Guide

The `azuread_conditional_access_policy` table provides insights into Conditional Access Policies within Azure Active Directory. As a security administrator, you can explore policy-specific details through this table, including conditions, grant controls, and associated metadata. Utilize it to uncover information about policies, such as those with specific conditions and controls, helping you to maintain security and compliance within your organization.

## Examples

### Basic info
Analyze the settings to understand the status and creation date of the built-in controls in your Azure Active Directory conditional access policy. This can help you assess the elements within your policy and make necessary adjustments.

```sql+postgres
select
  id,
  display_name,
  state,
  created_date_time,
  built_in_controls
from
  azuread_conditional_access_policy;
```

```sql+sqlite
select
  id,
  display_name,
  state,
  created_date_time,
  built_in_controls
from
  azuread_conditional_access_policy;
```

### List conditional access policies with mfa enabled
Uncover the details of conditional access policies that have multi-factor authentication enabled. This is useful for enhancing security by identifying policies that require an additional layer of verification.

```sql+postgres
select
  id,
  display_name,
  built_in_controls
from
  azuread_conditional_access_policy
where
  built_in_controls ?& array['mfa'];
```

```sql+sqlite
Error: SQLite does not support array operations and '?&' operator.
```