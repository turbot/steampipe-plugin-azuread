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
