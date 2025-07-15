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

func tableAzureAdDirectoryRoleDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_definition",
		Description: "Represents the role definitions for Azure AD directory resources.",
		Get: &plugin.GetConfig{
			Hydrate: getAdDirectoryRoleDefinition,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDirectoryRoleDefinitions,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the role definition.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the role definition.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description for the role definition.", Transform: transform.FromMethod("GetDescription")},
			{Name: "template_id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the role template.", Transform: transform.FromMethod("GetTemplateId")},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "The version of the role definition.", Transform: transform.FromMethod("GetVersion")},
			{Name: "is_built_in", Type: proto.ColumnType_BOOL, Description: "Flag indicating if the role definition is built in.", Transform: transform.FromMethod("GetIsBuiltIn")},
			{Name: "is_enabled", Type: proto.ColumnType_BOOL, Description: "Flag indicating whether the role is enabled for assignment.", Transform: transform.FromMethod("GetIsEnabled")},

			// JSON fields
			{Name: "resource_scopes", Type: proto.ColumnType_JSON, Description: "List of scopes that the role definition applies to.", Transform: transform.FromMethod("GetResourceScopes")},
			{Name: "role_permissions", Type: proto.ColumnType_JSON, Description: "List of permissions included in this role.", Transform: transform.FromMethod("GetRolePermissions")},
			{Name: "inherits_permissions_from", Type: proto.ColumnType_JSON, Description: "Read-only collection of role definitions that the given role definition inherits from.", Transform: transform.FromMethod("GetInheritsPermissionsFrom")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adDirectoryRoleDefinitionTitle)},
		}),
	}
}

//// LIST FUNCTION

func listAdDirectoryRoleDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_definition.listAdDirectoryRoleDefinitions", "connection_error", err)
		return nil, err
	}

	// List operations
	result, err := client.RoleManagement().Directory().RoleDefinitions().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdDirectoryRoleDefinitions", "list_directory_role_definition_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.UnifiedRoleDefinitionable](result, adapter, models.CreateUnifiedRoleDefinitionCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleDefinitions", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.UnifiedRoleDefinitionable) bool {
		d.StreamListItem(ctx, &ADDirectoryRoleDefinitionInfo{pageItem})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleDefinitions", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDirectoryRoleDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	roleDefinitionId := d.EqualsQuals["id"].GetStringValue()
	if roleDefinitionId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_definition.getAdDirectoryRoleDefinition", "connection_error", err)
		return nil, err
	}

	roleDefinition, err := client.RoleManagement().Directory().RoleDefinitions().ByUnifiedRoleDefinitionId(roleDefinitionId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdDirectoryRoleDefinition", "get_directory_role_definition_error", errObj)
		return nil, errObj
	}

	return &ADDirectoryRoleDefinitionInfo{roleDefinition}, nil
}

//// TRANSFORM FUNCTIONS

func adDirectoryRoleDefinitionTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADDirectoryRoleDefinitionInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}
