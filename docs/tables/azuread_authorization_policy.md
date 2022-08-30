# Table: azuread_application

Represents a policy that can control Azure Active Directory authorization settings.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  allow_invites_from
from
  azuread_authorization_policy;
```

### Check user consent to apps accessing company data on their behalf is not allowed

```sql
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  default_user_role_permissions ->> 'permissionGrantPoliciesAssigned' = '[]';
```
