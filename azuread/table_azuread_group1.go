package azuread

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/groups/item"
	"github.com/microsoftgraph/msgraph-sdk-go/groups/item/members"
	"github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdGroupTest() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_group_test",
		Description: "Represents an Azure AD user account.",
		Get: &plugin.GetConfig{
			Hydrate:           getAdGroupTest,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdGroupsTest,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "display_name", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "mail", Require: plugin.Optional},
				{Name: "mail_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "on_premises_sync_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "security_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the group.", Transform: transform.FromGo()},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "An optional description for the group."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for groups."},

			// Other fields
			{Name: "classification", Type: proto.ColumnType_STRING, Description: "Describes a classification for the group (such as low, medium or high business impact)."},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time at which the group was created."},
			{Name: "expiration_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the group is set to expire."},
			{Name: "is_assignable_to_role", Type: proto.ColumnType_BOOL, Description: "Indicates whether this group can be assigned to an Azure Active Directory role or not."},

			// Getting below error while requesting value for isSubscribedByMail
			// 	{
			// 		"error": {
			// 				"code": "ErrorInvalidGroup",
			// 				"message": "The requested group '[id@tenantId]' is invalid.",
			// 				"innerError": {
			// 						"date": "2022-07-13T11:06:23",
			// 						"request-id": "63a83d86-a007-4c68-be75-21cea61d830e",
			// 						"client-request-id": "d69d6667-e818-a322-c694-1fec40b438a8"
			// 				}
			// 		}
			// }
			// {Name: "is_subscribed_by_mail", Type: proto.ColumnType_BOOL, Description: "Indicates whether the signed-in user is subscribed to receive email conversations. Default value is true."},

			{Name: "mail", Type: proto.ColumnType_STRING, Description: "The SMTP address for the group, for example, \"serviceadmins@contoso.onmicrosoft.com\"."},
			{Name: "mail_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies whether the group is mail-enabled."},
			{Name: "mail_nickname", Type: proto.ColumnType_STRING, Description: "The mail alias for the user."},
			{Name: "membership_rule", Type: proto.ColumnType_STRING, Description: "The mail alias for the group, unique in the organization."},
			{Name: "membership_rule_processing_state", Type: proto.ColumnType_STRING, Description: "Indicates whether the dynamic membership processing is on or paused. Possible values are On or Paused."},
			{Name: "on_premises_domain_name", Type: proto.ColumnType_STRING, Description: "Contains the on-premises Domanin name synchronized from the on-premises directory."},
			{Name: "on_premises_last_sync_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Indicates the last time at which the group was synced with the on-premises directory."},
			{Name: "on_premises_net_bios_name", Type: proto.ColumnType_STRING, Description: "Contains the on-premises NetBiosName synchronized from the on-premises directory."},
			{Name: "on_premises_sam_account_name", Type: proto.ColumnType_STRING, Description: "Contains the on-premises SAM account name synchronized from the on-premises directory."},
			{Name: "on_premises_security_identifier", Type: proto.ColumnType_STRING, Description: "Contains the on-premises security identifier (SID) for the group that was synchronized from on-premises to the cloud."},
			{Name: "on_premises_sync_enabled", Type: proto.ColumnType_BOOL, Description: "True if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default)."},
			{Name: "renewed_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the group was last renewed. This cannot be modified directly and is only updated via the renew service action."},
			{Name: "security_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies whether the group is a security group."},
			{Name: "security_identifier", Type: proto.ColumnType_STRING, Description: "Security identifier of the group, used in Windows scenarios."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or Hiddenmembership."},

			// Json fields
			{Name: "assigned_labels", Type: proto.ColumnType_JSON, Description: "The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group."},
			{Name: "group_types", Type: proto.ColumnType_JSON, Description: "Specifies the group type and its membership. If the collection contains Unified, the group is a Microsoft 365 group; otherwise, it's either a security group or distribution group. For details, see [groups overview](https://docs.microsoft.com/en-us/graph/api/resources/groups-overview?view=graph-rest-1.0)."},
			{Name: "member_ids", Type: proto.ColumnType_JSON, Hydrate: getAdGroupMembers, Transform: transform.FromValue(), Description: "Id of Users and groups that are members of this group."},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getAdGroupOwners, Transform: transform.FromValue(), Description: "Id od the owners of the group. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "proxy_addresses", Type: proto.ColumnType_JSON, Description: "Email addresses for the group that direct to the same group mailbox. For example: [\"SMTP: bob@contoso.com\", \"smtp: bob@sales.contoso.com\"]. The any operator is required to filter expressions on multi-valued properties."},
			// {Name: "resource_behavior_options", Type: proto.ColumnType_JSON, Description: "Specifies the group behaviors that can be set for a Microsoft 365 group during creation. Possible values are AllowOnlyMembersToPost, HideGroupInOutlook, SubscribeNewGroupMembers, WelcomeEmailDisabled."},
			// {Name: "resource_provisioning_options", Type: proto.ColumnType_JSON, Description: "Specifies the group resources that are provisioned as part of Microsoft 365 group creation, that are not normally part of default group creation. Possible value is Team."},

			// Standard columns
			{Name: "tags", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTags, Transform: transform.From(adGroupTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdGroupsTest(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		panic(fmt.Errorf("error creating credentials: %w", err))
	}

	// List operations
	input := &groups.GroupsRequestBuilderGetQueryParameters{}

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
	selectColumns := buildGroupRequestFields(ctx, givenColumns)

	input.Select = selectColumns

	var queryFilter string
	filter := buildGroupQueryFilter(equalQuals)
	filter = append(filter, buildGroupBoolNEFilter(quals)...)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = &queryFilter
	} else if len(filter) > 0 {
		joinStr := strings.Join(filter, " and ")
		input.Filter = &joinStr
	}

	options := &groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Groups().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errors.New(fmt.Sprintf("failed to list groups. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateGroupCollectionResponseFromDiscriminatorValue)

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		group := pageItem.(models.Groupable)

		result := map[string]interface{}{
			"DisplayName":        group.GetDisplayName(),
			"ID":                 group.GetId(),
			"Description":        group.GetDescription(),
			"Classification":     group.GetClassification(),
			"CreatedDateTime":    group.GetCreatedDateTime(),
			"ExpirationDateTime": group.GetExpirationDateTime(),
			"IsAssignableToRole": group.GetIsAssignableToRole(),
			//"IsSubscribedByMail":            group.GetIsSubscribedByMail(),
			"Mail":                          group.GetMail(),
			"MailEnabled":                   group.GetMailEnabled(),
			"MailNickname":                  group.GetMailNickname(),
			"MembershipRule":                group.GetMembershipRule(),
			"MembershipRuleProcessingState": group.GetMembershipRuleProcessingState(),
			"OnPremisesDomainName":          group.GetOnPremisesDomainName(),
			"OnPremisesLastSyncDateTime":    group.GetOnPremisesLastSyncDateTime(),
			"OnPremisesNetBiosName":         group.GetOnPremisesNetBiosName(),
			"OnPremisesSamAccountName":      group.GetOnPremisesSamAccountName(),
			"OnPremisesSecurityIdentifier":  group.GetOnPremisesSecurityIdentifier(),
			"OnPremisesSyncEnabled":         group.GetOnPremisesSyncEnabled(),
			"RenewedDateTime":               group.GetRenewedDateTime(),
			"SecurityEnabled":               group.GetSecurityEnabled(),
			"SecurityIdentifier":            group.GetSecurityIdentifier(),
			"Visibility":                    group.GetVisibility(),
			"AssignedLabels":                group.GetAssignedLabels(),
			"GroupTypes":                    group.GetGroupTypes(),
			"ProxyAddresses":                group.GetProxyAddresses(),
		}

		d.StreamListItem(ctx, result)

		return true
	})

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdGroupTest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	groupId := d.KeyColumnQuals["id"].GetStringValue()
	if groupId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		panic(fmt.Errorf("error creating credentials: %w", err))
	}

	// Check for query context and requests only for queried columns
	givenColumns := d.QueryContext.Columns
	selectColumns := buildGroupRequestFields(ctx, givenColumns)

	input := &item.GroupItemRequestBuilderGetQueryParameters{}
	input.Select = selectColumns

	options := &item.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	group, err := client.GroupsById(groupId).GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		if isResourceNotFound(errObj) {
			return nil, nil
		}

		return nil, errors.New(fmt.Sprintf("failed to get group. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	result := map[string]interface{}{
		"DisplayName":        group.GetDisplayName(),
		"ID":                 group.GetId(),
		"Description":        group.GetDescription(),
		"Classification":     group.GetClassification(),
		"CreatedDateTime":    group.GetCreatedDateTime(),
		"ExpirationDateTime": group.GetExpirationDateTime(),
		"IsAssignableToRole": group.GetIsAssignableToRole(),
		//"IsSubscribedByMail":            group.GetIsSubscribedByMail(),
		"Mail":                          group.GetMail(),
		"MailEnabled":                   group.GetMailEnabled(),
		"MailNickname":                  group.GetMailNickname(),
		"MembershipRule":                group.GetMembershipRule(),
		"MembershipRuleProcessingState": group.GetMembershipRuleProcessingState(),
		"OnPremisesDomainName":          group.GetOnPremisesDomainName(),
		"OnPremisesLastSyncDateTime":    group.GetOnPremisesLastSyncDateTime(),
		"OnPremisesNetBiosName":         group.GetOnPremisesNetBiosName(),
		"OnPremisesSamAccountName":      group.GetOnPremisesSamAccountName(),
		"OnPremisesSecurityIdentifier":  group.GetOnPremisesSecurityIdentifier(),
		"OnPremisesSyncEnabled":         group.GetOnPremisesSyncEnabled(),
		"RenewedDateTime":               group.GetRenewedDateTime(),
		"SecurityEnabled":               group.GetSecurityEnabled(),
		"SecurityIdentifier":            group.GetSecurityIdentifier(),
		"Visibility":                    group.GetVisibility(),
		"AssignedLabels":                group.GetAssignedLabels(),
		"GroupTypes":                    group.GetGroupTypes(),
		"ProxyAddresses":                group.GetProxyAddresses(),
	}

	return result, nil
}

func getAdGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		panic(fmt.Errorf("error creating credentials: %w", err))
	}

	group := h.Item.(map[string]interface{})
	groupID := group["ID"].(*string)

	headers := map[string]string{
		"ConsistencyLevel": "eventual",
	}

	includeCount := true
	requestParameters := &members.MembersRequestBuilderGetQueryParameters{
		Count: &includeCount,
	}

	config := &members.MembersRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParameters,
	}

	memberIds := []*string{}
	members, err := client.GroupsById(*groupID).Members().GetWithRequestConfigurationAndResponseHandler(config, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errors.New(fmt.Sprintf("failed to list group members. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	pageIterator, err := msgraphcore.NewPageIterator(members, adapter, models.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue)
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		member := pageItem.(models.DirectoryObjectable)
		memberIds = append(memberIds, member.GetId())

		return true
	})

	return memberIds, nil
}

func getAdGroupOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		panic(fmt.Errorf("error creating credentials: %w", err))
	}

	group := h.Item.(map[string]interface{})
	groupID := group["ID"].(*string)

	headers := map[string]string{
		"ConsistencyLevel": "eventual",
	}

	includeCount := true
	requestParameters := &owners.OwnersRequestBuilderGetQueryParameters{
		Count: &includeCount,
	}

	config := &owners.OwnersRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParameters,
	}

	ownerIds := []*string{}
	owners, err := client.GroupsById(*groupID).Owners().GetWithRequestConfigurationAndResponseHandler(config, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errors.New(fmt.Sprintf("failed to list group owners. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	pageIterator, err := msgraphcore.NewPageIterator(owners, adapter, models.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue)
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		member := pageItem.(models.DirectoryObjectable)
		ownerIds = append(ownerIds, member.GetId())

		return true
	})

	return ownerIds, nil
}

//// TRANSFORM FUNCTIONS

func adGroupTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(map[string]interface{})
	assignedLabels := group["AssignedLabels"].([]models.AssignedLabelable)

	if assignedLabels == nil {
		return nil, nil
	}
	var tags = map[string]string{}
	for _, i := range assignedLabels {
		tags[*i.GetLabelId()] = *i.GetDisplayName()
	}
	return tags, nil
}

func buildGroupRequestFields(ctx context.Context, queryColumns []string) []string {
	var selectColumns []string

	if !helpers.StringSliceContains(queryColumns, "id") {
		queryColumns = append(queryColumns, "id")
	}

	if !helpers.StringSliceContains(queryColumns, "assignedLabels") && helpers.StringSliceContains(queryColumns, "tags") {
		queryColumns = append(queryColumns, "assignedLabels")
	}

	for _, columnName := range queryColumns {
		if columnName == "title" || columnName == "tags" || columnName == "filter" || columnName == "tenant_id" {
			continue
		}

		// Uses separate hydrate functions
		if columnName == "member_ids" || columnName == "owner_ids" {
			continue
		}

		selectColumns = append(selectColumns, strcase.ToLowerCamel(columnName))
	}

	return selectColumns
}

func buildGroupQueryFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := map[string]string{
		"display_name":             "string",
		"mail":                     "string",
		"mail_enabled":             "bool",
		"on_premises_sync_enabled": "bool",
		"security_enabled":         "bool",
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

func buildGroupBoolNEFilter(quals plugin.KeyColumnQualMap) []string {
	filters := []string{}

	filterQuals := []string{
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
