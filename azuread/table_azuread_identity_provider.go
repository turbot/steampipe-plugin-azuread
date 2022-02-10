package azuread

import (
	"context"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAdIdentityProvider() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_identity_provider",
		Description: "Represents an Azure Active Directory (Azure AD) identity provider",
		Get: &plugin.GetConfig{
			Hydrate:           getAdIdentityProvider,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdIdentityProviders,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the identity provider.", Transform: transform.FromGo()},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the identity provider."},

			// Other fields
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The identity provider type is a required field. For B2B scenario: Google, Facebook. For B2C scenario: Microsoft, Google, Amazon, LinkedIn, Facebook, GitHub, Twitter, Weibo, QQ, WeChat, OpenIDConnect."},
			{Name: "client_id", Type: proto.ColumnType_STRING, Description: "The client ID for the application. This is the client ID obtained when registering the application with the identity provider."},
			{Name: "client_secret", Type: proto.ColumnType_STRING, Description: "The client secret for the application. This is the client secret obtained when registering the application with the identity provider. This is write-only. A read operation will return ****."},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdIdentityProviders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewIdentityProvidersClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	identityProviders, _, err := client.List(ctx)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, identityProviders := range *identityProviders {
		d.StreamListItem(ctx, identityProviders)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// Hydrate Functions

func getAdIdentityProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var identityProviderId string
	if h.Item != nil {
		identityProviderId = *h.Item.(msgraph.IdentityProvider).ID
	} else {
		identityProviderId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if identityProviderId == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewIdentityProvidersClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer
	client.BaseClient.DisableRetries = true

	identityProvider, _, err := client.Get(ctx, identityProviderId)
	if err != nil {
		return nil, err
	}
	return *identityProvider, nil
}
