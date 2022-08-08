package azuread

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"

	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/groups/item"
	"github.com/microsoftgraph/msgraph-sdk-go/groups/item/members"
	"github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureAdGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_group",
		Description: "Represents an Azure AD user account.",
		Get: &plugin.GetConfig{
			Hydrate: getAdGroup,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdGroups,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery", "Invalid filter clause"}),
			},
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
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the group.", Transform: transform.FromMethod("GetId")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "An optional description for the group.", Transform: transform.FromMethod("GetDescription")},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Odata query to search for groups."},

			// Other fields
			{Name: "classification", Type: proto.ColumnType_STRING, Description: "Describes a classification for the group (such as low, medium or high business impact).", Transform: transform.FromMethod("GetClassification")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time at which the group was created.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "expiration_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the group is set to expire.", Transform: transform.FromMethod("GetExpirationDateTime")},
			{Name: "is_assignable_to_role", Type: proto.ColumnType_BOOL, Description: "Indicates whether this group can be assigned to an Azure Active Directory role or not.", Transform: transform.FromMethod("GetIsAssignableToRole")},
			{Name: "is_subscribed_by_mail", Type: proto.ColumnType_BOOL, Description: "Indicates whether the signed-in user is subscribed to receive email conversations. Default value is true.", Hydrate: getAdGroupIsSubscribedByMail, Transform: transform.FromValue()},
			{Name: "mail", Type: proto.ColumnType_STRING, Description: "The SMTP address for the group, for example, \"serviceadmins@contoso.onmicrosoft.com\".", Transform: transform.FromMethod("GetMail")},
			{Name: "mail_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies whether the group is mail-enabled.", Transform: transform.FromMethod("GetMailEnabled")},
			{Name: "mail_nickname", Type: proto.ColumnType_STRING, Description: "The mail alias for the user.", Transform: transform.FromMethod("GetMailNickname")},
			{Name: "membership_rule", Type: proto.ColumnType_STRING, Description: "The mail alias for the group, unique in the organization.", Transform: transform.FromMethod("GetMembershipRule")},
			{Name: "membership_rule_processing_state", Type: proto.ColumnType_STRING, Description: "Indicates whether the dynamic membership processing is on or paused. Possible values are On or Paused.", Transform: transform.FromMethod("GetMembershipRuleProcessingState")},
			{Name: "on_premises_domain_name", Type: proto.ColumnType_STRING, Description: "Contains the on-premises Domanin name synchronized from the on-premises directory.", Transform: transform.FromMethod("GetOnPremisesDomainName")},
			{Name: "on_premises_last_sync_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Indicates the last time at which the group was synced with the on-premises directory.", Transform: transform.FromMethod("GetOnPremisesLastSyncDateTime")},
			{Name: "on_premises_net_bios_name", Type: proto.ColumnType_STRING, Description: "Contains the on-premises NetBiosName synchronized from the on-premises directory.", Transform: transform.FromMethod("GetOnPremisesNetBiosName")},
			{Name: "on_premises_sam_account_name", Type: proto.ColumnType_STRING, Description: "Contains the on-premises SAM account name synchronized from the on-premises directory.", Transform: transform.FromMethod("GetOnPremisesSamAccountName")},
			{Name: "on_premises_security_identifier", Type: proto.ColumnType_STRING, Description: "Contains the on-premises security identifier (SID) for the group that was synchronized from on-premises to the cloud.", Transform: transform.FromMethod("GetOnPremisesSecurityIdentifier")},
			{Name: "on_premises_sync_enabled", Type: proto.ColumnType_BOOL, Description: "True if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default).", Transform: transform.FromMethod("GetOnPremisesSyncEnabled")},
			{Name: "renewed_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the group was last renewed. This cannot be modified directly and is only updated via the renew service action.", Transform: transform.FromMethod("GetRenewedDateTime")},
			{Name: "security_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies whether the group is a security group.", Transform: transform.FromMethod("GetSecurityEnabled")},
			{Name: "security_identifier", Type: proto.ColumnType_STRING, Description: "Security identifier of the group, used in Windows scenarios.", Transform: transform.FromMethod("GetSecurityIdentifier")},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or Hiddenmembership.", Transform: transform.FromMethod("GetVisibility")},

			// JSON fields
			{Name: "assigned_labels", Type: proto.ColumnType_JSON, Description: "The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group.", Transform: transform.FromMethod("GroupAssignedLabels")},
			{Name: "group_types", Type: proto.ColumnType_JSON, Description: "Specifies the group type and its membership. If the collection contains Unified, the group is a Microsoft 365 group; otherwise, it's either a security group or distribution group. For details, see [groups overview](https://docs.microsoft.com/en-us/graph/api/resources/groups-overview?view=graph-rest-1.0).", Transform: transform.FromMethod("GetGroupTypes")},
			{Name: "member_ids", Type: proto.ColumnType_JSON, Hydrate: getAdGroupMembers, Transform: transform.FromValue(), Description: "Id of Users and groups that are members of this group."},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getAdGroupOwners, Transform: transform.FromValue(), Description: "Id od the owners of the group. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "proxy_addresses", Type: proto.ColumnType_JSON, Description: "Email addresses for the group that direct to the same group mailbox. For example: [\"SMTP: bob@contoso.com\", \"smtp: bob@sales.contoso.com\"]. The any operator is required to filter expressions on multi-valued properties.", Transform: transform.FromMethod("GetProxyAddresses")},
			{Name: "resource_behavior_options", Type: proto.ColumnType_JSON, Description: "Specifies the group behaviors that can be set for a Microsoft 365 group during creation. Possible values are AllowOnlyMembersToPost, HideGroupInOutlook, SubscribeNewGroupMembers, WelcomeEmailDisabled."},
			{Name: "resource_provisioning_options", Type: proto.ColumnType_JSON, Description: "Specifies the group resources that are provisioned as part of Microsoft 365 group creation, that are not normally part of default group creation. Possible value is Team."},

			// Standard columns
			{Name: "tags", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTags, Transform: transform.From(adGroupTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adGroupTitle)},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_group.listAdGroups", "connection_error", err)
		return nil, err
	}

	// List operations
	input := &groups.GroupsRequestBuilderGetQueryParameters{}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	// unexpected status 400 with OData error: Request_UnsupportedQuery: Invalid page size specified: '1000'. Must be between 1 and 999 inclusive.
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit <= 999 {
			l := int32(*limit)
			input.Top = &l
		}
	}

	equalQuals := d.KeyColumnQuals
	quals := d.Quals

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
		plugin.Logger(ctx).Error("listAdGroups", "list_group_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateGroupCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("listAdGroups", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		group := pageItem.(models.Groupable)

		resourceBehaviorOptions := formatResourceBehaviorOptions(ctx, group)
		resourceProvisioningOptions := formatResourceProvisioningOptions(ctx, group)

		d.StreamListItem(ctx, &ADGroupInfo{group, resourceBehaviorOptions, resourceProvisioningOptions})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		return d.QueryStatus.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("listAdGroups", "paging_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	groupId := d.KeyColumnQuals["id"].GetStringValue()
	if groupId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_group.getAdGroup", "connection_error", err)
		return nil, err
	}

	input := &item.GroupItemRequestBuilderGetQueryParameters{}

	options := &item.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	group, err := client.GroupsById(groupId).GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdGroup", "get_group_error", errObj)
		return nil, errObj
	}
	resourceBehaviorOptions := formatResourceBehaviorOptions(ctx, group)
	resourceProvisioningOptions := formatResourceProvisioningOptions(ctx, group)

	return &ADGroupInfo{group, resourceBehaviorOptions, resourceProvisioningOptions}, nil
}

// Returned only on $select. Supported only on the Get group API (GET /groups/{ID}).

func getAdGroupIsSubscribedByMail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var groupId string
	if h.Item != nil {
		groupId = *h.Item.(*ADGroupInfo).GetId()
	} else {
		groupId = d.KeyColumnQuals["id"].GetStringValue()
	}
	if groupId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_group.getAdGroupIsSubscribedByMail", "connection_error", err)
		return nil, err
	}

	input := &item.GroupItemRequestBuilderGetQueryParameters{
		Select: []string{"isSubscribedByMail"},
	}

	options := &item.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	group, err := client.GroupsById(groupId).GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdGroupIsSubscribedByMail", "get_group_error", errObj)
		return nil, nil
	}

	return group.GetIsSubscribedByMail(), nil
}

func getAdGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_group.getAdGroupMembers", "connection_error", err)
		return nil, err
	}

	group := h.Item.(*ADGroupInfo)
	groupID := group.GetId()

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
		plugin.Logger(ctx).Error("getAdGroupMembers", "get_group_members_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(members, adapter, models.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("getAdGroupMembers", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		member := pageItem.(models.DirectoryObjectable)
		memberIds = append(memberIds, member.GetId())

		return true
	})
	if err != nil {
		plugin.Logger(ctx).Error("getAdGroupMembers", "paging_error", err)
		return nil, err
	}

	return memberIds, nil
}

func getAdGroupOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_group.getAdGroupOwners", "connection_error", err)
		return nil, err
	}

	group := h.Item.(*ADGroupInfo)
	groupID := group.GetId()

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
		plugin.Logger(ctx).Error("getAdGroupOwners", "get_group_owners_error", errObj)
		return nil, errObj
	}

	pageIterator, err := msgraphcore.NewPageIterator(owners, adapter, models.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue)
	if err != nil {
		plugin.Logger(ctx).Error("getAdGroupOwners", "create_iterator_instance_error", err)
		return nil, err
	}

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		member := pageItem.(models.DirectoryObjectable)
		ownerIds = append(ownerIds, member.GetId())

		return true
	})
	if err != nil {
		plugin.Logger(ctx).Error("getAdGroupMembers", "paging_error", err)
		return nil, err
	}

	return ownerIds, nil
}

//// TRANSFORM FUNCTIONS

func adGroupTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(*ADGroupInfo)
	if group == nil {
		return nil, nil
	}

	assignedLabels := group.GroupAssignedLabels()
	if len(assignedLabels) == 0 {
		return nil, nil
	}

	var tags = map[*string]*string{}
	for _, i := range assignedLabels {
		tags[i["labelId"]] = i["displayName"]
	}

	return tags, nil
}

func adGroupTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ADGroupInfo)
	if data == nil {
		return nil, nil
	}

	title := data.GetDisplayName()
	if title == nil {
		title = data.GetId()
	}

	return title, nil
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

func formatResourceBehaviorOptions(ctx context.Context, group models.Groupable) []string {
	var resourceBehaviorOptions []string
	data := group.GetAdditionalData()["resourceBehaviorOptions"]
	if data != nil {
		parsedData := group.GetAdditionalData()["resourceBehaviorOptions"].([]*jsonserialization.JsonParseNode)

		for _, r := range parsedData {
			val, err := r.GetStringValue()
			if err != nil {
				plugin.Logger(ctx).Error("failed to parse resourceBehaviorOptions: %v", err)
				val = nil
			}

			if val != nil {
				resourceBehaviorOptions = append(resourceBehaviorOptions, *val)
			}
		}
	}
	return resourceBehaviorOptions
}

func formatResourceProvisioningOptions(ctx context.Context, group models.Groupable) []string {
	var resourceProvisioningOptions []string
	data := group.GetAdditionalData()["resourceProvisioningOptions"]
	if data != nil {
		parsedData := data.([]*jsonserialization.JsonParseNode)

		for _, r := range parsedData {
			val, err := r.GetStringValue()
			if err != nil {
				plugin.Logger(ctx).Error("failed to parse resourceProvisioningOptions: %v", err)
				val = nil
			}

			if val != nil {
				resourceProvisioningOptions = append(resourceProvisioningOptions, *val)
			}
		}
	}
	return resourceProvisioningOptions
}
