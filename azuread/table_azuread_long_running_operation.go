package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	betamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func tableAzureAdLongRunningOperation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_long_running_operation",
		Description: "Represents long-running operations in Azure AD, including role management alert refresh and authentication method operations.",
		Get: &plugin.GetConfig{
			Hydrate: getAdLongRunningOperation,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdLongRunningOperations,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the long-running operation.", Transform: transform.FromMethod("GetId")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the long-running operation (notStarted, running, succeeded, failed, unknownFutureValue).", Transform: transform.FromMethod("GetStatus")},
			{Name: "status_detail", Type: proto.ColumnType_STRING, Description: "Additional details about the operation status.", Transform: transform.FromMethod("GetStatusDetail")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the operation was created.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "last_action_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the last action was performed on the operation.", Transform: transform.FromMethod("GetLastActionDateTime")},
			{Name: "resource_location", Type: proto.ColumnType_STRING, Description: "The URL of the resource created or modified by the operation (if successful).", Transform: transform.FromMethod("GetResourceLocation")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adLongRunningOperationTitle)},
		}),
	}
}

type ADLongRunningOperationInfo struct {
	betamodels.LongRunningOperationable
}

func (operation *ADLongRunningOperationInfo) GetId() *string {
	return operation.LongRunningOperationable.GetId()
}

func (operation *ADLongRunningOperationInfo) GetStatus() *string {
	if operation.LongRunningOperationable.GetStatus() != nil {
		status := operation.LongRunningOperationable.GetStatus().String()
		return &status
	}
	return nil
}

func (operation *ADLongRunningOperationInfo) GetStatusDetail() *string {
	return operation.LongRunningOperationable.GetStatusDetail()
}

func (operation *ADLongRunningOperationInfo) GetCreatedDateTime() *string {
	if operation.LongRunningOperationable.GetCreatedDateTime() != nil {
		dateTime := operation.LongRunningOperationable.GetCreatedDateTime().String()
		return &dateTime
	}
	return nil
}

func (operation *ADLongRunningOperationInfo) GetLastActionDateTime() *string {
	if operation.LongRunningOperationable.GetLastActionDateTime() != nil {
		dateTime := operation.LongRunningOperationable.GetLastActionDateTime().String()
		return &dateTime
	}
	return nil
}

func (operation *ADLongRunningOperationInfo) GetResourceLocation() *string {
	return operation.LongRunningOperationable.GetResourceLocation()
}

func adLongRunningOperationTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	operation := d.HydrateItem.(*ADLongRunningOperationInfo)
	if operation.GetId() != nil {
		return *operation.GetId(), nil
	}
	return nil, nil
}

func listAdLongRunningOperations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Note: The Microsoft Graph API doesn't provide a direct way to list all long-running operations
	// They are typically accessed by specific operation ID or through specific contexts
	// This is a placeholder implementation - in practice, you would need to know the operation IDs
	// or access them through specific contexts like role management alerts or user authentication operations

	plugin.Logger(ctx).Warn("azuread_long_running_operation.listAdLongRunningOperations", "info", "Long-running operations are typically accessed by specific operation ID. Use the Get operation with a known operation ID.")

	return nil, nil
}

func getAdLongRunningOperation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access long-running operations
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_long_running_operation.getAdLongRunningOperation", "connection_error", err)
		return nil, err
	}

	operationId := d.EqualsQuals["id"].GetStringValue()
	if operationId == "" {
		return nil, nil
	}

	// Try to get the operation from role management alerts first
	result, err := client.IdentityGovernance().RoleManagementAlerts().Operations().ByLongRunningOperationId(operationId).Get(ctx, nil)
	if err != nil {
		// If not found in role management alerts, it might be a user authentication operation
		// Note: User authentication operations require a user ID, which we don't have in this context
		plugin.Logger(ctx).Debug("azuread_long_running_operation.getAdLongRunningOperation", "operation_not_found_in_role_management", operationId)
		return nil, err
	}

	return &ADLongRunningOperationInfo{result}, nil
}
