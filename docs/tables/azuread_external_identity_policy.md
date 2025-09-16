---
title: "Steampipe Table: azuread_external_identity_policy - Query Azure Active Directory External Identity Policy using SQL"
description: "Allows users to query the tenant-wide external identity policy in Azure Active Directory, providing insights into external user management settings."
---

# Table: azuread_external_identity_policy - Query Azure Active Directory External Identity Policy using SQL

Azure Active Directory (Azure AD) External Identity Policy represents the tenant-wide policy that controls whether external users can leave a Microsoft Entra tenant via self-service controls. This policy helps organizations manage external user lifecycle and data retention settings.

## Table Usage Guide

The `azuread_external_identity_policy` table provides insights into the External Identity Policy within Azure Active Directory. As a DevOps engineer or IT professional, you can explore policy-specific details through this table, including self-service controls, data retention settings, and external user management configurations. Utilize it to monitor external user policies, ensure compliance with data governance requirements, and audit external identity management settings.

## Examples

### Basic info

Explore the external identity policy to understand the current tenant-wide settings for external user management.

```sql+postgres
select
  id,
  display_name,
  allow_external_identities_to_leave,
  allow_deleted_identities_data_removal,
  deleted_date_time
from
  azuread_external_identity_policy;
```

```sql+sqlite
select
  id,
  display_name,
  allow_external_identities_to_leave,
  allow_deleted_identities_data_removal,
  deleted_date_time
from
  azuread_external_identity_policy;
```

### Check external user self-service settings

Verify whether external users are allowed to leave the tenant via self-service controls.

```sql+postgres
select
  display_name,
  allow_external_identities_to_leave,
  case
    when allow_external_identities_to_leave then 'External users can self-remove'
    else 'External users cannot self-remove'
  end as self_service_status
from
  azuread_external_identity_policy;
```

```sql+sqlite
select
  display_name,
  allow_external_identities_to_leave,
  case
    when allow_external_identities_to_leave then 'External users can self-remove'
    else 'External users cannot self-remove'
  end as self_service_status
from
  azuread_external_identity_policy;
```

### Check data retention settings

Review the data retention policy for deleted external identities.

```sql+postgres
select
  display_name,
  allow_deleted_identities_data_removal,
  case
    when allow_deleted_identities_data_removal then 'Deleted identity data can be removed'
    else 'Deleted identity data cannot be removed'
  end as data_retention_status
from
  azuread_external_identity_policy;
```

```sql+sqlite
select
  display_name,
  allow_deleted_identities_data_removal,
  case
    when allow_deleted_identities_data_removal then 'Deleted identity data can be removed'
    else 'Deleted identity data cannot be removed'
  end as data_retention_status
from
  azuread_external_identity_policy;
```

### Policy compliance check

Perform a compliance check to ensure external identity policies meet organizational requirements.

```sql+postgres
select
  id,
  display_name,
  allow_external_identities_to_leave,
  allow_deleted_identities_data_removal,
  case
    when not allow_external_identities_to_leave and not allow_deleted_identities_data_removal then 'Strict Policy - No self-service, No data removal'
    when not allow_external_identities_to_leave and allow_deleted_identities_data_removal then 'Moderate Policy - No self-service, Data removal allowed'
    when allow_external_identities_to_leave and not allow_deleted_identities_data_removal then 'Moderate Policy - Self-service allowed, No data removal'
    else 'Permissive Policy - Self-service and data removal allowed'
  end as policy_type
from
  azuread_external_identity_policy;
```

```sql+sqlite
select
  id,
  display_name,
  allow_external_identities_to_leave,
  allow_deleted_identities_data_removal,
  case
    when not allow_external_identities_to_leave and not allow_deleted_identities_data_removal then 'Strict Policy - No self-service, No data removal'
    when not allow_external_identities_to_leave and allow_deleted_identities_data_removal then 'Moderate Policy - No self-service, Data removal allowed'
    when allow_external_identities_to_leave and not allow_deleted_identities_data_removal then 'Moderate Policy - Self-service allowed, No data removal'
    else 'Permissive Policy - Self-service and data removal allowed'
  end as policy_type
from
  azuread_external_identity_policy;
```

### Security assessment

Assess the security posture of external identity management policies.

```sql+postgres
select
  display_name,
  allow_external_identities_to_leave,
  allow_deleted_identities_data_removal,
  case
    when allow_external_identities_to_leave then 'HIGH - External users can self-remove'
    else 'LOW - External users cannot self-remove'
  end as self_service_risk,
  case
    when allow_deleted_identities_data_removal then 'MEDIUM - Data can be permanently removed'
    else 'LOW - Data retention enforced'
  end as data_retention_risk
from
  azuread_external_identity_policy;
```

```sql+sqlite
select
  display_name,
  allow_external_identities_to_leave,
  allow_deleted_identities_data_removal,
  case
    when allow_external_identities_to_leave then 'HIGH - External users can self-remove'
    else 'LOW - External users cannot self-remove'
  end as self_service_risk,
  case
    when allow_deleted_identities_data_removal then 'MEDIUM - Data can be permanently removed'
    else 'LOW - Data retention enforced'
  end as data_retention_risk
from
  azuread_external_identity_policy;
```

### Policy audit trail

Check if the policy has been deleted and when.

```sql+postgres
select
  id,
  display_name,
  deleted_date_time,
  case
    when deleted_date_time is not null then 'Policy has been deleted'
    else 'Policy is active'
  end as policy_status
from
  azuread_external_identity_policy;
```

```sql+sqlite
select
  id,
  display_name,
  deleted_date_time,
  case
    when deleted_date_time is not null then 'Policy has been deleted'
    else 'Policy is active'
  end as policy_status
from
  azuread_external_identity_policy;
```

## Schema

| Name                                    | Type        | Description                                                                            |
| --------------------------------------- | ----------- | -------------------------------------------------------------------------------------- |
| `id`                                    | `text`      | The unique identifier for the external identity policy.                                |
| `display_name`                          | `text`      | The display name for the external identity policy.                                     |
| `allow_external_identities_to_leave`    | `boolean`   | Flag indicating whether external users can leave the tenant via self-service controls. |
| `allow_deleted_identities_data_removal` | `boolean`   | Flag indicating whether deleted identities data can be removed.                        |
| `deleted_date_time`                     | `timestamp` | The date and time when the policy was deleted.                                         |
| `title`                                 | `text`      | Title of the resource.                                                                 |
| `tenant_id`                             | `text`      | The Azure Tenant ID where the resource is located.                                     |
| `sp_connection_name`                    | `text`      | Steampipe connection name.                                                             |
| `sp_ctx`                                | `jsonb`     | Steampipe context in JSON form.                                                        |
