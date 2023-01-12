package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

//// TABLE DEFINITION

func tableAzureAdLicense(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_license",
		Description: "Azure Active Directory comes in four editions—Free, Office 365 apps, Premium P1, and Premium P2. The Free edition is included with a subscription of a commercial online service, e.g. Azure, Dynamics 365, Intune and Power Platform.",
		List: &plugin.ListConfig{
			Hydrate: listAdLicenses,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique ID representing the license.", Transform: transform.FromMethod("GetId")},
			{Name: "sku_id", Type: proto.ColumnType_STRING, Description: "The SKU id of the license.", Transform: transform.FromMethod("GetSkuId")},
			{Name: "sku_part_number", Type: proto.ColumnType_STRING, Description: "The SKU name of the license..", Transform: transform.FromMethod("GetSkuPartNumber")},
			{Name: "service_plans", Type: proto.ColumnType_JSON, Description: "Service plan details of the license.", Transform: transform.FromMethod("ServicePlans")},
		},
	}
}

//// LIST FUNCTION

func listAdLicenses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_license.listAdLicenses", "connection_error", err)
		return nil, err
	}

	result, err := client.UsersById("steampipe@turbotoffice.onmicrosoft.com").LicenseDetails().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_license.listAdLicenses", "api_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateLicenseDetailsCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_license.listAdLicenses", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		license := pageItem.(models.LicenseDetailsable)
		d.StreamListItem(ctx, &ADLicenseInfo{license})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.QueryStatus.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuread_license.listAdLicenses", "paging_error", err)
		return nil, err
	}

	return nil, nil
}
