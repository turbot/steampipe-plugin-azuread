---
title: "Steampipe Table: azuread_identity_provider - Query Azure Active Directory Identity Providers using SQL"
description: "Allows users to query Identity Providers in Azure Active Directory, specifically configuration details and entity information, providing insights into single sign-on setup and identity federation."
---

# Table: azuread_identity_provider - Query Azure Active Directory Identity Providers using SQL

An Azure Active Directory Identity Provider is a service that authenticates users for access to applications and services. It provides a way to configure federation and single sign-on, enabling users to use their existing credentials to sign-in to multiple applications. Azure AD supports a variety of identity providers, including Microsoft Active Directory, Facebook, Google, and more.

## Table Usage Guide

The `azuread_identity_provider` table provides insights into Identity Providers within Azure Active Directory. As a system administrator or security analyst, explore provider-specific details through this table, including provider type, client id, and client secret. Utilize it to uncover information about providers, such as their configuration details and the applications they are linked to.

## Examples

### Basic info
Discover the identities that are registered within your Azure Active Directory. This can assist in managing access and authentication within your organization.

```sql
select
  name,
  id
from
  azuread_identity_provider;
```