package azuread

import (
	"context"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdDirectoryRoleTemplate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_template",
		Description: "Represents a directory role template in Azure Active Directory (Azure AD). A directory role template specifies the property values of a directory role.",
		Get: &plugin.GetConfig{
			Hydrate: getAdDirectoryRoleTemplate,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDirectoryRoleTemplates,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the directory role template.", Transform: transform.FromMethod("GetId")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description to set for the directory role.", Transform: transform.FromMethod("GetDescription")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name to set for the directory role.", Transform: transform.FromMethod("GetDisplayName")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adDirectoryRoleTemplateTitle)},
		}),
	}
}

//// LIST FUNCTION

func listAdDirectoryRoleTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_template.listAdDirectoryRoleTemplates", "connection_error", err)
		return nil, err
	}

	result, err := client.DirectoryRoleTemplates().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdDirectoryRoleTemplates", "list_directory_role_template_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.DirectoryRoleTemplateable](result, adapter, models.CreateDirectoryRoleTemplateCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleTemplates", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.DirectoryRoleTemplateable) bool {
		d.StreamListItem(ctx, &ADDirectoryRoleTemplateInfo{pageItem})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleTemplates", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDirectoryRoleTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	directoryRoleTemplateId := d.EqualsQuals["id"].GetStringValue()
	if directoryRoleTemplateId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_template.getAdDirectoryRoleTemplate", "connection_error", err)
		return nil, err
	}

	directoryRoleTemplate, err := client.DirectoryRoleTemplates().ByDirectoryRoleTemplateId(directoryRoleTemplateId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdDirectoryRoleTemplate", "get_directory_role_template_error", errObj)
		return nil, errObj
	}

	return &ADDirectoryRoleTemplateInfo{directoryRoleTemplate}, nil
}

//// TRANSFORM FUNCTIONS

func adDirectoryRoleTemplateTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADDirectoryRoleTemplateInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}
