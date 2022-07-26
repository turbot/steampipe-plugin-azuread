package azuread

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdApplicationTest() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_application_test",
		Description: "Represents an Azure Active Directory (Azure AD) application",
		Get: &plugin.GetConfig{
			Hydrate: getAdApplicationTest,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdApplicationsTest,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "app_id", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "publisher_domain", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the application.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the application.", Transform: transform.FromMethod("GetId")},
			{Name: "app_id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the application that is assigned to an application by Azure AD.", Transform: transform.FromMethod("GetAppId")},

			// Other fields
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Free text field to provide a description of the application object to end users.", Transform: transform.FromMethod("GetDescription")},
			// {Name: "is_authorization_service_enabled", Type: proto.ColumnType_BOOL, Description: "Is authorization service enabled."},
			{Name: "oauth2_require_post_response", Type: proto.ColumnType_BOOL, Description: "Specifies whether, as part of OAuth 2.0 token requests, Azure AD allows POST requests, as opposed to GET requests. The default is false, which specifies that only GET requests are allowed.", Transform: transform.FromMethod("GetOauth2RequirePostResponse"), Default: false},
			{Name: "publisher_domain", Type: proto.ColumnType_STRING, Description: "The verified publisher domain for the application.", Transform: transform.FromMethod("GetPublisherDomain")},
			{Name: "sign_in_audience", Type: proto.ColumnType_STRING, Description: "Specifies the Microsoft accounts that are supported for the current application.", Transform: transform.FromMethod("GetSignInAudience")},

			// Json fields
			{Name: "api", Type: proto.ColumnType_JSON, Description: "Specifies settings for an application that implements a web API.", Transform: transform.FromMethod("ApplicationAPI")},
			{Name: "identifier_uris", Type: proto.ColumnType_JSON, Description: "The URIs that identify the application within its Azure AD tenant, or within a verified custom domain if the application is multi-tenant.", Transform: transform.FromMethod("GetIdentifierUris")},
			{Name: "info", Type: proto.ColumnType_JSON, Description: "Basic profile information of the application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience.", Transform: transform.FromMethod("ApplicationInfo")},
			{Name: "key_credentials", Type: proto.ColumnType_JSON, Description: "The collection of key credentials associated with the application.", Transform: transform.FromMethod("ApplicationKeyCredentials")},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getAdApplicationOwners, Transform: transform.FromValue(), Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "parental_control_settings", Type: proto.ColumnType_JSON, Description: "Specifies parental control settings for an application.", Transform: transform.FromMethod("ApplicationParentalControlSettings")},
			{Name: "password_credentials", Type: proto.ColumnType_JSON, Description: "The collection of password credentials associated with the application.", Transform: transform.FromMethod("ApplicationPasswordCredentials")},
			{Name: "spa", Type: proto.ColumnType_JSON, Description: "Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.", Transform: transform.FromMethod("ApplicationSpa")},
			{Name: "tags_src", Type: proto.ColumnType_JSON, Description: "Custom strings that can be used to categorize and identify the application.", Transform: transform.FromMethod("GetTags")},
			{Name: "web", Type: proto.ColumnType_JSON, Description: "Specifies settings for a web application.", Transform: transform.FromMethod("ApplicationWeb")},

			// Standard columns
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(adApplicationTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adApplicationTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdApplicationsTest(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	// List operations
	input := &applications.ApplicationsRequestBuilderGetQueryParameters{}

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
	filter := buildApplicationQueryFilter(equalQuals)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = &queryFilter
	} else if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &applications.ApplicationsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Applications().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateApplicationCollectionResponseFromDiscriminatorValue)

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		application := pageItem.(models.Applicationable)

		d.StreamListItem(ctx, &ADApplicationInfo{application})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return false
		}

		return true
	})

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAdApplicationTest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	applicationId := d.KeyColumnQuals["id"].GetStringValue()
	if applicationId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	application, err := client.ApplicationsById(applicationId).Get()
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errObj
	}

	return &ADApplicationInfo{application}, nil
}

func getAdApplicationOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	application := h.Item.(*ADApplicationInfo)
	applicationID := application.GetId()

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
	owners, err := client.ApplicationsById(*applicationID).Owners().GetWithRequestConfigurationAndResponseHandler(config, nil)
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

func adApplicationTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	application := d.HydrateItem.(*ADApplicationInfo)
	tags := application.GetTags()
	return TagsToMap(tags)
}

func adApplicationTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADApplicationInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}

func buildApplicationQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name":     "string",
		"app_id":           "string",
		"publisher_domain": "string",
	}

	for qual := range filterQuals {
		if equalQuals[qual] != nil {
			filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
		}
	}

	return filters
}
