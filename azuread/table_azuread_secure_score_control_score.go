package azuread

import (
	"context"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/security"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdSecureScoreControlScore(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_secure_score_control_score",
		Description: "Represents the per-control scores of the most recent Microsoft Secure Score for a tenant.",
		List: &plugin.ListConfig{
			Hydrate: listAdSecureScoreControlScores,
		},

		Columns: commonColumns([]*plugin.Column{
			{Name: "control_name", Type: proto.ColumnType_STRING, Description: "The name of the control.", Transform: transform.FromMethod("GetControlName")},
			{Name: "control_category", Type: proto.ColumnType_STRING, Description: "The category that the control belongs to, for example Identity, Data, Device, Apps, or Infrastructure.", Transform: transform.FromMethod("GetControlCategory")},
			{Name: "score", Type: proto.ColumnType_DOUBLE, Description: "The tenant's current achieved score for this control.", Transform: transform.FromMethod("GetScore")},
			{Name: "max_score", Type: proto.ColumnType_DOUBLE, Description: "The maximum achievable score for the control, taken from the matching secureScoreControlProfile. Null when no profile exposes it.", Transform: transform.FromField("MaxScore")},
			{Name: "implementation_status", Type: proto.ColumnType_STRING, Description: "The implementation status of the control, as reported in the control score. Null when Graph does not provide it.", Transform: transform.FromMethod("ImplementationStatus")},
			{Name: "last_synced_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the control score was last synced, as reported in the control score. Null when Graph does not provide it.", Transform: transform.FromMethod("LastSyncedDateTime")},
		}),
	}
}

//// LIST FUNCTION

func listAdSecureScoreControlScores(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_secure_score_control_score.listAdSecureScoreControlScores", "connection_error", err)
		return nil, err
	}

	options := &security.SecureScoresRequestBuilderGetRequestConfiguration{
		QueryParameters: &security.SecureScoresRequestBuilderGetQueryParameters{
			Top: Int32(1),
		},
	}

	result, err := client.Security().SecureScores().Get(ctx, options)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("azuread_secure_score_control_score.listAdSecureScoreControlScores", "list_secure_score_error", errObj)
		return nil, errObj
	}

	scores := result.GetValue()
	if len(scores) == 0 {
		return nil, nil
	}
	latest := scores[0]

	maxScores := getSecureScoreControlMaxScores(ctx, d)

	for _, controlScore := range latest.GetControlScores() {
		if controlScore == nil {
			continue
		}
		info := &ADSecureScoreControlScoreInfo{ControlScoreable: controlScore}
		if controlScore.GetControlName() != nil {
			if maxScore, ok := maxScores[*controlScore.GetControlName()]; ok {
				info.MaxScore = maxScore
			}
		}
		d.StreamListItem(ctx, info)

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// getSecureScoreControlMaxScores maps each control name to the maxScore from its
// secureScoreControlProfile, the only place Graph exposes a control's maximum. A read failure
// yields an empty map so the control scores still stream with a null max_score.
func getSecureScoreControlMaxScores(ctx context.Context, d *plugin.QueryData) map[string]*float64 {
	maxScores := map[string]*float64{}

	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_secure_score_control_score.getSecureScoreControlMaxScores", "connection_error", err)
		return maxScores
	}

	result, err := client.Security().SecureScoreControlProfiles().Get(ctx, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_secure_score_control_score.getSecureScoreControlMaxScores", "list_profile_error", getErrorObject(err))
		return maxScores
	}

	pageIterator, err := msgraphcore.NewPageIterator[interface{}](result, adapter, models.CreateSecureScoreControlProfileCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_secure_score_control_score.getSecureScoreControlMaxScores", "create_iterator_instance_error", err)
		return maxScores
	}

	err = pageIterator.Iterate(ctx, func(pageItem interface{}) bool {
		if profile, ok := pageItem.(models.SecureScoreControlProfileable); ok {
			if profile.GetId() != nil && profile.GetMaxScore() != nil {
				maxScores[*profile.GetId()] = profile.GetMaxScore()
			}
		}
		return true
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuread_secure_score_control_score.getSecureScoreControlMaxScores", "paging_error", err)
	}

	return maxScores
}
