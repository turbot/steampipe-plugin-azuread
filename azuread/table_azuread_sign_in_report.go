package azuread

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/microsoftgraph/msgraph-sdk-go/auditlogs"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v6/rate_limiter"
)

// signInReportListLimiter keeps sign-in list requests, every page included, within Graph's
// identity and access reports limit of five requests per 10 seconds per tenant: the SDK runs an IN
// list on client_app_used as one concurrent list call per value.
var signInReportListLimiter = &rate_limiter.Definition{
	Name:       "azuread_sign_in_report_list",
	FillRate:   0.5,
	BucketSize: 5,
	Scope:      []string{rate_limiter.RateLimiterScopeConnection},
	Where:      "service = 'auditLogs' and action = 'signIns.list'",
}

//// TABLE DEFINITION

func tableAzureAdSignInReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_sign_in_report",
		Description: "Represents an Azure Active Directory (Azure AD) sign-in report.",
		Get: &plugin.GetConfig{
			Hydrate: getAdSignInReport,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdSignInReports,
			Tags:    map[string]string{"service": "auditLogs", "action": "signIns.list"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "created_date_time", Require: plugin.Optional, Operators: []string{">", ">=", "=", "<", "<="}},
				{Name: "client_app_used", Require: plugin.Optional},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID representing the sign-in activity.", Transform: transform.FromMethod("GetId")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time (UTC) the sign-in was initiated.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "user_display_name", Type: proto.ColumnType_STRING, Description: "Display name of the user that initiated the sign-in.", Transform: transform.FromMethod("GetUserDisplayName")},
			{Name: "user_principal_name", Type: proto.ColumnType_STRING, Description: "User principal name of the user that initiated the sign-in.", Transform: transform.FromMethod("GetUserPrincipalName")},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "ID of the user that initiated the sign-in.", Transform: transform.FromMethod("GetUserId")},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "Unique GUID representing the app ID in the Azure Active Directory.", Transform: transform.FromMethod("GetAppId")},
			{Name: "app_display_name", Type: proto.ColumnType_STRING, Description: "App name displayed in the Azure Portal.", Transform: transform.FromMethod("GetAppDisplayName")},
			{Name: "ip_address", Type: proto.ColumnType_STRING, Description: "IP address of the client used to sign in.", Transform: transform.FromMethod("GetIpAddress")},
			{Name: "client_app_used", Type: proto.ColumnType_STRING, Description: "Identifies the legacy client used for sign-in activity.", Transform: transform.FromMethod("GetClientAppUsed")},
			{Name: "correlation_id", Type: proto.ColumnType_STRING, Description: "The request ID sent from the client when the sign-in is initiated; used to troubleshoot sign-in activity.", Transform: transform.FromMethod("GetCorrelationId")},
			{Name: "conditional_access_status", Type: proto.ColumnType_STRING, Description: "Reports status of an activated conditional access policy. Possible values are: success, failure, notApplied, and unknownFutureValue.", Transform: transform.FromMethod("GetConditionalAccessStatus")},
			{Name: "is_interactive", Type: proto.ColumnType_BOOL, Description: "Indicates if a sign-in is interactive or not.", Transform: transform.FromMethod("GetIsInteractive")},
			{Name: "risk_detail", Type: proto.ColumnType_STRING, Description: "Provides the 'reason' behind a specific state of a risky user, sign-in or a risk event. The possible values are: none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, unknownFutureValue.", Transform: transform.FromMethod("GetRiskDetail")},
			{Name: "risk_level_aggregated", Type: proto.ColumnType_STRING, Description: "Aggregated risk level. The possible values are: none, low, medium, high, hidden, and unknownFutureValue.", Transform: transform.FromMethod("GetRiskLevelAggregated")},
			{Name: "risk_level_during_sign_in", Type: proto.ColumnType_STRING, Description: "Risk level during sign-in. The possible values are: none, low, medium, high, hidden, and unknownFutureValue.", Transform: transform.FromMethod("GetRiskLevelDuringSignIn")},
			{Name: "risk_state", Type: proto.ColumnType_STRING, Description: "Reports status of the risky user, sign-in, or a risk event. The possible values are: none, confirmedSafe, remediated, dismissed, atRisk, confirmedCompromised, unknownFutureValue.", Transform: transform.FromMethod("GetRiskState")},
			{Name: "resource_display_name", Type: proto.ColumnType_STRING, Description: "Name of the resource the user signed into.", Transform: transform.FromMethod("GetResourceDisplayName")},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "ID of the resource that the user signed into.", Transform: transform.FromMethod("GetResourceId")},

			// JSON fields
			{Name: "risk_event_types", Type: proto.ColumnType_JSON, Description: "Risk event types associated with the sign-in. The possible values are: unlikelyTravel, anonymizedIPAddress, maliciousIPAddress, unfamiliarFeatures, malwareInfectedIPAddress, suspiciousIPAddress, leakedCredentials, investigationsThreatIntelligence, generic, and unknownFutureValue.", Transform: transform.FromMethod("GetRiskEventTypes").Transform(formatSignInReportRiskEventTypes)},
			{Name: "status", Type: proto.ColumnType_JSON, Description: "Sign-in status. Includes the error code and description of the error (in case of a sign-in failure).", Transform: transform.FromMethod("SignInStatus")},
			{Name: "device_detail", Type: proto.ColumnType_JSON, Description: "Device information from where the sign-in occurred; includes device ID, operating system, and browser.", Transform: transform.FromMethod("SignInDeviceDetail")},
			{Name: "location", Type: proto.ColumnType_JSON, Description: "Provides the city, state, and country code where the sign-in originated.", Transform: transform.FromMethod("SignInLocation")},
			{Name: "applied_conditional_access_policies", Type: proto.ColumnType_JSON, Description: "Provides a list of conditional access policies that are triggered by the corresponding sign-in activity.", Transform: transform.FromMethod("SignInAppliedConditionalAccessPolicies")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetId")},
		}),
	}
}

//// LIST FUNCTION

func listAdSignInReports(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_sign_in_report.listAdSignInReports", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &auditlogs.SignInsRequestBuilderGetQueryParameters{
		Top: Int32(999),
	}

	pushed := parseSignInQuals(d.Quals)
	if filter := pushed.graphFilter(); len(filter) > 0 {
		joined := strings.Join(filter, " and ")
		input.Filter = &joined
	} else if limit := d.QueryContext.Limit; limit != nil && *limit > 0 && *limit < 999 {
		// Shrunk only without a $filter: a filtered page can hold rows admits drops, which a short page would spend a paced request on.
		input.Top = Int32(int32(*limit))
	}

	options := &auditlogs.SignInsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.AuditLogs().SignIns().Get(ctx, options)
	if err != nil {
		if ctx.Err() != nil {
			return nil, nil
		}
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdSignInReports", "list_sign_in_report_error", errObj)
		return nil, errObj
	}

	fetchPage := func(ctx context.Context, nextLink string) (models.SignInCollectionResponseable, error) {
		return client.AuditLogs().SignIns().WithUrl(nextLink).Get(ctx, nil)
	}
	err = iterateSignInPages(ctx, result, d.WaitForListRateLimit, fetchPage, func(signIn models.SignInable) bool {
		// Graph matches more than Postgres (whole-second bounds, case-insensitive eq); streamed, those rows would use up a LIMIT before Postgres rejects them.
		if pushed.admits(signIn) {
			d.StreamListItem(ctx, &ADSignInReportInfo{signIn})
		}

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		if ctx.Err() != nil {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listAdSignInReports", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

// iterateSignInPages calls emit for each sign-in on the first page and on every @odata.nextLink
// page, calling wait before each next-page fetch so every request is paced, not just the first.
func iterateSignInPages(ctx context.Context, page models.SignInCollectionResponseable, wait func(context.Context), fetchPage func(context.Context, string) (models.SignInCollectionResponseable, error), emit func(models.SignInable) bool) error {
	for page != nil {
		for _, signIn := range page.GetValue() {
			if signIn == nil {
				continue
			}
			if !emit(signIn) {
				return nil
			}
		}
		nextLink := page.GetOdataNextLink()
		if nextLink == nil || *nextLink == "" {
			return nil
		}
		wait(ctx)
		if ctx.Err() != nil {
			return nil
		}
		var err error
		if page, err = fetchPage(ctx, *nextLink); err != nil {
			return err
		}
	}
	return nil
}

type timeBound struct {
	operator string
	at       time.Time
}

// signInQuals holds the quals this table sends to Graph, parsed once per list call.
type signInQuals struct {
	created    []timeBound
	clientApps []string
}

// parseSignInQuals reads the created_date_time and client_app_used quals. A list value (an IN the
// SDK did not split) is skipped and left to Postgres.
func parseSignInQuals(keyQuals plugin.KeyColumnQualMap) signInQuals {
	var pushed signInQuals
	if keyQuals["created_date_time"] != nil {
		for _, q := range keyQuals["created_date_time"].Quals {
			if ts := q.Value.GetTimestampValue(); ts != nil {
				pushed.created = append(pushed.created, timeBound{q.Operator, ts.AsTime().UTC()})
			}
		}
	}
	if keyQuals["client_app_used"] != nil {
		for _, q := range keyQuals["client_app_used"].Quals {
			if q.Value.GetListValue() == nil {
				pushed.clientApps = append(pushed.clientApps, q.Value.GetStringValue())
			}
		}
	}
	return pushed
}

// graphFilter renders the quals as OData $filter terms, so Graph pages through only the matching
// sign-ins instead of the whole log. Graph filters createdDateTime with eq, le and ge only, in whole
// seconds, so each bound is widened to the second outside it and strict operators become inclusive.
func (q signInQuals) graphFilter() []string {
	filters := []string{}
	for _, b := range q.created {
		// An infinite or far-off bound would format as a year Graph rejects; Postgres applies it alone.
		if b.at.Year() < 1970 || b.at.Year() > 9999 {
			continue
		}
		floor := b.at.Truncate(time.Second)
		lower := fmt.Sprintf("createdDateTime ge %s", floor.Format(time.RFC3339))
		upper := fmt.Sprintf("createdDateTime le %s", floor.Add(time.Second).Format(time.RFC3339))
		switch b.operator {
		case ">", ">=":
			filters = append(filters, lower)
		case "<", "<=":
			filters = append(filters, upper)
		case "=":
			filters = append(filters, lower, upper)
		}
	}
	for _, app := range q.clientApps {
		filters = append(filters, fmt.Sprintf("clientAppUsed eq '%s'", strings.ReplaceAll(app, "'", "''")))
	}
	return filters
}

// admits reports whether Postgres may keep the sign-in under the pushed quals. The FDW hands
// Postgres the full 100 ns timestamp as RFC 3339 text, which timestamptz input rounds to the
// microsecond; an exact half microsecond may go either way, so both neighbours are tried there.
func (q signInQuals) admits(signIn models.SignInable) bool {
	for _, app := range q.clientApps {
		if got := signIn.GetClientAppUsed(); got == nil || *got != app {
			return false
		}
	}
	if len(q.created) == 0 {
		return true
	}
	created := signIn.GetCreatedDateTime()
	if created == nil {
		return false
	}
	if q.createdWithin(created.Round(time.Microsecond)) {
		return true
	}
	return created.Sub(created.Truncate(time.Microsecond)) == 500*time.Nanosecond && q.createdWithin(created.Truncate(time.Microsecond))
}

func (q signInQuals) createdWithin(at time.Time) bool {
	for _, b := range q.created {
		var ok bool
		switch b.operator {
		case ">":
			ok = at.After(b.at)
		case ">=":
			ok = !at.Before(b.at)
		case "=":
			ok = at.Equal(b.at)
		case "<=":
			ok = !at.After(b.at)
		case "<":
			ok = at.Before(b.at)
		default:
			ok = true
		}
		if !ok {
			return false
		}
	}
	return true
}

//// HYDRATE FUNCTIONS

func getAdSignInReport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	signInID := d.EqualsQuals["id"].GetStringValue()
	if signInID == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_sign_in_report.getAdSignInReport", "connection_error", err)
		return nil, err
	}

	signIn, err := client.AuditLogs().SignIns().BySignInId(signInID).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdSignInReport", "get_sign_in_report_error", errObj)
		return nil, errObj
	}

	return &ADSignInReportInfo{signIn}, nil
}

//// TRANSFORM FUNCTIONS

func formatSignInReportRiskEventTypes(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADSignInReportInfo)
	riskEventTypes := data.GetRiskEventTypes()
	if len(riskEventTypes) == 0 {
		return nil, nil
	}

	return riskEventTypes, nil
}
