---
title: "Steampipe Table: azuread_governance_role_setting - Query Azure AD Governance Role Settings using SQL"
description: "Allows users to query Azure AD Governance Role Settings, providing detailed information about role settings for Azure resources in Privileged Identity Management (PIM)."
---

# Table: azuread_governance_role_setting - Query Azure AD Governance Role Settings using SQL

Azure AD Governance Role Settings are configurations that define how roles are managed for Azure resources in Privileged Identity Management (PIM). These settings control various aspects of role assignments including expiration rules, MFA requirements, justification requirements, and approval workflows for both admin and user assignments.

## Table Usage Guide

The `azuread_governance_role_setting` table provides insights into governance role settings within Azure Active Directory Privileged Identity Management. As a security administrator or IT professional, you can explore role setting-specific details through this table, including admin and user settings, resource and role definition information, and policy configurations. Utilize it to understand role management policies, audit role settings, and ensure proper governance configurations are in place.

## Examples

### Basic info
Analyze the governance role settings to understand the current configuration and settings for Azure resources in your tenant.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  is_default,
  last_updated_by,
  last_updated_date_time
from
  azuread_governance_role_setting;
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  is_default,
  last_updated_by,
  last_updated_date_time
from
  azuread_governance_role_setting;
```

### List admin eligible settings
Explore admin eligible settings to understand the rules and configurations for admin role assignments.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  admin_eligible_settings
from
  azuread_governance_role_setting
where
  admin_eligible_settings is not null;
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  admin_eligible_settings
from
  azuread_governance_role_setting
where
  admin_eligible_settings is not null;
```

### List admin member settings
Explore admin member settings to understand the rules and configurations for admin role memberships.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  admin_member_settings
from
  azuread_governance_role_setting
where
  admin_member_settings is not null;
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  admin_member_settings
from
  azuread_governance_role_setting
where
  admin_member_settings is not null;
```

### List user eligible settings
Explore user eligible settings to understand the rules and configurations for user role assignments.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  user_eligible_settings
from
  azuread_governance_role_setting
where
  user_eligible_settings is not null;
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  user_eligible_settings
from
  azuread_governance_role_setting
where
  user_eligible_settings is not null;
```

### List user member settings
Explore user member settings to understand the rules and configurations for user role memberships.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  user_member_settings
from
  azuread_governance_role_setting
where
  user_member_settings is not null;
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  user_member_settings
from
  azuread_governance_role_setting
where
  user_member_settings is not null;
```

### Get resource details
Explore resource details to understand the Azure resources associated with role settings.

```sql+postgres
select
  id,
  resource_id,
  resource_details->>'displayName' as resource_display_name,
  resource_details->>'status' as resource_status,
  resource_details->>'type' as resource_type,
  resource_details->>'externalId' as resource_external_id
from
  azuread_governance_role_setting
where
  resource_details is not null;
```

```sql+sqlite
select
  id,
  resource_id,
  json_extract(resource_details, '$.displayName') as resource_display_name,
  json_extract(resource_details, '$.status') as resource_status,
  json_extract(resource_details, '$.type') as resource_type,
  json_extract(resource_details, '$.externalId') as resource_external_id
from
  azuread_governance_role_setting
where
  resource_details is not null;
```

### Get role definition details
Explore role definition details to understand the roles associated with the settings.

```sql+postgres
select
  id,
  role_definition_id,
  role_definition_details->>'displayName' as role_display_name,
  role_definition_details->>'externalId' as role_external_id,
  role_definition_details->>'templateId' as role_template_id
from
  azuread_governance_role_setting
where
  role_definition_details is not null;
```

```sql+sqlite
select
  id,
  role_definition_id,
  json_extract(role_definition_details, '$.displayName') as role_display_name,
  json_extract(role_definition_details, '$.externalId') as role_external_id,
  json_extract(role_definition_details, '$.templateId') as role_template_id
from
  azuread_governance_role_setting
where
  role_definition_details is not null;
```

### Find default role settings
Identify default role settings to understand which settings are applied by default.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  is_default,
  last_updated_by,
  last_updated_date_time
from
  azuread_governance_role_setting
where
  is_default = true;
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  is_default,
  last_updated_by,
  last_updated_date_time
from
  azuread_governance_role_setting
where
  is_default = 1;
```

### Analyze expiration rules
Analyze expiration rules in admin member settings to understand role assignment durations.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  jsonb_array_elements(admin_member_settings) as admin_setting
from
  azuread_governance_role_setting
where
  admin_member_settings is not null
  and jsonb_array_elements(admin_member_settings)->>'ruleIdentifier' = 'ExpirationRule';
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  json_extract(admin_member_settings.value, '$') as admin_setting
from
  azuread_governance_role_setting,
  json_each(admin_member_settings) as admin_member_settings
where
  admin_member_settings is not null
  and json_extract(admin_member_settings.value, '$.ruleIdentifier') = 'ExpirationRule';
```

### Analyze MFA rules
Analyze MFA rules in user member settings to understand multi-factor authentication requirements.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  jsonb_array_elements(user_member_settings) as user_setting
from
  azuread_governance_role_setting
where
  user_member_settings is not null
  and jsonb_array_elements(user_member_settings)->>'ruleIdentifier' = 'MfaRule';
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  json_extract(user_member_settings.value, '$') as user_setting
from
  azuread_governance_role_setting,
  json_each(user_member_settings) as user_member_settings
where
  user_member_settings is not null
  and json_extract(user_member_settings.value, '$.ruleIdentifier') = 'MfaRule';
```

### Analyze approval rules
Analyze approval rules in user member settings to understand approval workflows.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  jsonb_array_elements(user_member_settings) as user_setting
from
  azuread_governance_role_setting
where
  user_member_settings is not null
  and jsonb_array_elements(user_member_settings)->>'ruleIdentifier' = 'ApprovalRule';
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  json_extract(user_member_settings.value, '$') as user_setting
from
  azuread_governance_role_setting,
  json_each(user_member_settings) as user_member_settings
where
  user_member_settings is not null
  and json_extract(user_member_settings.value, '$.ruleIdentifier') = 'ApprovalRule';
```

### Find recently updated settings
Identify recently updated role settings to track changes and modifications.

```sql+postgres
select
  id,
  resource_id,
  role_definition_id,
  last_updated_by,
  last_updated_date_time
from
  azuread_governance_role_setting
where
  last_updated_date_time is not null
order by
  last_updated_date_time desc
limit 10;
```

```sql+sqlite
select
  id,
  resource_id,
  role_definition_id,
  last_updated_by,
  last_updated_date_time
from
  azuread_governance_role_setting
where
  last_updated_date_time is not null
order by
  last_updated_date_time desc
limit 10;
```

### Count settings by resource
Count the number of role settings per resource to understand role management complexity.

```sql+postgres
select
  resource_id,
  resource_details->>'displayName' as resource_display_name,
  count(*) as setting_count
from
  azuread_governance_role_setting
where
  resource_details is not null
group by
  resource_id,
  resource_details->>'displayName'
order by
  setting_count desc;
```

```sql+sqlite
select
  resource_id,
  json_extract(resource_details, '$.displayName') as resource_display_name,
  count(*) as setting_count
from
  azuread_governance_role_setting
where
  resource_details is not null
group by
  resource_id,
  json_extract(resource_details, '$.displayName')
order by
  setting_count desc;
```

### Count settings by role definition
Count the number of role settings per role definition to understand role usage.

```sql+postgres
select
  role_definition_id,
  role_definition_details->>'displayName' as role_display_name,
  count(*) as setting_count
from
  azuread_governance_role_setting
where
  role_definition_details is not null
group by
  role_definition_id,
  role_definition_details->>'displayName'
order by
  setting_count desc;
```

```sql+sqlite
select
  role_definition_id,
  json_extract(role_definition_details, '$.displayName') as role_display_name,
  count(*) as setting_count
from
  azuread_governance_role_setting
where
  role_definition_details is not null
group by
  role_definition_id,
  json_extract(role_definition_details, '$.displayName')
order by
  setting_count desc;
```