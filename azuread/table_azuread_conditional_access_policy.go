package azuread

import (
	"context"
	"strings"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureAdConditionalAccessPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_conditional_access_policy",
		Description: "Represents an Azure Active Directory (Azure AD) Conditional Access Policy",
		Get: &plugin.GetConfig{
			Hydrate:           getAdConditionalAccessPolicy,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate:           listAdConditionalAccessPolicies,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Request_UnsupportedQuery"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "display_name", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Specifies the identifier of a conditionalAccessPolicy object.", Transform: transform.FromGo()},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Specifies a display name for the conditionalAccessPolicy object."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "Specifies the state of the conditionalAccessPolicy object. Possible values are: enabled, disabled, enabledForReportingButNotEnforced."},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The create date of the conditional access policy."},
			{Name: "modified_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The modification date of the conditional access policy."},
			{Name: "operator", Type: proto.ColumnType_STRING, Description: "Defines the relationship of the grant controls. Possible values: AND, OR."},

			// Json fields
			{Name: "applications", Type: proto.ColumnType_JSON, Description: "Applications and user actions included in and excluded from the policy.", Transform: transform.FromField("Conditions.Applications")},
			{Name: "application_enforced_restrictions", Type: proto.ColumnType_JSON, Description: "Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control.", Transform: transform.FromField("SessionControls.ApplicationEnforcedRestrictions")},
			{Name: "built_in_controls", Type: proto.ColumnType_JSON, Description: "List of values of built-in controls required by the policy. Possible values: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.", Transform: transform.FromField("GrantControls.BuiltInControls")},
			{Name: "client_app_types", Type: proto.ColumnType_JSON, Description: "Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, easSupported, other.", Transform: transform.FromField("Conditions.ClientAppTypes")},
			{Name: "custom_authentication_factors", Type: proto.ColumnType_JSON, Description: "List of custom controls IDs required by the policy.", Transform: transform.FromField("GrantControls.CustomAuthenticationFactors")},
			{Name: "cloud_app_security", Type: proto.ColumnType_JSON, Description: "Session control to apply cloud app security.", Transform: transform.FromField("SessionControls.CloudAppSecurity")},
			{Name: "locations", Type: proto.ColumnType_JSON, Description: "Locations included in and excluded from the policy.", Transform: transform.FromField("Conditions.Locations")},
			{Name: "persistent_browser", Type: proto.ColumnType_JSON, Description: "Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly.", Transform: transform.FromField("SessionControls.PersistentBrowser")},
			{Name: "platforms", Type: proto.ColumnType_JSON, Description: "Platforms included in and excluded from the policy.", Transform: transform.FromField("Conditions.Platforms")},
			{Name: "sign_in_frequency", Type: proto.ColumnType_JSON, Description: "Session control to enforce signin frequency.", Transform: transform.FromField("SessionControls.SignInFrequency")},
			{Name: "sign_in_risk_levels", Type: proto.ColumnType_JSON, Description: "Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.", Transform: transform.FromField("Conditions.SignInRiskLevels")},
			{Name: "terms_of_use", Type: proto.ColumnType_JSON, Description: "List of terms of use IDs required by the policy.", Transform: transform.FromField("GrantControls.TermsOfUse")},
			{Name: "users", Type: proto.ColumnType_JSON, Description: "Users, groups, and roles included in and excluded from the policy.", Transform: transform.FromField("Conditions.Users")},
			{Name: "user_risk_levels", Type: proto.ColumnType_JSON, Description: "User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.", Transform: transform.FromField("Conditions.UserRiskLevels")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdConditionalAccessPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewConditionalAccessPolicyClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	input := odata.Query{}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	// unexpected status 400 with OData error: Request_UnsupportedQuery: Invalid page size specified: '1000'. Must be between 1 and 999 inclusive.
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit <= 999 {
			input.Top = int(*limit)
		}
	}
	
	qualsColumnMap := []QualsColumn{
		{"display_name", "string", "displayName"},
		{"state", "string", "state"},
	}

	filter := buildCommaonQueryFilter(qualsColumnMap, d.Quals)
	if len(filter) > 0 {
		input.Filter = strings.Join(filter, " and ")
	}

	conditionalAccessPolicies, _, err := client.List(ctx, input)
	if err != nil {
		if isNotFoundError(err) {
			plugin.Logger(ctx).Error("listAdConditionalAccessPolicies", "Resource not found", err)
			return nil, nil
		}
		return nil, err
	}

	for _, conditionalAccesPolicy := range *conditionalAccessPolicies {
		d.StreamListItem(ctx, conditionalAccesPolicy)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// Hydrate Functions

func getAdConditionalAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var conditionalAccessPolicyId string
	if h.Item != nil {
		conditionalAccessPolicyId = *h.Item.(msgraph.ConditionalAccessPolicy).ID
	} else {
		conditionalAccessPolicyId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if conditionalAccessPolicyId == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	client := msgraph.NewConditionalAccessPolicyClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer
	client.BaseClient.DisableRetries = true

	conditionalAccessPolicy, _, err := client.Get(ctx, conditionalAccessPolicyId, odata.Query{})
	if err != nil {
		if isNotFoundError(err) {
			plugin.Logger(ctx).Error("getAdConditionalAccessPolicy", "Resource not found", err)
			return nil, nil
		}
		return nil, err
	}
	return conditionalAccessPolicy, nil
}
