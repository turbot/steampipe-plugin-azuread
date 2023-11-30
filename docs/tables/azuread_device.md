---
title: "Steampipe Table: azuread_device - Query Azure Active Directory Devices using SQL"
description: "Allows users to query Azure Active Directory Devices, specifically providing details about the registered devices in an organization."
---

# Table: azuread_device - Query Azure Active Directory Devices using SQL

Azure Active Directory (Azure AD) is a Microsoft cloud-based identity and access management service. It provides an identity platform with enhanced security, access management, scalability, and reliability for connecting users with all the apps they need. Azure AD Devices are the registered devices within an organization that can access resources in the directory.

## Table Usage Guide

The `azuread_device` table provides insights into registered devices within Azure Active Directory. As a system administrator, explore device-specific details through this table, including device id, device type, and associated metadata. Utilize it to uncover information about devices, such as their operating system, physical device id, and the user who registered the device.

## Examples

### Basic info
Explore which devices in your Azure Active Directory are managed and compliant, as well as their group memberships. This is useful for maintaining security standards and managing device access within your organization.

```sql
select
  display_name,
  is_managed,
  is_compliant,
  member_of
from
  azuread_device;
```

### List managed devices
Explore which devices are managed within your Azure Active Directory. This allows you to gain insights into the device profiles, including their operating system versions, for better management and security compliance.

```sql
select
  display_name,
  profile_type,
  id,
  operating_system,
  operating_system_version
from
  azuread_device
where
  is_managed;
```

### List non-compliant devices
Explore which devices in your Azure Active Directory are not compliant with your organization's standards. This can help you identify potential security risks and take necessary corrective actions.

```sql
select
  display_name,
  profile_type,
  id,
  operating_system,
  operating_system_version
from
  azuread_device
where
  not is_compliant;
```