package azuread

import (
	"context"
	"strings"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureAdServicePrincipal() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_service_principal",
		Description: "Represents an Azure Active Directory (Azure AD) service principal",
		Get: &plugin.GetConfig{
			Hydrate:           getAdServicePrincipal,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate:           listAdServicePrincipals,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Request_UnsupportedQuery"}),
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "display_name", Require: plugin.Optional},
				{Name: "account_enabled", Require: plugin.Optional},
				{Name: "service_principal_type", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the service principal.", Transform: transform.FromGo()},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the service principal."},

			// Other fields
			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "true if the service principal account is enabled; otherwise, false."},
			{Name: "app_display_name", Type: proto.ColumnType_STRING, Description: "The display name exposed by the associated application."},
			{Name: "app_owner_organization_id", Type: proto.ColumnType_STRING, Description: "Contains the tenant id where the application is registered. This is applicable only to service principals backed by applications."},
			{Name: "app_role_assignment_required", Type: proto.ColumnType_BOOL, Description: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false."},
			{Name: "service_principal_type", Type: proto.ColumnType_STRING, Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Azure AD internally."},
			{Name: "sign_in_audience", Type: proto.ColumnType_STRING, Description: "Specifies the Microsoft accounts that are supported for the current application. Supported values are: AzureADMyOrg, AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, PersonalMicrosoftAccount."},

			// // Json fields
			{Name: "add_ins", Type: proto.ColumnType_JSON, Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts."},
			{Name: "alternative_names", Type: proto.ColumnType_JSON, Description: "Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities."},
			{Name: "app_roles", Type: proto.ColumnType_JSON, Description: "The roles exposed by the application which this service principal represents."},
			{Name: "info", Type: proto.ColumnType_JSON, Description: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs."},
			{Name: "keyCredentials", Type: proto.ColumnType_JSON, Description: "The collection of key credentials associated with the service principal."},
			{Name: "notification_email_addresses", Type: proto.ColumnType_JSON, Description: "Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications."},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getAdServicePrincipalOwners, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "password_credentials", Type: proto.ColumnType_JSON, Description: "Represents a password credential associated with a service principal."},
			{Name: "published_permission_scopes", Type: proto.ColumnType_JSON, Description: "The published permission scopes."},
			{Name: "reply_urls", Type: proto.ColumnType_JSON, Description: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application."},
			{Name: "service_principal_names", Type: proto.ColumnType_JSON, Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD."},
			{Name: "tags_src", Type: proto.ColumnType_JSON, Description: "Custom strings that can be used to categorize and identify the service principal.", Transform: transform.FromField("Tags")},
			{Name: "verified_publisher", Type: proto.ColumnType_JSON, Description: "Specifies the verified publisher of the application which this service principal represents."},

			// // Standard columns
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(servicePrincipalTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdServicePrincipals(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewServicePrincipalsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	input := odata.Query{
		Top: 999,
	}

	qualsColumnMap := []QualsColumn{
		{"display_name", "string", "displayName"},
		{"account_enabled", "bool", "accountEnabled"},
		{"service_principal_type", "string", "servicePrincipalType"},
	}

	filter := buildCommaonQueryFilter(qualsColumnMap, d.Quals)
	if len(filter) > 0 {
		input.Filter = strings.Join(filter, " and ")
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < 999 {
			input.Top = int(*limit)
		}
	}

	servicePrincipals, _, err := client.List(ctx, input)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, servicePrincipal := range *servicePrincipals {
		d.StreamListItem(ctx, servicePrincipal)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// Hydrate Functions

func getAdServicePrincipal(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var servicePrincipalId string
	if h.Item != nil {
		servicePrincipalId = *h.Item.(msgraph.ServicePrincipal).ID
	} else {
		servicePrincipalId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if servicePrincipalId == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewServicePrincipalsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer
	client.BaseClient.DisableRetries = true

	servicePrincipal, _, err := client.Get(ctx, servicePrincipalId, odata.Query{})
	if err != nil {
		return nil, err
	}
	return *servicePrincipal, nil
}

func getAdServicePrincipalOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	servicePrincipalId := *h.Item.(msgraph.ServicePrincipal).ID
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewServicePrincipalsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	owners, _, err := client.ListOwners(ctx, servicePrincipalId)
	if err != nil {
		return nil, err
	}
	return owners, nil
}

//// Transform Function

func servicePrincipalTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	servicePrincipal := d.HydrateItem.(msgraph.ServicePrincipal)
	return TagsToMap(*servicePrincipal.Tags)
}
