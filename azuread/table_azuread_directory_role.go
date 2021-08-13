package azuread

import (
	"context"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAdDirectoryRole() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role",
		Description: "Represents an Azure Active Directory (Azure AD) directory role",
		Get: &plugin.GetConfig{
			Hydrate:           getAdDirectoryRole,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError,
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDirectoryRoles,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "id", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the directory role.", Transform: transform.FromGo()},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description for the directory role."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the directory role."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for groups."},

			// Other fields
			{Name: "role_template_id", Type: proto.ColumnType_STRING, Description: "The id of the directoryRoleTemplate that this role is based on. The property must be specified when activating a directory role in a tenant with a POST operation. After the directory role has been activated, the property is read only."},

			// // Json fields
			{Name: "member_ids", Type: proto.ColumnType_JSON, Hydrate: getAdDirectoryRoleMembers, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},
			// // Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdDirectoryRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewDirectoryRolesClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	directoryRoles, _, err := client.List(ctx)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, directoryRoles := range *directoryRoles {
		d.StreamListItem(ctx, directoryRoles)
	}

	return nil, err
}

////Hydrate Functions

func getAdDirectoryRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	var directoryRoleId string
	if h.Item != nil {
		directoryRoleId = *h.Item.(msgraph.DirectoryRole).ID
	} else {
		directoryRoleId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if directoryRoleId == "" {
		return nil, nil
	}

	client := msgraph.NewDirectoryRolesClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	directoryRole, _, err := client.Get(ctx, directoryRoleId)
	if err != nil {
		return nil, err
	}
	return *directoryRole, nil
}

func getAdDirectoryRoleMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	var directoryRoleId string
	if h.Item != nil {
		directoryRoleId = *h.Item.(msgraph.DirectoryRole).ID
	} else {
		directoryRoleId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if directoryRoleId == "" {
		return nil, nil
	}

	client := msgraph.NewDirectoryRolesClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	members, _, err := client.ListMembers(ctx, directoryRoleId)
	if err != nil {
		return nil, err
	}
	return members, nil
}
