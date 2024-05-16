package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdUserAppRoleAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user_app_role_assignment",
		Description: "Represents an application role assigned to a user. Also includes application role assignments granted to groups that the user is a direct member of.",
		Get: &plugin.GetConfig{
			Hydrate: getAdUserAppRoleAssignment,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "user_id", Require: plugin.Required},
				{Name: "id", Require: plugin.Required},
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAdUserAppRoleAssignments,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "user_id", Require: plugin.Required},

				// Other fields for filtering OData
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "principal_display_name", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier for the appRoleAssignment key.", Transform: transform.FromMethod("GetId")},
			{Name: "app_role_id", Type: proto.ColumnType_STRING, Description: "The identifier (id) for the app role which is assigned to the principal. This app role must be exposed in the appRoles property on the resource application's service principal (resourceId). If the resource application has not declared any app roles, a default app role ID of 00000000-0000-0000-0000-000000000000 can be specified to signal that the principal is assigned to the resource app without any specific app roles.", Transform: transform.FromMethod("GetAppRoleId")},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "The unique identifier (id) for the resource service principal for which the assignment is made.", Transform: transform.FromMethod("GetResourceId")},
			{Name: "resource_display_name", Type: proto.ColumnType_STRING, Description: "The display name of the resource app's service principal to which the assignment is made.", Transform: transform.FromMethod("GetResourceDisplayName")},

			// Other fields
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "deleted_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the app role assignment was deleted. Always null for an appRoleAssignment object that hasn't been deleted.", Transform: transform.FromMethod("GetDeletedDateTime")},

			{Name: "principal_id", Type: proto.ColumnType_STRING, Description: "The unique identifier (id) for the user, security group, or service principal being granted the app role.", Transform: transform.FromMethod("GetPrincipalId")},
			{Name: "principal_display_name", Type: proto.ColumnType_STRING, Description: "The display name of the user, group, or service principal that was granted the app role assignment.", Transform: transform.FromMethod("GetPrincipalDisplayName")},
			{Name: "principal_type", Type: proto.ColumnType_STRING, Description: "The type of the assigned principal. This can either be User, Group, or ServicePrincipal.", Transform: transform.FromMethod("GetPrincipalType")},

			// Standard columns
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "The identifier (id) of the user principal.", Transform: transform.From(adUserAppRoleAssignmentUserId)},
		},
	}
}

//// LIST FUNCTION

func listAdUserAppRoleAssignments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	userId := d.EqualsQuals["user_id"].GetStringValue()
	if userId == "" {
		return nil, nil
	}

	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_user_app_role_assignment.listAdUserAppRoleAssignments", "connection_error", err)
		return nil, err
	}

	// List operations
	headers := &abstractions.RequestHeaders{}
	headers.Add("ConsistencyLevel", "eventual")

	input := &users.ItemAppRoleAssignmentsRequestBuilderGetQueryParameters{
		Top:   Int32(999),
		Count: Bool(true),
	}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	// unexpected status 400 with OData error: Request_UnsupportedQuery: Invalid page size specified: '1000'. Must be between 1 and 999 inclusive.
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit < 999 {
			l := int32(*limit)
			input.Top = Int32(l)
		}
	}

	// Apply optional filters
	filters := buildAdUserAppRoleAssignmentQueryFilter(d.Quals)
	if len(filters) > 0 {
		filterString := strings.Join(filters, " and ")
		input.Filter = &filterString
	}

	options := &users.ItemAppRoleAssignmentsRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: input,
	}

	result, err := client.Users().ByUserId(userId).AppRoleAssignments().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdUserAppRoleAssignments", "list_user_app_role_assignment_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.AppRoleAssignmentable](result, adapter, models.CreateAppRoleAssignmentCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdUserAppRoleAssignments", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.AppRoleAssignmentable) bool {
		d.StreamListItem(ctx, &ADUserAppRoleAssignmentInfo{pageItem, &userId})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})

	if err != nil {
		plugin.Logger(ctx).Error("listAdUserAppRoleAssignments", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdUserAppRoleAssignment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	userId := d.EqualsQuals["user_id"].GetStringValue()
	appRoleAssignmentId := d.EqualsQuals["id"].GetStringValue()
	if userId == "" || appRoleAssignmentId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_user_app_role_assignment.getAdUserAppRoleAssignment", "connection_error", err)
		return nil, err
	}

	options := &users.ItemAppRoleAssignmentsAppRoleAssignmentItemRequestBuilderGetRequestConfiguration{}

	appRoleAssignment, err := client.Users().ByUserId(userId).AppRoleAssignments().ByAppRoleAssignmentId(appRoleAssignmentId).Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdUserAppRoleAssignment", "get_user_app_role_assignment_error", errObj)
		return nil, errObj
	}

	return &ADUserAppRoleAssignmentInfo{appRoleAssignment, &userId}, nil
}

//// TRANSFORM FUNCTIONS

func adUserAppRoleAssignmentUserId(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADUserAppRoleAssignmentInfo)
	if data == nil {
		return nil, nil
	}

	return data.UserId, nil
}

func buildAdUserAppRoleAssignmentQueryFilter(quals plugin.KeyColumnQualMap) []string {
	filters := []string{}

	filterTypes := map[string]string{
		"resource_id":            "guid",
		"principal_display_name": "string",
	}

	for k, v := range quals {
		if filterType, ok := filterTypes[k]; ok {
			for _, q := range v.Quals {
				switch filterType {
				case "guid":
					filters = append(filters, fmt.Sprintf("%s eq %s", strcase.ToCamel(k), q.Value.GetStringValue()))
				case "string":
					filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(k), q.Value.GetStringValue()))
				}
			}
		}
	}

	return filters
}
