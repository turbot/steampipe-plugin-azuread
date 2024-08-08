---
title: "Steampipe Table: azuread_conditional_access_named_location - Query Microsoft Entra Named Locations using SQL"
description: "Allows users to query Microsoft Entra Named Locations, providing information about custom definitions of Named Locations"
---

# Table: azuread_conditional_access_named_location - Query Microsoft Entra Named Locations using SQL

Microsoft Entra Named Locations is a feature in Azure Active Directory (Microsoft Entra) that allows administrators to define custom Named Locations. These Custom named locations can be included in Conditional Access Policies and restrict user access to this specific locations. There are two types of Named Locations - IP based Named locations and Country based Named Locations, the table supports both types.

## Table Usage Guide

The `azuread_conditional_access_named_location` table provides insights into Named Locations within Azure Active Directory (Microsoft Entra). As a security administrator, you can understand policies based on Named Locations better through this table, including display name, type, and detailed location information. Utilize it to uncover information about custom Named Locations, understand Conditional Access policies better, and maintain security and compliance within your organization.

## Examples

### Basic info
Analyze the settings to understand the status and creation date of the Named Locations in your Microsoft Entra Named Locations. This can help you assess the locations elements within your Conditional Access Policy and make necessary adjustments.

```sql+postgres
select
  id,
  display_name,
  type,
  created_date_time,
  modified_date_time
from
  azuread_conditional_access_named_location;
```

```sql+sqlite
select
  id,
  display_name,
  type,
  created_date_time,
  modified_date_time
from
  azuread_conditional_access_named_location;
```

### Detailed information about the Namedl Location definitions
Analyze detailed information about the definition of Named Locations in your Microsoft Entra Named Locations. This can help you understand the locations elements within your Conditional Access Policy and assure the definitions are compliance within your organization policies.

```sql+postgres
select
  id,
  display_name,
  type,
  location_info
from
  azuread_conditional_access_named_location;
```

```sql+sqlite
select
  id,
  display_name,
  type,
  location_info
from
  azuread_conditional_access_named_location;
```