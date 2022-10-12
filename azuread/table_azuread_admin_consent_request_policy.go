package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdAdminConsentRequestPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_admin_consent_request_policy",
		Description: "Represents the policy for enabling or disabling the Azure AD admin consent workflow.",
		List: &plugin.ListConfig{
			Hydrate: listAdAdminConsentRequestPolicies,
		},

		Columns: []*plugin.Column{
			{Name: "is_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies whether the admin consent request feature is enabled or disabled.", Transform: transform.FromMethod("GetIsEnabled")},

			// Other fields
			{Name: "notify_reviewers", Type: proto.ColumnType_BOOL, Description: "Specifies whether reviewers will receive notifications.", Transform: transform.FromMethod("GetNotifyReviewers")},
			{Name: "reminders_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies whether reviewers will receive reminder emails.", Transform: transform.FromMethod("GetRemindersEnabled")},
			{Name: "request_duration_in_days", Type: proto.ColumnType_INT, Description: "Specifies the duration the request is active before it automatically expires if no decision is applied.", Transform: transform.FromMethod("GetRequestDurationInDays")},
			{Name: "version", Type: proto.ColumnType_INT, Description: "Specifies the version of this policy. When the policy is updated, this version is updated.", Transform: transform.FromMethod("GetVersion")},

			// JSON fields
			{Name: "reviewers", Type: proto.ColumnType_JSON, Description: "The list of reviewers for the admin consent.", Transform: transform.FromMethod("AdminConsentRequestPolicyReviewers")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Default: "adminConsentRequestPolicy"},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdAdminConsentRequestPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_admin_consent_request_policy.listAdAdminConsentRequestPolicies", "connection_error", err)
		return nil, err
	}

	result, err := client.Policies().AdminConsentRequestPolicy().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdAdminConsentRequestPolicies", "list_application_error", errObj)
		return nil, errObj
	}
	d.StreamListItem(ctx, &ADAdminConsentRequestPolicyInfo{result})

	return nil, nil
}
