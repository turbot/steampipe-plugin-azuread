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

func tableAzureAdDirectoryRoleEligibilityScheduleInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_eligibility_schedule_instance",
		Description: "Represents the schedule instances for role eligibility operations on Azure AD resources.",
		Get: &plugin.GetConfig{
			Hydrate: getAdDirectoryRoleEligibilityScheduleInstance,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDirectoryRoleEligibilityScheduleInstances,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the role eligibility schedule instance.", Transform: transform.FromMethod("GetId")},
			{Name: "principal_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the principal that's in the scope of the role eligibility.", Transform: transform.FromMethod("GetPrincipalId")},
			{Name: "role_definition_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the role definition that's in the scope of the role eligibility.", Transform: transform.FromMethod("GetRoleDefinitionId")},
			{Name: "directory_scope_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the directory scope that's in the scope of the role eligibility.", Transform: transform.FromMethod("GetDirectoryScopeId")},
			{Name: "app_scope_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the app scope that's in the scope of the role eligibility.", Transform: transform.FromMethod("GetAppScopeId")},
			{Name: "start_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The start date and time of the role eligibility schedule instance.", Transform: transform.FromMethod("GetStartDateTime")},
			{Name: "end_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The end date and time of the role eligibility schedule instance.", Transform: transform.FromMethod("GetEndDateTime")},
			{Name: "member_type", Type: proto.ColumnType_STRING, Description: "How the role eligibility is inherited. It can be Inherited, Direct, or Group.", Transform: transform.FromMethod("GetMemberType")},
			{Name: "role_eligibility_schedule_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the role eligibility schedule.", Transform: transform.FromMethod("GetRoleEligibilityScheduleId")},

			// JSON columns for complex objects
			{Name: "app_scope", Type: proto.ColumnType_JSON, Description: "The app scope that's in the scope of the role eligibility.", Transform: transform.FromMethod("DirectoryRoleEligibilityScheduleInstanceAppScope")},
			{Name: "directory_scope", Type: proto.ColumnType_JSON, Description: "The directory scope that's in the scope of the role eligibility.", Transform: transform.FromMethod("DirectoryRoleEligibilityScheduleInstanceDirectoryScope")},
			{Name: "principal", Type: proto.ColumnType_JSON, Description: "The principal that's in the scope of the role eligibility.", Transform: transform.FromMethod("DirectoryRoleEligibilityScheduleInstancePrincipal")},
			{Name: "role_definition", Type: proto.ColumnType_JSON, Description: "The role definition that's in the scope of the role eligibility.", Transform: transform.FromMethod("DirectoryRoleEligibilityScheduleInstanceRoleDefinition")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adDirectoryRoleEligibilityScheduleInstanceTitle)},
		}),
	}
}

//// LIST FUNCTION

func listAdDirectoryRoleEligibilityScheduleInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_eligibility_schedule_instance.listAdDirectoryRoleEligibilityScheduleInstances", "connection_error", err)
		return nil, err
	}

	// List operations
	result, err := client.RoleManagement().Directory().RoleEligibilityScheduleInstances().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdDirectoryRoleEligibilityScheduleInstances", "list_directory_role_eligibility_schedule_instance_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.UnifiedRoleEligibilityScheduleInstanceable](result, adapter, models.CreateUnifiedRoleEligibilityScheduleInstanceCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleEligibilityScheduleInstances", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.UnifiedRoleEligibilityScheduleInstanceable) bool {
		d.StreamListItem(ctx, &ADDirectoryRoleEligibilityScheduleInstanceInfo{pageItem})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdDirectoryRoleEligibilityScheduleInstances", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDirectoryRoleEligibilityScheduleInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	roleEligibilityScheduleInstanceId := d.EqualsQuals["id"].GetStringValue()
	if roleEligibilityScheduleInstanceId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_directory_role_eligibility_schedule_instance.getAdDirectoryRoleEligibilityScheduleInstance", "connection_error", err)
		return nil, err
	}

	roleEligibilityScheduleInstance, err := client.RoleManagement().Directory().RoleEligibilityScheduleInstances().ByUnifiedRoleEligibilityScheduleInstanceId(roleEligibilityScheduleInstanceId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdDirectoryRoleEligibilityScheduleInstance", "get_directory_role_eligibility_schedule_instance_error", errObj)
		return nil, errObj
	}

	return &ADDirectoryRoleEligibilityScheduleInstanceInfo{roleEligibilityScheduleInstance}, nil
}

//// TRANSFORM FUNCTIONS

func adDirectoryRoleEligibilityScheduleInstanceTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADDirectoryRoleEligibilityScheduleInstanceInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetId()
	if title == nil {
		title = data.GetPrincipalId()
	}

	return title, nil
}
