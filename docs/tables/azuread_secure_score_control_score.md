---
title: "Steampipe Table: azuread_secure_score_control_score - Query Microsoft Secure Score control scores using SQL"
description: "Allows users to query the per-control scores of the most recent Microsoft Secure Score for an Azure AD tenant."
---

# Table: azuread_secure_score_control_score - Query Microsoft Secure Score control scores using SQL

Microsoft Secure Score measures an organization's security posture as a single figure made up of many individual security controls. Each daily `secureScore` snapshot carries a `controlScores` collection - one entry per control - describing how the tenant scored against that control. This table returns one row per control of the **latest** secure score snapshot.

## Table Usage Guide

The `azuread_secure_score_control_score` table provides insights into the current security posture of a tenant, control by control. Use it to see which controls are implemented, the achieved versus maximum score for each, and when each control was last evaluated.

Reading this table requires the `SecurityEvents.Read.All` Microsoft Graph application permission.

### Column sources

Microsoft Graph splits a control's data across two resources, so this table draws each column from the place Graph actually populates it:

- `control_name`, `control_category`, `score` - read directly from the `controlScore` entry of the latest `secureScore`.
- `implementation_status`, `last_synced_date_time` - not formal properties of `controlScore`; Graph returns them only inside the entry's additional data (`implementationStatus`, `lastSynced`). They are `NULL` when Graph omits them.
- `max_score` - a control's maximum is not carried on the `secureScore` at all; it lives on the matching `secureScoreControlProfile`. This table joins the profile (by control name) to populate it, and returns `NULL` when no profile exposes it.

## Examples

### List the latest secure score controls
Review the whole latest snapshot, control by control.

```sql+postgres
select
  control_name,
  control_category,
  score,
  max_score,
  implementation_status
from
  azuread_secure_score_control_score
order by
  control_category,
  control_name;
```

```sql+sqlite
select
  control_name,
  control_category,
  score,
  max_score,
  implementation_status
from
  azuread_secure_score_control_score
order by
  control_category,
  control_name;
```

### Find controls that are not fully implemented
Surface the controls scoring below their maximum.

```sql+postgres
select
  control_name,
  control_category,
  score,
  max_score
from
  azuread_secure_score_control_score
where
  max_score is not null
  and score < max_score;
```

```sql+sqlite
select
  control_name,
  control_category,
  score,
  max_score
from
  azuread_secure_score_control_score
where
  max_score is not null
  and score < max_score;
```
