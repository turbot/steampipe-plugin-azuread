package azuread

import (
	"context"

	// msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	// "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdExternalIdentityPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_external_identity_policy",
		Description: "Represents the tenant-wide external identity policy that controls whether external users can leave a Microsoft Entra tenant via self-service controls.",
		List: &plugin.ListConfig{
			Hydrate: listAdExternalIdentityPolicies,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the external identity policy.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name for the external identity policy.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "allow_external_identities_to_leave", Type: proto.ColumnType_BOOL, Description: "Flag indicating whether external users can leave the tenant via self-service controls.", Transform: transform.FromMethod("GetAllowExternalIdentitiesToLeave")},
			{Name: "allow_deleted_identities_data_removal", Type: proto.ColumnType_BOOL, Description: "Flag indicating whether deleted identities data can be removed.", Transform: transform.FromMethod("GetAllowDeletedIdentitiesDataRemoval")},
			{Name: "deleted_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the policy was deleted.", Transform: transform.FromMethod("GetDeletedDateTime")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetDisplayName")},
		}),
	}
}

//// LIST FUNCTION

func listAdExternalIdentityPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create beta client
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_external_identity_policy.listAdExternalIdentityPolicies", "connection_error", err)
		return nil, err
	}

	// Get the external identities policy
	// Note: This is a singleton policy, so we return it as a single item
	policy, err := client.Policies().ExternalIdentitiesPolicy().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdExternalIdentityPolicies", "get_external_identity_policy_error", errObj)
		return nil, errObj
	}

	d.StreamListItem(ctx, &ADExternalIdentitiesPolicyInfo{policy})

	return nil, nil
}
