package azuread

import (
	"context"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureAdSignInReport() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_sign_in_report",
		Description: "Represents an Azure Active Directory (Azure AD) sign in report",
		Get: &plugin.GetConfig{
			Hydrate:           getAdSigninReport,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdSigninReports,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID representing the sign-in activity."},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time (UTC) the sign-in was initiated."},
			{Name: "user_display_name", Type: proto.ColumnType_STRING, Description: "Display name of the user that initiated the sign-in."},
			{Name: "user_principal_name", Type: proto.ColumnType_STRING, Description: "User principal name of the user that initiated the sign-in."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "ID of the user that initiated the sign-in."},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "Unique GUID representing the app ID in the Azure Active Directory."},
			{Name: "app_display_name", Type: proto.ColumnType_STRING, Description: "App name displayed in the Azure Portal."},
			{Name: "ip_address", Type: proto.ColumnType_STRING, Description: "IP address of the client used to sign in.", Transform: transform.FromField("IPAddress")},
			{Name: "client_app_used", Type: proto.ColumnType_STRING, Description: "Identifies the legacy client used for sign-in activity."},
			{Name: "correlation_id", Type: proto.ColumnType_STRING, Description: "The request ID sent from the client when the sign-in is initiated; used to troubleshoot sign-in activity."},
			{Name: "conditional_access_status", Type: proto.ColumnType_STRING, Description: "Reports status of an activated conditional access policy. Possible values are: success, failure, notApplied, and unknownFutureValue."},
			{Name: "is_interactive", Type: proto.ColumnType_BOOL, Description: "Indicates if a sign-in is interactive or not."},
			{Name: "risk_detail", Type: proto.ColumnType_STRING, Description: "Provides the 'reason' behind a specific state of a risky user, sign-in or a risk event. The possible values are: none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, unknownFutureValue."},
			{Name: "risk_level_aggregated", Type: proto.ColumnType_STRING, Description: "Aggregated risk level. The possible values are: none, low, medium, high, hidden, and unknownFutureValue."},
			{Name: "risk_level_during_sign_in", Type: proto.ColumnType_STRING, Description: "Risk level during sign-in. The possible values are: none, low, medium, high, hidden, and unknownFutureValue."},
			{Name: "risk_state", Type: proto.ColumnType_STRING, Description: "Reports status of the risky user, sign-in, or a risk event. The possible values are: none, confirmedSafe, remediated, dismissed, atRisk, confirmedCompromised, unknownFutureValue."},
			{Name: "resource_display_name", Type: proto.ColumnType_STRING, Description: "Name of the resource the user signed into."},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "ID of the resource that the user signed into."},

			// Json fields
			{Name: "risk_event_types", Type: proto.ColumnType_JSON, Description: "Risk event types associated with the sign-in. The possible values are: unlikelyTravel, anonymizedIPAddress, maliciousIPAddress, unfamiliarFeatures, malwareInfectedIPAddress, suspiciousIPAddress, leakedCredentials, investigationsThreatIntelligence, generic, and unknownFutureValue."},
			{Name: "status", Type: proto.ColumnType_JSON, Description: "Sign-in status. Includes the error code and description of the error (in case of a sign-in failure)."},
			{Name: "device_detail", Type: proto.ColumnType_JSON, Description: "Device information from where the sign-in occurred; includes device ID, operating system, and browser."},
			{Name: "location", Type: proto.ColumnType_JSON, Description: "	Provides the city, state, and country code where the sign-in originated."},
			{Name: "applied_conditional_access_policies", Type: proto.ColumnType_JSON, Description: "Provides a list of conditional access policies that are triggered by the corresponding sign-in activity."},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("Id")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdSigninReports(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewSignInLogsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	// As per our test result we have set the max limit to 999
	input := odata.Query{
		Top: 999,
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < 999 {
			input.Top = int(*limit)
		}
	}

	signInReports, _, err := client.List(ctx, input)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		logger.Error("listAdSigninReports", "list", err)
		return nil, err
	}

	for _, report := range *signInReports {
		d.StreamListItem(ctx, report)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// Hydrate Functions

func getAdSigninReport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	var signInId string
	if h.Item != nil {
		signInId = *h.Item.(msgraph.ServicePrincipal).ID
	} else {
		signInId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if signInId == "" {
		return nil, nil
	}
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewSignInLogsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer
	client.BaseClient.DisableRetries = true

	signInReport, _, err := client.Get(ctx, signInId, odata.Query{})
	if err != nil {
		logger.Error("getAdSigninReport", "get", err)
		return nil, err
	}
	return *signInReport, nil
}
