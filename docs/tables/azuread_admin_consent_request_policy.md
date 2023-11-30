---
title: "Steampipe Table: azuread_admin_consent_request_policy - Query Azure Active Directory Admin Consent Request Policies using SQL"
description: "Allows users to query Admin Consent Request Policies in Azure Active Directory, specifically the policy settings that control the workflow of admin consent requests, providing insights into the permissions and access management."
---

# Table: azuread_admin_consent_request_policy - Query Azure Active Directory Admin Consent Request Policies using SQL

An Azure Active Directory Admin Consent Request Policy is a feature within Microsoft Azure that controls the workflow of admin consent requests. It provides a centralized way to manage and review admin consent requests for applications requiring access to data they do not have permissions for. Azure AD admin consent request policy helps you stay informed about the access requests and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `azuread_admin_consent_request_policy` table provides insights into admin consent request policies within Azure Active Directory. As a security engineer, explore policy-specific details through this table, including policy settings, approval steps, and associated metadata. Utilize it to uncover information about policies, such as those with specific approval steps, the workflow of admin consent requests, and the verification of policy settings.

## Examples

### Basic info
Explore which Azure Active Directory admin consent request policies are enabled and their respective versions. This is useful for assessing the current status and versioning of your policies.

```sql
select
  title,
  is_enabled,
  version
from
  azuread_admin_consent_request_policy;
```

### Check admin consent workflow is enabled
Determine if the admin consent workflow is active in Azure Active Directory, which is essential for enhancing security by ensuring that admins explicitly approve access requests to specific resources.

```sql
select
  title,
  is_enabled,
  version
from
  azuread_admin_consent_request_policy
where
  is_enabled;
```

### List users who can review new admin consent requests
Determine the users who have the authority to review new administrative consent requests. This is useful for managing permissions and ensuring only appropriate personnel are able to handle these requests.

```sql
select
  p.title,
  p.is_enabled,
  u.display_name as user_display_name,
  u.user_principal_name
from
  azuread_admin_consent_request_policy as p,
  jsonb_array_elements(reviewers) as r
  left join azuread_user as u on split_part(r ->> 'query', '/', 4) = u.id
where
  is_enabled;
```