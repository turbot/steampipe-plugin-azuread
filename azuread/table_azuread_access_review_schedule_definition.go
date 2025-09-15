package azuread

import (
	"context"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	betamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func tableAzureAdAccessReviewScheduleDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_access_review_schedule_definition",
		Description: "Represents an access review schedule definition in Azure AD, which defines the settings and scope for recurring access reviews.",
		Get: &plugin.GetConfig{
			Hydrate: getAdAccessReviewScheduleDefinition,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdAccessReviewScheduleDefinitions,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the access review schedule definition.", Transform: transform.FromMethod("GetId")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the access review schedule definition.", Transform: transform.FromMethod("GetDisplayName")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the access review schedule definition (NotStarted, InProgress, Completed, etc.).", Transform: transform.FromMethod("GetStatus")},
			{Name: "description_for_admins", Type: proto.ColumnType_STRING, Description: "The description provided to administrators.", Transform: transform.FromMethod("GetDescriptionForAdmins")},
			{Name: "description_for_reviewers", Type: proto.ColumnType_STRING, Description: "The description provided to reviewers.", Transform: transform.FromMethod("GetDescriptionForReviewers")},
			{Name: "created_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the access review schedule definition was created.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "last_modified_date_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the access review schedule definition was last modified.", Transform: transform.FromMethod("GetLastModifiedDateTime")},
			{Name: "scope", Type: proto.ColumnType_JSON, Description: "The scope of the access review (what is being reviewed).", Transform: transform.FromMethod("GetScope")},
			{Name: "instance_enumeration_scope", Type: proto.ColumnType_JSON, Description: "The scope used to enumerate instances of the access review.", Transform: transform.FromMethod("GetInstanceEnumerationScope")},
			{Name: "reviewers", Type: proto.ColumnType_JSON, Description: "The reviewers assigned to the access review.", Transform: transform.FromMethod("GetReviewers")},
			{Name: "backup_reviewers", Type: proto.ColumnType_JSON, Description: "The backup reviewers assigned to the access review.", Transform: transform.FromMethod("GetBackupReviewers")},
			{Name: "fallback_reviewers", Type: proto.ColumnType_JSON, Description: "The fallback reviewers assigned to the access review.", Transform: transform.FromMethod("GetFallbackReviewers")},
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "The settings for the access review including notifications, recurrence, and other options.", Transform: transform.FromMethod("GetSettings")},
			{Name: "stage_settings", Type: proto.ColumnType_JSON, Description: "The stage settings for multi-stage access reviews.", Transform: transform.FromMethod("GetStageSettings")},
			{Name: "additional_notification_recipients", Type: proto.ColumnType_JSON, Description: "Additional recipients who will receive notifications about the access review.", Transform: transform.FromMethod("GetAdditionalNotificationRecipients")},
			{Name: "created_by", Type: proto.ColumnType_JSON, Description: "Information about the user who created the access review schedule definition.", Transform: transform.FromMethod("GetCreatedBy")},
			{Name: "instances", Type: proto.ColumnType_JSON, Description: "The instances of the access review schedule definition.", Transform: transform.FromMethod("GetInstances")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adAccessReviewScheduleDefinitionTitle)},
		}),
	}
}

type ADAccessReviewScheduleDefinitionInfo struct {
	betamodels.AccessReviewScheduleDefinitionable
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetId() *string {
	return definition.AccessReviewScheduleDefinitionable.GetId()
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetDisplayName() *string {
	return definition.AccessReviewScheduleDefinitionable.GetDisplayName()
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetStatus() *string {
	return definition.AccessReviewScheduleDefinitionable.GetStatus()
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetDescriptionForAdmins() *string {
	return definition.AccessReviewScheduleDefinitionable.GetDescriptionForAdmins()
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetDescriptionForReviewers() *string {
	return definition.AccessReviewScheduleDefinitionable.GetDescriptionForReviewers()
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetCreatedDateTime() *string {
	if definition.AccessReviewScheduleDefinitionable.GetCreatedDateTime() != nil {
		dateTime := definition.AccessReviewScheduleDefinitionable.GetCreatedDateTime().Format(time.RFC3339)
		return &dateTime
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetLastModifiedDateTime() *string {
	if definition.AccessReviewScheduleDefinitionable.GetLastModifiedDateTime() != nil {
		dateTime := definition.AccessReviewScheduleDefinitionable.GetLastModifiedDateTime().Format(time.RFC3339)
		return &dateTime
	}
	return nil
}
func (definition *ADAccessReviewScheduleDefinitionInfo) GetScope() map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetScope() != nil {
		scope := definition.AccessReviewScheduleDefinitionable.GetScope()
		scopeInfo := make(map[string]interface{})

		// Check if it's a query scope
		if queryScope, ok := scope.(betamodels.AccessReviewQueryScopeable); ok {
			if queryScope.GetQuery() != nil {
				scopeInfo["query"] = *queryScope.GetQuery()
			}
			if queryScope.GetQueryType() != nil {
				scopeInfo["queryType"] = *queryScope.GetQueryType()
			}
			if queryScope.GetQueryRoot() != nil {
				scopeInfo["queryRoot"] = *queryScope.GetQueryRoot()
			}
			scopeInfo["@odata.type"] = "#microsoft.graph.accessReviewQueryScope"
		} else if principalResourceScope, ok := scope.(betamodels.PrincipalResourceMembershipsScopeable); ok {
			// Handle principalResourceMembershipsScope
			scopeInfo["@odata.type"] = "#microsoft.graph.principalResourceMembershipsScope"
			
			// Extract principal scopes
			if principalResourceScope.GetPrincipalScopes() != nil {
				principalScopes := make([]map[string]interface{}, 0, len(principalResourceScope.GetPrincipalScopes()))
				for _, principalScope := range principalResourceScope.GetPrincipalScopes() {
					if principalScope != nil {
						principalScopeInfo := make(map[string]interface{})
						if queryScope, ok := principalScope.(betamodels.AccessReviewQueryScopeable); ok {
							if queryScope.GetQuery() != nil {
								principalScopeInfo["query"] = *queryScope.GetQuery()
							}
							if queryScope.GetQueryType() != nil {
								principalScopeInfo["queryType"] = *queryScope.GetQueryType()
							}
							if queryScope.GetQueryRoot() != nil {
								principalScopeInfo["queryRoot"] = *queryScope.GetQueryRoot()
							}
							principalScopeInfo["@odata.type"] = "#microsoft.graph.accessReviewQueryScope"
						}
						principalScopes = append(principalScopes, principalScopeInfo)
					}
				}
				scopeInfo["principalScopes"] = principalScopes
			}
			
			// Extract resource scopes
			if principalResourceScope.GetResourceScopes() != nil {
				resourceScopes := make([]map[string]interface{}, 0, len(principalResourceScope.GetResourceScopes()))
				for _, resourceScope := range principalResourceScope.GetResourceScopes() {
					if resourceScope != nil {
						resourceScopeInfo := make(map[string]interface{})
						if queryScope, ok := resourceScope.(betamodels.AccessReviewQueryScopeable); ok {
							if queryScope.GetQuery() != nil {
								resourceScopeInfo["query"] = *queryScope.GetQuery()
							}
							if queryScope.GetQueryType() != nil {
								resourceScopeInfo["queryType"] = *queryScope.GetQueryType()
							}
							if queryScope.GetQueryRoot() != nil {
								resourceScopeInfo["queryRoot"] = *queryScope.GetQueryRoot()
							}
							resourceScopeInfo["@odata.type"] = "#microsoft.graph.accessReviewQueryScope"
						}
						resourceScopes = append(resourceScopes, resourceScopeInfo)
					}
				}
				scopeInfo["resourceScopes"] = resourceScopes
			}
		}

		return scopeInfo
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetInstanceEnumerationScope() map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetInstanceEnumerationScope() != nil {
		scope := definition.AccessReviewScheduleDefinitionable.GetInstanceEnumerationScope()
		scopeInfo := make(map[string]interface{})

		// Check if it's a query scope
		if queryScope, ok := scope.(betamodels.AccessReviewQueryScopeable); ok {
			if queryScope.GetQuery() != nil {
				scopeInfo["query"] = *queryScope.GetQuery()
			}
			if queryScope.GetQueryType() != nil {
				scopeInfo["queryType"] = *queryScope.GetQueryType()
			}
			if queryScope.GetQueryRoot() != nil {
				scopeInfo["queryRoot"] = *queryScope.GetQueryRoot()
			}
		}

		return scopeInfo
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetReviewers() []map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetReviewers() != nil {
		reviewers := make([]map[string]interface{}, 0, len(definition.AccessReviewScheduleDefinitionable.GetReviewers()))
		for _, reviewer := range definition.AccessReviewScheduleDefinitionable.GetReviewers() {
			if reviewer != nil {
				reviewerInfo := make(map[string]interface{})
				// Check if it's a query scope
				if queryScope, ok := reviewer.(betamodels.AccessReviewQueryScopeable); ok {
					if queryScope.GetQuery() != nil {
						reviewerInfo["query"] = *queryScope.GetQuery()
					}
					if queryScope.GetQueryType() != nil {
						reviewerInfo["queryType"] = *queryScope.GetQueryType()
					}
					if queryScope.GetQueryRoot() != nil {
						reviewerInfo["queryRoot"] = *queryScope.GetQueryRoot()
					}
				}
				reviewers = append(reviewers, reviewerInfo)
			}
		}
		return reviewers
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetBackupReviewers() []map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetBackupReviewers() != nil {
		reviewers := make([]map[string]interface{}, 0, len(definition.AccessReviewScheduleDefinitionable.GetBackupReviewers()))
		for _, reviewer := range definition.AccessReviewScheduleDefinitionable.GetBackupReviewers() {
			if reviewer != nil {
				reviewerInfo := make(map[string]interface{})
				// Check if it's a query scope
				if queryScope, ok := reviewer.(betamodels.AccessReviewQueryScopeable); ok {
					if queryScope.GetQuery() != nil {
						reviewerInfo["query"] = *queryScope.GetQuery()
					}
					if queryScope.GetQueryType() != nil {
						reviewerInfo["queryType"] = *queryScope.GetQueryType()
					}
					if queryScope.GetQueryRoot() != nil {
						reviewerInfo["queryRoot"] = *queryScope.GetQueryRoot()
					}
				}
				reviewers = append(reviewers, reviewerInfo)
			}
		}
		return reviewers
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetFallbackReviewers() []map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetFallbackReviewers() != nil {
		reviewers := make([]map[string]interface{}, 0, len(definition.AccessReviewScheduleDefinitionable.GetFallbackReviewers()))
		for _, reviewer := range definition.AccessReviewScheduleDefinitionable.GetFallbackReviewers() {
			if reviewer != nil {
				reviewerInfo := make(map[string]interface{})
				// Check if it's a query scope
				if queryScope, ok := reviewer.(betamodels.AccessReviewQueryScopeable); ok {
					if queryScope.GetQuery() != nil {
						reviewerInfo["query"] = *queryScope.GetQuery()
					}
					if queryScope.GetQueryType() != nil {
						reviewerInfo["queryType"] = *queryScope.GetQueryType()
					}
					if queryScope.GetQueryRoot() != nil {
						reviewerInfo["queryRoot"] = *queryScope.GetQueryRoot()
					}
				}
				reviewers = append(reviewers, reviewerInfo)
			}
		}
		return reviewers
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetSettings() map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetSettings() != nil {
		settings := definition.AccessReviewScheduleDefinitionable.GetSettings()
		settingsInfo := make(map[string]interface{})

		if settings.GetMailNotificationsEnabled() != nil {
			settingsInfo["mailNotificationsEnabled"] = *settings.GetMailNotificationsEnabled()
		}
		if settings.GetReminderNotificationsEnabled() != nil {
			settingsInfo["reminderNotificationsEnabled"] = *settings.GetReminderNotificationsEnabled()
		}
		if settings.GetJustificationRequiredOnApproval() != nil {
			settingsInfo["justificationRequiredOnApproval"] = *settings.GetJustificationRequiredOnApproval()
		}
		if settings.GetDefaultDecisionEnabled() != nil {
			settingsInfo["defaultDecisionEnabled"] = *settings.GetDefaultDecisionEnabled()
		}
		if settings.GetDefaultDecision() != nil {
			settingsInfo["defaultDecision"] = *settings.GetDefaultDecision()
		}
		if settings.GetInstanceDurationInDays() != nil {
			settingsInfo["instanceDurationInDays"] = *settings.GetInstanceDurationInDays()
		}
		if settings.GetAutoApplyDecisionsEnabled() != nil {
			settingsInfo["autoApplyDecisionsEnabled"] = *settings.GetAutoApplyDecisionsEnabled()
		}
		if settings.GetRecommendationsEnabled() != nil {
			settingsInfo["recommendationsEnabled"] = *settings.GetRecommendationsEnabled()
		}
		if settings.GetRecurrence() != nil {
			recurrence := settings.GetRecurrence()
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
			settingsInfo["recurrence"] = recurrenceInfo
		}

		return settingsInfo
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetStageSettings() []map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetStageSettings() != nil {
		stageSettings := make([]map[string]interface{}, 0, len(definition.AccessReviewScheduleDefinitionable.GetStageSettings()))
		for _, stage := range definition.AccessReviewScheduleDefinitionable.GetStageSettings() {
			if stage != nil {
				stageInfo := make(map[string]interface{})
				if stage.GetStageId() != nil {
					stageInfo["stageId"] = *stage.GetStageId()
				}
				if stage.GetDurationInDays() != nil {
					stageInfo["durationInDays"] = *stage.GetDurationInDays()
				}
				if stage.GetFallbackReviewers() != nil {
					fallbackReviewers := make([]map[string]interface{}, 0, len(stage.GetFallbackReviewers()))
					for _, reviewer := range stage.GetFallbackReviewers() {
						if reviewer != nil {
							reviewerInfo := make(map[string]interface{})
							// Check if it's a query scope
							if queryScope, ok := reviewer.(betamodels.AccessReviewQueryScopeable); ok {
								if queryScope.GetQuery() != nil {
									reviewerInfo["query"] = *queryScope.GetQuery()
								}
								if queryScope.GetQueryType() != nil {
									reviewerInfo["queryType"] = *queryScope.GetQueryType()
								}
								if queryScope.GetQueryRoot() != nil {
									reviewerInfo["queryRoot"] = *queryScope.GetQueryRoot()
								}
							}
							fallbackReviewers = append(fallbackReviewers, reviewerInfo)
						}
					}
					stageInfo["fallbackReviewers"] = fallbackReviewers
				}
				if stage.GetReviewers() != nil {
					reviewers := make([]map[string]interface{}, 0, len(stage.GetReviewers()))
					for _, reviewer := range stage.GetReviewers() {
						if reviewer != nil {
							reviewerInfo := make(map[string]interface{})
							// Check if it's a query scope
							if queryScope, ok := reviewer.(betamodels.AccessReviewQueryScopeable); ok {
								if queryScope.GetQuery() != nil {
									reviewerInfo["query"] = *queryScope.GetQuery()
								}
								if queryScope.GetQueryType() != nil {
									reviewerInfo["queryType"] = *queryScope.GetQueryType()
								}
								if queryScope.GetQueryRoot() != nil {
									reviewerInfo["queryRoot"] = *queryScope.GetQueryRoot()
								}
							}
							reviewers = append(reviewers, reviewerInfo)
						}
					}
					stageInfo["reviewers"] = reviewers
				}
				stageSettings = append(stageSettings, stageInfo)
			}
		}
		return stageSettings
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetAdditionalNotificationRecipients() []map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetAdditionalNotificationRecipients() != nil {
		recipients := make([]map[string]interface{}, 0, len(definition.AccessReviewScheduleDefinitionable.GetAdditionalNotificationRecipients()))
		for _, recipient := range definition.AccessReviewScheduleDefinitionable.GetAdditionalNotificationRecipients() {
			if recipient != nil {
				recipientInfo := make(map[string]interface{})
				if recipient.GetNotificationTemplateType() != nil {
					recipientInfo["notificationTemplateType"] = *recipient.GetNotificationTemplateType()
				}
				if recipient.GetNotificationRecipientScope() != nil {
					recipientInfo["notificationRecipientScope"] = recipient.GetNotificationRecipientScope()
				}
				recipients = append(recipients, recipientInfo)
			}
		}
		return recipients
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetCreatedBy() map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetCreatedBy() != nil {
		createdBy := definition.AccessReviewScheduleDefinitionable.GetCreatedBy()
		createdByInfo := make(map[string]interface{})

		if createdBy.GetId() != nil {
			createdByInfo["id"] = *createdBy.GetId()
		}
		if createdBy.GetDisplayName() != nil {
			createdByInfo["displayName"] = *createdBy.GetDisplayName()
		}
		if createdBy.GetUserPrincipalName() != nil {
			createdByInfo["userPrincipalName"] = *createdBy.GetUserPrincipalName()
		}
		if createdBy.GetIpAddress() != nil {
			createdByInfo["ipAddress"] = *createdBy.GetIpAddress()
		}

		return createdByInfo
	}
	return nil
}

func (definition *ADAccessReviewScheduleDefinitionInfo) GetInstances() []map[string]interface{} {
	if definition.AccessReviewScheduleDefinitionable.GetInstances() != nil {
		instances := make([]map[string]interface{}, 0, len(definition.AccessReviewScheduleDefinitionable.GetInstances()))
		for _, instance := range definition.AccessReviewScheduleDefinitionable.GetInstances() {
			if instance != nil {
				instanceInfo := make(map[string]interface{})
				if instance.GetId() != nil {
					instanceInfo["id"] = *instance.GetId()
				}
				if instance.GetStartDateTime() != nil {
					instanceInfo["startDateTime"] = instance.GetStartDateTime().String()
				}
				if instance.GetEndDateTime() != nil {
					instanceInfo["endDateTime"] = instance.GetEndDateTime().String()
				}
				if instance.GetStatus() != nil {
					instanceInfo["status"] = *instance.GetStatus()
				}
				if instance.GetScope() != nil {
					scope := instance.GetScope()
					scopeInfo := make(map[string]interface{})
					// Check if it's a query scope
					if queryScope, ok := scope.(betamodels.AccessReviewQueryScopeable); ok {
						if queryScope.GetQuery() != nil {
							scopeInfo["query"] = *queryScope.GetQuery()
						}
						if queryScope.GetQueryType() != nil {
							scopeInfo["queryType"] = *queryScope.GetQueryType()
						}
						if queryScope.GetQueryRoot() != nil {
							scopeInfo["queryRoot"] = *queryScope.GetQueryRoot()
						}
					}
					instanceInfo["scope"] = scopeInfo
				}
				instances = append(instances, instanceInfo)
			}
		}
		return instances
	}
	return nil
}

func adAccessReviewScheduleDefinitionTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	definition := d.HydrateItem.(*ADAccessReviewScheduleDefinitionInfo)
	if definition.GetDisplayName() != nil {
		return *definition.GetDisplayName(), nil
	}
	if definition.GetId() != nil {
		return *definition.GetId(), nil
	}
	return nil, nil
}

func listAdAccessReviewScheduleDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access access review schedule definitions
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_access_review_schedule_definition.listAdAccessReviewScheduleDefinitions", "connection_error", err)
		return nil, err
	}

	result, err := client.IdentityGovernance().AccessReviews().Definitions().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdAccessReviewScheduleDefinitions", "list_access_review_schedule_definitions_error", errObj)
		return nil, errObj
	}

	if result.GetValue() != nil {
		for _, definition := range result.GetValue() {
			d.StreamListItem(ctx, &ADAccessReviewScheduleDefinitionInfo{definition})
		}
	}

	return nil, nil
}

func getAdAccessReviewScheduleDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access access review schedule definitions
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_access_review_schedule_definition.getAdAccessReviewScheduleDefinition", "connection_error", err)
		return nil, err
	}

	definitionId := d.EqualsQuals["id"].GetStringValue()
	if definitionId == "" {
		return nil, nil
	}

	result, err := client.IdentityGovernance().AccessReviews().Definitions().ByAccessReviewScheduleDefinitionId(definitionId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdAccessReviewScheduleDefinition", "get_access_review_schedule_definition_error", errObj)
		return nil, errObj
	}

	return &ADAccessReviewScheduleDefinitionInfo{result}, nil
}
