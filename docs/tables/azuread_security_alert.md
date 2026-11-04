---
title: "Steampipe Table: azuread_security_alert - Query Microsoft 365 Defender security alerts using SQL"
description: "Allows users to query the security alerts (alerts_v2) surfaced across Microsoft 365 Defender for an Azure AD tenant."
---

# Table: azuread_security_alert - Query Microsoft 365 Defender security alerts using SQL

Microsoft 365 Defender consolidates the alerts raised by its security products - Defender for Endpoint, Defender for Office 365, Defender for Identity and Defender for Cloud Apps - into a single unified `alerts_v2` collection. Each alert tracks a suspicious or malicious activity, its severity, its investigation state, and the evidence and analyst comments gathered while it is worked.

## Table Usage Guide

The `azuread_security_alert` table provides insights into security alerts within Microsoft 365 Defender. As a security analyst, explore alert-specific details through this table, including severity, status, classification, the service that raised the alert, and the evidence and comments attached to it. The `service_source`, `status`, `severity` and `last_update_date_time` quals are pushed down into the Graph `$filter`, so scope your query with them to keep it fast on busy tenants.

Reading this table requires the `SecurityAlert.Read.All` Microsoft Graph application permission.

## Examples

### List open Defender for Office 365 alerts
Focus on the email threat alerts that still need attention.

```sql+postgres
select
  id,
  title,
  severity,
  status,
  created_date_time
from
  azuread_security_alert
where
  service_source = 'microsoftDefenderForOffice365'
  and status = 'new';
```

```sql+sqlite
select
  id,
  title,
  severity,
  status,
  created_date_time
from
  azuread_security_alert
where
  service_source = 'microsoftDefenderForOffice365'
  and status = 'new';
```

### High severity alerts updated since a point in time
Use the pushed-down `last_update_date_time` qual to pull only recently changed alerts.

```sql+postgres
select
  id,
  title,
  severity,
  last_update_date_time
from
  azuread_security_alert
where
  severity = 'high'
  and last_update_date_time > now() - interval '1 day';
```

```sql+sqlite
select
  id,
  title,
  severity,
  last_update_date_time
from
  azuread_security_alert
where
  severity = 'high'
  and last_update_date_time > datetime('now', '-1 day');
```

### Inspect the evidence attached to an alert
The `evidence` column preserves the per-item verdict, remediation status and type-specific fields.

```sql+postgres
select
  id,
  title,
  jsonb_array_elements(evidence) -> 'verdict' as evidence_verdict
from
  azuread_security_alert
where
  id = 'da637551227677560813_-961444813';
```

```sql+sqlite
select
  id,
  title,
  json_extract(evidence.value, '$.verdict') as evidence_verdict
from
  azuread_security_alert,
  json_each(evidence)
where
  id = 'da637551227677560813_-961444813';
```
