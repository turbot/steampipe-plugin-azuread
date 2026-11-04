---
title: "Steampipe Table: azuread_subscribed_sku - Query Azure AD subscribed SKUs using SQL"
description: "Allows users to query the commercial subscriptions (SKUs) that an organization has acquired."
---

# Table: azuread_subscribed_sku - Query Azure AD subscribed SKUs using SQL

A subscribed SKU represents a commercial subscription an organization has acquired, for example Office 365 E5 or Microsoft Defender for Office 365 Plan 2. Each SKU bundles a set of service plans, and its capability status shows whether it is currently enabled. Inspecting the service plans is how a tenant's entitlement to a feature - such as Defender for Office 365 - is determined.

## Table Usage Guide

The `azuread_subscribed_sku` table provides insights into an organization's licensing. Use it to see which SKUs are held, their capability status, and the individual service plans each SKU includes.

Reading this table requires the `Organization.Read.All` Microsoft Graph application permission.

## Examples

### List all subscribed SKUs
Review every subscription the tenant holds.

```sql+postgres
select
  sku_id,
  sku_part_number,
  capability_status
from
  azuread_subscribed_sku;
```

```sql+sqlite
select
  sku_id,
  sku_part_number,
  capability_status
from
  azuread_subscribed_sku;
```

### Find SKUs that include a Defender for Office 365 service plan
Expand the service plans to check for a specific entitlement.

```sql+postgres
select
  sku_part_number,
  plan ->> 'servicePlanName' as service_plan_name,
  plan ->> 'provisioningStatus' as provisioning_status
from
  azuread_subscribed_sku,
  jsonb_array_elements(service_plans) as plan
where
  plan ->> 'servicePlanName' like 'THREAT_INTELLIGENCE%';
```

```sql+sqlite
select
  sku_part_number,
  json_extract(plan.value, '$.servicePlanName') as service_plan_name,
  json_extract(plan.value, '$.provisioningStatus') as provisioning_status
from
  azuread_subscribed_sku,
  json_each(service_plans) as plan
where
  json_extract(plan.value, '$.servicePlanName') like 'THREAT_INTELLIGENCE%';
```
