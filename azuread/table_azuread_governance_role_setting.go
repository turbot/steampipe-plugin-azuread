package azuread

import (
	"context"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	betamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func tableAzureAdGovernanceRoleSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_governance_role_setting",
		Description: "Represents governance role settings for Azure resources in Privileged Identity Management (PIM).",
		Get: &plugin.GetConfig{
			Hydrate: getAdGovernanceRoleSetting,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdGovernanceRoleSettings,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the governance role setting.", Transform: transform.FromMethod("GetId")},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the resource.", Transform: transform.FromMethod("GetResourceId")},
			{Name: "role_definition_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the role definition.", Transform: transform.FromMethod("GetRoleDefinitionId")},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Indicates whether this is the default role setting.", Transform: transform.FromMethod("GetIsDefault")},
			{Name: "last_updated_by", Type: proto.ColumnType_STRING, Description: "The user who last updated the role setting.", Transform: transform.FromMethod("GetLastUpdatedBy")},
			{Name: "last_updated_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the role setting was last updated.", Transform: transform.FromMethod("GetLastUpdatedDateTime")},
			{Name: "admin_eligible_settings", Type: proto.ColumnType_JSON, Description: "Settings for admin eligible assignments.", Transform: transform.FromMethod("GetAdminEligibleSettings")},
			{Name: "admin_member_settings", Type: proto.ColumnType_JSON, Description: "Settings for admin member assignments.", Transform: transform.FromMethod("GetAdminMemberSettings")},
			{Name: "user_eligible_settings", Type: proto.ColumnType_JSON, Description: "Settings for user eligible assignments.", Transform: transform.FromMethod("GetUserEligibleSettings")},
			{Name: "user_member_settings", Type: proto.ColumnType_JSON, Description: "Settings for user member assignments.", Transform: transform.FromMethod("GetUserMemberSettings")},
			{Name: "resource_details", Type: proto.ColumnType_JSON, Description: "Details about the resource.", Transform: transform.FromMethod("GetResourceDetails")},
			{Name: "role_definition_details", Type: proto.ColumnType_JSON, Description: "Details about the role definition.", Transform: transform.FromMethod("GetRoleDefinitionDetails")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adGovernanceRoleSettingTitle)},
		}),
	}
}

type ADGovernanceRoleSettingInfo struct {
	betamodels.GovernanceRoleSettingable
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetId() *string {
	return roleSetting.GovernanceRoleSettingable.GetId()
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetResourceId() *string {
	return roleSetting.GovernanceRoleSettingable.GetResourceId()
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetRoleDefinitionId() *string {
	return roleSetting.GovernanceRoleSettingable.GetRoleDefinitionId()
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetIsDefault() *bool {
	return roleSetting.GovernanceRoleSettingable.GetIsDefault()
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetLastUpdatedBy() *string {
	return roleSetting.GovernanceRoleSettingable.GetLastUpdatedBy()
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetLastUpdatedDateTime() *string {
	if roleSetting.GovernanceRoleSettingable.GetLastUpdatedDateTime() != nil {
		dateTime := roleSetting.GovernanceRoleSettingable.GetLastUpdatedDateTime().Format(time.RFC3339)
		return &dateTime
	}
	return nil
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetAdminEligibleSettings() []map[string]interface{} {
	if roleSetting.GovernanceRoleSettingable.GetAdminEligibleSettings() != nil {
		settings := make([]map[string]interface{}, 0, len(roleSetting.GovernanceRoleSettingable.GetAdminEligibleSettings()))
		for _, setting := range roleSetting.GovernanceRoleSettingable.GetAdminEligibleSettings() {
			if setting != nil {
				settingInfo := make(map[string]interface{})
				if setting.GetRuleIdentifier() != nil {
					settingInfo["ruleIdentifier"] = *setting.GetRuleIdentifier()
				}
				if setting.GetSetting() != nil {
					settingInfo["setting"] = *setting.GetSetting()
				}
				if setting.GetOdataType() != nil {
					settingInfo["@odata.type"] = *setting.GetOdataType()
				}
				settings = append(settings, settingInfo)
			}
		}
		return settings
	}
	return nil
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetAdminMemberSettings() []map[string]interface{} {
	if roleSetting.GovernanceRoleSettingable.GetAdminMemberSettings() != nil {
		settings := make([]map[string]interface{}, 0, len(roleSetting.GovernanceRoleSettingable.GetAdminMemberSettings()))
		for _, setting := range roleSetting.GovernanceRoleSettingable.GetAdminMemberSettings() {
			if setting != nil {
				settingInfo := make(map[string]interface{})
				if setting.GetRuleIdentifier() != nil {
					settingInfo["ruleIdentifier"] = *setting.GetRuleIdentifier()
				}
				if setting.GetSetting() != nil {
					settingInfo["setting"] = *setting.GetSetting()
				}
				if setting.GetOdataType() != nil {
					settingInfo["@odata.type"] = *setting.GetOdataType()
				}
				settings = append(settings, settingInfo)
			}
		}
		return settings
	}
	return nil
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetUserEligibleSettings() []map[string]interface{} {
	if roleSetting.GovernanceRoleSettingable.GetUserEligibleSettings() != nil {
		settings := make([]map[string]interface{}, 0, len(roleSetting.GovernanceRoleSettingable.GetUserEligibleSettings()))
		for _, setting := range roleSetting.GovernanceRoleSettingable.GetUserEligibleSettings() {
			if setting != nil {
				settingInfo := make(map[string]interface{})
				if setting.GetRuleIdentifier() != nil {
					settingInfo["ruleIdentifier"] = *setting.GetRuleIdentifier()
				}
				if setting.GetSetting() != nil {
					settingInfo["setting"] = *setting.GetSetting()
				}
				if setting.GetOdataType() != nil {
					settingInfo["@odata.type"] = *setting.GetOdataType()
				}
				settings = append(settings, settingInfo)
			}
		}
		return settings
	}
	return nil
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetUserMemberSettings() []map[string]interface{} {
	if roleSetting.GovernanceRoleSettingable.GetUserMemberSettings() != nil {
		settings := make([]map[string]interface{}, 0, len(roleSetting.GovernanceRoleSettingable.GetUserMemberSettings()))
		for _, setting := range roleSetting.GovernanceRoleSettingable.GetUserMemberSettings() {
			if setting != nil {
				settingInfo := make(map[string]interface{})
				if setting.GetRuleIdentifier() != nil {
					settingInfo["ruleIdentifier"] = *setting.GetRuleIdentifier()
				}
				if setting.GetSetting() != nil {
					settingInfo["setting"] = *setting.GetSetting()
				}
				if setting.GetOdataType() != nil {
					settingInfo["@odata.type"] = *setting.GetOdataType()
				}
				settings = append(settings, settingInfo)
			}
		}
		return settings
	}
	return nil
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetResourceDetails() map[string]interface{} {
	if roleSetting.GovernanceRoleSettingable.GetResource() != nil {
		resource := roleSetting.GovernanceRoleSettingable.GetResource()
		details := make(map[string]interface{})

		if resource.GetId() != nil {
			details["id"] = *resource.GetId()
		}
		if resource.GetDisplayName() != nil {
			details["displayName"] = *resource.GetDisplayName()
		}
		if resource.GetExternalId() != nil {
			details["externalId"] = *resource.GetExternalId()
		}
		if resource.GetStatus() != nil {
			details["status"] = *resource.GetStatus()
		}
		if resource.GetTypeEscaped() != nil {
			details["type"] = *resource.GetTypeEscaped()
		}
		if resource.GetRegisteredRoot() != nil {
			details["registeredRoot"] = *resource.GetRegisteredRoot()
		}
		if resource.GetRegisteredDateTime() != nil {
			details["registeredDateTime"] = resource.GetRegisteredDateTime().Format(time.RFC3339)
		}

		return details
	}
	return nil
}

func (roleSetting *ADGovernanceRoleSettingInfo) GetRoleDefinitionDetails() map[string]interface{} {
	if roleSetting.GovernanceRoleSettingable.GetRoleDefinition() != nil {
		roleDefinition := roleSetting.GovernanceRoleSettingable.GetRoleDefinition()
		details := make(map[string]interface{})

		if roleDefinition.GetId() != nil {
			details["id"] = *roleDefinition.GetId()
		}
		if roleDefinition.GetDisplayName() != nil {
			details["displayName"] = *roleDefinition.GetDisplayName()
		}
		if roleDefinition.GetExternalId() != nil {
			details["externalId"] = *roleDefinition.GetExternalId()
		}
		if roleDefinition.GetTemplateId() != nil {
			details["templateId"] = *roleDefinition.GetTemplateId()
		}
		if roleDefinition.GetResourceId() != nil {
			details["resourceId"] = *roleDefinition.GetResourceId()
		}

		return details
	}
	return nil
}

func adGovernanceRoleSettingTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	roleSetting := d.HydrateItem.(*ADGovernanceRoleSettingInfo)
	if roleSetting.GetId() != nil {
		return *roleSetting.GetId(), nil
	}
	return nil, nil
}

func listAdGovernanceRoleSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access governance role settings
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_governance_role_setting.listAdGovernanceRoleSettings", "connection_error", err)
		return nil, err
	}

	// Note: The API documentation shows this is accessed through privilegedAccess/azureResources/roleSettings
	// However, we'll use the GovernanceRoleSettings() method which should provide access to the same data
	result, err := client.GovernanceRoleSettings().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdGovernanceRoleSettings", "list_governance_role_settings_error", errObj)
		return nil, errObj
	}

	if result.GetValue() != nil {
		for _, roleSetting := range result.GetValue() {
			d.StreamListItem(ctx, &ADGovernanceRoleSettingInfo{roleSetting})
		}
	}

	return nil, nil
}

func getAdGovernanceRoleSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access governance role settings
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_governance_role_setting.getAdGovernanceRoleSetting", "connection_error", err)
		return nil, err
	}

	roleSettingId := d.EqualsQuals["id"].GetStringValue()
	if roleSettingId == "" {
		return nil, nil
	}

	result, err := client.GovernanceRoleSettings().ByGovernanceRoleSettingId(roleSettingId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdGovernanceRoleSetting", "get_governance_role_setting_error", errObj)
		return nil, errObj
	}

	return &ADGovernanceRoleSettingInfo{result}, nil
}
