package azuread

import (
	"context"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureAdDirectoryRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role",
		Description: "Represents an Azure Active Directory (Azure AD) directory role.",
		Get: &plugin.GetConfig{
			Hydrate: getAdDirectoryRole,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDirectoryRoles,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the directory role.", Transform: transform.FromMethod("GetId")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description for the directory role.", Transform: transform.FromMethod("GetDescription")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the directory role.", Transform: transform.FromMethod("GetDisplayName")},

			// Other fields
			{Name: "role_template_id", Type: proto.ColumnType_STRING, Description: "The id of the directoryRoleTemplate that this role is based on. The property must be specified when activating a directory role in a tenant with a POST operation. After the directory role has been activated, the property is read only.", Transform: transform.FromMethod("GetRoleTemplateId")},

			// Json fields
			{Name: "member_ids", Type: proto.ColumnType_JSON, Hydrate: getDirectoryRoleMembers, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adDirectoryRoleTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

type ADDirectoryRoleInfo struct {
	models.DirectoryRoleable
}

//// LIST FUNCTION

func listAdDirectoryRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role.listAdDirectoryRoles", "connection_error", err)
		return nil, err
	}

	result, err := client.DirectoryRoles().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdDirectoryRoles", "list_directory_role_error", errObj)
		return nil, errObj
	}

	for _, directoryRole := range result.GetValue() {
		d.StreamListItem(ctx, &ADDirectoryRoleInfo{directoryRole})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDirectoryRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	directoryRoleId := d.EqualsQuals["id"].GetStringValue()
	if directoryRoleId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role.getAdDirectoryRole", "connection_error", err)
		return nil, err
	}

	directoryRole, err := client.DirectoryRolesById(directoryRoleId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdDirectoryRole", "get_directory_role_error", errObj)
		return nil, errObj
	}

	return &ADDirectoryRoleInfo{directoryRole}, nil
}

func getDirectoryRoleMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role.getDirectoryRoleMembers", "connection_error", err)
		return nil, err
	}

	directoryRole := h.Item.(*ADDirectoryRoleInfo)
	directoryRoleID := directoryRole.GetId()

	headers := map[string]string{
		"ConsistencyLevel": "eventual",
	}

	includeCount := true
	requestParameters := &members.MembersRequestBuilderGetQueryParameters{
		Count: &includeCount,
	}

	config := &members.MembersRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParameters,
	}

	memberIds := []*string{}
	members, err := client.DirectoryRolesById(*directoryRoleID).Members().Get(ctx, config)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getDirectoryRoleMembers", "get_directory_role_members_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(members, adapter, models.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("getDirectoryRoleMembers", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		member := pageItem.(models.DirectoryObjectable)
		memberIds = append(memberIds, member.GetId())

		return true
	})
	if err != nil {
		plugin.Logger(ctx).Error("getDirectoryRoleMembers", "paging_error", err)
		return nil, err
	}

	return memberIds, nil
}

//// TRANSFORM FUNCTIONS

func adDirectoryRoleTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADDirectoryRoleInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}
