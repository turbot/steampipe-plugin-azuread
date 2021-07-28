package azuread

import (
	"context"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAdUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user",
		Description: "Azure AD User",
		// Get: &plugin.GetConfig{
		// 	KeyColumns:        plugin.SingleColumn("object_id"),
		// 	Hydrate:           getAdUser,
		// 	ShouldIgnoreError: isNotFoundError([]string{"Request_ResourceNotFound", "Request_BadRequest"}),
		// },
		List: &plugin.ListConfig{
			Hydrate: listAdUsers,
		},

		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "A friendly name that identifies an active directory user."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique ID that identifies an active directory user.", Transform: transform.FromGo()},
			{Name: "user_principal_name", Type: proto.ColumnType_STRING, Description: "Principal email of the active directory user."},
			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies the account status of the active directory user."},
			{Name: "user_type", Type: proto.ColumnType_STRING, Description: "A string value that can be used to classify user types in your directory."},
			{Name: "given_name", Type: proto.ColumnType_STRING, Description: "The given name(first name) of the active directory user."},
			{Name: "surname", Type: proto.ColumnType_STRING, Description: "Family name or last name of the active directory user."},

			// Other fields

			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: " The time at which the user was created."},
			// {Name: "deleted_date_time", Type: proto.ColumnType_TIMESTAMP, Description: " The time at which the directory object was deleted."},
			{Name: "is_management_restricted", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "is_resource_account", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "mail", Type: proto.ColumnType_STRING, Description: "The SMTP address for the user."},
			{Name: "mail_nickname", Type: proto.ColumnType_STRING, Description: "The mail alias for the user."},
			{Name: "password_policies", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "refresh_tokens_valid_from_date_time", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "sign_in_sessions_valid_from_date_time", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "usage_location", Type: proto.ColumnType_STRING, Description: "A two letter country code (ISO standard 3166), required for users that will be assigned licenses due to legal requirement to check for availability of services in countries."},

			// Json fields
			{Name: "additional_properties", Type: proto.ColumnType_JSON, Description: "A list of unmatched properties from the message are deserialized this collection."},
			{Name: "im_addresses", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "other_mails", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "password_profile", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "sign_in_activity", Type: proto.ColumnType_JSON, Description: ""},

			// {Name: "data", Type: proto.ColumnType_JSON, Description: "The unique ID that identifies an active directory user.", Transform: transform.FromValue()}, // For debugging

			// // Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("display_name", "user_principal_name"),
			},
			// {
			// 	Name:        "akas",
			// 	Description: ColumnDescriptionAkas,
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromP(getAdUserTurbotData, "TurbotAkas"),
			// },
		},
	}
}

//// LIST FUNCTION

func listAdUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}
	tenantID := session.TenantID

	client := msgraph.NewUsersClient(tenantID)
	client.BaseClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		users, _, err := client.List(ctx, odata.Query{})
		if err != nil {
			return nil, err
		}

		for _, user := range *users {
			d.StreamListItem(ctx, user)
		}
		pagesLeft = false
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

// func getAdUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getAdUser")

// 	session, err := GetNewSession(ctx, d, "GRAPH")
// 	if err != nil {
// 		return nil, err
// 	}
// 	tenantID := session.TenantID
// 	objectID := d.KeyColumnQuals["object_id"].GetStringValue()

// 	graphClient := graphrbac.NewUsersClient(tenantID)
// 	graphClient.Authorizer = session.Authorizer

// 	op, err := graphClient.Get(ctx, objectID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return op, nil
// }

// //// TRANSFORM FUNCTIONS

// func getAdUserTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
// 	data := d.HydrateItem.(graphrbac.User)
// 	param := d.Param.(string)

// 	// Get resource title
// 	title := data.ObjectID
// 	if data.DisplayName != nil {
// 		title = data.DisplayName
// 	}

// 	// Get resource tags
// 	akas := []string{"azure:///user/" + *data.ObjectID}

// 	if param == "TurbotTitle" {
// 		return title, nil
// 	}
// 	return akas, nil
// }
