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

func tableAzureAdApplication() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_application",
		Description: "Represents an Azure Active Directory (Azure AD) application",
		Get: &plugin.GetConfig{
			Hydrate:    getAdApplication,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdApplications,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Request_UnsupportedQuery"}),
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "app_id", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "publisher_domain", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the application."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the application.", Transform: transform.FromGo()},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the application that is assigned to an application by Azure AD."},

			// // Other fields
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time."},
			{Name: "is_authorization_service_enabled", Type: proto.ColumnType_BOOL, Description: "Is authorization service enabled."},
			{Name: "oauth2_require_post_response", Type: proto.ColumnType_BOOL, Description: "Specifies whether, as part of OAuth 2.0 token requests, Azure AD allows POST requests, as opposed to GET requests. The default is false, which specifies that only GET requests are allowed."},
			{Name: "publisher_domain", Type: proto.ColumnType_STRING, Description: "The verified publisher domain for the application."},
			{Name: "sign_in_audience", Type: proto.ColumnType_STRING, Description: "Specifies the Microsoft accounts that are supported for the current application."},

			// // Json fields
			{Name: "api", Type: proto.ColumnType_JSON, Description: "Specifies settings for an application that implements a web API."},
			{Name: "identifier_uris", Type: proto.ColumnType_JSON, Description: "The URIs that identify the application within its Azure AD tenant, or within a verified custom domain if the application is multi-tenant. "},
			{Name: "info", Type: proto.ColumnType_JSON, Description: "Basic profile information of the application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience."},
			{Name: "key_credentials", Type: proto.ColumnType_JSON, Description: "The collection of key credentials associated with the application."},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getApplicationOwners, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "parental_control_settings", Type: proto.ColumnType_JSON, Description: "Specifies parental control settings for an application."},
			{Name: "password_credentials", Type: proto.ColumnType_JSON, Description: "The collection of password credentials associated with the application."},
			{Name: "spa", Type: proto.ColumnType_JSON, Description: "Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens."},
			{Name: "tags_src", Type: proto.ColumnType_JSON, Description: "Custom strings that can be used to categorize and identify the application.", Transform: transform.FromField("Tags")},
			{Name: "web", Type: proto.ColumnType_JSON, Description: "Specifies settings for a web application."},

			// // Standard columns
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(applicationTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewApplicationsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	input := odata.Query{}

	qualsColumnMap := []QualsColumn{
		{"app_id", "string", "appId"},
		{"display_name", "string", "displayName"},
		{"publisher_domain", "string", "publisherDomain"},
	}

	filter := buildCommaonQueryFilter(qualsColumnMap, d.Quals)
	if len(filter) > 0 {
		input.Filter = strings.Join(filter, " and ")
	}

	applications, _, err := client.List(ctx, input)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, application := range *applications {
		d.StreamListItem(ctx, application)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// Hydrate Functions

func getAdApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var applicationId string
	if h.Item != nil {
		applicationId = *h.Item.(msgraph.Application).ID
	} else {
		applicationId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if applicationId == "" {
		return nil, nil
	}
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewApplicationsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer
	client.BaseClient.DisableRetries = true

	application, _, err := client.Get(ctx, applicationId, odata.Query{})
	if err != nil {
		if isNotFoundError(err){
			return nil, nil
		}
		return nil, err
	}
	return *application, nil
}

func getApplicationOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	application := h.Item.(msgraph.Application)
	session, err := GetNewSession(ctx, d)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	client := msgraph.NewApplicationsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	owners, _, err := client.ListOwners(ctx, *application.ID)
	if err != nil {
		return nil, err
	}
	return owners, nil
}

//// Transform Function

func applicationTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	application := d.HydrateItem.(msgraph.Application)
	return TagsToMap(*application.Tags)
}
