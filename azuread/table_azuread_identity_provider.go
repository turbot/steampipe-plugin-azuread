package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/identity/identityproviders"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAzureAdIdentityProvider(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_identity_provider",
		Description: "Represents an Azure Active Directory (Azure AD) identity provider.",
		List: &plugin.ListConfig{
			Hydrate: listAdIdentityProviders,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery", "Invalid filter clause"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "id", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the identity provider.", Transform: transform.FromMethod("GetId")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the identity provider.", Transform: transform.FromMethod("GetDisplayName")},

			// Other fields
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The identity provider type is a required field. For B2B scenario: Google, Facebook. For B2C scenario: Microsoft, Google, Amazon, LinkedIn, Facebook, GitHub, Twitter, Weibo, QQ, WeChat, OpenIDConnect.", Transform: transform.FromMethod("GetIdentityProviderType")},
			{Name: "client_id", Type: proto.ColumnType_STRING, Description: "The client ID for the application. This is the client ID obtained when registering the application with the identity provider."},
			{Name: "client_secret", Type: proto.ColumnType_STRING, Description: "The client secret for the application. This is the client secret obtained when registering the application with the identity provider. This is write-only. A read operation will return ****."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for resources."},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adIdentityProviderTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdIdentityProviders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_identity_provider.listAdIdentityProviders", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &identityproviders.IdentityProvidersRequestBuilderGetQueryParameters{}

	limit := d.QueryContext.Limit
	if limit != nil {
		l := int32(*limit)
		input.Top = &l
	}

	var queryFilter string
	equalQuals := d.KeyColumnQuals
	filter := buildIdentityProviderQueryFilter(equalQuals)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = &queryFilter
	} else if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &identityproviders.IdentityProvidersRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Identity().IdentityProviders().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdIdentityProviders", "list_identity_provider_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateBuiltInIdentityProviderFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdIdentityProviders", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		identityProvider := pageItem.(*models.BuiltInIdentityProvider)

		clientID := identityProvider.GetAdditionalData()["clientId"]
		clientSecret := identityProvider.GetAdditionalData()["clientSecret"]

		d.StreamListItem(ctx, &ADIdentityProviderInfo{*identityProvider, clientID, clientSecret})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.QueryStatus.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdIdentityProviders", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func buildIdentityProviderQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"id":   "string",
		"name": "string",
	}

	for qual := range filterQuals {
		if equalQuals[qual] != nil {
			if qual == "name" {
				filters = append(filters, fmt.Sprintf("displayName eq '%s'", equalQuals[qual].GetStringValue()))
			} else {
				filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
			}
		}
	}

	return filters
}

func adIdentityProviderTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(models.IdentityProviderBaseable)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}
