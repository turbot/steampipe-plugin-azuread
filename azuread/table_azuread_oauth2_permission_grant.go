package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/oauth2permissiongrants"
)

//// TABLE DEFINITION

func tableAzureAdOAuth2PermissionGrant(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_oauth2_permission_grant",
		Description: "Represents the delegated permissions that have been granted to an application's service principal, either for all users (admin consent) or for a single user.",
		Get: &plugin.GetConfig{
			Hydrate: getAdOAuth2PermissionGrant,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier", "Request_UnsupportedQuery"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdOAuth2PermissionGrants,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "client_id", Require: plugin.Optional},
				{Name: "consent_type", Require: plugin.Optional},
				{Name: "principal_id", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the delegated permission grant.", Transform: transform.FromMethod("GetId")},
			{Name: "client_id", Type: proto.ColumnType_STRING, Description: "The object ID (not the app ID) of the client service principal that is authorized to act on behalf of a signed-in user when accessing an API.", Transform: transform.FromMethod("GetClientId")},
			{Name: "consent_type", Type: proto.ColumnType_STRING, Description: "Whether the grant applies to all users (AllPrincipals, admin consent) or to a single user (Principal).", Transform: transform.FromMethod("GetConsentType")},
			{Name: "principal_id", Type: proto.ColumnType_STRING, Description: "The ID of the user on whose behalf the client is authorized to access the resource, when consent_type is Principal. Null when consent_type is AllPrincipals.", Transform: transform.FromMethod("GetPrincipalId")},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "The object ID of the resource service principal (the API, for example Microsoft Graph) to which access is authorized.", Transform: transform.FromMethod("GetResourceId")},
			{Name: "scope", Type: proto.ColumnType_STRING, Description: "A space-separated list of the claim values for delegated permissions granted, for example 'openid User.Read Files.Read.All'.", Transform: transform.FromMethod("GetScope")},

			// JSON fields
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "The delegated permissions granted, as an array.", Transform: transform.FromMethod("GetScope").Transform(splitScopes)},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetId")},
		}),
	}
}

type ADOAuth2PermissionGrantInfo struct {
	models.OAuth2PermissionGrantable
}

//// LIST FUNCTION

func listAdOAuth2PermissionGrants(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_oauth2_permission_grant.listAdOAuth2PermissionGrants", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &oauth2permissiongrants.Oauth2PermissionGrantsRequestBuilderGetQueryParameters{
		Top: Int32(999),
	}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	// unexpected status 400 with OData error: Request_UnsupportedQuery: Invalid page size specified: '1000'. Must be between 1 and 999 inclusive.
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit < 999 {
			l := int32(*limit)
			input.Top = Int32(l)
		}
	}

	// The API supports $filter with eq on clientId, consentType, principalId and resourceId
	filterColumns := map[string]string{
		"client_id":    "clientId",
		"consent_type": "consentType",
		"principal_id": "principalId",
		"resource_id":  "resourceId",
	}
	var filter []string
	for column, property := range filterColumns {
		if d.EqualsQuals[column] != nil {
			value := strings.ReplaceAll(d.EqualsQuals[column].GetStringValue(), "'", "''")
			filter = append(filter, fmt.Sprintf("%s eq '%s'", property, value))
		}
	}
	if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &oauth2permissiongrants.Oauth2PermissionGrantsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Oauth2PermissionGrants().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdOAuth2PermissionGrants", "list_oauth2_permission_grant_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.OAuth2PermissionGrantable](result, adapter, models.CreateOAuth2PermissionGrantCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdOAuth2PermissionGrants", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.OAuth2PermissionGrantable) bool {
		d.StreamListItem(ctx, &ADOAuth2PermissionGrantInfo{pageItem})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdOAuth2PermissionGrants", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdOAuth2PermissionGrant(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	grantId := d.EqualsQuals["id"].GetStringValue()
	if grantId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_oauth2_permission_grant.getAdOAuth2PermissionGrant", "connection_error", err)
		return nil, err
	}

	grant, err := client.Oauth2PermissionGrants().ByOAuth2PermissionGrantId(grantId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdOAuth2PermissionGrant", "get_oauth2_permission_grant_error", errObj)
		return nil, errObj
	}

	return &ADOAuth2PermissionGrantInfo{grant}, nil
}

//// TRANSFORM FUNCTIONS

// splitScopes turns the space-separated scope string into an array.
func splitScopes(_ context.Context, d *transform.TransformData) (interface{}, error) {
	scope, ok := d.Value.(*string)
	if !ok || scope == nil {
		return nil, nil
	}
	return strings.Fields(*scope), nil
}
