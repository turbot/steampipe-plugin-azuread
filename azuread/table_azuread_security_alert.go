package azuread

import (
	"context"
	"fmt"
	"strings"
	"time"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	securitymodels "github.com/microsoftgraph/msgraph-sdk-go/models/security"
	"github.com/microsoftgraph/msgraph-sdk-go/security"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdSecurityAlert(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_security_alert",
		Description: "Represents the security alerts (alerts_v2) surfaced across Microsoft 365 Defender for a tenant.",
		Get: &plugin.GetConfig{
			Hydrate: getAdSecurityAlert,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "ResourceNotFound", "NotFound"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdSecurityAlerts,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "service_source", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "last_update_date_time", Require: plugin.Optional, Operators: []string{">", ">="}},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the alert.", Transform: transform.FromMethod("GetId")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Brief identifying string value describing the alert.", Transform: transform.FromMethod("GetTitle")},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "Indicates the possible impact on assets. Possible values are: unknown, informational, low, medium, high.", Transform: transform.FromMethod("Severity")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the alert. Possible values are: new, inProgress, resolved.", Transform: transform.FromMethod("Status")},
			{Name: "classification", Type: proto.ColumnType_STRING, Description: "Specifies whether the alert represents a true threat. Possible values are: unknown, falsePositive, truePositive, informationalExpectedActivity.", Transform: transform.FromMethod("Classification")},
			{Name: "determination", Type: proto.ColumnType_STRING, Description: "Specifies the result of the investigation, whether the alert represents a true attack and, if so, the nature of it.", Transform: transform.FromMethod("Determination")},
			{Name: "service_source", Type: proto.ColumnType_STRING, Description: "The service or product that created this alert. Possible values include: microsoftDefenderForEndpoint, microsoftDefenderForOffice365, microsoftDefenderForIdentity, microsoftDefenderForCloudApps.", Transform: transform.FromMethod("ServiceSource")},
			{Name: "assigned_to", Type: proto.ColumnType_STRING, Description: "Owner of the alert, or null if no owner is assigned.", Transform: transform.FromMethod("GetAssignedTo")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time when Microsoft 365 Defender created the alert.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "last_update_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the alert was last updated at Microsoft 365 Defender.", Transform: transform.FromMethod("GetLastUpdateDateTime")},
			{Name: "resolved_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the alert was resolved.", Transform: transform.FromMethod("GetResolvedDateTime")},
			{Name: "alert_web_url", Type: proto.ColumnType_STRING, Description: "URL for the Microsoft 365 Defender view of the alert.", Transform: transform.FromMethod("GetAlertWebUrl")},
			{Name: "recommended_actions", Type: proto.ColumnType_STRING, Description: "Recommended response and remediation actions to take in the event this alert was generated.", Transform: transform.FromMethod("GetRecommendedActions")},

			// JSON fields
			{Name: "mitre_techniques", Type: proto.ColumnType_JSON, Description: "The attack techniques, as aligned with the MITRE ATT&CK framework.", Transform: transform.FromMethod("GetMitreTechniques")},
			{Name: "evidence", Type: proto.ColumnType_JSON, Description: "Collection of evidence related to the alert.", Transform: transform.FromMethod("Evidence")},
			{Name: "comments", Type: proto.ColumnType_JSON, Description: "Array of comments created by the Security Operations (SecOps) team during the alert management process.", Transform: transform.FromMethod("Comments")},
		}),
	}
}

//// LIST FUNCTION

func listAdSecurityAlerts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_security_alert.listAdSecurityAlerts", "connection_error", err)
		return nil, err
	}

	input := &security.Alerts_v2RequestBuilderGetQueryParameters{
		Top: Int32(1000),
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit < 1000 {
			input.Top = Int32(int32(*limit))
		}
	}

	filter := buildSecurityAlertQueryFilter(d)
	if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &security.Alerts_v2RequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Security().Alerts_v2().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_security_alert.listAdSecurityAlerts", "list_security_alert_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator[interface{}](result, adapter, securitymodels.CreateAlertCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_security_alert.listAdSecurityAlerts", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		if alert, ok := pageItem.(securitymodels.Alertable); ok {
			d.StreamListItem(ctx, &ADSecurityAlertInfo{alert})
		}
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuread_security_alert.listAdSecurityAlerts", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdSecurityAlert(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	alertID := d.EqualsQuals["id"].GetStringValue()
	if alertID == "" {
		return nil, nil
	}

	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_security_alert.getAdSecurityAlert", "connection_error", err)
		return nil, err
	}

	alert, err := client.Security().Alerts_v2().ByAlertId(alertID).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_security_alert.getAdSecurityAlert", "get_security_alert_error", errObj)
		return nil, errObj
	}

	return &ADSecurityAlertInfo{alert}, nil
}

// buildSecurityAlertQueryFilter renders the pushed-down quals as an OData $filter. The
// service_source/status/severity values are the Graph wire enum strings the columns already expose.
func buildSecurityAlertQueryFilter(d *plugin.QueryData) []string {
	filters := []string{}

	equalQuals := d.EqualsQuals
	stringQuals := map[string]string{
		"service_source": "serviceSource",
		"status":         "status",
		"severity":       "severity",
	}
	for qual, field := range stringQuals {
		if equalQuals[qual] != nil {
			filters = append(filters, fmt.Sprintf("%s eq '%s'", field, equalQuals[qual].GetStringValue()))
		}
	}

	if d.Quals["last_update_date_time"] != nil {
		for _, q := range d.Quals["last_update_date_time"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">":
				filters = append(filters, fmt.Sprintf("lastUpdateDateTime gt %s", givenTime.Format(time.RFC3339)))
			case ">=":
				filters = append(filters, fmt.Sprintf("lastUpdateDateTime ge %s", givenTime.Format(time.RFC3339)))
			}
		}
	}

	return filters
}
