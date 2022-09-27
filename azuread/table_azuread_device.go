package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/devices"
	"github.com/microsoftgraph/msgraph-sdk-go/devices/item"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
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
			Hydrate:    getAdDevice,
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDevices,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "display_name", Require: plugin.Optional},
				{Name: "account_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "operating_system", Require: plugin.Optional},
				{Name: "operating_system_version", Require: plugin.Optional},
				{Name: "profile_type", Require: plugin.Optional},
				{Name: "trust_type", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{

			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the device. Inherited from directoryObject.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The name displayed for the device.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "True if the account is enabled; otherwise, false.", Transform: transform.FromMethod("GetAccountEnabled")},
			{Name: "device_id", Type: proto.ColumnType_STRING, Description: "Unique identifier set by Azure Device Registration Service at the time of registration.", Transform: transform.FromMethod("GetDeviceId")},
			{Name: "approximate_last_sign_in_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp type represents date and time information using ISO 8601 format and is always in UTC time.", Transform: transform.FromMethod("GetApproximateLastSignInDateTime")},

			// Other fields
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for resources."},
			{Name: "is_compliant", Type: proto.ColumnType_BOOL, Description: "True if the device is compliant; otherwise, false.", Transform: transform.FromMethod("GetIsCompliant")},
			{Name: "is_managed", Type: proto.ColumnType_BOOL, Description: "True if the device is managed; otherwise, false.", Transform: transform.FromMethod("GetIsManaged")},
			{Name: "mdm_app_id", Type: proto.ColumnType_STRING, Description: "Application identifier used to register device into MDM.", Transform: transform.FromMethod("GetMdmAppId")},
			{Name: "operating_system", Type: proto.ColumnType_STRING, Description: "The type of operating system on the device.", Transform: transform.FromMethod("GetOperatingSystem")},
			{Name: "operating_system_version", Type: proto.ColumnType_STRING, Description: "The version of the operating system on the device.", Transform: transform.FromMethod("GetOperatingSystemVersion")},
			{Name: "profile_type", Type: proto.ColumnType_STRING, Description: "A string value that can be used to classify device types.", Transform: transform.FromMethod("GetProfileType")},
			{Name: "trust_type", Type: proto.ColumnType_STRING, Description: "Type of trust for the joined device. Possible values: Workplace (indicates bring your own personal devices), AzureAd (Cloud only joined devices), ServerAd (on-premises domain joined devices joined to Azure AD).", Transform: transform.FromMethod("GetTrustType")},

			// JSON fields
			{Name: "extension_attributes", Type: proto.ColumnType_JSON, Description: "Contains extension attributes 1-15 for the device. The individual extension attributes are not selectable. These properties are mastered in cloud and can be set during creation or update of a device object in Azure AD.", Transform: transform.FromMethod("GetExtensions")},
			{Name: "member_of", Type: proto.ColumnType_JSON, Description: "A list the groups and directory roles that the device is a direct member of.", Transform: transform.FromMethod("DeviceMemberOf")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adDeviceTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdDevices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, adapter, err := GetGraphClient(ctx, d)

	if err != nil {
		plugin.Logger(ctx).Error("azuread_device.listAdDevices", "connection_error", err)
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
	filter := buildDeviceQueryFilter(equalQuals)
	filter = append(filter, buildDeviceBoolNEFilter(quals)...)

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

	result, err := client.Devices().Get(ctx, options)

	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_device.listAdDevices", "list_device_error", errObj)
		return nil, errObj
	}

	if result.GetOdataNextLink() != nil {

		pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateDeviceCollectionResponseFromDiscriminatorValue)
		if err != nil {
			plugin.Logger(ctx).Error("azuread_device.listAdDevices", "create_iterator_instance_error", err)
			return nil, err
		}

		err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
			device := pageItem.(models.Deviceable)

			d.StreamListItem(ctx, &ADDeviceInfo{device})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			return d.QueryStatus.RowsRemaining(ctx) != 0
		})

		if err != nil {
			plugin.Logger(ctx).Error("azuread_device.listAdDevices", "paging_error", err)
			return nil, err
		}
	}

	for _, device := range result.GetValue() {
		d.StreamListItem(ctx, &ADDeviceInfo{device})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDevice(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_device.getAdDevice", "connection_error", err)
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

	device, err := client.DevicesById(deviceId).Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdDevice", "get_device_error", errObj)
		return nil, errObj
	}

	return &ADDeviceInfo{device}, nil
}

//// TRANSFORM FUNCTIONS

func adDeviceTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
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

func buildDeviceQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name":             "string",
		"id":                       "string",
		"operating_system":         "string",
		"operating_system_version": "string",
		"profile_type":             "string",
		"trust_type":               "string",
		"account_enabled":          "bool",
	}

	for qual, qualType := range filterQuals {
		switch qualType {
		case "string":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
			}
		case "bool":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), equalQuals[qual].GetBoolValue()))
			}
		}
	}

	return filters
}

func buildDeviceBoolNEFilter(quals plugin.KeyColumnQualMap) []string {
	filters := []string{}

	filterQuals := []string{
		"account_enabled",
	}

	for _, qual := range filterQuals {
		if quals[qual] != nil {
			for _, q := range quals[qual].Quals {
				value := q.Value.GetBoolValue()
				if q.Operator == "<>" {
					filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), !value))
					break
				}
			}
		}
	}

	return filters
}
