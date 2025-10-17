---
title: "Steampipe Table: azuread_device_registration_policy - Query Azure Active Directory Device Registration Policies using SQL"
description: "Allows users to query Azure Active Directory Device Registration Policies, specifically the policy settings that control device registration, providing insights into device enrollment and management configurations."
---

# Table: azuread_device_registration_policy - Query Azure Active Directory Device Registration Policies using SQL

Azure Active Directory (Azure AD) Device Registration Policy is a feature of Microsoft Azure that manages initial provisioning controls using quota restrictions, additional authentication and authorization checks for device registration. It provides a centralized way to control which users can register or join devices to your organization, the maximum number of devices per user, and whether multi-factor authentication is required for device registration.

## Table Usage Guide

The `azuread_device_registration_policy` table provides insights into device registration policies within Azure Active Directory. As a security administrator or IT manager, explore policy-specific details through this table, including user device quotas, multi-factor authentication requirements, and device registration authorization settings. Utilize it to uncover information about device registration controls, such as who can register devices, whether MFA is enforced, and local administrator settings for Azure AD joined devices.

## Examples

### Basic info
Analyze the settings to understand the display name, ID, user device quota, and multi-factor authentication configuration for your Azure AD device registration policy. This can be useful for auditing and managing device enrollment controls within your Azure environment.

```sql+postgres
select
  display_name,
  id,
  user_device_quota,
  multi_factor_auth_configuration,
  description
from
  azuread_device_registration_policy;
```

```sql+sqlite
select
  display_name,
  id,
  user_device_quota,
  multi_factor_auth_configuration,
  description
from
  azuread_device_registration_policy;
```

### Check if multi-factor authentication is required for device registration
Determine if your organization requires users to complete multi-factor authentication when registering devices. This is important for maintaining strong security controls.

```sql+postgres
select
  display_name,
  id,
  multi_factor_auth_configuration
from
  azuread_device_registration_policy
where
  multi_factor_auth_configuration = 'required';
```

```sql+sqlite
select
  display_name,
  id,
  multi_factor_auth_configuration
from
  azuread_device_registration_policy
where
  multi_factor_auth_configuration = 'required';
```

### List device registration settings for Azure AD Join
Explore the configuration settings for Azure AD Join, including who is allowed to join devices and local administrator settings. This helps understand device enrollment permissions in your organization.

```sql+postgres
select
  display_name,
  id,
  azure_ad_join
from
  azuread_device_registration_policy;
```

```sql+sqlite
select
  display_name,
  id,
  azure_ad_join
from
  azuread_device_registration_policy;
```

### Check if Azure AD device registration is configurable by admins
Determine if administrators can configure Azure AD device registration settings. This is useful for understanding the level of administrative control over device registration.

```sql+postgres
select
  display_name,
  id,
  azure_ad_registration -> 'isAdminConfigurable' as is_admin_configurable
from
  azuread_device_registration_policy;
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(azure_ad_registration, '$.isAdminConfigurable') as is_admin_configurable
from
  azuread_device_registration_policy;
```

### Check if Local Admin Password Solution (LAPS) is enabled
Verify whether the Local Admin Password Solution is enabled for your organization. LAPS helps manage local administrator passwords on domain-joined computers.

```sql+postgres
select
  display_name,
  id,
  local_admin_password -> 'isEnabled' as laps_enabled
from
  azuread_device_registration_policy;
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(local_admin_password, '$.isEnabled') as laps_enabled
from
  azuread_device_registration_policy;
```

### Check user device quota limits
Analyze the maximum number of devices each user is allowed to register. This helps in managing device sprawl and ensuring compliance with organizational policies.

```sql+postgres
select
  display_name,
  id,
  user_device_quota
from
  azuread_device_registration_policy
where
  user_device_quota > 5;
```

```sql+sqlite
select
  display_name,
  id,
  user_device_quota
from
  azuread_device_registration_policy
where
  user_device_quota > 5;
```