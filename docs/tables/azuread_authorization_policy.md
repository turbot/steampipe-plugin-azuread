# Table: azuread_authorization_policy

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

### Check if user consent to apps accessing company data on their behalf is not allowed

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

### Check if all members are allowed to invite external users to the organization

```sql
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  allow_invites_from = 'everyone';
```

### Check if email validation is not required to join the tenant

```sql
select
  display_name,
  id,
  default_user_role_permissions
from
  azuread_authorization_policy
where
  not allowed_email_verified_users_to_join_organization;
```
