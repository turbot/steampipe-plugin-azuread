package azuread

import (
	"context"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdDirectorySetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_setting",
		Description: "Represents the configurations that can be used to customize the tenant-wide and object-specific restrictions and allowed behavior",
		Get: &plugin.GetConfig{
			Hydrate:    getAdDirectorySetting,
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDirectorySetting,
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Display name of this group of settings, which comes from the associated template.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for these settings.", Transform: transform.FromMethod("GetId")},
			{Name: "template_id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the template used to create this group of settings.", Transform: transform.FromMethod("GetTemplateId")},

			// JSON fields
			{Name: "values", Type: proto.ColumnType_JSON, Description: "Collection of name-value pairs corresponding to the name and defaultValue properties in the referenced object.", Transform: transform.FromMethod("DirectorySettingValues")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetDisplayName")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdDirectorySetting(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_setting.listAdDirectorySetting", "connection_error", err)
		return nil, err
	}

	result, err := client.GroupSettings().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdDirectorySetting", "list_directory_setting_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateGroupSettingCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectorySetting", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		setting := pageItem.(models.GroupSettingable)

		d.StreamListItem(ctx, &ADDirectorySettingInfo{setting})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.QueryStatus.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectorySetting", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDirectorySetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	directorySettingID := d.KeyColumnQuals["id"].GetStringValue()
	if directorySettingID == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_setting.getAdDirectorySetting", "connection_error", err)
		return nil, err
	}

	setting, err := client.GroupSettingsById(directorySettingID).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_directory_setting.getAdDirectorySetting", "get_directory_setting_error", errObj)
		return nil, errObj
	}
	return &ADDirectorySettingInfo{setting}, nil
}
