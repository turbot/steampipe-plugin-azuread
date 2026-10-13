---
title: "Steampipe Table: azuread_intune_managed_device - Query Microsoft Intune managed devices using SQL"
description: "Allows users to query the devices enrolled in Microsoft Intune, including their compliance state, encryption, ownership and last sync time."
---

# Table: azuread_intune_managed_device - Query Microsoft Intune managed devices using SQL

A managed device is a laptop, desktop or mobile device enrolled in Microsoft Intune. Intune records each device's hardware identity (serial number, manufacturer, model), its primary user, its operating system, whether its storage is encrypted, the compliance state Intune's compliance policies assign it, and when it last checked in.

## Table Usage Guide

The `azuread_intune_managed_device` table provides an inventory of the devices Intune manages. Use it to find non-compliant or unencrypted devices, devices that have stopped syncing, and the user each device belongs to.

Reading this table requires the `DeviceManagementManagedDevices.Read.All` Microsoft Graph application permission, and the tenant must hold an Intune licence. Without either, Graph rejects the request and the query returns an error rather than an empty result.

## Examples

### Basic info
List every managed device with its user and operating system.

```sql+postgres
select
  device_name,
  serial_number,
  user_principal_name,
  operating_system,
  os_version,
  compliance_state
from
  azuread_intune_managed_device;
```

```sql+sqlite
select
  device_name,
  serial_number,
  user_principal_name,
  operating_system,
  os_version,
  compliance_state
from
  azuread_intune_managed_device;
```

### List devices that are not compliant
Find devices Intune's compliance policies flag as non-compliant or still in their grace period.

```sql+postgres
select
  device_name,
  user_principal_name,
  compliance_state,
  compliance_grace_period_expiration_date_time
from
  azuread_intune_managed_device
where
  compliance_state in ('noncompliant', 'inGracePeriod');
```

```sql+sqlite
select
  device_name,
  user_principal_name,
  compliance_state,
  compliance_grace_period_expiration_date_time
from
  azuread_intune_managed_device
where
  compliance_state in ('noncompliant', 'inGracePeriod');
```

### List devices whose storage is not encrypted
Identify devices that report unencrypted storage.

```sql+postgres
select
  device_name,
  serial_number,
  operating_system,
  user_principal_name
from
  azuread_intune_managed_device
where
  not is_encrypted;
```

```sql+sqlite
select
  device_name,
  serial_number,
  operating_system,
  user_principal_name
from
  azuread_intune_managed_device
where
  is_encrypted = 0;
```

### List devices that have not synced in the last 30 days
Find stale devices that may have been lost, retired or wiped without being removed from Intune.

```sql+postgres
select
  device_name,
  user_principal_name,
  last_sync_date_time
from
  azuread_intune_managed_device
where
  last_sync_date_time < now() - interval '30 days';
```

```sql+sqlite
select
  device_name,
  user_principal_name,
  last_sync_date_time
from
  azuread_intune_managed_device
where
  last_sync_date_time < datetime('now', '-30 days');
```

### List personally owned devices
Review devices enrolled as personal rather than company owned.

```sql+postgres
select
  device_name,
  user_principal_name,
  operating_system,
  managed_device_owner_type
from
  azuread_intune_managed_device
where
  managed_device_owner_type = 'personal';
```

```sql+sqlite
select
  device_name,
  user_principal_name,
  operating_system,
  managed_device_owner_type
from
  azuread_intune_managed_device
where
  managed_device_owner_type = 'personal';
```
