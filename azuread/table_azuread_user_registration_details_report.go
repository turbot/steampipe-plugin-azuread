package azuread

import (
	"context"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/reports"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdUserRegistrationDetailsReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user_registration_details_report",
		Description: "Represents an Azure Active Directory (Azure AD) user-registration-details report.",
		Get: &plugin.GetConfig{
			Hydrate: getAdUserRegistrationDetailsReport,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdUserRegistrationDetailsReport,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID representing the sign-in activity.", Transform: transform.FromMethod("GetId")},
			{Name: "user_principal_name", Type: proto.ColumnType_STRING, Description: "User principal name of the user that initiated the sign-in.", Transform: transform.FromMethod("GetUserPrincipalName")},
			{Name: "user_display_name", Type: proto.ColumnType_STRING, Description: "Display name of the user that initiated the sign-in.", Transform: transform.FromMethod("GetUserDisplayName")},
			{Name: "user_type", Type: proto.ColumnType_STRING, Description: "A string value that can be used to classify user types in your directory.", Transform: transform.FromMethod("GetUserType")},
			{Name: "is_admin", Type: proto.ColumnType_BOOL, Description: "Indicates whether the user has an admin role in the tenant", Transform: transform.FromMethod("GetIsAdmin")},
			{Name: "is_mfa_registered", Type: proto.ColumnType_BOOL, Description: "Indicates whether the user has registered a strong authentication method for multifactor authentication", Transform: transform.FromMethod("GetIsMfaRegistered")},
			{Name: "is_mfa_capable", Type: proto.ColumnType_BOOL, Description: "Indicates whether the user has registered a strong authentication method for multifactor authentication", Transform: transform.FromMethod("GetIsMfaCapable")},
			{Name: "is_sspr_registered", Type: proto.ColumnType_BOOL, Description: "Indicates whether the user has registered the required number of authentication methods for self-service password reset", Transform: transform.FromMethod("GetIsSsprRegistered")},
			{Name: "is_sspr_enabled", Type: proto.ColumnType_BOOL, Description: "Indicates whether the user is allowed to perform self-service password reset by policy", Transform: transform.FromMethod("GetIsSsprEnabled")},
			{Name: "is_sspr_capable", Type: proto.ColumnType_BOOL, Description: "Indicates whether the user has registered the required number of authentication methods for self-service password reset and the user is allowed to perform self-service password reset by policy", Transform: transform.FromMethod("GetIsSsprCapable")},
			{Name: "is_passwordless_capable", Type: proto.ColumnType_BOOL, Description: "Indicates whether the user has registered a passwordless strong authentication method", Transform: transform.FromMethod("GetIsPasswordlessCapable")},
			{Name: "is_system_preferred_authentication_method_enabled", Type: proto.ColumnType_BOOL, Description: "Indicates whether system preferred authentication method is enabled", Transform: transform.FromMethod("GetIsSystemPreferredAuthenticationMethodEnabled")},
			{Name: "user_preferred_method_for_secondary_authentication", Type: proto.ColumnType_STRING, Description: "The method the user selected as the default second-factor for performing multifactor authentication", Transform: transform.FromMethod("GetUserPreferredMethodForSecondaryAuthentication")},
			{Name: "last_updated_datetime", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time (UTC) when the report was last updated", Transform: transform.FromMethod("GetLastUpdatedDateTime")},

			// JSON fields
			{Name: "methods_registered", Type: proto.ColumnType_JSON, Description: "Collection of authentication methods registered", Transform: transform.FromMethod("GetMethodsRegistered")},
			{Name: "system_preferred_authentication_methods", Type: proto.ColumnType_JSON, Description: "Collection of authentication methods that the system determined to be the most secure authentication methods among the registered methods for second factor authentication.", Transform: transform.FromMethod("GetSystemPreferredAuthenticationMethods")},
		}),
	}
}

//// LIST FUNCTION

func listAdUserRegistrationDetailsReport(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("table_azuread_user_registration_details_report.listAdUserRegistrationDetailsReport", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &reports.AuthenticationMethodsUserRegistrationDetailsRequestBuilderGetQueryParameters{
		Top: Int32(999),
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

	conf := &reports.AuthenticationMethodsUserRegistrationDetailsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}
	result, err := client.Reports().AuthenticationMethods().UserRegistrationDetails().Get(ctx, conf)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdUserRegistrationDetailsReport", "list_user_registration_details_report_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[interface{}](result, adapter, models.CreateUserRegistrationDetailsCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdUserRegistrationDetailsReport", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		// To prevent errors during type conversion caused by inconsistent API responses (especially with larger data sets), we may get a different type of response (models.DirectoryAuditable). We need to include the following check.
		if details, ok := pageItem.(models.UserRegistrationDetailsable); ok {
			d.StreamListItem(ctx, &ADUserRegistrationDetailsReport{details})
		}

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdUserRegistrationDetailsReport", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdUserRegistrationDetailsReport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	signInID := d.EqualsQuals["id"].GetStringValue()
	if signInID == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_sign_in_report.getAdUserRegistrationDetailsReport", "connection_error", err)
		return nil, err
	}

	result, err := client.Reports().AuthenticationMethods().UserRegistrationDetails().ByUserRegistrationDetailsId(signInID).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdUserRegistrationDetailsReport", "get_user_registration_details_report_error", errObj)
		return nil, errObj
	}

	return &ADUserRegistrationDetailsReport{result}, nil
}

// //// TRANSFORM FUNCTIONS

// func formatSignInReportRiskEventTypes(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	data := d.HydrateItem.(*ADSignInReportInfo)
// 	riskEventTypes := data.GetRiskEventTypes()
// 	if len(riskEventTypes) == 0 {
// 		return nil, nil
// 	}

// 	return riskEventTypes, nil
// }
