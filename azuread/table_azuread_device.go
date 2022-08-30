package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/microsoftgraph/msgraph-sdk-go/devices"
	"github.com/microsoftgraph/msgraph-sdk-go/devices/item"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_device",
		Description: "Represents an Azure AD device.",
		Get: &plugin.GetConfig{
			Hydrate:    getDevice,
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listDevices,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "display_name", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{

			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the device. Inherited from directoryObject.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The name displayed for the device.", Transform: transform.FromMethod("GetDisplayName")},

			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "True if the account is enabled; otherwise, false.", Transform: transform.FromMethod("GetAccountEnabled")},
			{Name: "approximate_last_sign_in_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp type represents date and time information using ISO 8601 format and is always in UTC time.", Transform: transform.FromMethod("GetApproximateLastSignInDateTime")},
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "Unique identifier set by Azure Device Registration Service at the time of registration.", Transform: transform.FromMethod("GetDeviceId")},
			{Name: "extension_attributes", Type: proto.ColumnType_JSON, Description: "Contains extension attributes 1-15 for the device. The individual extension attributes are not selectable. These properties are mastered in cloud and can be set during creation or update of a device object in Azure AD.", Transform: transform.FromMethod("GetExtensions")},
			{Name: "is_compliant", Type: proto.ColumnType_BOOL, Description: "True if the device is compliant; otherwise, false.", Transform: transform.FromMethod("GetIsCompliant")},
			{Name: "is_managed", Type: proto.ColumnType_BOOL, Description: "True if the device is managed; otherwise, false.", Transform: transform.FromMethod("GetIsManaged")},
			{Name: "mdm_app_id", Type: proto.ColumnType_STRING, Description: "Application identifier used to register device into MDM.", Transform: transform.FromMethod("GetMdmAppId")},
			{Name: "member_of", Type: proto.ColumnType_JSON, Description: "A list the groups and directory roles that the device is a direct member of.", Transform: transform.FromMethod("DeviceMemberOf")},
			{Name: "operating_system", Type: proto.ColumnType_STRING, Description: "The type of operating system on the device.", Transform: transform.FromMethod("GetOperatingSystem")},
			{Name: "operating_system_version", Type: proto.ColumnType_STRING, Description: "The version of the operating system on the device.", Transform: transform.FromMethod("GetOperatingSystemVersion")},
			{Name: "profile_type", Type: proto.ColumnType_STRING, Description: "A string value that can be used to classify device types.", Transform: transform.FromMethod("GetProfileType")},
			{Name: "trust_type", Type: proto.ColumnType_STRING, Description: "Type of trust for the joined device. Read-only. Possible values: Workplace (indicates bring your own personal devices), AzureAd (Cloud only joined devices), ServerAd (on-premises domain joined devices joined to Azure AD).", Transform: transform.FromMethod("GetTrustType")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(deviceTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listDevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, _, err := GetGraphClient(ctx, d)

	if err != nil {
		plugin.Logger(ctx).Error("azuread_device.listDevices", "connection_error", err)
		return nil, err
	}

	input := &devices.DevicesRequestBuilderGetQueryParameters{}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit <= 999 {
			l := int32(*limit)
			input.Top = &l
		}
	}

	equalQuals := d.KeyColumnQuals
	quals := d.Quals

	givenColumns := d.QueryContext.Columns
	selectColumns, expandColumns := buildDeviceRequestFields(ctx, givenColumns)

	input.Select = selectColumns
	input.Expand = expandColumns

	var queryFilter string
	filter := buildQueryFilter(equalQuals)
	filter = append(filter, buildBoolNEFilter(quals)...)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = &queryFilter
	} else if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &devices.DevicesRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Devices().GetWithRequestConfigurationAndResponseHandler(options, nil)

	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listDevices", "list_device_error", errObj)
		return nil, errObj
	}

	for _, device := range result.GetValue() {
		d.StreamLeafListItem(ctx, &ADDeviceInfo{device})

		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func buildDeviceRequestFields(ctx context.Context, queryColumns []string) ([]string, []string) {
	var selectColumns, expandColumns []string

	for _, columnName := range queryColumns {
		if columnName == "title" || columnName == "filter" || columnName == "tenant_id" {
			continue
		}

		if columnName == "member_of" {
			expandColumns = append(expandColumns, fmt.Sprintf("%s($select=id,displayName)", strcase.ToLowerCamel(columnName)))
			continue
		}

		selectColumns = append(selectColumns, strcase.ToLowerCamel(columnName))
	}

	return selectColumns, expandColumns
}

func deviceTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADDeviceInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetDeviceId()
	}

	return title, nil
}

func getDevice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_device.getDevice", "connection_error", err)
		return nil, err
	}

	deviceId := d.KeyColumnQuals["id"].GetStringValue()
	if deviceId == "" {
		return nil, nil
	}

	// Check for query context and requests only for queried columns
	givenColumns := d.QueryContext.Columns
	selectColumns, expandColumns := buildDeviceRequestFields(ctx, givenColumns)

	input := &item.DeviceItemRequestBuilderGetQueryParameters{}
	input.Select = selectColumns
	input.Expand = expandColumns

	options := &item.DeviceItemRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	device, err := client.DevicesById(deviceId).GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getDevice", "get_device_error", errObj)
		return nil, errObj
	}

	return &ADDeviceInfo{device}, nil
}
