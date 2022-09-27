package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/identity/conditionalaccess/policies"
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
			{Name: "operator", Type: proto.ColumnType_STRING, Description: "Defines the relationship of the grant controls. Possible values: AND, OR.", Transform: transform.FromMethod("ConditionalAccessPolicyGrantControlsOperator")},

			// Json fields
			{Name: "applications", Type: proto.ColumnType_JSON, Description: "Applications and user actions included in and excluded from the policy.", Transform: transform.FromMethod("ConditionalAccessPolicyConditionsApplications")},
			{Name: "application_enforced_restrictions", Type: proto.ColumnType_JSON, Description: "Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control.", Transform: transform.FromMethod("ConditionalAccessPolicySessionControlsApplicationEnforcedRestrictions")},
			{Name: "built_in_controls", Type: proto.ColumnType_JSON, Description: "List of values of built-in controls required by the policy. Possible values: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.", Transform: transform.FromMethod("ConditionalAccessPolicyGrantControlsBuiltInControls")},
			{Name: "client_app_types", Type: proto.ColumnType_JSON, Description: "Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, easSupported, other.", Transform: transform.FromMethod("ConditionalAccessPolicyConditionsClientAppTypes")},
			{Name: "custom_authentication_factors", Type: proto.ColumnType_JSON, Description: "List of custom controls IDs required by the policy.", Transform: transform.FromMethod("ConditionalAccessPolicyGrantControlsCustomAuthenticationFactors")},
			{Name: "cloud_app_security", Type: proto.ColumnType_JSON, Description: "Session control to apply cloud app security.", Transform: transform.FromMethod("ConditionalAccessPolicySessionControlsCloudAppSecurity")},
			{Name: "locations", Type: proto.ColumnType_JSON, Description: "Locations included in and excluded from the policy.", Transform: transform.FromMethod("ConditionalAccessPolicyConditionsLocations")},
			{Name: "persistent_browser", Type: proto.ColumnType_JSON, Description: "Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly.", Transform: transform.FromMethod("ConditionalAccessPolicySessionControlsPersistentBrowser")},
			{Name: "platforms", Type: proto.ColumnType_JSON, Description: "Platforms included in and excluded from the policy.", Transform: transform.FromMethod("ConditionalAccessPolicyConditionsPlatforms")},
			{Name: "sign_in_frequency", Type: proto.ColumnType_JSON, Description: "Session control to enforce signin frequency.", Transform: transform.FromMethod("ConditionalAccessPolicySessionControlsSignInFrequency")},
			{Name: "sign_in_risk_levels", Type: proto.ColumnType_JSON, Description: "Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.", Transform: transform.FromMethod("ConditionalAccessPolicyConditionsSignInRiskLevels")},
			{Name: "terms_of_use", Type: proto.ColumnType_JSON, Description: "List of terms of use IDs required by the policy.", Transform: transform.FromMethod("ConditionalAccessPolicyGrantControlsTermsOfUse")},
			{Name: "users", Type: proto.ColumnType_JSON, Description: "Users, groups, and roles included in and excluded from the policy.", Transform: transform.FromMethod("ConditionalAccessPolicyConditionsUsers")},
			{Name: "user_risk_levels", Type: proto.ColumnType_JSON, Description: "User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.", Transform: transform.FromMethod("ConditionalAccessPolicyConditionsUserRiskLevels")},

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
	input := &policies.PoliciesRequestBuilderGetQueryParameters{}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	// unexpected status 400 with OData error: Request_UnsupportedQuery: Invalid page size specified: '1000'. Must be between 1 and 999 inclusive.
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit <= 999 {
			l := int32(*limit)
			input.Top = &l
		}
	}

	equalQuals := d.KeyColumnQuals
	filter := buildConditionalAccessPolicyQueryFilter(equalQuals)

	if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &policies.PoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Identity().ConditionalAccess().Policies().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdConditionalAccessPolicies", "list_conditional_access_policy_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateConditionalAccessPolicyCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdConditionalAccessPolicies", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		policy := pageItem.(models.ConditionalAccessPolicyable)

		d.StreamListItem(ctx, &ADConditionalAccessPolicyInfo{policy})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.QueryStatus.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdConditionalAccessPolicies", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdConditionalAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	conditionalAccessPolicyId := d.KeyColumnQuals["id"].GetStringValue()
	if conditionalAccessPolicyId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_conditional_access_policy.getAdConditionalAccessPolicy", "connection_error", err)
		return nil, err
	}

	policy, err := client.Identity().ConditionalAccess().PoliciesById(conditionalAccessPolicyId).Get(ctx, nil)
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
