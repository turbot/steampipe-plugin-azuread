package azuread

import (
	"context"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/identity/identityproviders"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureAdIdentityProviderTest() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_identity_provider_test",
		Description: "Represents an Azure Active Directory (Azure AD) identity provider",
		Get: &plugin.GetConfig{
			Hydrate: getAdIdentityProviderTest,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdIdentityProvidersTest,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the identity provider.", Transform: transform.FromMethod("GetId")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the identity provider.", Transform: transform.FromMethod("GetDisplayName")},

			// Other fields
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The identity provider type is a required field. For B2B scenario: Google, Facebook. For B2C scenario: Microsoft, Google, Amazon, LinkedIn, Facebook, GitHub, Twitter, Weibo, QQ, WeChat, OpenIDConnect.", Transform: transform.FromMethod("GetIdentityProviderType")},
			// {Name: "client_id", Type: proto.ColumnType_STRING, Description: "The client ID for the application. This is the client ID obtained when registering the application with the identity provider.", Transform: transform.FromMethod("GetClientId")},
			// {Name: "client_secret", Type: proto.ColumnType_STRING, Description: "The client secret for the application. This is the client secret obtained when registering the application with the identity provider. This is write-only. A read operation will return ****.", Transform: transform.FromMethod("GetClientSecret")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adIdentityProviderTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdIdentityProvidersTest(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	// List operations
	input := &identityproviders.IdentityProvidersRequestBuilderGetQueryParameters{}

	limit := d.QueryContext.Limit
	if limit != nil {
		l := int32(*limit)
		input.Top = &l
	}

	options := &identityproviders.IdentityProvidersRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Identity().IdentityProviders().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateBuiltInIdentityProviderFromDiscriminatorValue)

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		identityProvider := pageItem.(*models.BuiltInIdentityProvider)

		d.StreamListItem(ctx, identityProvider)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return false
		}

		return true
	})

	return nil, nil
}

//// HYDRATE FUNCTIONS
// https://docs.microsoft.com/en-us/graph/api/identityproviderbase-get?view=graph-rest-1.0&tabs=go

func getAdIdentityProviderTest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityProviderId := d.KeyColumnQuals["id"].GetStringValue()
	if identityProviderId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	identityProvider, err := client.Identity().IdentityProvidersById(identityProviderId).Get()
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errObj
	}

	return identityProvider, nil
}

//// TRANSFORM FUNCTIONS

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
