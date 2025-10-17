package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdDeviceRegistrationPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_device_registration_policy",
		Description: "Represents the Azure Active Directory (Azure AD) device registration policy that manages initial provisioning controls using quota restrictions, additional authentication and authorization checks.",
		List: &plugin.ListConfig{
			Hydrate: listAdDeviceRegistrationPolicy,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The identifier of the device registration policy.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the device registration policy.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the device registration policy.", Transform: transform.FromMethod("GetDescription")},
			{Name: "user_device_quota", Type: proto.ColumnType_INT, Description: "Specifies the maximum number of devices that a user can have within your organization before blocking new device registrations.", Transform: transform.FromMethod("GetUserDeviceQuota")},
			{Name: "multi_factor_auth_configuration", Type: proto.ColumnType_STRING, Description: "Specifies the authentication policy for a user to complete registration using Azure AD Join or Azure AD register within your organization. Possible values are: notRequired, required, unknownFutureValue.", Transform: transform.FromMethod("DeviceRegistrationPolicyMultiFactorAuthConfiguration")},

			// JSON fields
			{Name: "azure_ad_registration", Type: proto.ColumnType_JSON, Description: "Specifies the authorization policy for controlling registration of new devices using Azure AD registered within your organization.", Transform: transform.FromMethod("DeviceRegistrationPolicyAzureADRegistration")},
			{Name: "azure_ad_join", Type: proto.ColumnType_JSON, Description: "Specifies the authorization policy for controlling registration of new devices using Azure AD join within your organization.", Transform: transform.FromMethod("DeviceRegistrationPolicyAzureADJoin")},
			{Name: "local_admin_password", Type: proto.ColumnType_JSON, Description: "Specifies the setting for Local Admin Password Solution (LAPS) within your organization.", Transform: transform.FromMethod("DeviceRegistrationPolicyLocalAdminPassword")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetDisplayName")},
		}),
	}
}

//// LIST FUNCTION

func listAdDeviceRegistrationPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_device_registration_policy.listAdDeviceRegistrationPolicy", "connection_error", err)
		return nil, err
	}

	result, err := client.Policies().DeviceRegistrationPolicy().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdDeviceRegistrationPolicy", "get_device_registration_policy_error", errObj)
		return nil, errObj
	}

	d.StreamListItem(ctx, &ADDeviceRegistrationPolicyInfo{result})

	return nil, nil
}
