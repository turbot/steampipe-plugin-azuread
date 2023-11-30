---
title: "Steampipe Table: azuread_security_defaults_policy - Query Azure AD Security Defaults Policies using SQL"
description: "Allows users to query Security Defaults Policies in Azure AD, specifically providing information about the default security settings in Azure Active Directory (Azure AD)."
---

# Table: azuread_security_defaults_policy - Query Azure AD Security Defaults Policies using SQL

Security Defaults in Azure AD is a set of basic identity security mechanisms recommended by Microsoft. It provides a level of protection to organizations that may not have dedicated security and identity professionals on their IT staff. Security Defaults include requiring all users to register for Azure AD Multi-Factor Authentication, requiring administrators to perform multi-factor authentication, blocking legacy authentication protocols, and more.

## Table Usage Guide

The `azuread_security_defaults_policy` table provides insights into the Security Defaults Policies within Azure Active Directory. As a security analyst, explore policy-specific details through this table, including the status of the policy and if it is enabled or not. Utilize it to monitor and manage your organization's basic identity security settings, ensuring that all users and administrators are adhering to recommended security practices.

## Examples

### Basic info
Explore which security policies are active within your Azure Active Directory. This can help in assessing your current security settings and identifying areas that might need reinforcement.

```sql
select
  display_name,
  id,
  is_enabled
from
  azuread_security_defaults_policy;
```