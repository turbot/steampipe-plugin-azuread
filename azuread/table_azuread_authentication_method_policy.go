package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdAuthenticationMethodPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_authentication_method_policy",
		Description: "Represents the authentication methods policy for the Microsoft Entra tenant.",
		List: &plugin.ListConfig{
			Hydrate: listAdAuthenticationMethodPolicy,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The identifier for the authentication methods policy.", Transform: transform.FromMethod("ResourceId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the authentication methods policy.", Transform: transform.FromMethod("DisplayName")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the authentication methods policy.", Transform: transform.FromMethod("Description")},
			{Name: "last_modified_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time of the last update to the policy.", Transform: transform.FromMethod("LastModifiedDateTime")},
			{Name: "policy_migration_state", Type: proto.ColumnType_STRING, Description: "The state of migration of the authentication methods policy from the legacy multifactor authentication and self-service password reset (SSPR) policies.", Transform: transform.FromMethod("PolicyMigrationState")},
			{Name: "policy_version", Type: proto.ColumnType_STRING, Description: "The version of the policy in use.", Transform: transform.FromMethod("PolicyVersion")},
			{Name: "reconfirmation_in_days", Type: proto.ColumnType_INT, Description: "The reconfirmation in days property.", Transform: transform.FromMethod("ReconfirmationInDays")},
			{Name: "registration_enforcement", Type: proto.ColumnType_JSON, Description: "Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods.", Transform: transform.FromMethod("RegistrationEnforcement")},
			{Name: "authentication_method_configurations", Type: proto.ColumnType_JSON, Description: "Represents the settings for each authentication method.", Transform: transform.FromMethod("AuthenticationMethodConfigurations")},
			{Name: "additional_data", Type: proto.ColumnType_JSON, Description: "Additional data associated with the authentication methods policy.", Transform: transform.FromMethod("AdditionalData")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("DisplayName")},
		}),
	}
}

//// LIST FUNCTION

func listAdAuthenticationMethodPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access additional properties like numberMatchingRequiredState
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_authentication_method_policy.listAdAuthenticationMethodPolicy", "connection_error", err)
		return nil, err
	}

	result, err := client.Policies().AuthenticationMethodsPolicy().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdAuthenticationMethodPolicy", "list_authentication_method_policy_error", errObj)
		return nil, errObj
	}

	d.StreamListItem(ctx, &ADAuthenticationMethodPolicyInfo{result})

	return nil, nil
}
