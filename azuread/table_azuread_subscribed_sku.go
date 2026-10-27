package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdSubscribedSku(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_subscribed_sku",
		Description: "Represents the commercial subscriptions (SKUs) that an organization has acquired.",
		List: &plugin.ListConfig{
			Hydrate: listAdSubscribedSkus,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "sku_id", Type: proto.ColumnType_STRING, Description: "The unique identifier (GUID) for the service SKU.", Transform: transform.FromMethod("SkuId")},
			{Name: "sku_part_number", Type: proto.ColumnType_STRING, Description: "The SKU part number, for example ENTERPRISEPACK.", Transform: transform.FromMethod("GetSkuPartNumber")},
			{Name: "capability_status", Type: proto.ColumnType_STRING, Description: "The status of the SKU. Possible values are Enabled, Warning, Suspended, Deleted, LockedOut.", Transform: transform.FromMethod("GetCapabilityStatus")},

			// JSON fields
			{Name: "service_plans", Type: proto.ColumnType_JSON, Description: "Information about the service plans that are available with the SKU.", Transform: transform.FromMethod("ServicePlans")},
		}),
	}
}

//// LIST FUNCTION

func listAdSubscribedSkus(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_subscribed_sku.listAdSubscribedSkus", "connection_error", err)
		return nil, err
	}

	result, err := client.SubscribedSkus().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_subscribed_sku.listAdSubscribedSkus", "list_subscribed_sku_error", errObj)
		return nil, errObj
	}

	for _, sku := range result.GetValue() {
		if sku == nil {
			continue
		}
		d.StreamListItem(ctx, &ADSubscribedSkuInfo{sku})

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
