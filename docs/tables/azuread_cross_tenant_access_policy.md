---
title: "Steampipe Table: azuread_cross_tenant_access_policy - Query Azure AD Cross-Tenant Access Policies using SQL"
description: "Allows users to query Azure AD Cross-Tenant Access Policies, providing detailed information about the policies that control cross-tenant collaboration and access."
---

# Table: azuread_cross_tenant_access_policy - Query Azure AD Cross-Tenant Access Policies using SQL

Azure AD Cross-Tenant Access Policy is a feature in Azure Active Directory that allows administrators to define policies that control cross-tenant collaboration and access. These policies manage how users from other organizations can access your resources via Azure AD B2B collaboration and B2B direct connect, and vice versa. This feature is crucial for managing secure collaboration between organizations.

## Table Usage Guide

The `azuread_cross_tenant_access_policy` table provides insights into Cross-Tenant Access Policies within Azure Active Directory. As a security administrator, you can explore policy-specific details through this table, including default configurations, partner-specific settings, and cloud endpoint configurations. Utilize it to uncover information about cross-tenant policies, helping you to maintain security and compliance in your organization's external collaborations.

## Examples

### Basic info
Analyze the settings to understand the configuration of your Azure Active Directory cross-tenant access policy. This can help you assess the elements within your policy and make necessary adjustments.

```sql+postgres
select
  id,
  display_name,
  allowed_cloud_endpoints,
  default_configuration
from
  azuread_cross_tenant_access_policy;
```

```sql+sqlite
select
  id,
  display_name,
  allowed_cloud_endpoints,
  default_configuration
from
  azuread_cross_tenant_access_policy;
```

### List cross-tenant access policy with default configuration
Uncover the details of cross-tenant access policies that have default configurations. This is useful for understanding how your organization collaborates with external partners by default.

```sql+postgres
select
  id,
  display_name,
  default_configuration
from
  azuread_cross_tenant_access_policy
where
  default_configuration is not null;
```

```sql+sqlite
select
  id,
  display_name,
  default_configuration
from
  azuread_cross_tenant_access_policy
where
  default_configuration is not null;
```

### List cross-tenant access policy with partner configurations
Examine the partner-specific configurations to understand how your organization interacts with specific external organizations.

```sql+postgres
select
  id,
  display_name,
  partners
from
  azuread_cross_tenant_access_policy
where
  partners is not null;
```

```sql+sqlite
select
  id,
  display_name,
  partners
from
  azuread_cross_tenant_access_policy
where
  partners is not null;
```

### List policies with specific cloud endpoints
Identify cross-tenant access policies that have specific cloud endpoint configurations, which determine which Microsoft clouds your organization collaborates with.

```sql+postgres
select
  id,
  display_name,
  allowed_cloud_endpoints
from
  azuread_cross_tenant_access_policy
where
  allowed_cloud_endpoints is not null;
```

```sql+sqlite
select
  id,
  display_name,
  allowed_cloud_endpoints
from
  azuread_cross_tenant_access_policy
where
  allowed_cloud_endpoints is not null;
```
