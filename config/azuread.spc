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
