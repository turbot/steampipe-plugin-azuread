package azuread

import (
	"context"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdIntuneManagedDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_intune_managed_device",
		Description: "Represents the devices enrolled in Microsoft Intune for a tenant.",
		Get: &plugin.GetConfig{
			Hydrate: getAdIntuneManagedDevice,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "ResourceNotFound", "NotFound"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdIntuneManagedDevices,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the managed device.", Transform: transform.FromMethod("GetId")},
			{Name: "device_name", Type: proto.ColumnType_STRING, Description: "The name of the device.", Transform: transform.FromMethod("GetDeviceName")},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "The serial number of the device.", Transform: transform.FromMethod("GetSerialNumber")},
			{Name: "azure_ad_device_id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the Microsoft Entra device.", Transform: transform.FromMethod("GetAzureADDeviceId")},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the user associated with the device.", Transform: transform.FromMethod("GetUserId")},
			{Name: "user_principal_name", Type: proto.ColumnType_STRING, Description: "The device user principal name.", Transform: transform.FromMethod("GetUserPrincipalName")},
			{Name: "email_address", Type: proto.ColumnType_STRING, Description: "The email address of the user associated with the device.", Transform: transform.FromMethod("GetEmailAddress")},
			{Name: "user_display_name", Type: proto.ColumnType_STRING, Description: "The display name of the user associated with the device.", Transform: transform.FromMethod("GetUserDisplayName")},
			{Name: "operating_system", Type: proto.ColumnType_STRING, Description: "The operating system of the device, for example Windows, iOS or macOS.", Transform: transform.FromMethod("GetOperatingSystem")},
			{Name: "os_version", Type: proto.ColumnType_STRING, Description: "The operating system version of the device.", Transform: transform.FromMethod("GetOsVersion")},
			{Name: "compliance_state", Type: proto.ColumnType_STRING, Description: "The compliance state of the device. Possible values are: unknown, compliant, noncompliant, conflict, error, inGracePeriod, configManager.", Transform: transform.FromMethod("ComplianceState")},
			{Name: "is_encrypted", Type: proto.ColumnType_BOOL, Description: "Whether the device storage is encrypted.", Transform: transform.FromMethod("GetIsEncrypted")},
			{Name: "management_agent", Type: proto.ColumnType_STRING, Description: "The management channel of the device, for example mdm, easMdm or configurationManagerClientMdm.", Transform: transform.FromMethod("ManagementAgent")},
			{Name: "managed_device_owner_type", Type: proto.ColumnType_STRING, Description: "The ownership of the device. Possible values are: unknown, company, personal.", Transform: transform.FromMethod("ManagedDeviceOwnerType")},
			{Name: "enrolled_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the device was enrolled.", Transform: transform.FromMethod("EnrolledDateTime")},
			{Name: "last_sync_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the device last completed a successful sync with Intune.", Transform: transform.FromMethod("LastSyncDateTime")},
			{Name: "compliance_grace_period_expiration_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the device compliance grace period expires.", Transform: transform.FromMethod("ComplianceGracePeriodExpirationDateTime")},
			{Name: "manufacturer", Type: proto.ColumnType_STRING, Description: "The manufacturer of the device.", Transform: transform.FromMethod("GetManufacturer")},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "The model of the device.", Transform: transform.FromMethod("GetModel")},
			{Name: "jail_broken", Type: proto.ColumnType_STRING, Description: "Whether the device is jail broken or rooted.", Transform: transform.FromMethod("GetJailBroken")},
			{Name: "wifi_mac_address", Type: proto.ColumnType_STRING, Description: "The Wi-Fi MAC address of the device.", Transform: transform.FromMethod("GetWiFiMacAddress")},
			{Name: "ethernet_mac_address", Type: proto.ColumnType_STRING, Description: "The Ethernet MAC address of the device.", Transform: transform.FromMethod("GetEthernetMacAddress")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetDeviceName")},
		}),
	}
}

//// LIST FUNCTION

func listAdIntuneManagedDevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_intune_managed_device.listAdIntuneManagedDevices", "connection_error", err)
		return nil, err
	}

	input := &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{}

	limit := d.QueryContext.Limit
	if limit != nil && *limit > 0 && *limit < 1000 {
		input.Top = Int32(int32(*limit))
	}

	options := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.DeviceManagement().ManagedDevices().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_intune_managed_device.listAdIntuneManagedDevices", "list_managed_device_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[interface{}](result, adapter, models.CreateManagedDeviceCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_intune_managed_device.listAdIntuneManagedDevices", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		if device, ok := pageItem.(models.ManagedDeviceable); ok {
			d.StreamListItem(ctx, &ADIntuneManagedDeviceInfo{device})
		}
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuread_intune_managed_device.listAdIntuneManagedDevices", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdIntuneManagedDevice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	deviceID := d.EqualsQuals["id"].GetStringValue()
	if deviceID == "" {
		return nil, nil
	}

	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_intune_managed_device.getAdIntuneManagedDevice", "connection_error", err)
		return nil, err
	}

	device, err := client.DeviceManagement().ManagedDevices().ByManagedDeviceId(deviceID).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_intune_managed_device.getAdIntuneManagedDevice", "get_managed_device_error", errObj)
		return nil, errObj
	}

	return &ADIntuneManagedDeviceInfo{device}, nil
}
