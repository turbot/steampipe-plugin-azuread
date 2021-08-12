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
	}

	return nil, err
}

// Hydrate Functions

// we didn't add the get function as it retries 5 times on 404 errors
