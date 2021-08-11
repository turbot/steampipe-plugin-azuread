package azuread

import (
	"context"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAdServicePrincipal() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_service_principal",
		Description: "Represents an Azure Active Directory (Azure AD) service principal",
		Get: &plugin.GetConfig{
			Hydrate:           getAdServicePrincipal,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError,
		},
		List: &plugin.ListConfig{
			Hydrate: listAdServicePrincipals,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "id", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},

				// Other fields for filtering OData
				// {Name: "mail", Require: plugin.Optional, Operators: []string{"<>", "="}},                     // $filter (eq, ne, NOT, ge, le, in, startsWith).
				// {Name: "mail_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},             // $filter (eq, ne, NOT).
				// {Name: "on_premises_sync_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}}, // $filter (eq, ne, NOT, in).
				// {Name: "security_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},         // $filter (eq, ne, NOT, in).

				// TODO
				// {Name: "created_date_time", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},    // Supports $filter (eq, ne, NOT, ge, le, in).
				// {Name: "expiration_date_time", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional}, // Supports $filter (eq, ne, NOT, ge, le, in).
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the service principal.", Transform: transform.FromGo()},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the service principal"},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for service principals."},

			// Other fields
			{Name: "app_display_name", Type: proto.ColumnType_STRING, Description: "The display name exposed by the associated application."},
			{Name: "sign_in_audience", Type: proto.ColumnType_STRING, Description: "Specifies the Microsoft accounts that are supported for the current application. Supported values are: AzureADMyOrg, AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, PersonalMicrosoftAccount"},
			{Name: "service_principal_type", Type: proto.ColumnType_STRING, Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Azure AD internally."},
			{Name: "app_owner_organization_id", Type: proto.ColumnType_STRING, Description: "Contains the tenant id where the application is registered. This is applicable only to service principals backed by applications."},
			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "true if the service principal account is enabled; otherwise, false"},
			{Name: "app_role_assignment_required", Type: proto.ColumnType_BOOL, Description: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false."},

			// // Json fields
			{Name: "info", Type: proto.ColumnType_JSON, Description: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "Custom strings that can be used to categorize and identify the service principal."},
			{Name: "add_ins", Type: proto.ColumnType_JSON, Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts."},
			{Name: "app_roles", Type: proto.ColumnType_JSON, Description: "The roles exposed by the application which this service principal represents."},
			{Name: "reply_urls", Type: proto.ColumnType_JSON, Description: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application."},
			{Name: "keyCredentials", Type: proto.ColumnType_JSON, Description: "The collection of key credentials associated with the service principal."},
			{Name: "alternative_names", Type: proto.ColumnType_JSON, Description: "Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities."},
			{Name: "verified_publisher", Type: proto.ColumnType_JSON, Description: "Specifies the verified publisher of the application which this service principal represents."},
			{Name: "password_credentials", Type: proto.ColumnType_JSON, Description: "Represents a password credential associated with a service principal."},
			{Name: "service_principal_names", Type: proto.ColumnType_JSON, Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD."},
			{Name: "published_permission_scopes", Type: proto.ColumnType_JSON, Description: "The published permission scopes."},
			{Name: "notification_email_addresses", Type: proto.ColumnType_JSON, Description: "Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications."},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getAdServicePrincipalOwners, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},

			// // Standard columns
			// {Name: "tags", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTags, Transform: transform.From(applicationTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
			// {Name: "data", Type: proto.ColumnType_JSON, Description: "The unique ID that identifies an active directory user.", Transform: transform.FromValue()}, // For debugging
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

	// TODO filters
	input := odata.Query{}
	filter := ""
	input.Filter = filter

	pagesLeft := true
	for pagesLeft {
		servicePrincipals, _, err := client.List(ctx, input)
		if err != nil {
			if isNotFoundError(err) {
				return nil, nil
			}
			return nil, err
		}

		for _, servicePrincipals := range *servicePrincipals {
			d.StreamListItem(ctx, servicePrincipals)
		}
		pagesLeft = false
	}

	return nil, err
}

// Hydrate Functions

func getAdServicePrincipal(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewServicePrincipalsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	var servicePrincipalId string
	if h.Item != nil {
		servicePrincipalId = *h.Item.(msgraph.ServicePrincipal).ID
	} else {
		servicePrincipalId = d.KeyColumnQuals["id"].GetStringValue()
	}

	// TODO filters
	input := odata.Query{}
	filter := ""
	input.Filter = filter


	servicePrincipal, _, err := client.Get(ctx, servicePrincipalId, input)
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
