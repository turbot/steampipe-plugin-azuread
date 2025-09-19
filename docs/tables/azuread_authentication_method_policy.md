---
title: "Steampipe Table: azuread_authentication_method_policy - Query Azure AD Authentication Methods Policy using SQL"
description: "Allows users to query Azure AD Authentication Methods Policy, providing detailed information about authentication method configurations and settings."
---

# Table: azuread_authentication_method_policy - Query Azure AD Authentication Methods Policy using SQL

Azure AD Authentication Methods Policy is a tenant-wide policy that controls which authentication methods are allowed in the tenant, authentication method registration requirements, and self-service password reset settings. This policy defines the available authentication methods and their configurations for users in the organization.

## Table Usage Guide

The `azuread_authentication_method_policy` table provides insights into Authentication Methods Policy within Azure Active Directory. As a security administrator or IT professional, you can explore policy-specific details through this table, including authentication method configurations, registration enforcement settings, and policy metadata. Utilize it to understand available authentication methods, audit authentication settings, and ensure proper security configurations are in place.

## Examples

### Basic info

Analyze the authentication methods policy to understand the current configuration and settings in your Azure Active Directory tenant.

```sql+postgres
select
  id,
  display_name,
  description,
  last_modified_date_time,
  policy_migration_state,
  policy_version,
  reconfirmation_in_days
from
  azuread_authentication_method_policy;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  last_modified_date_time,
  policy_migration_state,
  policy_version,
  reconfirmation_in_days
from
  azuread_authentication_method_policy;
```

### List all authentication method configurations

Explore all available authentication methods and their key properties to understand what authentication options are configured in your tenant.

```sql+postgres
select
  auth_method->>'id' as method_id,
  auth_method->>'state' as state,
  auth_method->>'@odata.type' as method_type,
  -- Common properties
  jsonb_array_length(auth_method->'includeTargets') as include_targets_count,
  jsonb_array_length(auth_method->'excludeTargets') as exclude_targets_count,
  -- FIDO2 specific properties
  auth_method->'isAttestationEnforced' as is_attestation_enforced,
  auth_method->'isSelfServiceRegistrationAllowed' as is_self_service_registration_allowed,
  -- Microsoft Authenticator specific properties
  auth_method->'isSoftwareOathEnabled' as is_software_oath_enabled,
  auth_method->'featureSettings'->'numberMatchingRequiredState'->>'state' as number_matching_state,
  auth_method->'featureSettings'->'companionAppAllowedState'->>'state' as companion_app_state,
  auth_method->'featureSettings'->'displayAppInformationRequiredState'->>'state' as display_app_info_state,
  auth_method->'featureSettings'->'displayLocationInformationRequiredState'->>'state' as display_location_state,
  -- Email specific properties
  auth_method->'allowExternalIdToUseEmailOtp' as allow_external_id_to_use_email_otp,
  -- Voice specific properties
  auth_method->'isOfficePhoneAllowed' as is_office_phone_allowed,
  -- Temporary Access Pass specific properties
  auth_method->'defaultLength' as default_length,
  auth_method->'defaultLifetimeInMinutes' as default_lifetime_in_minutes,
  auth_method->'isUsableOnce' as is_usable_once
from
  azuread_authentication_method_policy,
  jsonb_array_elements(authentication_method_configurations) as auth_method
order by
  auth_method->>'id';
```

```sql+sqlite
select
  json_extract(auth_config.value, '$.id') as method_id,
  json_extract(auth_config.value, '$.state') as state,
  json_extract(auth_config.value, '$."@odata.type"') as method_type,
  -- Common properties
  json_array_length(json_extract(auth_config.value, '$.includeTargets')) as include_targets_count,
  json_array_length(json_extract(auth_config.value, '$.excludeTargets')) as exclude_targets_count,
  -- FIDO2 specific properties
  json_extract(auth_config.value, '$.isAttestationEnforced') as is_attestation_enforced,
  json_extract(auth_config.value, '$.isSelfServiceRegistrationAllowed') as is_self_service_registration_allowed,
  -- Microsoft Authenticator specific properties
  json_extract(auth_config.value, '$.isSoftwareOathEnabled') as is_software_oath_enabled,
  json_extract(auth_config.value, '$.featureSettings.numberMatchingRequiredState.state') as number_matching_state,
  json_extract(auth_config.value, '$.featureSettings.companionAppAllowedState.state') as companion_app_state,
  json_extract(auth_config.value, '$.featureSettings.displayAppInformationRequiredState.state') as display_app_info_state,
  json_extract(auth_config.value, '$.featureSettings.displayLocationInformationRequiredState.state') as display_location_state,
  -- Email specific properties
  json_extract(auth_config.value, '$.allowExternalIdToUseEmailOtp') as allow_external_id_to_use_email_otp,
  -- Voice specific properties
  json_extract(auth_config.value, '$.isOfficePhoneAllowed') as is_office_phone_allowed,
  -- Temporary Access Pass specific properties
  json_extract(auth_config.value, '$.defaultLength') as default_length,
  json_extract(auth_config.value, '$.defaultLifetimeInMinutes') as default_lifetime_in_minutes,
  json_extract(auth_config.value, '$.isUsableOnce') as is_usable_once
from
  azuread_authentication_method_policy,
  json_each(authentication_method_configurations) as auth_config
order by
  json_extract(auth_config.value, '$.id');
```

### Find enabled authentication methods

Identify which authentication methods are currently enabled in your tenant to assess the available authentication options.

```sql+postgres
select
  auth_method->>'id' as method_id,
  auth_method->>'state' as state,
  auth_method->>'@odata.type' as method_type
from
  azuread_authentication_method_policy,
  jsonb_array_elements(authentication_method_configurations) as auth_method
where
  auth_method->>'state' = 'enabled';
```

```sql+sqlite
select
  json_extract(auth_config.value, '$.id') as method_id,
  json_extract(auth_config.value, '$.state') as state,
  json_extract(auth_config.value, '$."@odata.type"') as method_type
from
  azuread_authentication_method_policy,
  json_each(authentication_method_configurations) as auth_config
where
  json_extract(auth_config.value, '$.state') = 'enabled';
```

### Analyze Microsoft Authenticator settings

Examine the Microsoft Authenticator configuration including feature settings like number matching and companion app access to understand the security requirements.

```sql+postgres
select
  auth_method->>'id' as method_id,
  auth_method->>'state' as state,
  auth_method->'featureSettings'->'numberMatchingRequiredState'->>'state' as number_matching_state,
  auth_method->'featureSettings'->'companionAppAllowedState'->>'state' as companion_app_state,
  auth_method->'featureSettings'->'displayAppInformationRequiredState'->>'state' as display_app_info_state,
  auth_method->'featureSettings'->'displayLocationInformationRequiredState'->>'state' as display_location_state
from
  azuread_authentication_method_policy,
  jsonb_array_elements(authentication_method_configurations) as auth_method
where
  auth_method->>'id' = 'MicrosoftAuthenticator';
```

```sql+sqlite
select
  json_extract(auth_config.value, '$.id') as method_id,
  json_extract(auth_config.value, '$.state') as state,
  json_extract(auth_config.value, '$.featureSettings.numberMatchingRequiredState.state') as number_matching_state,
  json_extract(auth_config.value, '$.featureSettings.companionAppAllowedState.state') as companion_app_state,
  json_extract(auth_config.value, '$.featureSettings.displayAppInformationRequiredState.state') as display_app_info_state,
  json_extract(auth_config.value, '$.featureSettings.displayLocationInformationRequiredState.state') as display_location_state
from
  azuread_authentication_method_policy,
  json_each(authentication_method_configurations) as auth_config
where
  json_extract(auth_config.value, '$.id') = 'MicrosoftAuthenticator';
```

### Check registration enforcement settings

Review the registration enforcement configuration to understand how authentication method registration is managed in your tenant.

```sql+postgres
select
  id,
  display_name,
  registration_enforcement
from
  azuread_authentication_method_policy;
```

```sql+sqlite
select
  id,
  display_name,
  registration_enforcement
from
  azuread_authentication_method_policy;
```

### List authentication methods with include targets

Identify authentication methods that have specific include targets configured to understand the scope of authentication method availability.

```sql+postgres
select
  auth_method->>'id' as method_id,
  auth_method->>'state' as state,
  jsonb_array_elements(auth_method->'includeTargets') as include_target
from
  azuread_authentication_method_policy,
  jsonb_array_elements(authentication_method_configurations) as auth_method
where
  auth_method ? 'includeTargets';
```

```sql+sqlite
select
  json_extract(auth_config.value, '$.id') as method_id,
  json_extract(auth_config.value, '$.state') as state,
  json_extract(target.value, '$') as include_target
from
  azuread_authentication_method_policy,
  json_each(authentication_method_configurations) as auth_config,
  json_each(json_extract(auth_config.value, '$.includeTargets')) as target
where
  json_extract(auth_config.value, '$.includeTargets') is not null;
```

### Policy compliance check

Assess the overall policy configuration to ensure proper authentication method settings are in place for security compliance.

```sql+postgres
select
  id,
  display_name,
  description,
  policy_migration_state,
  policy_version,
  registration_enforcement,
  authentication_method_configurations
from
  azuread_authentication_method_policy;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  policy_migration_state,
  policy_version,
  registration_enforcement,
  authentication_method_configurations
from
  azuread_authentication_method_policy;
```
