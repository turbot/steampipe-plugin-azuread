package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/identity"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

//// TABLE DEFINITION

func tableAzureAdConditionalAccessNamedLocation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_conditional_access_named_location",
		Description: "Represents an Azure Active Directory (Azure AD) Conditional Access Named Location.",
		Get: &plugin.GetConfig{
			Hydrate: getAdConditionalAccessNamedLocation,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdConditionalAccessNamedLocations,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "display_name", Require: plugin.Optional},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Specifies the identifier of a Named Location object.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Specifies a display name for the Named Location object.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The create date of the Named Location object.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "modified_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The modification date of Named Location object.", Transform: transform.FromMethod("GetModifiedDateTime")},
		}),
	}
}

//// LIST FUNCTION

func listAdConditionalAccessNamedLocations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_conditional_access_named_location.listAdConditionalAccessNamedLocations", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &identity.ConditionalAccessNamedLocationsRequestBuilderGetQueryParameters{
		Top: Int32(1000),
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit < 1000 {
			l := int32(*limit)
			input.Top = Int32(l)
		}
	}

	equalQuals := d.EqualsQuals
	filter := buildConditionalAccessNamedLocationQueryFilter(equalQuals)

	if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &identity.ConditionalAccessNamedLocationsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Identity().ConditionalAccess().NamedLocations().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdConditionalAccessNamedLocations", "list_conditional_access_named_location_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.NamedLocationable](result, adapter, models.CreateNamedLocationCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdConditionalAccessNamedLocations", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.NamedLocationable) bool {
		d.StreamListItem(ctx, &ADLocationInfo{pageItem})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdConditionalAccessNamedLocations", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdConditionalAccessNamedLocation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	conditionalAccessNamedLocationId := d.EqualsQuals["id"].GetStringValue()
	if conditionalAccessNamedLocationId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_conditional_access_named_location.getAdConditionalAccessNamedLocation", "connection_error", err)
		return nil, err
	}

	location, err := client.Identity().ConditionalAccess().NamedLocations().ByNamedLocationId(conditionalAccessNamedLocationId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdConditionalAccessNamedLocation", "get_conditional_access_location_error", errObj)
		return nil, errObj
	}
	return &ADLocationInfo{location}, nil
}

func buildConditionalAccessNamedLocationQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name": "string",
		"state":        "string",
	}

	for qual, qualType := range filterQuals {
		switch qualType {
		case "string":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
			}
		}
	}

	return filters
}

//// TRANSFORM FUNCTIONS

func adConditionalAccessNamedLocationTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADLocationInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
}
