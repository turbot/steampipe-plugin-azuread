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
+----------------+--------------------------------------------+-----------+
| display_name   | user_principal_name                        | user_type |
+----------------+--------------------------------------------+-----------+
| Dwight Schrute | dwight@contoso.onmicrosoft.com             | Member    |
| Jim Halpert    | jim_gmail.com#EXT#@contoso.onmicrosoft.com | Guest     |
| Pam Beesly     | pam_beesly@contoso.onmicrosoft.com         | Member    |
| Michael Scott  | michael@contoso.onmicrosoft.com            | Member    |
+----------------+--------------------------------------------+-----------+
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
| Permissions | Grant the following permissions to your user: <br /><li> `Application.Read.All` </li><li> `AuditLog.Read.All` </li><li> `Directory.Read.All` </li><li> `Domain.Read.All` </li><li> `Group.Read.All` </li><li> `IdentityProvider.Read.All` </li><li> `Policy.Read.All` </li><li> `User.Read.All` </li>                                                                                                                                                            |
| Radius      | Each connection represents a single Azure Tenant.                                                                                                                                                                       |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/azuread.spc`).<br />2. Credentials specified in [environment variables](#credentials-from-environment-variables) e.g. `AZURE_TENANT_ID`. |

### Configuration

Installing the latest azuread plugin will create a config file (~/.steampipe/config/azuread.spc) with a single connection named azuread:

```hcl
connection "azuread" {
  plugin = "azuread"

  # Defaults to "AZUREPUBLICCLOUD". Valid environments are "AZUREPUBLICCLOUD", "AZURECHINACLOUD" and "AZUREUSGOVERNMENTCLOUD"
  # environment = "AZUREPUBLICCLOUD"

  # You can connect to Azure using one of options below:

  # Use client secret authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#option-2-create-a-new-application-secret)
  # tenant_id     = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id     = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # client_secret = "ZZZZZZZZZZZZZZZZZZZZZZZZ"

  # Use client certificate authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#option-1-upload-a-certificate)
  # required options:
  # tenant_id             = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id             = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # certificate_path      = "~/home/azure_cert.pem"
  # certificate_password  = "notreal~pwd"

  # Use a managed identity (https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview)
  # This method is useful with Azure virtual machines
  # tenant_id  = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id  = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # enable_msi = true
  # msi_endpoint = "http://169.254.169.254/metadata/identity/oauth2/token"

  # If no credentials are specified, the plugin will use Azure CLI authentication
}
```

By default, all options are commented out in the default connection, thus Steampipe will resolve your credentials using the same order as mentioned in [Credentials](#credentials). This provides a quick way to get started with Steampipe, but you will probably want to customize your experience using configuration options for querying multiple tenants, [configuring credentials](#configuring-active-directory-credentials) from your Azure CLI, Client Certificate, etc.

## Multi-Tenant Connections

You may create multiple azuread connections:

```hcl
connection "azuread_all" {
  type        = "aggregator"
  plugin      = "azuread"
  connections = ["azuread_*"]
}

connection "azuread_ten_1" {
  plugin        = "azuread"
  tenant_id     = "crfsd708-7da0-4cea-abeb-0a4c334d0f90"
  client_id     = "ea4v6490-c9b5-41db-a942-40a3eaba7053"
  client_secret = "oIb8Q~2hrXFNvWkvtKRtSriRg-kAM3CWEn0g0aGn"
}

connection "azuread_ten_2" {
  plugin        = "azuread"
  tenant_id     = "crfsd708-7da0-4cea-abeb-0a4c334d0f80"
  client_id     = "ea4v6490-c9b5-41db-a942-40a3eaba7088"
  client_secret = "oIb8Q~2hrXFNvWkvtKRtSriRg-kAM3NEEn0g0aGn"
}

connection "azuread_ten_3" {
  plugin        = "azuread"
  tenant_id     = "crfsd708-7da0-4cea-abeb-0a4c334d0f70"
  client_id     = "ea4v6490-c9b5-41db-a942-40a3eaba7000"
  client_secret = "oIb8Q~2hrXFNvWkvtKRtSriRg-kAM3VREn0g0aGn"
}
```

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from azuread_ten_1.azuread_user
```

Alternatively, you can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/using-steampipe/managing-connections#setting-the-search-path):

```sql
select * from azuread_user
```

You can create multi-tenant connections by using an [**aggregator** connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators). Aggregators allow you to query data from multiple connections for a plugin as if they are a single connection:

```hcl
connection "azuread_all" {
  plugin      = "azuread"
  type        = "aggregator"
  connections = ["azuread_ten_1", "azuread_ten_2", "azuread_ten_3"]
}
```

Querying tables from this connection will return results from the `azuread_ten_1`, `azuread_ten_2`, and `azuread_ten_3` connections:

```sql
select * from azuread_all.azuread_user
```

Steampipe supports the `*` wildcard in the connection names. For example, to aggregate all the Azuread plugin connections whose names begin with `azuread_`:

```hcl
connection "azuread_all" {
  type        = "aggregator"
  plugin      = "azuread"
  connections = ["azuread_*"]
}
```

## Configuring Azure Active Directory Credentials

The Azure AD plugin support multiple formats and authentication mechanisms, and they are tried in the below order:

1. [Client Secret Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-saml-bearer-assertion#prerequisites) if set; otherwise
2. [Client Certificate Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-microsoft-identity-platform) if set; otherwise
3. Azure [Managed System Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-managed-identities-work-vm#system-assigned-managed-identity) (useful with virtual machines) if set; otherwise
4. If no credentials are supplied, then the [az cli](https://docs.microsoft.com/en-us/cli/azure/) credentials are used

### Client Secret Credentials

You may specify the tenant ID, client ID, and client secret to authenticate:

- `tenant_id`: Specify the tenant to authenticate with.
- `client_id`: Specify the app client ID to use.
- `client_secret`: Specify the app secret to use.

```hcl
  connection "azuread_via_sp_secret" {
    plugin        = "azuread"
    tenant_id     = "00000000-0000-0000-0000-000000000000"
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "my plaintext password"
  }
```

### Client Certificate Credentials

You may specify the tenant ID, client ID, certificate path, and certificate password to authenticate:

- `tenant_id`: Specify the tenant to authenticate with.
- `client_id`: Specify the app client ID to use.
- `certificate_path`: Specify the certificate path to use.
- `certificate_password`: Specify the certificate password to use.

```hcl
  connection "azuread_via_sp_cert" {
    plugin               = "azuread"
    tenant_id            = "00000000-0000-0000-0000-000000000000"
    client_id            = "00000000-0000-0000-0000-000000000000"
    certificate_path     = "path/to/file.pem"
    certificate_password = "my plaintext password"
  }
```

### Azure Managed Identity

Steampipe works with managed identities (formerly known as Managed Service Identity), provided it is running in Azure, e.g., on a VM. All configuration is handled by Azure. See [Azure Managed Identities](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview) for more details.

- `enable_msi`: Specify `true` to use managed identity credentials.
- `tenant_id`: Specify the tenant to authenticate with.
- `client_id`: Specify the app client ID of managed identity to use.
- `msi_endpoint`: Specify the MSI endpoint to connect to, otherwise use the default Azure Instance Metadata Service (IMDS) endpoint.

```hcl
connection "azure_msi" {
  plugin       = "azuread"
  tenant_id    = "00000000-0000-0000-0000-000000000000"
  client_id    = "00000000-0000-0000-0000-000000000000"
  enable_msi   = true
  msi_endpoint = "http://169.254.169.254/metadata/identity/oauth2/token"
}
```

### Azure CLI

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
export AZURE_ENVIRONMENT="AZUREPUBLICCLOUD" # Defaults to "AZUREPUBLICCLOUD". Valid environments are "AZUREPUBLICCLOUD", "AZURECHINACLOUD" and "AZUREUSGOVERNMENTCLOUD"
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

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-azuread
- Community: [Slack Channel](https://steampipe.io/community/join)
