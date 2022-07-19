package azuread

import (
	"context"
	"errors"
	"fmt"
	"os"
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

func tableAzureAdUserTest() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user_test",
		Description: "Represents an Azure AD user account.",
		Get: &plugin.GetConfig{
			Hydrate:    getAdUserTest,
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdUsersTest,
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
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the user. Should be treated as an opaque identifier.", Transform: transform.FromGo()},
			{Name: "user_principal_name", Type: proto.ColumnType_STRING, Description: "Principal email of the active directory user."},
			{Name: "account_enabled", Type: proto.ColumnType_BOOL, Description: "True if the account is enabled; otherwise, false."},
			{Name: "user_type", Type: proto.ColumnType_STRING, Description: "A string value that can be used to classify user types in your directory."},
			{Name: "given_name", Type: proto.ColumnType_STRING, Description: "The given name (first name) of the user."},
			{Name: "surname", Type: proto.ColumnType_STRING, Description: "Family name or last name of the active directory user."},

			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for resources."},

			// // Other fields
			{Name: "on_premises_immutable_id", Type: proto.ColumnType_STRING, Description: "Used to associate an on-premises Active Directory user account with their Azure AD user object."},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time at which the user was created."},
			{Name: "mail", Type: proto.ColumnType_STRING, Description: "The SMTP address for the user, for example, jeff@contoso.onmicrosoft.com."},
			{Name: "mail_nickname", Type: proto.ColumnType_STRING, Description: "The mail alias for the user."},
			{Name: "password_policies", Type: proto.ColumnType_STRING, Description: "Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two may be specified together; for example: DisablePasswordExpiration, DisableStrongPassword."},
			// {Name: "refresh_tokens_valid_from_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph)."},
			{Name: "sign_in_sessions_valid_from_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph)."},
			{Name: "usage_location", Type: proto.ColumnType_STRING, Description: "A two letter country code (ISO standard 3166), required for users that will be assigned licenses due to legal requirement to check for availability of services in countries."},

			// // Json fields
			{Name: "member_of", Type: proto.ColumnType_JSON, Description: "A list the groups and directory roles that the user is a direct member of."},
			// {Name: "additional_properties", Type: proto.ColumnType_JSON, Description: "A list of unmatched properties from the message are deserialized this collection."},
			{Name: "im_addresses", Type: proto.ColumnType_JSON, Description: "The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user."},
			{Name: "other_mails", Type: proto.ColumnType_JSON, Description: "A list of additional email addresses for the user."},
			{Name: "password_profile", Type: proto.ColumnType_JSON, Description: "Specifies the password profile for the user. The profile contains the user’s password. This property is required when a user is created."},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "UserPrincipalName")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdUsersTest(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		panic(fmt.Errorf("error creating credentials: %w", err))
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
		return nil, errors.New(fmt.Sprintf("failed to list users. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateUserCollectionResponseFromDiscriminatorValue)

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		user := pageItem.(models.Userable)

		result := map[string]interface{}{
			"DisplayName":                     user.GetDisplayName(),
			"ID":                              user.GetId(),
			"UserPrincipalName":               user.GetUserPrincipalName(),
			"AccountEnabled":                  user.GetAccountEnabled(),
			"UserType":                        user.GetUserType(),
			"GivenName":                       user.GetGivenName(),
			"Surname":                         user.GetSurname(),
			"OnPremisesImmutableId":           user.GetOnPremisesImmutableId(),
			"CreatedDateTime":                 user.GetCreatedDateTime(),
			"Mail":                            user.GetMail(),
			"MailNickname":                    user.GetMailNickname(),
			"PasswordPolicies":                user.GetPasswordPolicies(),
			"SignInSessionsValidFromDateTime": user.GetSignInSessionsValidFromDateTime(),
			"UsageLocation":                   user.GetUsageLocation(),
			"ImAddresses":                     user.GetImAddresses(),
			"OtherMails":                      user.GetOtherMails(),
			"PasswordProfile":                 user.GetPasswordProfile(),
		}

		memberIds := []string{}
		for _, i := range user.GetMemberOf() {
			memberIds = append(memberIds, *i.GetId())
		}
		result["MemberOf"] = memberIds

		d.StreamListItem(ctx, result)

		return true
	})

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdUserTest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		panic(fmt.Errorf("error creating credentials: %w", err))
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
		if isResourceNotFound(errObj) {
			return nil, nil
		}

		return nil, errors.New(fmt.Sprintf("failed to get user. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	result := map[string]interface{}{
		"DisplayName":                     user.GetDisplayName(),
		"ID":                              user.GetId(),
		"UserPrincipalName":               user.GetUserPrincipalName(),
		"AccountEnabled":                  user.GetAccountEnabled(),
		"UserType":                        user.GetUserType(),
		"GivenName":                       user.GetGivenName(),
		"Surname":                         user.GetSurname(),
		"OnPremisesImmutableId":           user.GetOnPremisesImmutableId(),
		"CreatedDateTime":                 user.GetCreatedDateTime(),
		"Mail":                            user.GetMail(),
		"MailNickname":                    user.GetMailNickname(),
		"PasswordPolicies":                user.GetPasswordPolicies(),
		"SignInSessionsValidFromDateTime": user.GetSignInSessionsValidFromDateTime(),
		"UsageLocation":                   user.GetUsageLocation(),
		"ImAddresses":                     user.GetImAddresses(),
		"OtherMails":                      user.GetOtherMails(),
		"PasswordProfile":                 user.GetPasswordProfile(),
	}

	memberIds := []string{}
	for _, i := range user.GetMemberOf() {
		memberIds = append(memberIds, *i.GetId())
	}
	result["MemberOf"] = memberIds

	return result, nil
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

func getTenant(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var tenantID string
	var err error

	// Read tenantID from config, or environment variables
	azureADConfig := GetConfig(d.Connection)
	if azureADConfig.TenantID != nil {
		tenantID = *azureADConfig.TenantID
	} else {
		tenantID = os.Getenv("AZURE_TENANT_ID")
	}

	// If not set in config, get tenantID from CLI
	if tenantID != "" {
		tenantID, err = getTenantFromCLI()
		if err != nil {
			return nil, err
		}
	}

	return tenantID, nil
}
