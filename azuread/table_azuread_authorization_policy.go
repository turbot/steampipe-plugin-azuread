package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdAuthorizationPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_authorization_policy",
		Description: "Represents a policy that can control Azure Active Directory authorization settings",
		List: &plugin.ListConfig{
			Hydrate: listAdAuthorizationPolicies,
		},

		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Display name for this policy.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the authorization policy.", Transform: transform.FromMethod("GetId")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of this policy.", Transform: transform.FromMethod("GetDescription")},

			// Other fields
			{Name: "allowed_to_sign_up_email_based_subscriptions", Type: proto.ColumnType_BOOL, Description: "Indicates whether users can sign up for email based subscriptions.", Transform: transform.FromMethod("GetAllowedToSignUpEmailBasedSubscriptions")},
			{Name: "allowed_to_use_sspr", Type: proto.ColumnType_BOOL, Description: "Indicates whether the Self-Serve Password Reset feature can be used by users on the tenant.", Transform: transform.FromMethod("GetAllowedToUseSSPR")},
			{Name: "allowed_email_verified_users_to_join_organization", Type: proto.ColumnType_BOOL, Description: "Indicates whether a user can join the tenant by email validation.", Transform: transform.FromMethod("GetAllowEmailVerifiedUsersToJoinOrganization")},
			{Name: "allow_invites_from", Type: proto.ColumnType_STRING, Description: "Indicates who can invite external users to the organization. Possible values are: none, adminsAndGuestInviters, adminsGuestInvitersAndAllMembers, everyone.", Transform: transform.FromMethod("AuthorizationPolicyAllowInvitesFrom")},
			{Name: "block_msol_powershell", Type: proto.ColumnType_BOOL, Description: "To disable the use of MSOL PowerShell set this property to true. This will also disable user-based access to the legacy service endpoint used by MSOL PowerShell. This does not affect Azure AD Connect or Microsoft Graph.", Transform: transform.FromMethod("GetBlockMsolPowerShell")},
			{Name: "guest_user_role_id", Type: proto.ColumnType_STRING, Description: "Represents role templateId for the role that should be granted to guest user.", Transform: transform.FromMethod("GetGuestUserRoleId")},

			// JSON fields
			{Name: "default_user_role_permissions", Type: proto.ColumnType_JSON, Description: "Specifies certain customizable permissions for default user role.", Transform: transform.FromMethod("AuthorizationPolicyDefaultUserRolePermissions")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetId")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdAuthorizationPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_authorization_policy.listAdAuthorizationPolicies", "connection_error", err)
		return nil, err
	}

	result, err := client.Policies().AuthorizationPolicy().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdAuthorizationPolicies", "list_application_error", errObj)
		return nil, errObj
	}
	d.StreamListItem(ctx, &ADAuthorizationPolicyInfo{result})

	return nil, nil
}
