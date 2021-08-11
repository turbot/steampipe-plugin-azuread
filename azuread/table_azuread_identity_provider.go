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
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError,
		},
		List: &plugin.ListConfig{
			Hydrate: listAdIdentityProviders,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "id", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the service principal.", Transform: transform.FromGo()},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name for the service principal"},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for service principals."},

			// Other fields
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The display name exposed by the associated application."},
			{Name: "client_id", Type: proto.ColumnType_STRING, Description: "Specifies the Microsoft accounts that are supported for the current application. Supported values are: AzureADMyOrg, AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, PersonalMicrosoftAccount"},
			{Name: "client_secret", Type: proto.ColumnType_STRING, Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Azure AD internally."},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
			// {Name: "data", Type: proto.ColumnType_JSON, Description: "The unique ID that identifies an active directory user.", Transform: transform.FromValue()}, // For debugging
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

	pagesLeft := true
	for pagesLeft {
		identityProviders, _, err := client.List(ctx)
		if err != nil {
			if isNotFoundError(err) {
				return nil, nil
			}
			return nil, err
		}

		for _, identityProviders := range *identityProviders {
			d.StreamListItem(ctx, identityProviders)
		}
		pagesLeft = false
	}

	return nil, err
}

// Hydrate Functions

func getAdIdentityProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewIdentityProvidersClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	var identityProviderId string
	if h.Item != nil {
		identityProviderId = *h.Item.(msgraph.IdentityProvider).ID
	} else {
		identityProviderId = d.KeyColumnQuals["id"].GetStringValue()
	}

	identityProvider, _, err := client.Get(ctx, identityProviderId)
	if err != nil {
		return nil, err
	}
	return *identityProvider, nil
}
