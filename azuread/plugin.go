package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-azuread"

// Plugin creates this (azuread) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"azuread_user":  tableAzureAdUser(),
			"azuread_group": tableAzureAdGroup(),
			"azuread_application": tableAzureAdApplication(),
			"azuread_domain": tableAzureAdDomain(),
			"azuread_directory_role": tableAzureAdDirectoryRole(),
			"azuread_service_principal": tableAzureAdServicePrincipal(),
			"azuread_identity_provider": tableAzureAdIdentityProvider(),
		},
	}

	return p
}
