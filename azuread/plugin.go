package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const pluginName = "steampipe-plugin-azuread"

// Plugin creates this (azuread) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "tenant_id",
				Hydrate: getTenant,
			},
		},
		DefaultGetConfig: &plugin.GetConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound"}),
			},
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"azuread_access_review_schedule_definition":            tableAzureAdAccessReviewScheduleDefinition(ctx),
			"azuread_admin_consent_request_policy":                 tableAzureAdAdminConsentRequestPolicy(ctx),
			"azuread_application_app_role_assigned_to":             tableAzureAdApplicationAppRoleAssignment(ctx),
			"azuread_application":                                  tableAzureAdApplication(ctx),
			"azuread_authentication_method_policy":                 tableAzureAdAuthenticationMethodPolicy(ctx),
			"azuread_authorization_policy":                         tableAzureAdAuthorizationPolicy(ctx),
			"azuread_conditional_access_named_location":            tableAzureAdConditionalAccessNamedLocation(ctx),
			"azuread_conditional_access_policy":                    tableAzureAdConditionalAccessPolicy(ctx),
			"azuread_cross_tenant_access_policy":                   tableAzureAdCrossTenantAccessPolicy(ctx),
			"azuread_device":                                       tableAzureAdDevice(ctx),
			"azuread_directory_audit_report":                       tableAzureAdDirectoryAuditReport(ctx),
			"azuread_directory_role_assignment":                    tableAzureAdDirectoryRoleAssignment(ctx),
			"azuread_directory_role_definition":                    tableAzureAdDirectoryRoleDefinition(ctx),
			"azuread_directory_role_eligibility_schedule_instance": tableAzureAdDirectoryRoleEligibilityScheduleInstance(ctx),
			"azuread_directory_role_template":                      tableAzureAdDirectoryRoleTemplate(ctx),
			"azuread_directory_role":                               tableAzureAdDirectoryRole(ctx),
			"azuread_directory_setting":                            tableAzureAdDirectorySetting(ctx),
			"azuread_domain":                                       tableAzureAdDomain(ctx),
			"azuread_external_identity_policy":                     tableAzureAdExternalIdentityPolicy(ctx),
			"azuread_group_app_role_assignment":                    tableAzureAdGroupAppRoleAssignment(ctx),
			"azuread_group":                                        tableAzureAdGroup(ctx),
			"azuread_identity_provider":                            tableAzureAdIdentityProvider(ctx),
			"azuread_long_running_operation":                       tableAzureAdLongRunningOperation(ctx),
			"azuread_security_defaults_policy":                     tableAzureAdSecurityDefaultsPolicy(ctx),
			"azuread_service_principal_app_role_assigned_to":       tableAzureAdServicePrincipalAppRoleAssignedTo(ctx),
			"azuread_service_principal_app_role_assignment":        tableAzureAdServicePrincipalAppRoleAssignment(ctx),
			"azuread_service_principal":                            tableAzureAdServicePrincipal(ctx),
			"azuread_sign_in_report":                               tableAzureAdSignInReport(ctx),
			"azuread_unified_role_assignment_schedule":             tableAzureAdUnifiedRoleAssignmentSchedule(ctx),
			"azuread_user_app_role_assignment":                     tableAzureAdUserAppRoleAssignment(ctx),
			"azuread_user_authentication_method":                   tableAzureAdUserAuthenticationMethod(ctx),
			"azuread_user_registration_details_report":             tableAzureAdUserRegistrationDetailsReport(ctx),
			"azuread_user":                                         tableAzureAdUser(ctx),
		},
	}

	return p
}
