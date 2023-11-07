package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/identity"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

//// TABLE DEFINITION

func tableAzureAdConditionalAccessPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_conditional_access_policy",
		Description: "Represents an Azure Active Directory (Azure AD) Conditional Access Policy.",
		Get: &plugin.GetConfig{
			Hydrate: getAdConditionalAccessPolicy,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdConditionalAccessPolicies,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "display_name", Require: plugin.Optional},
				{Name: "state", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Specifies the identifier of a conditionalAccessPolicy object.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Specifies a display name for the conditionalAccessPolicy object.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "Specifies the state of the conditionalAccessPolicy object. Possible values are: enabled, disabled, enabledForReportingButNotEnforced.", Transform: transform.FromMethod("GetState")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The create date of the conditional access policy.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "modified_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The modification date of the conditional access policy.", Transform: transform.FromMethod("GetModifiedDateTime")},
			{Name: "operator", Type: proto.ColumnType_STRING, Description: "Defines the relationship of the grant controls. Possible values: AND, OR.", Transform: transform.FromP(extractGrantControls, "operator")},

			// Json fields
			{Name: "applications", Type: proto.ColumnType_JSON, Description: "Applications and user actions included in and excluded from the policy.", Transform: transform.FromP(extractConditions, "applications")},
			{Name: "application_enforced_restrictions", Type: proto.ColumnType_JSON, Description: "Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control.", Transform: transform.FromP(extractSessionControls, "application_enforced_restrictions")},
			{Name: "built_in_controls", Type: proto.ColumnType_JSON, Description: "List of values of built-in controls required by the policy. Possible values: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.", Transform: transform.FromP(extractGrantControls, "built_in_controls")},
			{Name: "client_app_types", Type: proto.ColumnType_JSON, Description: "Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, easSupported, other.", Transform: transform.FromP(extractConditions, "client_app_types")},
			{Name: "custom_authentication_factors", Type: proto.ColumnType_JSON, Description: "List of custom controls IDs required by the policy.", Transform: transform.FromP(extractGrantControls, "custom_authentication_factors")},
			{Name: "cloud_app_security", Type: proto.ColumnType_JSON, Description: "Session control to apply cloud app security.", Transform: transform.FromP(extractSessionControls, "cloud_app_security")},
			{Name: "locations", Type: proto.ColumnType_JSON, Description: "Locations included in and excluded from the policy.", Transform: transform.FromP(extractConditions, "locations")},
			{Name: "persistent_browser", Type: proto.ColumnType_JSON, Description: "Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly.", Transform: transform.FromP(extractSessionControls, "persistent_browser")},
			{Name: "platforms", Type: proto.ColumnType_JSON, Description: "Platforms included in and excluded from the policy.", Transform: transform.FromP(extractConditions, "platforms")},
			{Name: "sign_in_frequency", Type: proto.ColumnType_JSON, Description: "Session control to enforce signin frequency.", Transform: transform.FromP(extractSessionControls, "sign_in_frequency")},
			{Name: "sign_in_risk_levels", Type: proto.ColumnType_JSON, Description: "Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.", Transform: transform.FromP(extractConditions, "sign_in_risk_levels")},
			{Name: "terms_of_use", Type: proto.ColumnType_JSON, Description: "List of terms of use IDs required by the policy.", Transform: transform.FromP(extractGrantControls, "terms_of_use")},
			{Name: "users", Type: proto.ColumnType_JSON, Description: "Users, groups, and roles included in and excluded from the policy.", Transform: transform.FromP(extractConditions, "users")},
			{Name: "user_risk_levels", Type: proto.ColumnType_JSON, Description: "User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.", Transform: transform.FromP(extractConditions, "user_risk_levels")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adConditionalAccessPolicyTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdConditionalAccessPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_conditional_access_policy.listAdConditionalAccessPolicies", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &identity.ConditionalAccessPoliciesRequestBuilderGetQueryParameters{
		Top: Int32(1000),
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit < 1000 {
			l := int32(*limit)
			input.Top = Int32(l)
		}
	}

	equalQuals := d.EqualsQuals
	filter := buildConditionalAccessPolicyQueryFilter(equalQuals)

	if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &identity.ConditionalAccessPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Identity().ConditionalAccess().Policies().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdConditionalAccessPolicies", "list_conditional_access_policy_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.ConditionalAccessPolicyable](result, adapter, models.CreateConditionalAccessPolicyCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdConditionalAccessPolicies", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(policy models.ConditionalAccessPolicyable) bool {
		d.StreamListItem(ctx, &ADConditionalAccessPolicyInfo{policy})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdConditionalAccessPolicies", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdConditionalAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	conditionalAccessPolicyId := d.EqualsQuals["id"].GetStringValue()
	if conditionalAccessPolicyId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_conditional_access_policy.getAdConditionalAccessPolicy", "connection_error", err)
		return nil, err
	}

	policy, err := client.Identity().ConditionalAccess().Policies().ByConditionalAccessPolicyId(conditionalAccessPolicyId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdConditionalAccessPolicy", "get_conditional_access_policy_error", errObj)
		return nil, errObj
	}
	return &ADConditionalAccessPolicyInfo{policy}, nil
}

func buildConditionalAccessPolicyQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name": "string",
		"state":        "string",
	}

	for qual, qualType := range filterQuals {
		switch qualType {
		case "string":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
			}
		}
	}

	return filters
}

//// TRANSFORM FUNCTIONS

func adConditionalAccessPolicyTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADConditionalAccessPolicyInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}

func extractGrantControls(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADConditionalAccessPolicyInfo)
	if data == nil {
		return nil, nil
	}
	param := d.Param.(string)
	grantControls := data.GetGrantControls()
	if grantControls == nil {
		return nil, nil
	}

	grantControlsMap := map[string]interface{}{
		"authentication_strength":       grantControls.GetAuthenticationStrength(),
		"backing_store":                 grantControls.GetBackingStore(),
		"built_in_controls":             grantControls.GetBuiltInControls(),
		"custom_authentication_factors": grantControls.GetCustomAuthenticationFactors(),
		"odata_type":                    grantControls.GetOdataType(),
		"operator":                      grantControls.GetOperator(),
		"terms_of_use":                  grantControls.GetTermsOfUse(),
	}

	return grantControlsMap[param], nil
}

func extractSessionControls(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADConditionalAccessPolicyInfo)
	if data == nil {
		return nil, nil
	}
	param := d.Param.(string)
	sessionControls := data.GetSessionControls()
	if sessionControls == nil {
		return nil, nil
	}

	sessionControlsMap := map[string]interface{}{
		"application_enforced_restrictions": sessionControls.GetApplicationEnforcedRestrictions(),
		"backing_store":                     sessionControls.GetBackingStore(),
		"cloud_app_security":                sessionControls.GetCloudAppSecurity(),
		"disable_resilience_defaults":       sessionControls.GetDisableResilienceDefaults(),
		"odata_type":                        sessionControls.GetOdataType(),
		"persistent_browser":                sessionControls.GetPersistentBrowser(),
		"sign_in_frequency":                 sessionControls.GetSignInFrequency(),
	}

	return sessionControlsMap[param], nil
}

func extractConditions(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADConditionalAccessPolicyInfo)
	if data == nil {
		return nil, nil
	}
	param := d.Param.(string)
	conditions := data.GetConditions()
	if conditions == nil {
		return nil, nil
	}

	conditionsMap := map[string]interface{}{
		"applications":                  conditions.GetApplications(),
		"backing_store":                 conditions.GetBackingStore(),
		"client_applications":           conditions.GetClientApplications(),
		"client_app_types":              conditions.GetClientAppTypes(),
		"devices":                       conditions.GetDevices(),
		"locations":                     conditions.GetLocations(),
		"odata_type":                    conditions.GetOdataType(),
		"platforms":                     conditions.GetPlatforms(),
		"service_principal_risk_levels": conditions.GetServicePrincipalRiskLevels(),
		"sign_in_risk_levels":           conditions.GetSignInRiskLevels(),
		"user_risk_levels":              conditions.GetUserRiskLevels(),
		"users":                         conditions.GetUsers(),
	}

	return conditionsMap[param], nil
}
