package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

const pluginName = "steampipe-plugin-azuread"

// Plugin creates this (azuread) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Request_ResourceNotFound"}),
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"azuread_application":               tableAzureAdApplication(),
			"azuread_conditional_access_policy": tableAzureAdConditionalAccessPolicy(),
			"azuread_directory_role":            tableAzureAdDirectoryRole(),
			"azuread_domain":                    tableAzureAdDomain(),
			"azuread_group":                     tableAzureAdGroup(),
			"azuread_identity_provider":         tableAzureAdIdentityProvider(),
			"azuread_service_principal":         tableAzureAdServicePrincipal(),
			"azuread_sign_in_report":            tableAzureAdSignInReport(),
			"azuread_user":                      tableAzureAdUser(),
			"azuread_identity_provider_test":    tableAzureAdIdentityProviderTest(),
			// "azuread_application_test":               tableAzureAdApplicationTest(),
			// "azuread_conditional_access_policy_test": tableAzureAdConditionalAccessPolicyTest(),
			// "azuread_directory_role_test":            tableAzureAdDirectoryRoleTest(),
			// "azuread_domain_test":                    tableAzureAdDomainTest(),
			// "azuread_group_test":                     tableAzureAdGroupTest(),
			// "azuread_service_principal_test":         tableAzureAdServicePrincipalTest(),
			// "azuread_sign_in_report_test":            tableAzureAdSignInReportTest(),
			// "azuread_user_test":                      tableAzureAdUserTest(),
		},
	}

	return p
}
