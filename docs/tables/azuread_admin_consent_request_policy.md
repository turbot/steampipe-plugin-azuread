# Table: azuread_admin_consent_request_policy

Represents the policy for enabling or disabling the Azure AD admin consent workflow. The admin consent workflow allows users to request access for apps that they wish to use and that require admin authorization before users can use the apps to access organizational data.

## Examples

### Basic info

```sql
select
  title,
  is_enabled,
  version
from
  azuread_admin_consent_request_policy;
```

### Check admin consent workflow is enabled

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
