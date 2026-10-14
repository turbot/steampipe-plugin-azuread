package azuread

import (
	"context"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdIntuneDetectedApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_intune_detected_app",
		Description: "Represents the applications Microsoft Intune has discovered on a managed device.",
		List: &plugin.ListConfig{
			Hydrate: listAdIntuneDetectedApps,
			// A device removed from Intune between reads must not fail a query that names many devices.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "ResourceNotFound", "NotFound"}),
			},
			KeyColumns: plugin.SingleColumn("managed_device_id"),
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "managed_device_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the managed device the application was discovered on.", Transform: transform.FromMethod("ManagedDeviceID")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the detected application.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The name of the detected application.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "The version of the detected application.", Transform: transform.FromMethod("GetVersion")},
			{Name: "publisher", Type: proto.ColumnType_STRING, Description: "The publisher of the detected application.", Transform: transform.FromMethod("GetPublisher")},
			{Name: "platform", Type: proto.ColumnType_STRING, Description: "The platform of the detected application, for example windows, macOS, ios, chromeOS or androidWorkProfile.", Transform: transform.FromMethod("Platform")},
			{Name: "size_in_byte", Type: proto.ColumnType_INT, Description: "The size of the detected application in bytes.", Transform: transform.FromMethod("GetSizeInByte")},
			{Name: "device_count", Type: proto.ColumnType_INT, Description: "The number of devices in the tenant that have installed this application.", Transform: transform.FromMethod("GetDeviceCount")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetDisplayName")},
		}),
	}
}

//// LIST FUNCTION

// listAdIntuneDetectedApps lists one device's apps; the SDK calls it once per value of an IN list.
func listAdIntuneDetectedApps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	deviceID := d.EqualsQuals["managed_device_id"].GetStringValue()
	if deviceID == "" {
		return nil, nil
	}

	// Graph v1.0 has no managedDevice detectedApps relationship; only beta exposes it.
	client, adapter, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_intune_detected_app.listAdIntuneDetectedApps", "connection_error", err)
		return nil, err
	}

	result, err := client.DeviceManagement().ManagedDevices().ByManagedDeviceId(deviceID).DetectedApps().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_intune_detected_app.listAdIntuneDetectedApps", "list_detected_app_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[interface{}](result, adapter, models.CreateDetectedAppCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_intune_detected_app.listAdIntuneDetectedApps", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		if app, ok := pageItem.(models.DetectedAppable); ok {
			d.StreamListItem(ctx, &ADIntuneDetectedAppInfo{app, deviceID})
		}
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_intune_detected_app.listAdIntuneDetectedApps", "paging_error", errObj)
		return nil, errObj
	}

	return nil, nil
}
