---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/azuread.svg"
brand_color: "#0089D6"
display_name: "Azure Active Directory"
name: "azuread"
description: "Steampipe plugin for querying resource users, groups, applications and more from Azure Active Directory."
og_description: "Query Azure Active Directory with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/azuread-social-graphic.png"
---

# Azure AD + Steampipe

[Azure AD](https://www.office.com/) provides access to data stored across Microsoft 365 services. Custom applications can use the Microsoft Graph API to connect to data and use it in custom applications to enhance organizational productivity.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  display_name,
  user_principal_name,
  user_type
from
  azuread_user;
```

```
+--------------------------------+--------------------------------------------+-----------+
| display_name                   | user_principal_name                        | user_type |
+--------------------------------+--------------------------------------------+-----------+
| Dwight Schrute                 | dwight@contoso.onmicrosoft.com             | Member    |
| Jim Halpert                    | jim_gmail.com#EXT#@contoso.onmicrosoft.com | Guest     |
| Pam Beesly                     | pam_beesly@contoso.onmicrosoft.com         | Member    |
| Michael Scott                  | michael@contoso.onmicrosoft.com            | Member    |
+--------------------------------+--------------------------------------------+-----------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/azuread/tables)**

## Get started

### Install

Download and install the latest Azure Active Directory plugin:

```bash
steampipe plugin install azuread
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Use the `az login` command to setup your [Azure AD Default Connection](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| Permissions | Grant the `Global Reader` permission to your user.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| Radius      | Each connection represents a single Azure Tenant.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| Resolution  | 1. [Client certificate authentication](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-microsoft-identity-platform)<br />2. [Client secret authentication](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-saml-bearer-assertion#prerequisites)<br />3. [Managed System Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-managed-identities-work-vm#system-assigned-managed-identity) useful for virtual machines. <br />4. If no credentials are supplied, then the [az cli](https://docs.microsoft.com/en-us/cli/azure/#:~:text=The%20Azure%20command%2Dline%20interface,with%20an%20emphasis%20on%20automation.) credentials are used |

### Configuration

Installing the latest azuread plugin will create a config file (~/.steampipe/config/azuread.spc) with a single connection named azuread:

```hcl
connection "azuread" {
  plugin = "azuread"

  # You may connect to azure using more than one option
  # 1. client certificate authentication, specify TenantID, ClientID and ClientCertData / ClientCertPath.
  # required options:
  # tenant_id             = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id             = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # certificate_path      = "~/home/azure_cert.pem"
  # certificate_password  = "notreal~pwd"
  #


  # 2. For client secret authentication, specify TenantID, ClientID and ClientSecret.
  # required options:
  # tenant_id             = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id             = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # client_secret         = "ZZZZZZZZZZZZZZZZZZZZZZZZ"

  # 3. MSI authentication (if enabled) using the Azure Metadata Service is then attempted
  # Useful for virtual machine hosted in azure
  # If applicable provide msi endpoint, otherwise default endpoiint will be used
  # required options:
  # enable_msi = true
  # msi_endpoint = "http://169.254.169.254/metadata/identity/oauth2/token"

  # 4. Azure CLI authentication (if enabled) is attempted last
}

```

By default, all options are commented out in the default connection, thus Steampipe will resolve your credentials using the same order as mentioned in [Credentials](#credentials). This provides a quick way to get started with Steampipe, but you will probably want to customize your experience using configuration options for querying multiple tenants, configuring credentials from your Azure CLI, Client Certificate, etc.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-azure
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
