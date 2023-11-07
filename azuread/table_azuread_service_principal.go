package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
)

//// TABLE DEFINITION

func tableAzureAdServicePrincipal(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_service_principal",
		Description: "Represents an Azure Active Directory (Azure AD) service principal.",
		Get: &plugin.GetConfig{
			Hydrate: getAdServicePrincipal,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdServicePrincipals,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "display_name", Require: plugin.Optional},
				{Name: "account_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "service_principal_type", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the service principal.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the service principal.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the associated application (its appId property).", Transform: transform.FromMethod("GetAppId")},

			// Other fields
			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "true if the service principal account is enabled; otherwise, false.", Transform: transform.FromMethod("GetAccountEnabled")},
			{Name: "app_display_name", Type: proto.ColumnType_STRING, Description: "The display name exposed by the associated application.", Transform: transform.FromMethod("GetAppDisplayName")},
			{Name: "app_owner_organization_id", Type: proto.ColumnType_STRING, Description: "Contains the tenant id where the application is registered. This is applicable only to service principals backed by applications.", Transform: transform.FromMethod("GetAppOwnerOrganizationId")},
			{Name: "app_role_assignment_required", Type: proto.ColumnType_BOOL, Description: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false.", Transform: transform.FromMethod("GetAppRoleAssignmentRequired")},
			{Name: "service_principal_type", Type: proto.ColumnType_STRING, Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Azure AD internally.", Transform: transform.FromMethod("GetServicePrincipalType")},
			{Name: "sign_in_audience", Type: proto.ColumnType_STRING, Description: "Specifies the Microsoft accounts that are supported for the current application. Supported values are: AzureADMyOrg, AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, PersonalMicrosoftAccount.", Transform: transform.FromMethod("GetSignInAudience")},
			{Name: "app_description", Type: proto.ColumnType_STRING, Description: "The description exposed by the associated application.", Transform: transform.FromMethod("GetAppDescription")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Free text field to provide an internal end-user facing description of the service principal.", Transform: transform.FromMethod("GetDescription")},
			{Name: "login_url", Type: proto.ColumnType_STRING, Description: "Specifies the URL where the service provider redirects the user to Azure AD to authenticate. Azure AD uses the URL to launch the application from Microsoft 365 or the Azure AD My Apps. When blank, Azure AD performs IdP-initiated sign-on for applications configured with SAML-based single sign-on.", Transform: transform.FromMethod("GetLoginUrl")},
			{Name: "logout_url", Type: proto.ColumnType_STRING, Description: "Specifies the URL that will be used by Microsoft's authorization service to logout an user using OpenId Connect front-channel, back-channel or SAML logout protocols.", Transform: transform.FromMethod("GetLogoutUrl")},

			// JSON fields
			{Name: "add_ins", Type: proto.ColumnType_JSON, Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts.", Transform: transform.FromMethod("GetAddIns")},
			{Name: "alternative_names", Type: proto.ColumnType_JSON, Description: "Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities.", Transform: transform.FromMethod("GetAlternativeNames")},
			{Name: "app_roles", Type: proto.ColumnType_JSON, Description: "The roles exposed by the application which this service principal represents.", Transform: transform.FromMethod("GetAppRoles")},
			{Name: "info", Type: proto.ColumnType_JSON, Description: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs.", Transform: transform.FromMethod("GetInfo")},
			{Name: "key_credentials", Type: proto.ColumnType_JSON, Description: "The collection of key credentials associated with the service principal.", Transform: transform.FromMethod("GetKeyCredentials")},
			{Name: "notification_email_addresses", Type: proto.ColumnType_JSON, Description: "Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications.", Transform: transform.FromMethod("GetNotificationEmailAddresses")},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getServicePrincipalOwners, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "password_credentials", Type: proto.ColumnType_JSON, Description: "Represents a password credential associated with a service principal.", Transform: transform.FromMethod("GetPasswordCredentials")},
			{Name: "oauth2_permission_scopes", Type: proto.ColumnType_JSON, Description: "The published permission scopes.", Transform: transform.FromMethod("GetOauth2PermissionScopes")},
			{Name: "reply_urls", Type: proto.ColumnType_JSON, Description: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application.", Transform: transform.FromMethod("GetReplyUrls")},
			{Name: "service_principal_names", Type: proto.ColumnType_JSON, Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD.", Transform: transform.FromMethod("GetServicePrincipalNames")},
			{Name: "tags_src", Type: proto.ColumnType_JSON, Description: "Custom strings that can be used to categorize and identify the service principal.", Transform: transform.FromMethod("GetTags")},

			// Standard columns
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(adServicePrincipalTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adServicePrincipalTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdServicePrincipals(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_service_principal.listAdServicePrincipals", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
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

	var queryFilter string
	equalQuals := d.EqualsQuals
	quals := d.Quals
	filter := buildServicePrincipalQueryFilter(equalQuals)
	filter = append(filter, buildServicePrincipalBoolNEFilter(quals)...)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = &queryFilter
	} else if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.ServicePrincipals().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdServicePrincipals", "list_service_principal_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.ServicePrincipalable](result, adapter, models.CreateServicePrincipalCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdServicePrincipals", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.ServicePrincipalable) bool {

		d.StreamListItem(ctx, &ADServicePrincipalInfo{pageItem})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdServicePrincipals", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdServicePrincipal(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	servicePrincipalID := d.EqualsQuals["id"].GetStringValue()
	if servicePrincipalID == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_service_principal.getAdServicePrincipal", "connection_error", err)
		return nil, err
	}

	servicePrincipal, err := client.ServicePrincipals().ByServicePrincipalId(servicePrincipalID).Get(ctx, &serviceprincipals.ServicePrincipalItemRequestBuilderGetRequestConfiguration{})
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdServicePrincipal", "get_service_principal_error", errObj)
		return nil, errObj
	}

	return &ADServicePrincipalInfo{servicePrincipal}, nil
}

func getServicePrincipalOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_service_principal.getServicePrincipalOwners", "connection_error", err)
		return nil, err
	}

	servicePrincipal := h.Item.(*ADServicePrincipalInfo)
	servicePrincipalID := servicePrincipal.GetId()

	if servicePrincipalID == nil {
		return nil, nil
	}

	headers := &abstractions.RequestHeaders{}
	headers.Add("ConsistencyLevel", "eventual")

	includeCount := true
	requestParameters := &serviceprincipals.ItemOwnersRequestBuilderGetQueryParameters{
		Count: &includeCount,
	}

	config := &serviceprincipals.ItemOwnersRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParameters,
	}

	ownerIds := []*string{}
	owners, err := client.ServicePrincipals().ByServicePrincipalId(*servicePrincipalID).Owners().Get(ctx, config)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getServicePrincipalOwners", "get_service_principal_owners_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.DirectoryObjectable](owners, adapter, models.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("getServicePrincipalOwners", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.DirectoryObjectable) bool {
		owner := pageItem.(models.DirectoryObjectable)
		ownerIds = append(ownerIds, owner.GetId())

		return true
	})
	if err != nil {
		plugin.Logger(ctx).Error("getServicePrincipalOwners", "paging_error", err)
		return nil, err
	}

	return ownerIds, nil
}

//// TRANSFORM FUNCTIONS

func adServicePrincipalTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	servicePrincipal := d.HydrateItem.(*ADServicePrincipalInfo)
	tags := servicePrincipal.GetTags()
	return TagsToMap(tags)
}

func adServicePrincipalTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADServicePrincipalInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}

func buildServicePrincipalQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name":           "string",
		"account_enabled":        "bool",
		"service_principal_type": "string",
	}

	for qual, qualType := range filterQuals {
		switch qualType {
		case "string":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
			}
		case "bool":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), equalQuals[qual].GetBoolValue()))
			}
		}
	}

	return filters
}

func buildServicePrincipalBoolNEFilter(quals plugin.KeyColumnQualMap) []string {
	filters := []string{}

	filterQuals := []string{
		"account_enabled",
	}

	for _, qual := range filterQuals {
		if quals[qual] != nil {
			for _, q := range quals[qual].Quals {
				value := q.Value.GetBoolValue()
				if q.Operator == "<>" {
					filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), !value))
					break
				}
			}
		}
	}

	return filters
}
