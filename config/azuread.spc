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
  # client_secret = "ZZZZZZZZZZZZZZZZZZZZZZZZ"

  # 3. MSI authentication (if enabled) using the Azure Metadata Service is then attempted
  # Useful for virtual machine hosted in azure
  # required options:
  # enable_msi = true

  # 4. // Azure CLI authentication (if enabled) is attempted last
}
