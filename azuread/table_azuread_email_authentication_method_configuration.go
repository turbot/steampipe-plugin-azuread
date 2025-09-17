package azuread

import (
	"context"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdEmailAuthenticationMethodConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_email_authentication_method_configuration",
		Description: "Represents the email OTP authentication method policy for the Microsoft Entra tenant.",
		List: &plugin.ListConfig{
			Hydrate: listAdEmailAuthenticationMethodConfiguration,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The identifier for the authentication method configuration.", Transform: transform.FromMethod("GetId")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the authentication method configuration. Possible values are: enabled, disabled.", Transform: transform.FromMethod("GetState")},
			{Name: "allow_external_id_to_use_email_otp", Type: proto.ColumnType_STRING, Description: "Determines whether email OTP is usable by external users for authentication. Possible values are: default, enabled, disabled.", Transform: transform.FromMethod("GetAllowExternalIdToUseEmailOtp")},
			{Name: "include_targets", Type: proto.ColumnType_JSON, Description: "A collection of users or groups who are enabled to use the authentication method.", Transform: transform.FromMethod("GetIncludeTargets").Transform(transformEmailAuthIncludeTargets)},
			{Name: "exclude_targets", Type: proto.ColumnType_JSON, Description: "A collection of users or groups who are excluded from using the authentication method.", Transform: transform.FromMethod("GetExcludeTargets").Transform(transformEmailAuthExcludeTargets)},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetId")},
		}),
	}
}

//// LIST FUNCTION

func listAdEmailAuthenticationMethodConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_email_authentication_method_configuration.listAdEmailAuthenticationMethodConfiguration", "connection_error", err)
		return nil, err
	}

	result, err := client.Policies().AuthenticationMethodsPolicy().AuthenticationMethodConfigurations().ByAuthenticationMethodConfigurationId("email").Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdEmailAuthenticationMethodConfiguration", "list_email_authentication_method_configuration_error", errObj)
		return nil, errObj
	}

	// Cast the result to the specific email authentication method configuration type
	emailConfig := result.(models.EmailAuthenticationMethodConfigurationable)
	d.StreamListItem(ctx, &ADEmailAuthenticationMethodConfigurationInfo{emailConfig})

	return nil, nil
}
