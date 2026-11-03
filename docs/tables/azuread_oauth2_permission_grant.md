---
title: "Steampipe Table: azuread_oauth2_permission_grant - Query Azure Active Directory Delegated Permission Grants using SQL"
description: "Allows users to query delegated permission grants (OAuth2 permission grants) in Azure Active Directory, showing which applications can act on behalf of users, for which APIs, and with which permissions."
---

# Table: azuread_oauth2_permission_grant - Query Azure Active Directory Delegated Permission Grants using SQL

A delegated permission grant (oAuth2PermissionGrant) records that a client application is allowed to call an API, such as Microsoft Graph, on behalf of a signed-in user. Grants are created when an administrator consents for the whole organization (consent type `AllPrincipals`) or when a single user consents for themselves (consent type `Principal`). Delegated permissions are separate from application permissions, which are app role assignments (see `azuread_service_principal_app_role_assigned_to`).

## Table Usage Guide

The `azuread_oauth2_permission_grant` table shows what third-party and in-house applications can do with your users' data. Use it to find apps that users have consented to on their own, apps holding broad permissions such as `Files.Read.All` or `Mail.Read`, and admin-consented grants that apply to everyone in the tenant.

`client_id` and `resource_id` are service principal object IDs (not application IDs). Join them to `azuread_service_principal.id` to get display names.

**Important Notes**
- This table requires the `Directory.Read.All` permission.
- For improved performance, it is advised that you use the optional qualifiers `client_id`, `consent_type`, `principal_id` and `resource_id` to limit the result set.

## Examples

### Basic info
List every delegated permission grant with its consent type and the permissions it includes.

```sql+postgres
select
  id,
  client_id,
  consent_type,
  principal_id,
  resource_id,
  scope
from
  azuread_oauth2_permission_grant;
```

```sql+sqlite
select
  id,
  client_id,
  consent_type,
  principal_id,
  resource_id,
  scope
from
  azuread_oauth2_permission_grant;
```

### List grants with application and API names
See which applications can act on behalf of users, and against which API.

```sql+postgres
select
  c.display_name as client_app,
  r.display_name as api,
  g.consent_type,
  g.scopes
from
  azuread_oauth2_permission_grant as g
  left join azuread_service_principal as c on c.id = g.client_id
  left join azuread_service_principal as r on r.id = g.resource_id
order by
  client_app;
```

```sql+sqlite
select
  c.display_name as client_app,
  r.display_name as api,
  g.consent_type,
  g.scopes
from
  azuread_oauth2_permission_grant as g
  left join azuread_service_principal as c on c.id = g.client_id
  left join azuread_service_principal as r on r.id = g.resource_id
order by
  client_app;
```

### List grants made by individual users
Find applications that users consented to for themselves rather than through an administrator.

```sql+postgres
select
  g.client_id,
  u.user_principal_name,
  g.scope
from
  azuread_oauth2_permission_grant as g
  left join azuread_user as u on u.id = g.principal_id
where
  g.consent_type = 'Principal';
```

```sql+sqlite
select
  g.client_id,
  u.user_principal_name,
  g.scope
from
  azuread_oauth2_permission_grant as g
  left join azuread_user as u on u.id = g.principal_id
where
  g.consent_type = 'Principal';
```

### List grants that include high-risk Microsoft Graph permissions
Identify applications that can read all files, mail or the directory, or write data, on behalf of users.

```sql+postgres
select
  c.display_name as client_app,
  g.consent_type,
  s as permission
from
  azuread_oauth2_permission_grant as g
  join azuread_service_principal as r on r.id = g.resource_id
  left join azuread_service_principal as c on c.id = g.client_id,
  jsonb_array_elements_text(g.scopes) as s
where
  r.app_id = '00000003-0000-0000-c000-000000000000'
  and (
    s ~ '^(Files|Sites|Mail|Directory|RoleManagement|Application|AppRoleAssignment)\.'
    or s like '%ReadWrite%'
  );
```

```sql+sqlite
select
  c.display_name as client_app,
  g.consent_type,
  s.value as permission
from
  azuread_oauth2_permission_grant as g
  join azuread_service_principal as r on r.id = g.resource_id
  left join azuread_service_principal as c on c.id = g.client_id,
  json_each(g.scopes) as s
where
  r.app_id = '00000003-0000-0000-c000-000000000000'
  and (
    s.value like 'Files.%'
    or s.value like 'Sites.%'
    or s.value like 'Mail.%'
    or s.value like 'Directory.%'
    or s.value like 'RoleManagement.%'
    or s.value like 'Application.%'
    or s.value like 'AppRoleAssignment.%'
    or s.value like '%ReadWrite%'
  );
```

### List grants for a specific application
Show everything one application has been granted, using its service principal object ID.

```sql+postgres
select
  consent_type,
  principal_id,
  resource_id,
  scopes
from
  azuread_oauth2_permission_grant
where
  client_id = 'af79d139-fd2e-4291-821d-b2841e5f53b2';
```

```sql+sqlite
select
  consent_type,
  principal_id,
  resource_id,
  scopes
from
  azuread_oauth2_permission_grant
where
  client_id = 'af79d139-fd2e-4291-821d-b2841e5f53b2';
```
