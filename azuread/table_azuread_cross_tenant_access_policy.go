package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureAdCrossTenantAccessPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_cross_tenant_access_policy",
		Description: "Represents an Azure Active Directory (Azure AD) Cross-Tenant Access Policy.",
		List: &plugin.ListConfig{
			Hydrate: listAdCrossTenantAccessPolicies,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery"}),
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Specifies the identifier of a crossTenantAccessPolicy object.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Specifies a display name for the crossTenantAccessPolicy object.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "allowed_cloud_endpoints", Type: proto.ColumnType_JSON, Description: "Used to specify which Microsoft clouds an organization would like to collaborate with. By default, this value is empty. Supported values are: microsoftonline.com, microsoftonline.us, and partner.microsoftonline.cn.", Transform: transform.FromMethod("CrossTenantAccessPolicyAllowedCloudEndpoints")},
			{Name: "default_configuration", Type: proto.ColumnType_JSON, Description: "Defines the default configuration for how your organization interacts with external Microsoft Entra organizations.", Transform: transform.FromMethod("CrossTenantAccessPolicyDefault")},
			{Name: "partners", Type: proto.ColumnType_JSON, Description: "Defines partner-specific configurations for external Microsoft Entra organizations.", Transform: transform.FromMethod("CrossTenantAccessPolicyPartners")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adCrossTenantAccessPolicyTitle)},
		}),
	}
}

//// LIST FUNCTION

func listAdCrossTenantAccessPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_cross_tenant_access_policy.listAdCrossTenantAccessPolicies", "connection_error", err)
		return nil, err
	}

	// Get the cross-tenant access policy
	policy, err := client.Policies().CrossTenantAccessPolicy().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdCrossTenantAccessPolicies", "get_cross_tenant_access_policy_error", errObj)
		return nil, errObj
	}

	// Stream the single policy item
	d.StreamListItem(ctx, &ADCrossTenantAccessPolicyInfo{policy})

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func adCrossTenantAccessPolicyTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADCrossTenantAccessPolicyInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}
