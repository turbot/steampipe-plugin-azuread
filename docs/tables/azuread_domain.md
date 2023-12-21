---
title: "Steampipe Table: azuread_domain - Query Azure Active Directory Domains using SQL"
description: "Allows users to query Azure Active Directory Domains, providing information about the various domains associated with an Azure Active Directory instance."
---

# Table: azuread_domain - Query Azure Active Directory Domains using SQL

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service. It helps organizations to sign in and access resources in external resources, such as Microsoft Office 365, the Azure portal, and thousands of other SaaS applications. Azure AD Domains represent the various domains that are associated with an Azure AD instance.

## Table Usage Guide

The `azuread_domain` table provides insights into the domains within Azure Active Directory. As an IT administrator, you can explore domain-specific details through this table, including verification status, type of domain, and associated metadata. Utilize it to uncover information about domains, such as their authentication type, availability status, and whether they are the primary domain.

## Examples

### Basic info
Determine the administrative status and verification state of your Azure Active Directory domains, and gain insights into the supported services. This is beneficial for managing access controls and understanding the capabilities of your domains.

```sql+postgres
select
  id,
  is_admin_managed,
  is_verified,
  supported_services
from
  azuread_domain;
```

```sql+sqlite
select
  id,
  is_admin_managed,
  is_verified,
  supported_services
from
  azuread_domain;
```

### List verified domains
Discover the segments that are verified within your Azure Active Directory (AD) domain. This can help ensure the integrity and security of your domain by identifying those that have been validated.

```sql+postgres
select
  id
from
  azuread_domain
where
  is_verified;
```

```sql+sqlite
select
  id
from
  azuread_domain
where
  is_verified;
```