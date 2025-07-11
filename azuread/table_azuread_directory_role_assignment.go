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

func tableAzureAdDirectoryRoleAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_assignment",
		Description: "Represents the role assignments for Azure AD resources.",
		Get: &plugin.GetConfig{
			Hydrate: getAdDirectoryRoleAssignment,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDirectoryRoleAssignments,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the role assignment.", Transform: transform.FromMethod("GetId")},
			{Name: "principal_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the principal that's in the scope of the role assignment.", Transform: transform.FromMethod("GetPrincipalId")},
			{Name: "role_definition_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the role definition that's in the scope of the role assignment.", Transform: transform.FromMethod("GetRoleDefinitionId")},
			{Name: "directory_scope_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the directory scope that's in the scope of the role assignment.", Transform: transform.FromMethod("GetDirectoryScopeId")},
			{Name: "app_scope_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the app scope that's in the scope of the role assignment.", Transform: transform.FromMethod("GetAppScopeId")},
			{Name: "condition", Type: proto.ColumnType_STRING, Description: "The condition which describes the circumstances under which the role assignment is valid.", Transform: transform.FromMethod("GetCondition")},

			// JSON fields
			{Name: "app_scope", Type: proto.ColumnType_JSON, Description: "The app scope of the role assignment. Contains information about the application scope.", Transform: transform.FromMethod("DirectoryRoleAssignmentAppScope")},
			{Name: "directory_scope", Type: proto.ColumnType_JSON, Description: "The directory scope of the role assignment. Contains information about the directory scope.", Transform: transform.FromMethod("DirectoryRoleAssignmentDirectoryScope")},
			{Name: "principal", Type: proto.ColumnType_JSON, Description: "The principal (user, group, or service principal) that the role is assigned to.", Transform: transform.FromMethod("DirectoryRoleAssignmentPrincipal")},
			{Name: "role_definition", Type: proto.ColumnType_JSON, Description: "The role definition that describes the role being assigned.", Transform: transform.FromMethod("DirectoryRoleAssignmentRoleDefinition")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adDirectoryRoleAssignmentTitle)},
		}),
	}
}

//// LIST FUNCTION

func listAdDirectoryRoleAssignments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_assignment.listAdDirectoryRoleAssignments", "connection_error", err)
		return nil, err
	}

	// List operations
	result, err := client.RoleManagement().Directory().RoleAssignments().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdDirectoryRoleAssignments", "list_directory_role_assignment_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.UnifiedRoleAssignmentable](result, adapter, models.CreateUnifiedRoleAssignmentCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleAssignments", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.UnifiedRoleAssignmentable) bool {
		d.StreamListItem(ctx, &ADDirectoryRoleAssignmentInfo{pageItem})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleAssignments", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDirectoryRoleAssignment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	roleAssignmentId := d.EqualsQuals["id"].GetStringValue()
	if roleAssignmentId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_assignment.getAdDirectoryRoleAssignment", "connection_error", err)
		return nil, err
	}

	roleAssignment, err := client.RoleManagement().Directory().RoleAssignments().ByUnifiedRoleAssignmentId(roleAssignmentId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdDirectoryRoleAssignment", "get_directory_role_assignment_error", errObj)
		return nil, errObj
	}

	return &ADDirectoryRoleAssignmentInfo{roleAssignment}, nil
}

//// TRANSFORM FUNCTIONS

func adDirectoryRoleAssignmentTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADDirectoryRoleAssignmentInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetId()
	if title == nil {
		title = data.GetPrincipalId()
	}

	return title, nil
}
