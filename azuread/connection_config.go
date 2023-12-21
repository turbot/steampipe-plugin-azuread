package azuread

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type azureADConfig struct {
	TenantID            *string `hcl:"tenant_id"`
	ClientID            *string `hcl:"client_id"`
	ClientSecret        *string `hcl:"client_secret"`
	CertificatePath     *string `hcl:"certificate_path"`
	CertificatePassword *string `hcl:"certificate_password"`
	EnableMsi           *bool   `hcl:"enable_msi"`
	MsiEndpoint         *string `hcl:"msi_endpoint"`
	Environment         *string `hcl:"environment"`
}

func ConfigInstance() interface{} {
	return &azureADConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) azureADConfig {
	if connection == nil || connection.Config == nil {
		return azureADConfig{}
	}
	config, _ := connection.Config.(azureADConfig)
	return config
}
