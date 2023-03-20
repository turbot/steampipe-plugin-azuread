package azuread

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type azureADConfig struct {
	TenantID            *string `cty:"tenant_id"`
	ClientID            *string `cty:"client_id"`
	ClientSecret        *string `cty:"client_secret"`
	CertificatePath     *string `cty:"certificate_path"`
	CertificatePassword *string `cty:"certificate_password"`
	EnableMsi           *bool   `cty:"enable_msi"`
	MsiEndpoint         *string `cty:"msi_endpoint"`
	Environment         *string `cty:"environment"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"tenant_id": {
		Type: schema.TypeString,
	},
	"client_id": {
		Type: schema.TypeString,
	},
	"client_secret": {
		Type: schema.TypeString,
	},
	"certificate_path": {
		Type: schema.TypeString,
	},
	"certificate_password": {
		Type: schema.TypeString,
	},
	"environment": {
		Type: schema.TypeString,
	},
	"enable_msi": {
		Type: schema.TypeBool,
	},
	"msi_endpoint": {
		Type: schema.TypeString,
	},
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
