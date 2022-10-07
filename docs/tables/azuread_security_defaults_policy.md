# Table: azuread_security_defaults_policy

Security defaults in Azure Active Directory (Azure AD) make it easier to be secure and help protect your organization. Security defaults contain pre-configured security settings for common attacks.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  is_enabled
from
  azuread_security_defaults_policy;
```
