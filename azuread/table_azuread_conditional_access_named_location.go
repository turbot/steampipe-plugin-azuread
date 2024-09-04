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
				{Name: "id", Require: plugin.Optional},
				{Name: "location_type", Require: plugin.Optional},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Specifies the identifier of a Named Location object.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Specifies a display name for the Named Location object.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "location_type", Type: proto.ColumnType_STRING, Description: "Specifies the type of the Named Location object: IP or Country.", Transform: transform.FromMethod("GetType")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The create date of the Named Location object.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "modified_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The modification date of Named Location object.", Transform: transform.FromMethod("GetModifiedDateTime")},
			{Name: "location_info", Type: proto.ColumnType_JSON, Description: "Specifies some location information for the Named Location object. Now supported: IP (v4/6 and CIDR/Range), odata_type, IsTrusted (for IP named locations only). Country (and regions, if exist), lookup method, UnkownCountriesAndRegions (for country named locations only).", Transform: transform.FromMethod("GetLocationInfo")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetDisplayName")},
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
		plugin.Logger(ctx).Error("azuread_conditional_access_named_location.listAdConditionalAccessNamedLocations", "list_conditional_access_named_location_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.NamedLocationable](result, adapter, models.CreateNamedLocationCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_conditional_access_named_location.listAdConditionalAccessNamedLocations", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem models.NamedLocationable) bool {
		d.StreamListItem(ctx, ADNamedLocationInfo{
			NamedLocationable: pageItem,
			NamedLocation:     getNamedLocationDetails(pageItem),
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuread_conditional_access_named_location.listAdConditionalAccessNamedLocations", "paging_error", err)
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
		plugin.Logger(ctx).Error("azuread_conditional_access_named_location.getAdConditionalAccessNamedLocation", "get_conditional_access_location_error", errObj)
		return nil, errObj
	}

	return &ADNamedLocationInfo{
		NamedLocationable: location,
		NamedLocation:     getNamedLocationDetails(location),
	}, nil
}

func buildConditionalAccessNamedLocationQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name": "string",
		"id":           "string",
	}

	for qual, qualType := range filterQuals {
		switch qualType {
		case "string":
			if equalQuals[qual] != nil {
				if qual == "location_type" {
					filters = append(filters, fmt.Sprintf("type eq '%s'", equalQuals[qual].GetStringValue()))
				} else {
					filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToLowerCamel(qual), equalQuals[qual].GetStringValue()))
				}
			}
		}
	}

	return filters
}

/// UTILITY FUNCTION

func getNamedLocationDetails(i interface{}) models.NamedLocationable {

	switch t := i.(type) {
	case *models.IpNamedLocation:
		return ADIpNamedLocationInfo{t}
	case *models.CountryNamedLocation:
		return ADCountryNamedLocationInfo{t}
	}

	return nil
}

//// TRANSFORM FUNCTIONS

func IpGetLocationInfo(ipLocationInfo *ADIpNamedLocationInfo) map[string]interface{} {
	ipRangesArray := ipLocationInfo.GetIpRanges()
	locationInfoJSON := map[string]interface{}{}

	IPv4CidrArr := []map[string]interface{}{}
	IPv4RangeArr := []map[string]interface{}{}
	IPv6CidrArr := []map[string]interface{}{}
	IPv6RangeArr := []map[string]interface{}{}

	for i := 0; i < len(ipRangesArray); i++ {
		switch t := ipRangesArray[i].(type) {
		case *models.IPv4CidrRange:
			IPv4CidrPair := map[string]interface{}{}
			IPv4CidrPair["Address"] = *t.GetCidrAddress()
			IPv4CidrArr = append(IPv4CidrArr, IPv4CidrPair)
		case *models.IPv4Range:
			IPv4AddressPair := map[string]interface{}{}
			IPv4AddressPair["Lower"] = *t.GetLowerAddress()
			IPv4AddressPair["Upper"] = *t.GetUpperAddress()
			IPv4RangeArr = append(IPv4RangeArr, IPv4AddressPair)
		case *models.IPv6CidrRange:
			IPv6CidrPair := map[string]interface{}{}
			IPv6CidrPair["Address"] = *t.GetCidrAddress()
			IPv6CidrArr = append(IPv6CidrArr, IPv6CidrPair)
		case *models.IPv6Range:
			IPv6AddressPair := map[string]interface{}{}
			IPv6AddressPair["Lower"] = *t.GetLowerAddress()
			IPv6AddressPair["Upper"] = *t.GetUpperAddress()
			IPv6RangeArr = append(IPv6RangeArr, IPv6AddressPair)
		}
	}

	locationInfoJSON["IPv4Cidr"] = IPv4CidrArr
	locationInfoJSON["IPv4Range"] = IPv4RangeArr
	locationInfoJSON["IPv6Cidr"] = IPv6CidrArr
	locationInfoJSON["IPv6Range"] = IPv6RangeArr
	locationInfoJSON["IsTrusted"] = ipLocationInfo.GetIsTrusted()
	return locationInfoJSON
}

func CountryGetLocationInfo(countryLocationInfo *ADCountryNamedLocationInfo) map[string]interface{} {
	locationInfoJSON := map[string]interface{}{}
	locationInfoJSON["Countries_and_Regions"] = countryLocationInfo.GetCountriesAndRegions()
	locationInfoJSON["Get_Unknown_Countries_and_Regions"] = countryLocationInfo.GetIncludeUnknownCountriesAndRegions()
	locationInfoJSON["Lookup_Method"] = countryLocationInfo.GetCountryLookupMethod().String()
	return locationInfoJSON
}

func (locationInfo *ADNamedLocationInfo) GetLocationInfo() map[string]interface{} {
	switch t := locationInfo.NamedLocation.(type) {
	case ADIpNamedLocationInfo:
		return IpGetLocationInfo(&ADIpNamedLocationInfo{t})
	case ADCountryNamedLocationInfo:
		return CountryGetLocationInfo(&ADCountryNamedLocationInfo{t})
	}
	return nil
}

func (locationInfo *ADNamedLocationInfo) GetType() string {
	switch locationInfo.NamedLocation.(type) {
	case ADIpNamedLocationInfo:
		return "IP"
	case ADCountryNamedLocationInfo:
		return "Country"
	}
	return "Unknown"
}
