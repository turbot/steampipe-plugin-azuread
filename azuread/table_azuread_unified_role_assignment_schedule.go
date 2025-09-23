package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	betamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func tableAzureAdUnifiedRoleAssignmentSchedule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_unified_role_assignment_schedule",
		Description: "Represents the schedule for an active role assignment operation in Azure AD.",
		Get: &plugin.GetConfig{
			Hydrate: getAdUnifiedRoleAssignmentSchedule,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdUnifiedRoleAssignmentSchedules,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the role assignment schedule.", Transform: transform.FromMethod("GetId")},
			{Name: "principal_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the principal (user, group, or service principal) to which the role is assigned.", Transform: transform.FromMethod("GetPrincipalId")},
			{Name: "role_definition_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the role definition.", Transform: transform.FromMethod("GetRoleDefinitionId")},
			{Name: "directory_scope_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the directory scope.", Transform: transform.FromMethod("GetDirectoryScopeId")},
			{Name: "app_scope_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the app scope.", Transform: transform.FromMethod("GetAppScopeId")},
			{Name: "created_using", Type: proto.ColumnType_STRING, Description: "The unique identifier of the object that was used to create this role assignment schedule.", Transform: transform.FromMethod("GetCreatedUsing")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the role assignment schedule was created.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "modified_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the role assignment schedule was last modified.", Transform: transform.FromMethod("GetModifiedDateTime")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the role assignment schedule (Provisioned, ProvisionedAsResource, etc.).", Transform: transform.FromMethod("GetStatus")},
			{Name: "assignment_type", Type: proto.ColumnType_STRING, Description: "The type of assignment (Assigned, Activated, etc.).", Transform: transform.FromMethod("GetAssignmentType")},
			{Name: "member_type", Type: proto.ColumnType_STRING, Description: "The type of member (Direct, Inherited, etc.).", Transform: transform.FromMethod("GetMemberType")},
			{Name: "schedule_info", Type: proto.ColumnType_JSON, Description: "The schedule information including start date, expiration, and recurrence.", Transform: transform.FromMethod("GetScheduleInfo")},
			{Name: "principal_details", Type: proto.ColumnType_JSON, Description: "Details about the principal (user, group, or service principal).", Transform: transform.FromMethod("GetPrincipalDetails")},
			{Name: "role_definition_details", Type: proto.ColumnType_JSON, Description: "Details about the role definition.", Transform: transform.FromMethod("GetRoleDefinitionDetails")},
			{Name: "directory_scope_details", Type: proto.ColumnType_JSON, Description: "Details about the directory scope.", Transform: transform.FromMethod("GetDirectoryScopeDetails")},
			{Name: "app_scope_details", Type: proto.ColumnType_JSON, Description: "Details about the app scope.", Transform: transform.FromMethod("GetAppScopeDetails")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adUnifiedRoleAssignmentScheduleTitle)},
		}),
	}
}

type ADUnifiedRoleAssignmentScheduleInfo struct {
	betamodels.UnifiedRoleAssignmentScheduleable
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetId() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetId()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetPrincipalId() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetPrincipalId()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetRoleDefinitionId() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetRoleDefinitionId()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetDirectoryScopeId() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetDirectoryScopeId()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetAppScopeId() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetAppScopeId()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetCreatedUsing() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetCreatedUsing()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetCreatedDateTime() *string {
	if schedule.UnifiedRoleAssignmentScheduleable.GetCreatedDateTime() != nil {
		dateTime := schedule.UnifiedRoleAssignmentScheduleable.GetCreatedDateTime().String()
		return &dateTime
	}
	return nil
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetModifiedDateTime() *string {
	if schedule.UnifiedRoleAssignmentScheduleable.GetModifiedDateTime() != nil {
		dateTime := schedule.UnifiedRoleAssignmentScheduleable.GetModifiedDateTime().String()
		return &dateTime
	}
	return nil
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetStatus() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetStatus()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetAssignmentType() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetAssignmentType()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetMemberType() *string {
	return schedule.UnifiedRoleAssignmentScheduleable.GetMemberType()
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetScheduleInfo() map[string]interface{} {
	if schedule.UnifiedRoleAssignmentScheduleable.GetScheduleInfo() != nil {
		scheduleInfo := schedule.UnifiedRoleAssignmentScheduleable.GetScheduleInfo()
		info := make(map[string]interface{})

		if scheduleInfo.GetStartDateTime() != nil {
			info["startDateTime"] = scheduleInfo.GetStartDateTime().String()
		}
		if scheduleInfo.GetExpiration() != nil {
			expiration := scheduleInfo.GetExpiration()
			expirationInfo := make(map[string]interface{})
			if expiration.GetTypeEscaped() != nil {
				expirationInfo["type"] = expiration.GetTypeEscaped().String()
			}
			if expiration.GetEndDateTime() != nil {
				expirationInfo["endDateTime"] = expiration.GetEndDateTime().String()
			}
			if expiration.GetDuration() != nil {
				expirationInfo["duration"] = *expiration.GetDuration()
			}
			info["expiration"] = expirationInfo
		}
		if scheduleInfo.GetRecurrence() != nil {
			recurrence := scheduleInfo.GetRecurrence()
			recurrenceInfo := make(map[string]interface{})
			if recurrence.GetPattern() != nil {
				pattern := recurrence.GetPattern()
				patternInfo := make(map[string]interface{})
				if pattern.GetTypeEscaped() != nil {
					patternInfo["type"] = pattern.GetTypeEscaped().String()
				}
				if pattern.GetInterval() != nil {
					patternInfo["interval"] = *pattern.GetInterval()
				}
				if pattern.GetMonth() != nil {
					patternInfo["month"] = *pattern.GetMonth()
				}
				if pattern.GetDayOfMonth() != nil {
					patternInfo["dayOfMonth"] = *pattern.GetDayOfMonth()
				}
				if pattern.GetDaysOfWeek() != nil {
					days := make([]string, 0, len(pattern.GetDaysOfWeek()))
					for _, day := range pattern.GetDaysOfWeek() {
						days = append(days, day.String())
					}
					patternInfo["daysOfWeek"] = days
				}
				if pattern.GetFirstDayOfWeek() != nil {
					patternInfo["firstDayOfWeek"] = pattern.GetFirstDayOfWeek().String()
				}
				if pattern.GetIndex() != nil {
					patternInfo["index"] = pattern.GetIndex().String()
				}
				recurrenceInfo["pattern"] = patternInfo
			}
			if recurrence.GetRangeEscaped() != nil {
				rangeInfo := recurrence.GetRangeEscaped()
				rangeData := make(map[string]interface{})
				if rangeInfo.GetTypeEscaped() != nil {
					rangeData["type"] = rangeInfo.GetTypeEscaped().String()
				}
				if rangeInfo.GetNumberOfOccurrences() != nil {
					rangeData["numberOfOccurrences"] = *rangeInfo.GetNumberOfOccurrences()
				}
				if rangeInfo.GetRecurrenceTimeZone() != nil {
					rangeData["recurrenceTimeZone"] = *rangeInfo.GetRecurrenceTimeZone()
				}
				if rangeInfo.GetStartDate() != nil {
					rangeData["startDate"] = rangeInfo.GetStartDate().String()
				}
				if rangeInfo.GetEndDate() != nil {
					rangeData["endDate"] = rangeInfo.GetEndDate().String()
				}
				recurrenceInfo["range"] = rangeData
			}
			info["recurrence"] = recurrenceInfo
		}

		return info
	}
	return nil
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetPrincipalDetails() map[string]interface{} {
	if schedule.UnifiedRoleAssignmentScheduleable.GetPrincipal() != nil {
		principal := schedule.UnifiedRoleAssignmentScheduleable.GetPrincipal()
		details := make(map[string]interface{})

		if principal.GetId() != nil {
			details["id"] = *principal.GetId()
		}
		if principal.GetOdataType() != nil {
			details["@odata.type"] = *principal.GetOdataType()
		}

		return details
	}
	return nil
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetRoleDefinitionDetails() map[string]interface{} {
	if schedule.UnifiedRoleAssignmentScheduleable.GetRoleDefinition() != nil {
		roleDefinition := schedule.UnifiedRoleAssignmentScheduleable.GetRoleDefinition()
		details := make(map[string]interface{})

		if roleDefinition.GetId() != nil {
			details["id"] = *roleDefinition.GetId()
		}
		if roleDefinition.GetDisplayName() != nil {
			details["displayName"] = *roleDefinition.GetDisplayName()
		}
		if roleDefinition.GetDescription() != nil {
			details["description"] = *roleDefinition.GetDescription()
		}
		if roleDefinition.GetIsBuiltIn() != nil {
			details["isBuiltIn"] = *roleDefinition.GetIsBuiltIn()
		}
		if roleDefinition.GetIsEnabled() != nil {
			details["isEnabled"] = *roleDefinition.GetIsEnabled()
		}
		if roleDefinition.GetTemplateId() != nil {
			details["templateId"] = *roleDefinition.GetTemplateId()
		}
		if roleDefinition.GetVersion() != nil {
			details["version"] = *roleDefinition.GetVersion()
		}

		return details
	}
	return nil
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetDirectoryScopeDetails() map[string]interface{} {
	if schedule.UnifiedRoleAssignmentScheduleable.GetDirectoryScope() != nil {
		directoryScope := schedule.UnifiedRoleAssignmentScheduleable.GetDirectoryScope()
		details := make(map[string]interface{})

		if directoryScope.GetId() != nil {
			details["id"] = *directoryScope.GetId()
		}
		if directoryScope.GetOdataType() != nil {
			details["@odata.type"] = *directoryScope.GetOdataType()
		}

		return details
	}
	return nil
}

func (schedule *ADUnifiedRoleAssignmentScheduleInfo) GetAppScopeDetails() map[string]interface{} {
	if schedule.UnifiedRoleAssignmentScheduleable.GetAppScope() != nil {
		appScope := schedule.UnifiedRoleAssignmentScheduleable.GetAppScope()
		details := make(map[string]interface{})

		if appScope.GetId() != nil {
			details["id"] = *appScope.GetId()
		}
		if appScope.GetDisplayName() != nil {
			details["displayName"] = *appScope.GetDisplayName()
		}
		if appScope.GetId() != nil {
			details["appScopeId"] = *appScope.GetId()
		}

		return details
	}
	return nil
}

func adUnifiedRoleAssignmentScheduleTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	schedule := d.HydrateItem.(*ADUnifiedRoleAssignmentScheduleInfo)
	if schedule.GetId() != nil {
		return *schedule.GetId(), nil
	}
	return nil, nil
}

func listAdUnifiedRoleAssignmentSchedules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access unified role assignment schedules
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_unified_role_assignment_schedule.listAdUnifiedRoleAssignmentSchedules", "connection_error", err)
		return nil, err
	}

	result, err := client.RoleManagement().Directory().RoleAssignmentSchedules().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdUnifiedRoleAssignmentSchedules", "list_unified_role_assignment_schedules_error", errObj)
		return nil, errObj
	}

	if result.GetValue() != nil {
		for _, schedule := range result.GetValue() {
			d.StreamListItem(ctx, &ADUnifiedRoleAssignmentScheduleInfo{schedule})
		}
	}

	return nil, nil
}

func getAdUnifiedRoleAssignmentSchedule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access unified role assignment schedules
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_unified_role_assignment_schedule.getAdUnifiedRoleAssignmentSchedule", "connection_error", err)
		return nil, err
	}

	scheduleId := d.EqualsQuals["id"].GetStringValue()
	if scheduleId == "" {
		return nil, nil
	}

	result, err := client.RoleManagement().Directory().RoleAssignmentSchedules().ByUnifiedRoleAssignmentScheduleId(scheduleId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdUnifiedRoleAssignmentSchedule", "get_unified_role_assignment_schedule_error", errObj)
		return nil, errObj
	}

	return &ADUnifiedRoleAssignmentScheduleInfo{result}, nil
}
