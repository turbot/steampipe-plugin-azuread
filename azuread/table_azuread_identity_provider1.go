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
			Hydrate:           getAdIdentityProviderTest,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdIdentityProvidersTest,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the identity provider.", Transform: transform.FromGo()},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the identity provider."},

			// Other fields
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The identity provider type is a required field. For B2B scenario: Google, Facebook. For B2C scenario: Microsoft, Google, Amazon, LinkedIn, Facebook, GitHub, Twitter, Weibo, QQ, WeChat, OpenIDConnect."},
			// {Name: "client_id", Type: proto.ColumnType_STRING, Description: "The client ID for the application. This is the client ID obtained when registering the application with the identity provider."},
			// {Name: "client_secret", Type: proto.ColumnType_STRING, Description: "The client secret for the application. This is the client secret obtained when registering the application with the identity provider. This is write-only. A read operation will return ****."},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
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
		return nil, errors.New(fmt.Sprintf("failed to list identity providers. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateIdentityProviderCollectionResponseFromDiscriminatorValue)

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		identityProvider := pageItem.(models.IdentityProviderBaseable)

		result := map[string]interface{}{
			"ID":   identityProvider.GetId(),
			"Name": identityProvider.GetDisplayName(),
			"Type": identityProvider.GetType(),
			// "ClientId":     identityProvider.GetClientId(),
			// "ClientSecret": identityProvider.GetClientSecret(),
		}

		d.StreamListItem(ctx, result)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return false
		}

		return true
	})

	return nil, nil
}

//// Hydrate Functions

// Need to validate.
// Getting following error
// Code: AADB2C90063 Message: There is a problem with the service.
func getAdIdentityProviderTest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	identityProviderId := d.KeyColumnQuals["id"].GetStringValue()
	if identityProviderId == "" {
		return nil, nil
	}

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
		if isResourceNotFound(errObj) {
			return nil, nil
		}

		return nil, errors.New(fmt.Sprintf("failed to get identity provider. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	result := map[string]interface{}{
		"ID":   identityProvider.GetId(),
		"Name": identityProvider.GetDisplayName(),
		"Type": identityProvider.GetType(),
		// "ClientId":     identityProvider.GetClientId(),
		// "ClientSecret": identityProvider.GetClientSecret(),
	}

	return result, nil
}
