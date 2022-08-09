package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/microsoftgraph/msgraph-sdk-go/users/item"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdUser() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user",
		Description: "Represents an Azure AD user account.",
		Get: &plugin.GetConfig{
			Hydrate: getAdUser,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdUsers,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "user_principal_name", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},

				// Other fields for filtering OData
				{Name: "user_type", Require: plugin.Optional},
				{Name: "account_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "surname", Require: plugin.Optional},
			},
		},

		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the user. Should be treated as an opaque identifier.", Transform: transform.FromMethod("GetId")},
			{Name: "user_principal_name", Type: proto.ColumnType_STRING, Description: "Principal email of the active directory user.", Transform: transform.FromMethod("GetUserPrincipalName")},
			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "True if the account is enabled; otherwise, false.", Transform: transform.FromMethod("GetAccountEnabled")},
			{Name: "user_type", Type: proto.ColumnType_STRING, Description: "A string value that can be used to classify user types in your directory.", Transform: transform.FromMethod("GetUserType")},
			{Name: "given_name", Type: proto.ColumnType_STRING, Description: "The given name (first name) of the user.", Transform: transform.FromMethod("GetGivenName")},
			{Name: "surname", Type: proto.ColumnType_STRING, Description: "Family name or last name of the active directory user.", Transform: transform.FromMethod("GetSurname")},

			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for resources."},

			// Other fields
			{Name: "on_premises_immutable_id", Type: proto.ColumnType_STRING, Description: "Used to associate an on-premises Active Directory user account with their Azure AD user object.", Transform: transform.FromMethod("GetOnPremisesImmutableId")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time at which the user was created.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "mail", Type: proto.ColumnType_STRING, Description: "The SMTP address for the user, for example, jeff@contoso.onmicrosoft.com.", Transform: transform.FromMethod("GetMail")},
			{Name: "mail_nickname", Type: proto.ColumnType_STRING, Description: "The mail alias for the user.", Transform: transform.FromMethod("GetMailNickname")},
			{Name: "password_policies", Type: proto.ColumnType_STRING, Description: "Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two may be specified together; for example: DisablePasswordExpiration, DisableStrongPassword.", Transform: transform.FromMethod("GetPasswordPolicies")},
			{Name: "sign_in_sessions_valid_from_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).", Transform: transform.FromMethod("GetSignInSessionsValidFromDateTime")},
			{Name: "usage_location", Type: proto.ColumnType_STRING, Description: "A two letter country code (ISO standard 3166), required for users that will be assigned licenses due to legal requirement to check for availability of services in countries.", Transform: transform.FromMethod("GetUsageLocation")},

			// Json fields
			{Name: "member_of", Type: proto.ColumnType_JSON, Description: "A list the groups and directory roles that the user is a direct member of.", Transform: transform.FromMethod("UserMemberOf")},
			{Name: "im_addresses", Type: proto.ColumnType_JSON, Description: "The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user.", Transform: transform.FromMethod("GetImAddresses")},
			{Name: "other_mails", Type: proto.ColumnType_JSON, Description: "A list of additional email addresses for the user.", Transform: transform.FromMethod("GetOtherMails")},
			{Name: "password_profile", Type: proto.ColumnType_JSON, Description: "Specifies the password profile for the user. The profile contains the user’s password. This property is required when a user is created.", Transform: transform.FromMethod("UserPasswordProfile")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adUserTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_user.listAdUsers", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &users.UsersRequestBuilderGetQueryParameters{}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit <= 999 {
			l := int32(*limit)
			input.Top = &l
		}
	}

	equalQuals := d.KeyColumnQuals
	quals := d.Quals

	// Check for query context and requests only for queried columns
	givenColumns := d.QueryContext.Columns
	selectColumns, expandColumns := buildUserRequestFields(ctx, givenColumns)

	input.Select = selectColumns
	input.Expand = expandColumns

	var queryFilter string
	filter := buildQueryFilter(equalQuals)
	filter = append(filter, buildBoolNEFilter(quals)...)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = &queryFilter
	} else if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Users().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdUsers", "list_user_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateUserCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdUsers", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		user := pageItem.(models.Userable)

		refreshTokensValidFromDateTime := user.GetAdditionalData()["refreshTokensValidFromDateTime"]

		d.StreamListItem(ctx, &ADUserInfo{user, refreshTokensValidFromDateTime})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.QueryStatus.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdUsers", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_user.getAdUser", "connection_error", err)
		return nil, err
	}

	userId := d.KeyColumnQuals["id"].GetStringValue()
	if userId == "" {
		return nil, nil
	}

	// Check for query context and requests only for queried columns
	givenColumns := d.QueryContext.Columns
	selectColumns, expandColumns := buildUserRequestFields(ctx, givenColumns)

	input := &item.UserItemRequestBuilderGetQueryParameters{}
	input.Select = selectColumns
	input.Expand = expandColumns

	options := &item.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	user, err := client.UsersById(userId).GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdUser", "get_user_error", errObj)
		return nil, errObj
	}
	refreshTokensValidFromDateTime := user.GetAdditionalData()["refreshTokensValidFromDateTime"]

	return &ADUserInfo{user, refreshTokensValidFromDateTime}, nil
}

func buildUserRequestFields(ctx context.Context, queryColumns []string) ([]string, []string) {
	var selectColumns, expandColumns []string

	for _, columnName := range queryColumns {
		if columnName == "title" || columnName == "filter" || columnName == "tenant_id" {
			continue
		}

		if columnName == "member_of" {
			expandColumns = append(expandColumns, fmt.Sprintf("%s($select=id,displayName)", strcase.ToLowerCamel(columnName)))
			continue
		}

		selectColumns = append(selectColumns, strcase.ToLowerCamel(columnName))
	}

	return selectColumns, expandColumns
}

//// TRANSFORM FUNCTIONS

func adUserTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADUserInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetUserPrincipalName()
	}

	return title, nil
}

func buildQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name":             "string",
		"id":                       "string",
		"surname":                  "string",
		"user_principal_name":      "string",
		"user_type":                "string",
		"account_enabled":          "bool",
		"mail_enabled":             "bool",
		"security_enabled":         "bool",
		"on_premises_sync_enabled": "bool",
	}

	for qual, qualType := range filterQuals {
		switch qualType {
		case "string":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq '%s'", strcase.ToCamel(qual), equalQuals[qual].GetStringValue()))
			}
		case "bool":
			if equalQuals[qual] != nil {
				filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), equalQuals[qual].GetBoolValue()))
			}
		}
	}

	return filters
}

func buildBoolNEFilter(quals plugin.KeyColumnQualMap) []string {
	filters := []string{}

	filterQuals := []string{
		"account_enabled",
		"mail_enabled",
		"on_premises_sync_enabled",
		"security_enabled",
	}

	for _, qual := range filterQuals {
		if quals[qual] != nil {
			for _, q := range quals[qual].Quals {
				value := q.Value.GetBoolValue()
				if q.Operator == "<>" {
					filters = append(filters, fmt.Sprintf("%s eq %t", strcase.ToCamel(qual), !value))
					break
				}
			}
		}
	}

	return filters
}
