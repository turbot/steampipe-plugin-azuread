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

# Azure Active Directory + Steampipe

[Azure Active Directory](https://docs.microsoft.com/en-in/azure/active-directory/fundamentals/active-directory-whatis) is Microsoft’s cloud-based identity and access management service, which helps your employees sign in and access resources in:

- External resources, such as Microsoft 365, the Azure portal, and thousands of other SaaS applications.
- Internal resources, such as apps on your corporate network and intranet, along with any cloud apps developed by your own organization.

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

| Item        | Description                                                                                                                                                                                                             |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Use the `az login` command to setup your [Azure AD Default Connection](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli)                                                                               |
| Permissions | Grant the `Global Reader` permission to your user.                                                                                                                                                                      |
| Radius      | Each connection represents a single Azure Tenant.                                                                                                                                                                       |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/azuread.spc`).<br />2. Credentials specified in [environment variables](#credentials-from-environment-variables) e.g. `AZURE_TENANT_ID`. |

### Configuration

Installing the latest azuread plugin will create a config file (~/.steampipe/config/azuread.spc) with a single connection named azuread:

```hcl
connection "azuread" {
  plugin = "azuread"

  # "Defaults to "AZUREPUBLICCLOUD". Can be one of "AZUREPUBLICCLOUD", "AZURECHINACLOUD", "AZUREGERMANCLOUD" and "AZUREUSGOVERNMENTCLOUD"
  # environment = "AZUREPUBLICCLOUD"

  # You may connect to azure using more than one option
  # 1. For client secret authentication, specify TenantID, ClientID and ClientSecret.
  # required options:
  # tenant_id             = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id             = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # client_secret         = "ZZZZZZZZZZZZZZZZZZZZZZZZ"


  # 2. client certificate authentication, specify TenantID, ClientID and ClientCertData / ClientCertPath.
  # required options:
  # tenant_id             = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id             = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # certificate_path      = "~/home/azure_cert.pem"
  # certificate_password  = "notreal~pwd"
  #

  # 3. MSI authentication (if enabled) using the Azure Metadata Service is then attempted
  # Useful for virtual machine hosted in azure
  # If applicable provide msi endpoint, otherwise default endpoint will be used
  # required options:
  # enable_msi = true
  # msi_endpoint = "http://169.254.169.254/metadata/identity/oauth2/token"

  # 4. Azure CLI authentication (if enabled) is attempted last
}

```

By default, all options are commented out in the default connection, thus Steampipe will resolve your credentials using the same order as mentioned in [Credentials](#credentials). This provides a quick way to get started with Steampipe, but you will probably want to customize your experience using configuration options for querying multiple tenants, [configuring credentials](#configuring-active-directory-credentials) from your Azure CLI, Client Certificate, etc.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-azuread
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)

## Configuring Active Directory Credentials

The Azure AD plugin support multiple formats / authentication mechanisms, and they are tried in the below order:

1. [Client Secret Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-saml-bearer-assertion#prerequisites) if set; otherwise
2. [Client Certificate Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-microsoft-identity-platform) if set; otherwise
3. Azure [Managed System Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-managed-identities-work-vm#system-assigned-managed-identity) useful for virtual machines if set; otherwise
4. If no credentials are supplied, then the [az cli](https://docs.microsoft.com/en-us/cli/azure/#:~:text=The%20Azure%20command%2Dline%20interface,with%20an%20emphasis%20on%20automation.) credentials are used

### Client Secret Credentials: Azure AD Application ID and Secret

- `tenant_id`: Specifies the Tenant to which to authenticate.
- `client_id`: Specifies the app client ID to use.
- `client_secret`: Specifies the app secret to use.

```hcl
  connection "azuread_via_sp_secret" {
    plugin        = "azuread"
    tenant_id     = "00000000-0000-0000-0000-000000000000"
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "my plaintext password"
  }
```

### Client Certificate Credentials: Azure AD Application ID and X.509 Certificate

- `tenant_id`: Specifies the Tenant to which to authenticate.
- `client_id`: Specifies the app client ID to use.
- `certificate_path`: Specifies the certificate path to use.
- `certificate_password`: Specifies the certificate password to use.

```hcl
  connection "azuread_via_sp_cert" {
    plugin               = "azuread"
    tenant_id            = "00000000-0000-0000-0000-000000000000"
    client_id            = "00000000-0000-0000-0000-000000000000"
    certificate_path     = "path/to/file.pem"
    certificate_password = "my plaintext password"
  }
```

### Azure Managed Service Identity: Delegate credential management to the platform.

Requires that code is running in Azure, e.g. on a VM. All configuration is handled by Azure. See [Azure Managed Service Identity](https://docs.microsoft.com/azure/active-directory/msi-overview) for more details.

- `enable_msi`: pass `true` to use msi credentials
- `msi_endpoint`: If applicable provide msi endpoint, otherwise default endpoint(i.e. http://169.254.169.254/metadata/identity/oauth2/token) will be used.

```hcl
connection "azure_msi" {
  plugin       = "azuread"
  enable_msi   = true
  msi_endpoint = "http://169.254.169.254/metadata/identity/oauth2/token"
}
```

### Azure CLI: If no credentials are specified

If no credentials are specified and the SDK environment variables are not set, the plugin will use the active credentials from the `az` cli. You can run `az login` to set up these credentials.

```hcl
connection "azuread" {
  plugin = "azuread"
}
```

### Credentials from Environment Variables

The Azure AD plugin will use the standard Azure environment variables to obtain credentials **only if other arguments (`tenant_id`, `client_id`, `client_secret`, `certificate_path`, etc..) are not specified** in the connection:

```sh
export AZURE_TENANT_ID="00000000-0000-0000-0000-000000000000"
export AZURE_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export AZURE_CLIENT_SECRET="my plaintext secret"
export AZURE_CERTIFICATE_PATH=path/to/file.pem
export AZURE_CERTIFICATE_PASSWORD="my plaintext password"
```

```hcl
connection "azuread" {
  plugin = "azuread"
}
```
