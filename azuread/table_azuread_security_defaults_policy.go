package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdSecurityDefaultsPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_security_defaults_policy",
		Description: "Represents the Azure Active Directory security defaults policy",
		List: &plugin.ListConfig{
			Hydrate: listAdSecurityDefaultPolicies,
		},

		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Display name for this policy.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier for this policy.", Transform: transform.FromMethod("GetId")},
			{Name: "is_enabled", Type: proto.ColumnType_BOOL, Description: "If set to true, Azure Active Directory security defaults is enabled for the tenant.", Transform: transform.FromMethod("GetIsEnabled")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description for this policy.", Transform: transform.FromMethod("GetDescription")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetDisplayName")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdSecurityDefaultPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_security_defaults_policy.listAdSecurityDefaultPolicies", "connection_error", err)
		return nil, err
	}

	result, err := client.Policies().IdentitySecurityDefaultsEnforcementPolicy().Get(ctx, nil)
	if err != nil {
		errObj := getBetaErrorObject(err)
		plugin.Logger(ctx).Error("listAdSecurityDefaultPolicies", "list_security_defaults_policy_error", errObj)
		return nil, errObj
	}
	d.StreamListItem(ctx, &ADSecurityDefaultsPolicyInfo{result})

	return nil, nil
}
