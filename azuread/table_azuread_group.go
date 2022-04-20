package azuread

import (
	"context"
	"strings"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureAdGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_group",
		Description: "Represents an Azure Active Directory (Azure AD) group, which can be a Microsoft 365 group, or a security group.",
		Get: &plugin.GetConfig{
			Hydrate:           getAdGroup,
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdGroups,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{Name: "id", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},

				// Other fields for filtering OData
				{Name: "mail", Require: plugin.Optional, Operators: []string{"<>", "="}},                     // $filter (eq, ne, NOT, ge, le, in, startsWith).
				{Name: "mail_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},             // $filter (eq, ne, NOT).
				{Name: "on_premises_sync_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}}, // $filter (eq, ne, NOT, in).
				{Name: "security_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},         // $filter (eq, ne, NOT, in).

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
			{Name: "is_subscribed_by_mail", Type: proto.ColumnType_BOOL, Description: "Indicates whether the signed-in user is subscribed to receive email conversations. Default value is true."},
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
			{Name: "member_ids", Type: proto.ColumnType_JSON, Hydrate: getGroupMembers, Transform: transform.FromValue(), Description: "Id of Users and groups that are members of this group."},
			{Name: "owner_ids", Type: proto.ColumnType_JSON, Hydrate: getGroupOwners, Transform: transform.FromValue(), Description: "Id od the owners of the group. The owners are a set of non-admin users who are allowed to modify this object."},
			{Name: "proxy_addresses", Type: proto.ColumnType_JSON, Description: "Email addresses for the group that direct to the same group mailbox. For example: [\"SMTP: bob@contoso.com\", \"smtp: bob@sales.contoso.com\"]. The any operator is required to filter expressions on multi-valued properties."},
			{Name: "resource_behavior_options", Type: proto.ColumnType_JSON, Description: "Specifies the group behaviors that can be set for a Microsoft 365 group during creation. Possible values are AllowOnlyMembersToPost, HideGroupInOutlook, SubscribeNewGroupMembers, WelcomeEmailDisabled."},
			{Name: "resource_provisioning_options", Type: proto.ColumnType_JSON, Description: "Specifies the group resources that are provisioned as part of Microsoft 365 group creation, that are not normally part of default group creation. Possible value is Team."},

			// Standard columns
			{Name: "tags", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTags, Transform: transform.From(groupTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewGroupsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	input := odata.Query{}

	equalQuals := d.KeyColumnQuals
	quals := d.Quals

	var queryFilter string
	filter := buildQueryFilter(equalQuals)
	filter = append(filter, buildBoolNEFilter(quals)...)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.Filter = queryFilter
	} else if len(filter) > 0 {
		input.Filter = strings.Join(filter, " and ")
	}

	groups, _, err := client.List(ctx, input)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, group := range *groups {
		d.StreamListItem(ctx, group)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, err
}

//// Hydrate Functions

func getAdGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var groupId string
	if h.Item != nil {
		groupId = *h.Item.(msgraph.Group).ID
	} else {
		groupId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if groupId == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewGroupsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer
	client.BaseClient.DisableRetries = true

	group, _, err := client.Get(ctx, groupId, odata.Query{})
	if err != nil {
		return nil, err
	}
	return *group, nil
}

func getGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(msgraph.Group)
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewGroupsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	members, _, err := client.ListMembers(ctx, *group.ID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func getGroupOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(msgraph.Group)
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewGroupsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	owners, _, err := client.ListOwners(ctx, *group.ID)
	if err != nil {
		return nil, err
	}
	return owners, nil
}

//// Transform Function

func groupTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(msgraph.Group)

	if group.AssignedLabels == nil {
		return nil, nil
	}
	var tags = map[string]string{}
	for _, i := range *group.AssignedLabels {
		tags[*i.LabelId] = *i.DisplayName
	}
	return tags, nil
}
