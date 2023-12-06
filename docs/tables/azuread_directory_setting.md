---
title: "Steampipe Table: azuread_directory_setting - Query Azure Active Directory Directory Settings using SQL"
description: "Allows users to query Directory Settings in Azure Active Directory, specifically the settings that define the behavior and functionality of the directory."
---

# Table: azuread_directory_setting - Query Azure Active Directory Directory Settings using SQL

Azure Active Directory (Azure AD) is Microsoft's multi-tenant, cloud-based directory, and identity management service. Directory Settings in Azure AD are configurable settings that define the behavior and functionality of the directory. These settings include the ability to manage features like self-service password reset, device settings, group settings, and more.

## Table Usage Guide

The `azuread_directory_setting` table provides insights into Directory Settings within Azure Active Directory. As an IT administrator, explore settings-specific details through this table, including the status of various features like self-service password reset, device settings, group settings, and more. Utilize it to uncover information about the configuration and behavior of your Azure AD directory.

## Examples

### Basic info
Explore the basic information in your Azure Active Directory settings to determine the areas where changes or updates may be needed. This can be especially useful in managing user access and permissions within your organization.

```sql+postgres
select
  display_name,
  id,
  value
from
  azuread_directory_setting;
```

```sql+sqlite
select
  display_name,
  id,
  value
from
  azuread_directory_setting;
```

### Check if user admin consent workflow is enabled
Determine if the workflow for user admin consent is activated. This is useful for managing and enforcing user permissions and access controls within your Azure Active Directory.

```sql+postgres
select
  display_name,
  id,
  name,
  value
from
  azuread_directory_setting
where
  display_name = 'Consent Policy Settings'
  and name = 'EnableAdminConsentRequests'
  and value::bool;
```

```sql+sqlite
select
  display_name,
  id,
  name,
  value
from
  azuread_directory_setting
where
  display_name = 'Consent Policy Settings'
  and name = 'EnableAdminConsentRequests'
  and value = '1';
```

### Check if banned password protection is enabled
Determine if your organization's password policy is effectively safeguarding against the use of commonly banned passwords. This query is beneficial in identifying potential vulnerabilities in your password protection settings, ensuring a robust security protocol.

```sql+postgres
select
  display_name,
  id,
  name,
  value
from
  azuread_directory_setting
where
  display_name = 'Password Rule Settings'
  and (
    name = 'EnableBannedPasswordCheckOnPremises'
    and value::bool
  ) and (
    name = 'BannedPasswordCheckOnPremisesMode'
    and value = 'Enforced'
  );
```

```sql+sqlite
select
  display_name,
  id,
  name,
  value
from
  azuread_directory_setting
where
  display_name = 'Password Rule Settings'
  and (
    name = 'EnableBannedPasswordCheckOnPremises'
    and value = 'true'
  ) and (
    name = 'BannedPasswordCheckOnPremisesMode'
    and value = 'Enforced'
  );
```