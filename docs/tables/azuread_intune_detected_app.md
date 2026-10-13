---
title: "Steampipe Table: azuread_intune_detected_app - Query applications Microsoft Intune detected on managed devices using SQL"
description: "Allows users to query the applications Microsoft Intune discovered on each managed device, including name, version, publisher and platform."
---

# Table: azuread_intune_detected_app - Query applications Microsoft Intune detected on managed devices using SQL

Microsoft Intune inventories the applications installed on the devices it manages. Each detected application carries its name, version, publisher, platform and size, and the number of devices in the tenant that have it installed.

## Table Usage Guide

The `azuread_intune_detected_app` table lists the software installed on Intune managed devices. Use it to build a per-device software inventory, find devices running an outdated version of an application, or find unapproved software.

Rows are read per device from the Microsoft Graph beta endpoint `/deviceManagement/managedDevices/{id}/detectedApps`. Reading this table requires the `DeviceManagementManagedDevices.Read.All` Microsoft Graph application permission, and the tenant must hold an Intune licence. Without either, Graph rejects the request and the query returns an error rather than an empty result.

**Important Notes**
- You must specify `managed_device_id` in the `where` clause to query this table. An `in (...)` list, or a join from `azuread_intune_managed_device`, reads each device in turn.
- A device with no detected applications, or one removed from Intune, returns no rows.

## Examples

### Basic info
List the applications detected on a device.

```sql+postgres
select
  display_name,
  version,
  publisher,
  platform
from
  azuread_intune_detected_app
where
  managed_device_id = '0a1b2c3d-0000-0000-0000-000000000001';
```

```sql+sqlite
select
  display_name,
  version,
  publisher,
  platform
from
  azuread_intune_detected_app
where
  managed_device_id = '0a1b2c3d-0000-0000-0000-000000000001';
```

### List the applications on every managed device
Build a software inventory across all devices Intune manages.

```sql+postgres
select
  d.device_name,
  d.user_principal_name,
  a.display_name,
  a.version,
  a.publisher
from
  azuread_intune_managed_device as d
  join azuread_intune_detected_app as a on a.managed_device_id = d.id
order by
  d.device_name,
  a.display_name;
```

```sql+sqlite
select
  d.device_name,
  d.user_principal_name,
  a.display_name,
  a.version,
  a.publisher
from
  azuread_intune_managed_device as d
  join azuread_intune_detected_app as a on a.managed_device_id = d.id
order by
  d.device_name,
  a.display_name;
```

### Find devices running a specific application
Identify the devices that have a given application installed, and which version.

```sql+postgres
select
  d.device_name,
  d.user_principal_name,
  a.version
from
  azuread_intune_managed_device as d
  join azuread_intune_detected_app as a on a.managed_device_id = d.id
where
  a.display_name ilike '%chrome%';
```

```sql+sqlite
select
  d.device_name,
  d.user_principal_name,
  a.version
from
  azuread_intune_managed_device as d
  join azuread_intune_detected_app as a on a.managed_device_id = d.id
where
  a.display_name like '%chrome%';
```

### Count the applications detected on each device
Spot devices with an unusually large or small number of installed applications.

```sql+postgres
select
  d.device_name,
  count(a.id) as app_count
from
  azuread_intune_managed_device as d
  join azuread_intune_detected_app as a on a.managed_device_id = d.id
group by
  d.device_name
order by
  app_count desc;
```

```sql+sqlite
select
  d.device_name,
  count(a.id) as app_count
from
  azuread_intune_managed_device as d
  join azuread_intune_detected_app as a on a.managed_device_id = d.id
group by
  d.device_name
order by
  app_count desc;
```
