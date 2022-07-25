package azuread

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners"
)

//// TABLE DEFINITION

func tableAzureAdServicePrincipalTest() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_service_principal_test",
		Description: "Represents an Azure Active Directory (Azure AD) service principal",
		Get: &plugin.GetConfig{
			Hydrate: getAdServicePrincipalTest,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdServicePrincipalsTest,
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

			// Json fields
			{Name: "add_ins", Type: proto.ColumnType_JSON, Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts.", Transform: transform.FromMethod("ServicePrincipalAddIns")},
			{Name: "alternative_names", Type: proto.ColumnType_JSON, Description: "Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities.", Transform: transform.FromMethod("GetAlternativeNames")},
			{Name: "app_roles", Type: proto.ColumnType_JSON, Description: "The roles exposed by the application which this service principal represents.", Transform: transform.FromMethod("ServicePrincipalAppRoles")},
			{Name: "info", Type: proto.ColumnType_JSON, Description: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs.", Transform: transform.FromMethod("ServicePrincipalInfo")},
			{Name: "key_credentials", Type: proto.ColumnType_JSON, Description: "The collection of key credentials associated with the service principal.", Transform: transform.FromMethod("ServicePrincipalKeyCredentials")},
			{Name: "notification_email_addresses", Type: proto.ColumnType_JSON, Description: "Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications.", Transform: transform.FromMethod("GetNotificationEmailAddresses")},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getServicePrincipalOwners, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "password_credentials", Type: proto.ColumnType_JSON, Description: "Represents a password credential associated with a service principal.", Transform: transform.FromMethod("ServicePrincipalPasswordCredentials")},
			// {Name: "published_permission_scopes", Type: proto.ColumnType_JSON, Description: "The published permission scopes."},
			{Name: "reply_urls", Type: proto.ColumnType_JSON, Description: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application.", Transform: transform.FromMethod("GetReplyUrls")},
			{Name: "service_principal_names", Type: proto.ColumnType_JSON, Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD.", Transform: transform.FromMethod("GetServicePrincipalNames")},
			{Name: "tags_src", Type: proto.ColumnType_JSON, Description: "Custom strings that can be used to categorize and identify the service principal.", Transform: transform.FromMethod("GetTags")},
			// {Name: "verified_publisher", Type: proto.ColumnType_JSON, Description: "Specifies the verified publisher of the application which this service principal represents."},

			// Standard columns
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(adServicePrincipalTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adServicePrincipalTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdServicePrincipalsTest(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	// List operations
	input := &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	// unexpected status 400 with OData error: Request_UnsupportedQuery: Invalid page size specified: '1000'. Must be between 1 and 999 inclusive.
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit <= 999 {
			l := int32(*limit)
			input.Top = &l
		}
	}

	var queryFilter string
	equalQuals := d.KeyColumnQuals
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

	result, err := client.ServicePrincipals().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateServicePrincipalCollectionResponseFromDiscriminatorValue)

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		servicePrincipal := pageItem.(models.ServicePrincipalable)

		d.StreamListItem(ctx, &ADServicePrincipalInfo{servicePrincipal})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return false
		}

		return true
	})

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdServicePrincipalTest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	servicePrincipalID := d.KeyColumnQuals["id"].GetStringValue()
	if servicePrincipalID == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	servicePrincipal, err := client.ServicePrincipalsById(servicePrincipalID).Get()
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errObj
	}

	return &ADServicePrincipalInfo{servicePrincipal}, nil
}

func getServicePrincipalOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	servicePrincipal := h.Item.(*ADServicePrincipalInfo)
	servicePrincipalID := servicePrincipal.GetId()

	if servicePrincipalID == nil {
		return nil, nil
	}

	headers := map[string]string{
		"ConsistencyLevel": "eventual",
	}

	includeCount := true
	requestParameters := &owners.OwnersRequestBuilderGetQueryParameters{
		Count: &includeCount,
	}

	config := &owners.OwnersRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParameters,
	}

	ownerIds := []*string{}
	owners, err := client.ServicePrincipalsById(*servicePrincipalID).Owners().GetWithRequestConfigurationAndResponseHandler(config, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(owners, adapter, models.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue)
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		owner := pageItem.(models.DirectoryObjectable)
		ownerIds = append(ownerIds, owner.GetId())

		return true
	})

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
