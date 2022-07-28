# Table: azuread_sign_in_report

Retrieve the Azure AD user sign-ins for your tenant. Sign-ins that are interactive in nature (where a username/password is passed as part of auth token) and successful federated sign-ins are currently included in the sign-in logs.

## Examples

### Basic info

```sql
select
  id,
  created_date_time,
  user_display_name,
  user_principal_name,
  ip_address,
  location ->> 'city' as city
from
  azuread_sign_in_report;
```

### List an user sign in details

```sql
select
  user_display_name,
  id,
  app_display_name,
  user_principal_name
from
  azuread_sign_in_report
where
  user_principal_name = 'abc@myacc.onmicrosoft.com';
```
