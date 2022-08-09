# Table: azuread_conditional_access_policy

Conditional Access is the tool used by Azure Active Directory to bring signals together, to make decisions, and enforce organizational policies. Conditional Access is at the heart of the new identity driven control plane.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  state,
  created_date_time,
  built_in_controls
from
  azuread_conditional_access_policy;
```

### List conditional access policies with mfa enabled

```sql
select
  id,
  display_name,
  built_in_controls
from
  azuread_conditional_access_policy
where
  built_in_controls ?& array['mfa'];
```
