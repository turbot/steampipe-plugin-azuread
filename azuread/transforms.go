package azuread

import (
	"context"
	"time"

	betamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type ADAdminConsentRequestPolicyInfo struct {
	models.AdminConsentRequestPolicyable
}

type ADApplicationInfo struct {
	models.Applicationable
	IsAuthorizationServiceEnabled interface{}
}

type ADApplicationAppRoleAssignmentInfo struct {
	models.AppRoleAssignmentable
	ApplicationId *string
}

type ADAppRoleAssignmentInfo struct {
	models.AppRoleAssignmentable
}

type ADAuthorizationPolicyInfo struct {
	models.AuthorizationPolicyable
}

type ADConditionalAccessPolicyInfo struct {
	models.ConditionalAccessPolicyable
}

type ADCrossTenantAccessPolicyInfo struct {
	models.CrossTenantAccessPolicyable
}

type ADDeviceRegistrationPolicyInfo struct {
	models.DeviceRegistrationPolicyable
}

type ADDeviceInfo struct {
	models.Deviceable
}

type ADDirectoryAuditReportInfo struct {
	models.DirectoryAuditable
}

type ADDirectorySettingInfo struct {
	// models.GroupSettingable
	DisplayName *string
	Id          *string
	TemplateId  *string
	Name        *string
	Value       *string
}

type ADGroupInfo struct {
	models.Groupable
	ResourceBehaviorOptions     []string
	ResourceProvisioningOptions []string
}

type ADIdentityProviderInfo struct {
	models.BuiltInIdentityProvider
	ClientId     interface{}
	ClientSecret interface{}
}

type ADNamedLocationInfo struct {
	models.NamedLocationable
	NamedLocation models.NamedLocationable
}

type ADIpNamedLocationInfo struct {
	models.IpNamedLocationable
}

type ADCountryNamedLocationInfo struct {
	models.CountryNamedLocationable
}

type ADSecurityDefaultsPolicyInfo struct {
	models.IdentitySecurityDefaultsEnforcementPolicyable
}

type ADAuthenticationMethodPolicyInfo struct {
	betamodels.AuthenticationMethodsPolicyable
}

type ADServicePrincipalInfo struct {
	models.ServicePrincipalable
}

type ADSignInReportInfo struct {
	models.SignInable
}

type ADUserRegistrationDetailsReport struct {
	models.UserRegistrationDetailsable
}

type ADUserInfo struct {
	models.Userable
	RefreshTokensValidFromDateTime interface{}
}

type ADUserAppRoleAssignmentInfo struct {
	models.AppRoleAssignmentable
	UserId *string
}

type ADDirectoryRoleTemplateInfo struct {
	models.DirectoryRoleTemplateable
}

type ADDirectoryRoleEligibilityScheduleInstanceInfo struct {
	models.UnifiedRoleEligibilityScheduleInstanceable
}

type ADDirectoryRoleAssignmentInfo struct {
	models.UnifiedRoleAssignmentable
}

type ADDirectoryRoleDefinitionInfo struct {
	models.UnifiedRoleDefinitionable
}

type ADExternalIdentitiesPolicyInfo struct {
	betamodels.ExternalIdentitiesPolicyable
}

func (roleAssignment *ADDirectoryRoleAssignmentInfo) DirectoryRoleAssignmentPrincipal() map[string]interface{} {
	if roleAssignment.GetPrincipal() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if roleAssignment.GetPrincipal().GetId() != nil {
		data["id"] = *roleAssignment.GetPrincipal().GetId()
	}
	if roleAssignment.GetPrincipal().GetOdataType() != nil {
		data["@odata.type"] = *roleAssignment.GetPrincipal().GetOdataType()
	}

	return data
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) ResourceId() *string {
	if methodPolicy.GetId() != nil {
		return methodPolicy.GetId()
	}
	return nil
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) DisplayName() *string {
	if methodPolicy.GetDisplayName() != nil {
		return methodPolicy.GetDisplayName()
	}
	return nil
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) Description() *string {
	if methodPolicy.GetDescription() != nil {
		return methodPolicy.GetDescription()
	}
	return nil
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) LastModifiedDateTime() *time.Time {
	if methodPolicy.GetLastModifiedDateTime() != nil {
		return methodPolicy.GetLastModifiedDateTime()
	}
	return nil
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) PolicyMigrationState() string {
	if methodPolicy.GetPolicyMigrationState() != nil {
		return methodPolicy.GetPolicyMigrationState().String()
	}
	return ""
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) PolicyVersion() *string {
	if methodPolicy.GetPolicyVersion() != nil {
		return methodPolicy.GetPolicyVersion()
	}
	return nil
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) ReconfirmationInDays() *int32 {
	if methodPolicy.GetReconfirmationInDays() != nil {
		return methodPolicy.GetReconfirmationInDays()
	}
	return nil
}
func (methodPolicy *ADAuthenticationMethodPolicyInfo) AdditionalData() map[string]interface{} {
	if methodPolicy.GetAdditionalData() != nil {
		data := map[string]interface{}{}
		for key, value := range methodPolicy.GetAdditionalData() {
			data[key] = value
		}
		return data
	}
	return nil
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) RegistrationEnforcement() map[string]interface{} {
	if methodPolicy.GetRegistrationEnforcement() != nil {
		data := map[string]interface{}{}

		// Get authentication methods registration campaign
		if methodPolicy.GetRegistrationEnforcement().GetAuthenticationMethodsRegistrationCampaign() != nil {
			campaign := methodPolicy.GetRegistrationEnforcement().GetAuthenticationMethodsRegistrationCampaign()
			campaignData := map[string]interface{}{}

			// Additional data
			if campaign.GetAdditionalData() != nil {
				data := map[string]interface{}{}
				for key, value := range campaign.GetAdditionalData() {
					data[key] = value
				}
				campaignData["additionalData"] = data
			}

			// Exclude targets
			if campaign.GetExcludeTargets() != nil {
				excludeTargets := make([]map[string]interface{}, 0, len(campaign.GetExcludeTargets()))
				for _, target := range campaign.GetExcludeTargets() {
					if target != nil {
						targetData := map[string]interface{}{}
						if target.GetId() != nil {
							targetData["id"] = *target.GetId()
						}
						if target.GetOdataType() != nil {
							targetData["@odata.type"] = *target.GetOdataType()
						}
						if target.GetTargetType() != nil {
							targetData["targetType"] = target.GetTargetType().String()
						}
						if target.GetAdditionalData() != nil {
							targetData["additionalData"] = target.GetAdditionalData()
						}
						excludeTargets = append(excludeTargets, targetData)
					}
				}
				campaignData["excludeTargets"] = excludeTargets
			}

			// Include targets
			if campaign.GetIncludeTargets() != nil {
				includeTargets := make([]map[string]interface{}, 0, len(campaign.GetIncludeTargets()))
				for _, target := range campaign.GetIncludeTargets() {
					if target != nil {
						targetData := map[string]interface{}{}
						if target.GetId() != nil {
							targetData["id"] = *target.GetId()
						}
						if target.GetOdataType() != nil {
							targetData["@odata.type"] = *target.GetOdataType()
						}
						if target.GetTargetType() != nil {
							targetData["targetType"] = target.GetTargetType().String()
						}
						if target.GetTargetedAuthenticationMethod() != nil {
							targetData["targetedAuthenticationMethod"] = *target.GetTargetedAuthenticationMethod()
						}
						if target.GetAdditionalData() != nil {
							targetData["additionalData"] = target.GetAdditionalData()
						}
						includeTargets = append(includeTargets, targetData)
					}
				}
				campaignData["includeTargets"] = includeTargets
			}

			// Snooze duration in days
			if campaign.GetSnoozeDurationInDays() != nil {
				campaignData["snoozeDurationInDays"] = *campaign.GetSnoozeDurationInDays()
			}

			// State
			if campaign.GetState() != nil {
				campaignData["state"] = campaign.GetState().String()
			}

			// OData type
			if campaign.GetOdataType() != nil {
				campaignData["@odata.type"] = *campaign.GetOdataType()
			}

			data["authenticationMethodsRegistrationCampaign"] = campaignData
		}

		// OData type
		if methodPolicy.GetRegistrationEnforcement().GetOdataType() != nil {
			data["@odata.type"] = *methodPolicy.GetRegistrationEnforcement().GetOdataType()
		}

		// Additional data
		if methodPolicy.GetRegistrationEnforcement().GetAdditionalData() != nil {
			addData := make(map[string]interface{}, len(methodPolicy.GetRegistrationEnforcement().GetAdditionalData()))
			for key, value := range methodPolicy.GetRegistrationEnforcement().GetAdditionalData() {
				addData[key] = value
			}
			data["additionalData"] = addData
		}

		return data
	}
	return nil
}

func (methodPolicy *ADAuthenticationMethodPolicyInfo) AuthenticationMethodConfigurations() []map[string]interface{} {
	if methodPolicy.GetAuthenticationMethodConfigurations() != nil {
		configurations := make([]map[string]interface{}, 0, len(methodPolicy.GetAuthenticationMethodConfigurations()))

		for _, config := range methodPolicy.GetAuthenticationMethodConfigurations() {
			if config == nil {
				continue
			}

			configData := map[string]interface{}{}

			// Base properties from AuthenticationMethodConfiguration
			if config.GetId() != nil {
				configData["id"] = *config.GetId()
			}
			if config.GetOdataType() != nil {
				configData["@odata.type"] = *config.GetOdataType()
			}
			if config.GetState() != nil {
				configData["state"] = config.GetState().String()
			}
			if config.GetAdditionalData() != nil {
				data := map[string]interface{}{}
				for key, value := range config.GetAdditionalData() {
					data[key] = value
				}
				configData["additionalData"] = data
			}

			// Exclude targets (common to all configurations)
			if config.GetExcludeTargets() != nil {
				excludeTargets := make([]map[string]interface{}, 0, len(config.GetExcludeTargets()))
				for _, target := range config.GetExcludeTargets() {
					if target != nil {
						targetData := map[string]interface{}{}
						if target.GetId() != nil {
							targetData["id"] = *target.GetId()
						}
						if target.GetOdataType() != nil {
							targetData["@odata.type"] = *target.GetOdataType()
						}
						if target.GetTargetType() != nil {
							targetData["targetType"] = target.GetTargetType().String()
						}
						if target.GetAdditionalData() != nil {
							data := map[string]interface{}{}
							for key, value := range target.GetAdditionalData() {
								data[key] = value
							}
							targetData["additionalData"] = data
						}
						excludeTargets = append(excludeTargets, targetData)
					}
				}
				configData["excludeTargets"] = excludeTargets
			}

			// Type-specific properties
			switch configType := config.(type) {
			case betamodels.EmailAuthenticationMethodConfigurationable:
				// Email-specific properties
				if configType.GetAllowExternalIdToUseEmailOtp() != nil {
					configData["allowExternalIdToUseEmailOtp"] = configType.GetAllowExternalIdToUseEmailOtp().String()
				}
				if configType.GetIncludeTargets() != nil {
					includeTargets := make([]map[string]interface{}, 0, len(configType.GetIncludeTargets()))
					for _, target := range configType.GetIncludeTargets() {
						if target != nil {
							targetData := map[string]interface{}{}
							if target.GetId() != nil {
								targetData["id"] = *target.GetId()
							}
							if target.GetOdataType() != nil {
								targetData["@odata.type"] = *target.GetOdataType()
							}
							if target.GetTargetType() != nil {
								targetData["targetType"] = target.GetTargetType().String()
							}
							if target.GetIsRegistrationRequired() != nil {
								targetData["isRegistrationRequired"] = *target.GetIsRegistrationRequired()
							}
							if target.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range target.GetAdditionalData() {
									data[key] = value
								}
								targetData["additionalData"] = data
							}
							includeTargets = append(includeTargets, targetData)
						}
					}
					configData["includeTargets"] = includeTargets
				}

			case betamodels.Fido2AuthenticationMethodConfigurationable:
				// FIDO2-specific properties
				if configType.GetIsAttestationEnforced() != nil {
					configData["isAttestationEnforced"] = *configType.GetIsAttestationEnforced()
				}
				if configType.GetIsSelfServiceRegistrationAllowed() != nil {
					configData["isSelfServiceRegistrationAllowed"] = *configType.GetIsSelfServiceRegistrationAllowed()
				}
				if configType.GetKeyRestrictions() != nil {
					keyRestrictions := configType.GetKeyRestrictions()
					keyRestrictionsData := map[string]interface{}{}
					if keyRestrictions.GetAdditionalData() != nil {
						data := map[string]interface{}{}
						for key, value := range keyRestrictions.GetAdditionalData() {
							data[key] = value
						}
						keyRestrictionsData["additionalData"] = data
					}
					if keyRestrictions.GetOdataType() != nil {
						keyRestrictionsData["@odata.type"] = *keyRestrictions.GetOdataType()
					}
					configData["keyRestrictions"] = keyRestrictionsData
				}
				if configType.GetIncludeTargets() != nil {
					includeTargets := make([]map[string]interface{}, 0, len(configType.GetIncludeTargets()))
					for _, target := range configType.GetIncludeTargets() {
						if target != nil {
							targetData := map[string]interface{}{}
							if target.GetId() != nil {
								targetData["id"] = *target.GetId()
							}
							if target.GetOdataType() != nil {
								targetData["@odata.type"] = *target.GetOdataType()
							}
							if target.GetTargetType() != nil {
								targetData["targetType"] = target.GetTargetType().String()
							}
							if target.GetIsRegistrationRequired() != nil {
								targetData["isRegistrationRequired"] = *target.GetIsRegistrationRequired()
							}
							if target.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range target.GetAdditionalData() {
									data[key] = value
								}
								targetData["additionalData"] = data
							}
							includeTargets = append(includeTargets, targetData)
						}
					}
					configData["includeTargets"] = includeTargets
				}

			case betamodels.SmsAuthenticationMethodConfigurationable:
				// SMS-specific properties
				if configType.GetIncludeTargets() != nil {
					includeTargets := make([]map[string]interface{}, 0, len(configType.GetIncludeTargets()))
					for _, target := range configType.GetIncludeTargets() {
						if target != nil {
							targetData := map[string]interface{}{}
							if target.GetId() != nil {
								targetData["id"] = *target.GetId()
							}
							if target.GetOdataType() != nil {
								targetData["@odata.type"] = *target.GetOdataType()
							}
							if target.GetTargetType() != nil {
								targetData["targetType"] = target.GetTargetType().String()
							}
							if target.GetIsRegistrationRequired() != nil {
								targetData["isRegistrationRequired"] = *target.GetIsRegistrationRequired()
							}
							if target.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range target.GetAdditionalData() {
									data[key] = value
								}
								targetData["additionalData"] = data
							}
							// SMS-specific property
							if target.GetIsUsableForSignIn() != nil {
								targetData["isUsableForSignIn"] = *target.GetIsUsableForSignIn()
							}
							includeTargets = append(includeTargets, targetData)
						}
					}
					configData["includeTargets"] = includeTargets
				}

			case betamodels.MicrosoftAuthenticatorAuthenticationMethodConfigurationable:
				// Microsoft Authenticator-specific properties
				if configType.GetIsSoftwareOathEnabled() != nil {
					configData["isSoftwareOathEnabled"] = *configType.GetIsSoftwareOathEnabled()
				}
				if configType.GetIncludeTargets() != nil {
					includeTargets := make([]map[string]interface{}, 0, len(configType.GetIncludeTargets()))
					for _, target := range configType.GetIncludeTargets() {
						if target != nil {
							targetData := map[string]interface{}{}
							if target.GetId() != nil {
								targetData["id"] = *target.GetId()
							}
							if target.GetOdataType() != nil {
								targetData["@odata.type"] = *target.GetOdataType()
							}
							if target.GetTargetType() != nil {
								targetData["targetType"] = target.GetTargetType().String()
							}
							if target.GetIsRegistrationRequired() != nil {
								targetData["isRegistrationRequired"] = *target.GetIsRegistrationRequired()
							}
							if target.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range target.GetAdditionalData() {
									data[key] = value
								}
								targetData["additionalData"] = data
							}
							includeTargets = append(includeTargets, targetData)
						}
					}
					configData["includeTargets"] = includeTargets
				}
				// Feature settings
				if configType.GetFeatureSettings() != nil {
					featureSettings := configType.GetFeatureSettings()
					featureSettingsData := map[string]interface{}{}
					if featureSettings.GetOdataType() != nil {
						featureSettingsData["@odata.type"] = *featureSettings.GetOdataType()
					}
					if featureSettings.GetAdditionalData() != nil {
						data := map[string]interface{}{}
						for key, value := range featureSettings.GetAdditionalData() {
							data[key] = value
						}
						featureSettingsData["additionalData"] = data
					}
					if featureSettings.GetDisplayAppInformationRequiredState() != nil {
						displayAppInfo := featureSettings.GetDisplayAppInformationRequiredState()
						displayAppInfoData := map[string]interface{}{}
						if displayAppInfo.GetOdataType() != nil {
							displayAppInfoData["@odata.type"] = *displayAppInfo.GetOdataType()
						}
						if displayAppInfo.GetState() != nil {
							displayAppInfoData["state"] = displayAppInfo.GetState().String()
						}
						if displayAppInfo.GetExcludeTarget() != nil {
							excludeTarget := displayAppInfo.GetExcludeTarget()
							excludeTargetData := map[string]interface{}{}
							if excludeTarget.GetId() != nil {
								excludeTargetData["id"] = *excludeTarget.GetId()
							}
							if excludeTarget.GetOdataType() != nil {
								excludeTargetData["@odata.type"] = *excludeTarget.GetOdataType()
							}
							if excludeTarget.GetTargetType() != nil {
								excludeTargetData["targetType"] = excludeTarget.GetTargetType().String()
							}
							if excludeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range excludeTarget.GetAdditionalData() {
									data[key] = value
								}
								excludeTargetData["additionalData"] = data
							}
							displayAppInfoData["excludeTarget"] = excludeTargetData
						}
						if displayAppInfo.GetIncludeTarget() != nil {
							includeTarget := displayAppInfo.GetIncludeTarget()
							includeTargetData := map[string]interface{}{}
							if includeTarget.GetId() != nil {
								includeTargetData["id"] = *includeTarget.GetId()
							}
							if includeTarget.GetOdataType() != nil {
								includeTargetData["@odata.type"] = *includeTarget.GetOdataType()
							}
							if includeTarget.GetTargetType() != nil {
								includeTargetData["targetType"] = includeTarget.GetTargetType().String()
							}
							if includeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range includeTarget.GetAdditionalData() {
									data[key] = value
								}
								includeTargetData["additionalData"] = data
							}
							displayAppInfoData["includeTarget"] = includeTargetData
						}
						featureSettingsData["displayAppInformationRequiredState"] = displayAppInfoData
					}
					if featureSettings.GetDisplayLocationInformationRequiredState() != nil {
						displayLocationInfo := featureSettings.GetDisplayLocationInformationRequiredState()
						displayLocationInfoData := map[string]interface{}{}
						if displayLocationInfo.GetOdataType() != nil {
							displayLocationInfoData["@odata.type"] = *displayLocationInfo.GetOdataType()
						}
						if displayLocationInfo.GetState() != nil {
							displayLocationInfoData["state"] = displayLocationInfo.GetState().String()
						}
						if displayLocationInfo.GetExcludeTarget() != nil {
							excludeTarget := displayLocationInfo.GetExcludeTarget()
							excludeTargetData := map[string]interface{}{}
							if excludeTarget.GetId() != nil {
								excludeTargetData["id"] = *excludeTarget.GetId()
							}
							if excludeTarget.GetOdataType() != nil {
								excludeTargetData["@odata.type"] = *excludeTarget.GetOdataType()
							}
							if excludeTarget.GetTargetType() != nil {
								excludeTargetData["targetType"] = excludeTarget.GetTargetType().String()
							}
							if excludeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range excludeTarget.GetAdditionalData() {
									data[key] = value
								}
								excludeTargetData["additionalData"] = data
							}
							displayLocationInfoData["excludeTarget"] = excludeTargetData
						}
						if displayLocationInfo.GetIncludeTarget() != nil {
							includeTarget := displayLocationInfo.GetIncludeTarget()
							includeTargetData := map[string]interface{}{}
							if includeTarget.GetId() != nil {
								includeTargetData["id"] = *includeTarget.GetId()
							}
							if includeTarget.GetOdataType() != nil {
								includeTargetData["@odata.type"] = *includeTarget.GetOdataType()
							}
							if includeTarget.GetTargetType() != nil {
								includeTargetData["targetType"] = includeTarget.GetTargetType().String()
							}
							if includeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range includeTarget.GetAdditionalData() {
									data[key] = value
								}
								includeTargetData["additionalData"] = data
							}
							displayLocationInfoData["includeTarget"] = includeTargetData
						}
						featureSettingsData["displayLocationInformationRequiredState"] = displayLocationInfoData
					}

					// Add numberMatchingRequiredState using beta SDK method
					if featureSettings.GetNumberMatchingRequiredState() != nil {
						numberMatchingInfo := featureSettings.GetNumberMatchingRequiredState()
						numberMatchingData := map[string]interface{}{}
						if numberMatchingInfo.GetOdataType() != nil {
							numberMatchingData["@odata.type"] = *numberMatchingInfo.GetOdataType()
						}
						if numberMatchingInfo.GetState() != nil {
							numberMatchingData["state"] = numberMatchingInfo.GetState().String()
						}
						if numberMatchingInfo.GetExcludeTarget() != nil {
							excludeTarget := numberMatchingInfo.GetExcludeTarget()
							excludeTargetData := map[string]interface{}{}
							if excludeTarget.GetId() != nil {
								excludeTargetData["id"] = *excludeTarget.GetId()
							}
							if excludeTarget.GetOdataType() != nil {
								excludeTargetData["@odata.type"] = *excludeTarget.GetOdataType()
							}
							if excludeTarget.GetTargetType() != nil {
								excludeTargetData["targetType"] = excludeTarget.GetTargetType().String()
							}
							if excludeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range excludeTarget.GetAdditionalData() {
									data[key] = value
								}
								excludeTargetData["additionalData"] = data
							}
							numberMatchingData["excludeTarget"] = excludeTargetData
						}
						if numberMatchingInfo.GetIncludeTarget() != nil {
							includeTarget := numberMatchingInfo.GetIncludeTarget()
							includeTargetData := map[string]interface{}{}
							if includeTarget.GetId() != nil {
								includeTargetData["id"] = *includeTarget.GetId()
							}
							if includeTarget.GetOdataType() != nil {
								includeTargetData["@odata.type"] = *includeTarget.GetOdataType()
							}
							if includeTarget.GetTargetType() != nil {
								includeTargetData["targetType"] = includeTarget.GetTargetType().String()
							}
							if includeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range includeTarget.GetAdditionalData() {
									data[key] = value
								}
								includeTargetData["additionalData"] = data
							}
							numberMatchingData["includeTarget"] = includeTargetData
						}
						featureSettingsData["numberMatchingRequiredState"] = numberMatchingData
					}

					// Add companionAppAllowedState using beta SDK method
					if featureSettings.GetCompanionAppAllowedState() != nil {
						companionAppInfo := featureSettings.GetCompanionAppAllowedState()
						companionAppData := map[string]interface{}{}
						if companionAppInfo.GetOdataType() != nil {
							companionAppData["@odata.type"] = *companionAppInfo.GetOdataType()
						}
						if companionAppInfo.GetState() != nil {
							companionAppData["state"] = companionAppInfo.GetState().String()
						}
						if companionAppInfo.GetExcludeTarget() != nil {
							excludeTarget := companionAppInfo.GetExcludeTarget()
							excludeTargetData := map[string]interface{}{}
							if excludeTarget.GetId() != nil {
								excludeTargetData["id"] = *excludeTarget.GetId()
							}
							if excludeTarget.GetOdataType() != nil {
								excludeTargetData["@odata.type"] = *excludeTarget.GetOdataType()
							}
							if excludeTarget.GetTargetType() != nil {
								excludeTargetData["targetType"] = excludeTarget.GetTargetType().String()
							}
							if excludeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range excludeTarget.GetAdditionalData() {
									data[key] = value
								}
								excludeTargetData["additionalData"] = data
							}
							companionAppData["excludeTarget"] = excludeTargetData
						}
						if companionAppInfo.GetIncludeTarget() != nil {
							includeTarget := companionAppInfo.GetIncludeTarget()
							includeTargetData := map[string]interface{}{}
							if includeTarget.GetId() != nil {
								includeTargetData["id"] = *includeTarget.GetId()
							}
							if includeTarget.GetOdataType() != nil {
								includeTargetData["@odata.type"] = *includeTarget.GetOdataType()
							}
							if includeTarget.GetTargetType() != nil {
								includeTargetData["targetType"] = includeTarget.GetTargetType().String()
							}
							if includeTarget.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range includeTarget.GetAdditionalData() {
									data[key] = value
								}
								includeTargetData["additionalData"] = data
							}
							companionAppData["includeTarget"] = includeTargetData
						}
						featureSettingsData["companionAppAllowedState"] = companionAppData
					}

					configData["featureSettings"] = featureSettingsData
				}

			case betamodels.X509CertificateAuthenticationMethodConfigurationable:
				// X509 Certificate-specific properties
				if configType.GetAuthenticationModeConfiguration() != nil {
					authModeConfig := configType.GetAuthenticationModeConfiguration()
					authModeConfigData := map[string]interface{}{}
					if authModeConfig.GetOdataType() != nil {
						authModeConfigData["@odata.type"] = *authModeConfig.GetOdataType()
					}
					if authModeConfig.GetAdditionalData() != nil {
						data := map[string]interface{}{}
						for key, value := range authModeConfig.GetAdditionalData() {
							data[key] = value
						}
						authModeConfigData["additionalData"] = data
					}
					configData["authenticationModeConfiguration"] = authModeConfigData
				}
				if configType.GetCertificateUserBindings() != nil {
					certUserBindings := make([]map[string]interface{}, 0, len(configType.GetCertificateUserBindings()))
					for _, binding := range configType.GetCertificateUserBindings() {
						if binding != nil {
							bindingData := map[string]interface{}{}
							if binding.GetOdataType() != nil {
								bindingData["@odata.type"] = *binding.GetOdataType()
							}
							if binding.GetPriority() != nil {
								bindingData["priority"] = *binding.GetPriority()
							}
							if binding.GetUserProperty() != nil {
								bindingData["userProperty"] = *binding.GetUserProperty()
							}
							if binding.GetX509CertificateField() != nil {
								bindingData["x509CertificateField"] = *binding.GetX509CertificateField()
							}
							if binding.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range binding.GetAdditionalData() {
									data[key] = value
								}
								bindingData["additionalData"] = data
							}
							certUserBindings = append(certUserBindings, bindingData)
						}
					}
					configData["certificateUserBindings"] = certUserBindings
				}
				if configType.GetIncludeTargets() != nil {
					includeTargets := make([]map[string]interface{}, 0, len(configType.GetIncludeTargets()))
					for _, target := range configType.GetIncludeTargets() {
						if target != nil {
							targetData := map[string]interface{}{}
							if target.GetId() != nil {
								targetData["id"] = *target.GetId()
							}
							if target.GetOdataType() != nil {
								targetData["@odata.type"] = *target.GetOdataType()
							}
							if target.GetTargetType() != nil {
								targetData["targetType"] = target.GetTargetType().String()
							}
							if target.GetIsRegistrationRequired() != nil {
								targetData["isRegistrationRequired"] = *target.GetIsRegistrationRequired()
							}
							if target.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range target.GetAdditionalData() {
									data[key] = value
								}
								targetData["additionalData"] = data
							}
							includeTargets = append(includeTargets, targetData)
						}
					}
					configData["includeTargets"] = includeTargets
				}

			case betamodels.TemporaryAccessPassAuthenticationMethodConfigurationable:
				// Temporary Access Pass-specific properties
				if configType.GetDefaultLength() != nil {
					configData["defaultLength"] = *configType.GetDefaultLength()
				}
				if configType.GetDefaultLifetimeInMinutes() != nil {
					configData["defaultLifetimeInMinutes"] = *configType.GetDefaultLifetimeInMinutes()
				}
				if configType.GetIsUsableOnce() != nil {
					configData["isUsableOnce"] = *configType.GetIsUsableOnce()
				}
				if configType.GetMaximumLifetimeInMinutes() != nil {
					configData["maximumLifetimeInMinutes"] = *configType.GetMaximumLifetimeInMinutes()
				}
				if configType.GetMinimumLifetimeInMinutes() != nil {
					configData["minimumLifetimeInMinutes"] = *configType.GetMinimumLifetimeInMinutes()
				}
				if configType.GetIncludeTargets() != nil {
					includeTargets := make([]map[string]interface{}, 0, len(configType.GetIncludeTargets()))
					for _, target := range configType.GetIncludeTargets() {
						if target != nil {
							targetData := map[string]interface{}{}
							if target.GetId() != nil {
								targetData["id"] = *target.GetId()
							}
							if target.GetOdataType() != nil {
								targetData["@odata.type"] = *target.GetOdataType()
							}
							if target.GetTargetType() != nil {
								targetData["targetType"] = target.GetTargetType().String()
							}
							if target.GetIsRegistrationRequired() != nil {
								targetData["isRegistrationRequired"] = *target.GetIsRegistrationRequired()
							}
							if target.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range target.GetAdditionalData() {
									data[key] = value
								}
								targetData["additionalData"] = data
							}
							includeTargets = append(includeTargets, targetData)
						}
					}
					configData["includeTargets"] = includeTargets
				}

			case betamodels.VoiceAuthenticationMethodConfigurationable:
				// Voice-specific properties
				if configType.GetIsOfficePhoneAllowed() != nil {
					configData["isOfficePhoneAllowed"] = *configType.GetIsOfficePhoneAllowed()
				}
				if configType.GetIncludeTargets() != nil {
					includeTargets := make([]map[string]interface{}, 0, len(configType.GetIncludeTargets()))
					for _, target := range configType.GetIncludeTargets() {
						if target != nil {
							targetData := map[string]interface{}{}
							if target.GetId() != nil {
								targetData["id"] = *target.GetId()
							}
							if target.GetOdataType() != nil {
								targetData["@odata.type"] = *target.GetOdataType()
							}
							if target.GetTargetType() != nil {
								targetData["targetType"] = target.GetTargetType().String()
							}
							if target.GetIsRegistrationRequired() != nil {
								targetData["isRegistrationRequired"] = *target.GetIsRegistrationRequired()
							}
							if target.GetAdditionalData() != nil {
								data := map[string]interface{}{}
								for key, value := range target.GetAdditionalData() {
									data[key] = value
								}
								targetData["additionalData"] = data
							}
							includeTargets = append(includeTargets, targetData)
						}
					}
					configData["includeTargets"] = includeTargets
				}
			}

			configurations = append(configurations, configData)
		}

		return configurations
	}
	return nil
}

func (roleAssignment *ADDirectoryRoleAssignmentInfo) DirectoryRoleAssignmentRoleDefinition() map[string]interface{} {
	if roleAssignment.GetRoleDefinition() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if roleAssignment.GetRoleDefinition().GetId() != nil {
		data["id"] = *roleAssignment.GetRoleDefinition().GetId()
	}
	if roleAssignment.GetRoleDefinition().GetDisplayName() != nil {
		data["display_name"] = *roleAssignment.GetRoleDefinition().GetDisplayName()
	}
	if roleAssignment.GetRoleDefinition().GetDescription() != nil {
		data["description"] = *roleAssignment.GetRoleDefinition().GetDescription()
	}
	if roleAssignment.GetRoleDefinition().GetOdataType() != nil {
		data["@odata.type"] = *roleAssignment.GetRoleDefinition().GetOdataType()
	}

	return data
}

func (roleEligibilityScheduleInstance *ADDirectoryRoleEligibilityScheduleInstanceInfo) DirectoryRoleEligibilityScheduleInstanceAppScope() map[string]interface{} {
	if roleEligibilityScheduleInstance.GetAppScope() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if roleEligibilityScheduleInstance.GetAppScope().GetId() != nil {
		data["id"] = *roleEligibilityScheduleInstance.GetAppScope().GetId()
	}
	if roleEligibilityScheduleInstance.GetAppScope().GetDisplayName() != nil {
		data["display_name"] = *roleEligibilityScheduleInstance.GetAppScope().GetDisplayName()
	}
	if roleEligibilityScheduleInstance.GetAppScope().GetOdataType() != nil {
		data["@odata.type"] = *roleEligibilityScheduleInstance.GetAppScope().GetOdataType()
	}

	return data
}

func (roleEligibilityScheduleInstance *ADDirectoryRoleEligibilityScheduleInstanceInfo) DirectoryRoleEligibilityScheduleInstanceDirectoryScope() map[string]interface{} {
	if roleEligibilityScheduleInstance.GetDirectoryScope() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if roleEligibilityScheduleInstance.GetDirectoryScope().GetId() != nil {
		data["id"] = *roleEligibilityScheduleInstance.GetDirectoryScope().GetId()
	}
	if roleEligibilityScheduleInstance.GetDirectoryScope().GetOdataType() != nil {
		data["@odata.type"] = *roleEligibilityScheduleInstance.GetDirectoryScope().GetOdataType()
	}

	return data
}

func (roleEligibilityScheduleInstance *ADDirectoryRoleEligibilityScheduleInstanceInfo) DirectoryRoleEligibilityScheduleInstancePrincipal() map[string]interface{} {
	if roleEligibilityScheduleInstance.GetPrincipal() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if roleEligibilityScheduleInstance.GetPrincipal().GetId() != nil {
		data["id"] = *roleEligibilityScheduleInstance.GetPrincipal().GetId()
	}
	if roleEligibilityScheduleInstance.GetPrincipal().GetOdataType() != nil {
		data["@odata.type"] = *roleEligibilityScheduleInstance.GetPrincipal().GetOdataType()
	}

	return data
}

func (roleEligibilityScheduleInstance *ADDirectoryRoleEligibilityScheduleInstanceInfo) DirectoryRoleEligibilityScheduleInstanceRoleDefinition() map[string]interface{} {
	if roleEligibilityScheduleInstance.GetRoleDefinition() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if roleEligibilityScheduleInstance.GetRoleDefinition().GetId() != nil {
		data["id"] = *roleEligibilityScheduleInstance.GetRoleDefinition().GetId()
	}
	if roleEligibilityScheduleInstance.GetRoleDefinition().GetDisplayName() != nil {
		data["display_name"] = *roleEligibilityScheduleInstance.GetRoleDefinition().GetDisplayName()
	}
	if roleEligibilityScheduleInstance.GetRoleDefinition().GetDescription() != nil {
		data["description"] = *roleEligibilityScheduleInstance.GetRoleDefinition().GetDescription()
	}
	if roleEligibilityScheduleInstance.GetRoleDefinition().GetOdataType() != nil {
		data["@odata.type"] = *roleEligibilityScheduleInstance.GetRoleDefinition().GetOdataType()
	}

	return data
}

func (adminConsentRequestPolicy *ADAdminConsentRequestPolicyInfo) AdminConsentRequestPolicyReviewers() []map[string]interface{} {
	if adminConsentRequestPolicy.GetReviewers() == nil {
		return nil
	}
	reviewers := []map[string]interface{}{}

	for _, a := range adminConsentRequestPolicy.GetReviewers() {
		data := map[string]interface{}{}
		if a.GetOdataType() != nil {
			data["@odata.type"] = *a.GetOdataType()
		}
		if a.GetQuery() != nil {
			data["query"] = *a.GetQuery()
		}
		if a.GetQueryRoot() != nil {
			data["queryRoot"] = *a.GetQueryRoot()
		}
		if a.GetQueryType() != nil {
			data["queryType"] = *a.GetQueryType()
		}
		reviewers = append(reviewers, data)
	}

	return reviewers
}

func (application *ADApplicationInfo) ApplicationAPI() map[string]interface{} {
	if application.GetApi() == nil {
		return nil
	}

	apiData := map[string]interface{}{
		"knownClientApplications": application.GetApi().GetKnownClientApplications(),
	}

	if application.GetApi().GetAcceptMappedClaims() != nil {
		apiData["acceptMappedClaims"] = *application.GetApi().GetAcceptMappedClaims()
	}
	if application.GetApi().GetRequestedAccessTokenVersion() != nil {
		apiData["requestedAccessTokenVersion"] = *application.GetApi().GetRequestedAccessTokenVersion()
	}

	oauth2PermissionScopes := []map[string]interface{}{}
	for _, p := range application.GetApi().GetOauth2PermissionScopes() {
		data := map[string]interface{}{}
		if p.GetAdminConsentDescription() != nil {
			data["adminConsentDescription"] = *p.GetAdminConsentDescription()
		}
		if p.GetAdminConsentDisplayName() != nil {
			data["adminConsentDisplayName"] = *p.GetAdminConsentDisplayName()
		}
		if p.GetId() != nil {
			data["id"] = *p.GetId()
		}
		if p.GetIsEnabled() != nil {
			data["isEnabled"] = *p.GetIsEnabled()
		}
		if p.GetOrigin() != nil {
			data["origin"] = *p.GetOrigin()
		}
		if p.GetTypeEscaped() != nil {
			data["type"] = *p.GetTypeEscaped()
		}
		if p.GetUserConsentDescription() != nil {
			data["userConsentDescription"] = p.GetUserConsentDescription()
		}
		if p.GetUserConsentDisplayName() != nil {
			data["userConsentDisplayName"] = p.GetUserConsentDisplayName()
		}
		if p.GetValue() != nil {
			data["value"] = *p.GetValue()
		}
		oauth2PermissionScopes = append(oauth2PermissionScopes, data)
	}
	apiData["oauth2PermissionScopes"] = oauth2PermissionScopes

	preAuthorizedApplications := []map[string]interface{}{}
	for _, p := range application.GetApi().GetPreAuthorizedApplications() {
		data := map[string]interface{}{
			"delegatedPermissionIds": p.GetDelegatedPermissionIds(),
		}
		if p.GetAppId() != nil {
			data["appId"] = *p.GetAppId()
		}
		preAuthorizedApplications = append(preAuthorizedApplications, data)
	}
	apiData["preAuthorizedApplications"] = preAuthorizedApplications

	return apiData
}

func (application *ADApplicationInfo) ApplicationInfo() map[string]interface{} {
	if application.GetInfo() == nil {
		return nil
	}

	return map[string]interface{}{
		"logoUrl":             application.GetInfo().GetLogoUrl(),
		"marketingUrl":        application.GetInfo().GetMarketingUrl(),
		"privacyStatementUrl": application.GetInfo().GetPrivacyStatementUrl(),
		"supportUrl":          application.GetInfo().GetSupportUrl(),
		"termsOfServiceUrl":   application.GetInfo().GetTermsOfServiceUrl(),
	}
}

func (application *ADApplicationInfo) ApplicationKeyCredentials() []map[string]interface{} {
	if application.GetKeyCredentials() == nil {
		return nil
	}

	keyCredentials := []map[string]interface{}{}
	for _, p := range application.GetKeyCredentials() {
		keyCredentialData := map[string]interface{}{}
		if p.GetDisplayName() != nil {
			keyCredentialData["displayName"] = *p.GetDisplayName()
		}
		if p.GetEndDateTime() != nil {
			keyCredentialData["endDateTime"] = *p.GetEndDateTime()
		}
		if p.GetKeyId() != nil {
			keyCredentialData["keyId"] = *p.GetKeyId()
		}
		if p.GetStartDateTime() != nil {
			keyCredentialData["startDateTime"] = *p.GetStartDateTime()
		}
		if p.GetTypeEscaped() != nil {
			keyCredentialData["type"] = *p.GetTypeEscaped()
		}
		if p.GetUsage() != nil {
			keyCredentialData["usage"] = *p.GetUsage()
		}
		if p.GetCustomKeyIdentifier() != nil {
			keyCredentialData["customKeyIdentifier"] = p.GetCustomKeyIdentifier()
		}
		if p.GetKey() != nil {
			keyCredentialData["key"] = p.GetKey()
		}
		keyCredentials = append(keyCredentials, keyCredentialData)
	}

	return keyCredentials
}

func (application *ADApplicationInfo) ApplicationParentalControlSettings() map[string]interface{} {
	if application.GetParentalControlSettings() == nil {
		return nil
	}

	parentalControlSettingData := map[string]interface{}{
		"countriesBlockedForMinors": application.GetParentalControlSettings().GetCountriesBlockedForMinors(),
	}
	if application.GetParentalControlSettings().GetLegalAgeGroupRule() != nil {
		parentalControlSettingData["legalAgeGroupRule"] = *application.GetParentalControlSettings().GetLegalAgeGroupRule()
	}

	return parentalControlSettingData
}

func (application *ADApplicationInfo) ApplicationPasswordCredentials() []map[string]interface{} {
	if application.GetPasswordCredentials() == nil {
		return nil
	}

	passwordCredentials := []map[string]interface{}{}
	for _, p := range application.GetPasswordCredentials() {
		passwordCredentialData := map[string]interface{}{}
		if p.GetDisplayName() != nil {
			passwordCredentialData["displayName"] = *p.GetDisplayName()
		}
		if p.GetHint() != nil {
			passwordCredentialData["hint"] = *p.GetHint()
		}
		if p.GetSecretText() != nil {
			passwordCredentialData["secretText"] = *p.GetSecretText()
		}
		if p.GetKeyId() != nil {
			passwordCredentialData["keyId"] = *p.GetKeyId()
		}
		if p.GetEndDateTime() != nil {
			passwordCredentialData["endDateTime"] = *p.GetEndDateTime()
		}
		if p.GetStartDateTime() != nil {
			passwordCredentialData["startDateTime"] = *p.GetStartDateTime()
		}
		if p.GetCustomKeyIdentifier() != nil {
			passwordCredentialData["customKeyIdentifier"] = p.GetCustomKeyIdentifier()
		}
		passwordCredentials = append(passwordCredentials, passwordCredentialData)
	}

	return passwordCredentials
}

func (application *ADApplicationInfo) ApplicationSpa() map[string]interface{} {
	if application.GetSpa() == nil {
		return nil
	}

	return map[string]interface{}{
		"redirectUris": application.GetSpa().GetRedirectUris(),
	}
}

func (application *ADApplicationInfo) ApplicationWeb() map[string]interface{} {
	if application.GetWeb() == nil {
		return nil
	}

	webData := map[string]interface{}{}
	if application.GetWeb().GetHomePageUrl() != nil {
		webData["homePageUrl"] = *application.GetWeb().GetHomePageUrl()
	}
	if application.GetWeb().GetLogoutUrl() != nil {
		webData["logoutUrl"] = *application.GetWeb().GetLogoutUrl()
	}
	if application.GetWeb().GetRedirectUris() != nil {
		webData["redirectUris"] = application.GetWeb().GetRedirectUris()
	}
	if application.GetWeb().GetImplicitGrantSettings() != nil {
		implicitGrantSettingsData := map[string]*bool{}

		if application.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance() != nil {
			implicitGrantSettingsData["enableAccessTokenIssuance"] = application.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance()
		}
		if application.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance() != nil {
			implicitGrantSettingsData["enableIdTokenIssuance"] = application.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance()
		}
		webData["implicitGrantSettings"] = implicitGrantSettingsData
	}

	return webData
}

func (authorizationPolicy *ADAuthorizationPolicyInfo) AuthorizationPolicyDefaultUserRolePermissions() map[string]interface{} {
	if authorizationPolicy.GetDefaultUserRolePermissions() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToCreateApps() != nil {
		data["allowedToCreateApps"] = *authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToCreateApps()
	}
	if authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToCreateTenants() != nil {
		data["allowedToCreateTenants"] = *authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToCreateTenants()
	}
	if authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToCreateSecurityGroups() != nil {
		data["allowedToCreateSecurityGroups"] = *authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToCreateSecurityGroups()
	}
	if authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToReadOtherUsers() != nil {
		data["allowedToReadOtherUsers"] = *authorizationPolicy.GetDefaultUserRolePermissions().GetAllowedToReadOtherUsers()
	}
	if authorizationPolicy.GetDefaultUserRolePermissions().GetPermissionGrantPoliciesAssigned() != nil {
		data["permissionGrantPoliciesAssigned"] = authorizationPolicy.GetDefaultUserRolePermissions().GetPermissionGrantPoliciesAssigned()
	}

	return data
}

func (authorizationPolicy *ADAuthorizationPolicyInfo) AuthorizationPolicyAllowInvitesFrom() string {
	if authorizationPolicy.GetAllowInvitesFrom() == nil {
		return ""
	}
	return authorizationPolicy.GetAllowInvitesFrom().String()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsAdditionalData() interface{} {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}
	additionalData := make(map[string]interface{})
	if conditionalAccessPolicy.GetConditions().GetAdditionalData() != nil {
		addData := make(map[string]interface{})
		for k, v := range conditionalAccessPolicy.GetConditions().GetAdditionalData() {
			addData[k] = v
		}
		additionalData = addData
	}

	return additionalData
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsApplications() map[string]interface{} {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}

	if conditionalAccessPolicy.GetConditions().GetApplications() == nil {
		return nil
	}

	return map[string]interface{}{
		"excludeApplications":                         conditionalAccessPolicy.GetConditions().GetApplications().GetExcludeApplications(),
		"includeApplications":                         conditionalAccessPolicy.GetConditions().GetApplications().GetIncludeApplications(),
		"includeAuthenticationContextClassReferences": conditionalAccessPolicy.GetConditions().GetApplications().GetIncludeAuthenticationContextClassReferences(),
		"includeUserActions":                          conditionalAccessPolicy.GetConditions().GetApplications().GetIncludeUserActions(),
	}
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsClientAppTypes() []models.ConditionalAccessClientApp {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetConditions().GetClientAppTypes()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsLocations() map[string]interface{} {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}

	if conditionalAccessPolicy.GetConditions().GetLocations() == nil {
		return nil
	}

	return map[string]interface{}{
		"excludeLocations": conditionalAccessPolicy.GetConditions().GetLocations().GetExcludeLocations(),
		"includeLocations": conditionalAccessPolicy.GetConditions().GetLocations().GetIncludeLocations(),
	}
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsPlatforms() map[string]interface{} {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}

	if conditionalAccessPolicy.GetConditions().GetPlatforms() == nil {
		return nil
	}

	return map[string]interface{}{
		"excludePlatforms": conditionalAccessPolicy.GetConditions().GetPlatforms().GetExcludePlatforms(),
		"includePlatforms": conditionalAccessPolicy.GetConditions().GetPlatforms().GetIncludePlatforms(),
	}
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsSignInRiskLevels() []string {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}

	riskLevels := conditionalAccessPolicy.GetConditions().GetSignInRiskLevels()
	if riskLevels == nil {
		return nil
	}

	// Convert RiskLevel enums to strings for better readability
	result := make([]string, 0, len(riskLevels))
	for _, riskLevel := range riskLevels {
		result = append(result, riskLevel.String())
	}

	return result
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsServicePrincipalRiskLevels() []string {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}

	riskLevels := conditionalAccessPolicy.GetConditions().GetSignInRiskLevels()
	if riskLevels == nil {
		return nil
	}

	// Convert RiskLevel enums to strings for better readability
	result := make([]string, 0, len(riskLevels))
	for _, riskLevel := range riskLevels {
		result = append(result, riskLevel.String())
	}

	return result
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsUsers() map[string]interface{} {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}

	if conditionalAccessPolicy.GetConditions().GetUsers() == nil {
		return nil
	}

	return map[string]interface{}{
		"excludeGroups": conditionalAccessPolicy.GetConditions().GetUsers().GetExcludeGroups(),
		"excludeRoles":  conditionalAccessPolicy.GetConditions().GetUsers().GetExcludeRoles(),
		"excludeUsers":  conditionalAccessPolicy.GetConditions().GetUsers().GetExcludeUsers(),
		"includeGroups": conditionalAccessPolicy.GetConditions().GetUsers().GetIncludeGroups(),
		"includeRoles":  conditionalAccessPolicy.GetConditions().GetUsers().GetIncludeRoles(),
		"includeUsers":  conditionalAccessPolicy.GetConditions().GetUsers().GetIncludeUsers(),
	}
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsUserRiskLevels() []models.RiskLevel {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetConditions().GetUserRiskLevels()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyGrantControlsBuiltInControls() []models.ConditionalAccessGrantControl {
	if conditionalAccessPolicy.GetGrantControls() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetGrantControls().GetBuiltInControls()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyGrantAuthenticationStrength() []models.AuthenticationMethodModes {
	if conditionalAccessPolicy.GetGrantControls() == nil {
		return nil
	}
	if conditionalAccessPolicy.GetGrantControls().GetAuthenticationStrength() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetGrantControls().GetAuthenticationStrength().GetAllowedCombinations()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyGrantControlsCustomAuthenticationFactors() []string {
	if conditionalAccessPolicy.GetGrantControls() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetGrantControls().GetCustomAuthenticationFactors()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyGrantControlsOperator() *string {
	if conditionalAccessPolicy.GetGrantControls() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetGrantControls().GetOperator()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyGrantControlsTermsOfUse() []string {
	if conditionalAccessPolicy.GetGrantControls() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetGrantControls().GetTermsOfUse()
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicySessionControlsApplicationEnforcedRestrictions() map[string]interface{} {
	if conditionalAccessPolicy.GetSessionControls() == nil {
		return nil
	}
	if conditionalAccessPolicy.GetSessionControls().GetApplicationEnforcedRestrictions() == nil {
		return nil
	}

	data := map[string]interface{}{}
	if conditionalAccessPolicy.GetSessionControls().GetApplicationEnforcedRestrictions().GetIsEnabled() != nil {
		data["isEnabled"] = conditionalAccessPolicy.GetSessionControls().GetApplicationEnforcedRestrictions().GetIsEnabled()
	}
	if conditionalAccessPolicy.GetSessionControls().GetApplicationEnforcedRestrictions().GetOdataType() != nil {
		data["@odata.type"] = conditionalAccessPolicy.GetSessionControls().GetApplicationEnforcedRestrictions().GetOdataType()
	}
	return data
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicySessionControlsCloudAppSecurity() map[string]interface{} {
	if conditionalAccessPolicy.GetSessionControls() == nil {
		return nil
	}
	if conditionalAccessPolicy.GetSessionControls().GetCloudAppSecurity() == nil {
		return nil
	}

	data := map[string]interface{}{}
	if conditionalAccessPolicy.GetSessionControls().GetCloudAppSecurity().GetIsEnabled() != nil {
		data["isEnabled"] = conditionalAccessPolicy.GetSessionControls().GetCloudAppSecurity().GetIsEnabled()
	}
	if conditionalAccessPolicy.GetSessionControls().GetCloudAppSecurity().GetCloudAppSecurityType() != nil {
		data["cloudAppSecurityType"] = conditionalAccessPolicy.GetSessionControls().GetCloudAppSecurity().GetCloudAppSecurityType()
	}
	return data
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicySessionControlsPersistentBrowser() map[string]interface{} {
	if conditionalAccessPolicy.GetSessionControls() == nil {
		return nil
	}
	if conditionalAccessPolicy.GetSessionControls().GetPersistentBrowser() == nil {
		return nil
	}

	data := map[string]interface{}{}
	if conditionalAccessPolicy.GetSessionControls().GetPersistentBrowser().GetIsEnabled() != nil {
		data["isEnabled"] = conditionalAccessPolicy.GetSessionControls().GetPersistentBrowser().GetIsEnabled()
	}
	if conditionalAccessPolicy.GetSessionControls().GetPersistentBrowser().GetMode() != nil {
		data["mode"] = conditionalAccessPolicy.GetSessionControls().GetPersistentBrowser().GetMode()
	}
	return data
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicySessionControlsSignInFrequency() map[string]interface{} {
	if conditionalAccessPolicy.GetSessionControls() == nil {
		return nil
	}
	if conditionalAccessPolicy.GetSessionControls().GetSignInFrequency() == nil {
		return nil
	}

	data := map[string]interface{}{}
	if conditionalAccessPolicy.GetSessionControls().GetSignInFrequency().GetIsEnabled() != nil {
		data["isEnabled"] = conditionalAccessPolicy.GetSessionControls().GetSignInFrequency().GetIsEnabled()
	}
	if conditionalAccessPolicy.GetSessionControls().GetSignInFrequency().GetValue() != nil {
		data["value"] = conditionalAccessPolicy.GetSessionControls().GetSignInFrequency().GetValue()
	}
	return data
}

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicySessionControlsDisableResilienceDefaults() *bool { //map[string]interface{} {
	sessionControls := conditionalAccessPolicy.GetSessionControls()
	if sessionControls == nil {
		return nil
	}
	if sessionControls.GetDisableResilienceDefaults() == nil {
		return nil
	}

	return sessionControls.GetDisableResilienceDefaults()
}

// Cross-Tenant Access Policy Transform Functions

func (crossTenantAccessPolicy *ADCrossTenantAccessPolicyInfo) CrossTenantAccessPolicyAllowedCloudEndpoints() []string {
	return crossTenantAccessPolicy.GetAllowedCloudEndpoints()
}

func (crossTenantAccessPolicy *ADCrossTenantAccessPolicyInfo) CrossTenantAccessPolicyDefault() map[string]interface{} {
	if crossTenantAccessPolicy.GetDefaultEscaped() == nil {
		return nil
	}

	defaultConfig := crossTenantAccessPolicy.GetDefaultEscaped()
	return map[string]interface{}{
		"b2bCollaborationInbound":  defaultConfig.GetB2bCollaborationInbound(),
		"b2bCollaborationOutbound": defaultConfig.GetB2bCollaborationOutbound(),
		"b2bDirectConnectInbound":  defaultConfig.GetB2bDirectConnectInbound(),
		"b2bDirectConnectOutbound": defaultConfig.GetB2bDirectConnectOutbound(),
		"inboundTrust":             defaultConfig.GetInboundTrust(),
	}
}

func (crossTenantAccessPolicy *ADCrossTenantAccessPolicyInfo) CrossTenantAccessPolicyPartners() []map[string]interface{} {
	partners := crossTenantAccessPolicy.GetPartners()
	if partners == nil {
		return nil
	}

	var result []map[string]interface{}
	for _, partner := range partners {
		partnerData := map[string]interface{}{
			"tenantId": partner.GetTenantId(),
		}

		if partner.GetB2bCollaborationInbound() != nil {
			partnerData["b2bCollaborationInbound"] = partner.GetB2bCollaborationInbound()
		}
		if partner.GetB2bCollaborationOutbound() != nil {
			partnerData["b2bCollaborationOutbound"] = partner.GetB2bCollaborationOutbound()
		}
		if partner.GetB2bDirectConnectInbound() != nil {
			partnerData["b2bDirectConnectInbound"] = partner.GetB2bDirectConnectInbound()
		}
		if partner.GetB2bDirectConnectOutbound() != nil {
			partnerData["b2bDirectConnectOutbound"] = partner.GetB2bDirectConnectOutbound()
		}
		if partner.GetInboundTrust() != nil {
			partnerData["inboundTrust"] = partner.GetInboundTrust()
		}

		result = append(result, partnerData)
	}

	return result
}

// Device Registration Policy Transform Functions

func (deviceRegistrationPolicy *ADDeviceRegistrationPolicyInfo) DeviceRegistrationPolicyMultiFactorAuthConfiguration() string {
	if deviceRegistrationPolicy.GetMultiFactorAuthConfiguration() == nil {
		return ""
	}
	return deviceRegistrationPolicy.GetMultiFactorAuthConfiguration().String()
}

func (deviceRegistrationPolicy *ADDeviceRegistrationPolicyInfo) DeviceRegistrationPolicyAzureADRegistration() map[string]interface{} {
	if deviceRegistrationPolicy.GetAzureADRegistration() == nil {
		return nil
	}

	data := map[string]interface{}{}
	azureADReg := deviceRegistrationPolicy.GetAzureADRegistration()

	if azureADReg.GetIsAdminConfigurable() != nil {
		data["isAdminConfigurable"] = *azureADReg.GetIsAdminConfigurable()
	}

	if azureADReg.GetAllowedToRegister() != nil {
		allowedData := map[string]interface{}{}
		if azureADReg.GetAllowedToRegister().GetOdataType() != nil {
			allowedData["@odata.type"] = *azureADReg.GetAllowedToRegister().GetOdataType()
		}
		if azureADReg.GetAllowedToRegister().GetAdditionalData() != nil {
			allowedData["additionalData"] = azureADReg.GetAllowedToRegister().GetAdditionalData()
		}
		data["allowedToRegister"] = allowedData
	}

	return data
}

func (deviceRegistrationPolicy *ADDeviceRegistrationPolicyInfo) DeviceRegistrationPolicyAzureADJoin() map[string]interface{} {
	if deviceRegistrationPolicy.GetAzureADJoin() == nil {
		return nil
	}

	data := map[string]interface{}{}
	azureADJoin := deviceRegistrationPolicy.GetAzureADJoin()

	if azureADJoin.GetIsAdminConfigurable() != nil {
		data["isAdminConfigurable"] = *azureADJoin.GetIsAdminConfigurable()
	}

	if azureADJoin.GetAllowedToJoin() != nil {
		allowedData := map[string]interface{}{}
		if azureADJoin.GetAllowedToJoin().GetOdataType() != nil {
			allowedData["@odata.type"] = *azureADJoin.GetAllowedToJoin().GetOdataType()
		}
		if azureADJoin.GetAllowedToJoin().GetAdditionalData() != nil {
			allowedData["additionalData"] = azureADJoin.GetAllowedToJoin().GetAdditionalData()
		}
		data["allowedToJoin"] = allowedData
	}

	// Get additional data which may contain localAdmins and other fields
	if azureADJoin.GetAdditionalData() != nil {
		for key, value := range azureADJoin.GetAdditionalData() {
			data[key] = value
		}
	}

	return data
}

func (deviceRegistrationPolicy *ADDeviceRegistrationPolicyInfo) DeviceRegistrationPolicyLocalAdminPassword() map[string]interface{} {
	if deviceRegistrationPolicy.GetLocalAdminPassword() == nil {
		return nil
	}

	data := map[string]interface{}{}
	localAdminPassword := deviceRegistrationPolicy.GetLocalAdminPassword()

	if localAdminPassword.GetIsEnabled() != nil {
		data["isEnabled"] = *localAdminPassword.GetIsEnabled()
	}

	return data
}

func (device *ADDeviceInfo) DeviceMemberOf() []map[string]interface{} {
	if device.GetMemberOf() == nil {
		return nil
	}

	members := []map[string]interface{}{}
	for _, i := range device.GetMemberOf() {
		member := map[string]interface{}{
			"@odata.type": i.GetOdataType(),
			"id":          i.GetId(),
		}
		members = append(members, member)
	}
	return members
}

func (directoryAuditReport *ADDirectoryAuditReportInfo) DirectoryAuditAdditionalDetails() []map[string]interface{} {
	if directoryAuditReport.GetAdditionalDetails() == nil {
		return nil
	}
	additionalDetails := []map[string]interface{}{}

	for _, i := range directoryAuditReport.GetAdditionalDetails() {
		data := map[string]interface{}{}
		if i.GetKey() != nil {
			data["key"] = *i.GetKey()
		}
		if i.GetValue() != nil {
			data["value"] = *i.GetValue()
		}
		if i.GetOdataType() != nil {
			data["@odata.type"] = *i.GetOdataType()
		}
		additionalDetails = append(additionalDetails, data)
	}

	return additionalDetails
}

func (directoryAuditReport *ADDirectoryAuditReportInfo) DirectoryAuditInitiatedBy() map[string]interface{} {
	if directoryAuditReport.GetInitiatedBy() == nil {
		return nil
	}
	data := map[string]interface{}{}

	if directoryAuditReport.GetInitiatedBy().GetOdataType() != nil {
		data["@odata.type"] = *directoryAuditReport.GetInitiatedBy().GetOdataType()
	}

	if directoryAuditReport.GetInitiatedBy().GetUser() != nil {
		userData := map[string]interface{}{}

		if directoryAuditReport.GetInitiatedBy().GetUser().GetDisplayName() != nil {
			userData["displayName"] = *directoryAuditReport.GetInitiatedBy().GetUser().GetDisplayName()
		}
		if directoryAuditReport.GetInitiatedBy().GetUser().GetId() != nil {
			userData["id"] = *directoryAuditReport.GetInitiatedBy().GetUser().GetId()
		}
		if directoryAuditReport.GetInitiatedBy().GetUser().GetUserPrincipalName() != nil {
			userData["userPrincipalName"] = *directoryAuditReport.GetInitiatedBy().GetUser().GetUserPrincipalName()
		}
		if directoryAuditReport.GetInitiatedBy().GetUser().GetIpAddress() != nil {
			userData["ipAddress"] = *directoryAuditReport.GetInitiatedBy().GetUser().GetIpAddress()
		}
		data["user"] = userData
	}

	if directoryAuditReport.GetInitiatedBy().GetApp() != nil {
		appData := map[string]interface{}{}

		if directoryAuditReport.GetInitiatedBy().GetApp().GetDisplayName() != nil {
			appData["displayName"] = *directoryAuditReport.GetInitiatedBy().GetApp().GetDisplayName()
		}
		if directoryAuditReport.GetInitiatedBy().GetApp().GetAppId() != nil {
			appData["appId"] = *directoryAuditReport.GetInitiatedBy().GetApp().GetAppId()
		}
		if directoryAuditReport.GetInitiatedBy().GetApp().GetServicePrincipalId() != nil {
			appData["servicePrincipalId"] = *directoryAuditReport.GetInitiatedBy().GetApp().GetServicePrincipalId()
		}
		if directoryAuditReport.GetInitiatedBy().GetApp().GetServicePrincipalName() != nil {
			appData["servicePrincipalName"] = *directoryAuditReport.GetInitiatedBy().GetApp().GetServicePrincipalName()
		}
		data["app"] = appData
	}

	return data
}

func (directoryAuditReport *ADDirectoryAuditReportInfo) DirectoryAuditResult() string {
	if directoryAuditReport.GetResult() == nil {
		return ""
	}
	return directoryAuditReport.GetResult().String()
}

func (directoryAuditReport *ADDirectoryAuditReportInfo) DirectoryAuditTargetResources() []map[string]interface{} {
	if directoryAuditReport.GetTargetResources() == nil {
		return nil
	}
	targetResources := []map[string]interface{}{}

	for _, i := range directoryAuditReport.GetTargetResources() {
		data := map[string]interface{}{}
		if i.GetDisplayName() != nil {
			data["displayName"] = *i.GetDisplayName()
		}
		if i.GetId() != nil {
			data["id"] = *i.GetId()
		}
		if i.GetOdataType() != nil {
			data["@odata.type"] = *i.GetOdataType()
		}
		if i.GetGroupType() != nil {
			data["groupType"] = i.GetGroupType().String()
		}
		if i.GetTypeEscaped() != nil {
			data["type"] = *i.GetTypeEscaped()
		}
		if i.GetUserPrincipalName() != nil {
			data["userPrincipalName"] = *i.GetUserPrincipalName()
		}

		modifiedProperties := []map[string]interface{}{}
		for _, m := range i.GetModifiedProperties() {
			prop := map[string]interface{}{}
			if m.GetDisplayName() != nil {
				prop["displayName"] = *m.GetDisplayName()
			}
			if m.GetNewValue() != nil {
				prop["newValue"] = *m.GetNewValue()
			}
			if m.GetOldValue() != nil {
				prop["oldValue"] = *m.GetOldValue()
			}
			if m.GetOdataType() != nil {
				prop["@odata.type"] = *m.GetOdataType()
			}
			modifiedProperties = append(modifiedProperties, prop)
		}
		data["modifiedProperties"] = modifiedProperties

		targetResources = append(targetResources, data)
	}

	return targetResources
}

// func (directorySetting *ADDirectorySettingInfo) DirectorySettingValues() []map[string]interface{} {
// 	if directorySetting.GetValues() == nil {
// 		return nil
// 	}
// 	values := []map[string]interface{}{}

// 	for _, v := range directorySetting.GetValues() {
// 		data := map[string]interface{}{}
// 		if v.GetName() != nil {
// 			data["name"] = *v.GetName()
// 		}
// 		if v.GetOdataType() != nil {
// 			data["@odata.type"] = *v.GetOdataType()
// 		}
// 		if v.GetValue() != nil {
// 			data["value"] = *v.GetValue()
// 		}
// 		values = append(values, data)
// 	}
// 	return values
// }

func (group *ADGroupInfo) GroupAssignedLabels() []map[string]*string {
	if group.GetAssignedLabels() == nil {
		return nil
	}

	assignedLabels := []map[string]*string{}
	for _, i := range group.GetAssignedLabels() {
		label := map[string]*string{
			"labelId":     i.GetLabelId(),
			"displayName": i.GetDisplayName(),
		}
		assignedLabels = append(assignedLabels, label)
	}
	return assignedLabels
}

func (servicePrincipal *ADServicePrincipalInfo) ServicePrincipalAddIns() []map[string]interface{} {
	if servicePrincipal.GetAddIns() == nil {
		return nil
	}

	addIns := []map[string]interface{}{}
	for _, p := range servicePrincipal.GetAddIns() {
		addInData := map[string]interface{}{}
		if p.GetId() != nil {
			addInData["id"] = *p.GetId()
		}
		if p.GetTypeEscaped() != nil {
			addInData["type"] = *p.GetTypeEscaped()
		}

		addInProperties := []map[string]interface{}{}
		for _, k := range p.GetProperties() {
			addInPropertyData := map[string]interface{}{}
			if k.GetKey() != nil {
				addInPropertyData["key"] = *k.GetKey()
			}
			if k.GetValue() != nil {
				addInPropertyData["value"] = *k.GetValue()
			}
			addInProperties = append(addInProperties, addInPropertyData)
		}
		addInData["properties"] = addInProperties

		addIns = append(addIns, addInData)
	}
	return addIns
}

func (servicePrincipal *ADServicePrincipalInfo) ServicePrincipalAppRoles() []map[string]interface{} {
	if servicePrincipal.GetAppRoles() == nil {
		return nil
	}

	appRoles := []map[string]interface{}{}
	for _, p := range servicePrincipal.GetAppRoles() {
		appRoleData := map[string]interface{}{
			"allowedMemberTypes": p.GetAllowedMemberTypes(),
		}
		if p.GetDescription() != nil {
			appRoleData["description"] = *p.GetDescription()
		}
		if p.GetDisplayName() != nil {
			appRoleData["displayName"] = *p.GetDisplayName()
		}
		if p.GetId() != nil {
			appRoleData["id"] = *p.GetId()
		}
		if p.GetIsEnabled() != nil {
			appRoleData["isEnabled"] = *p.GetIsEnabled()
		}
		if p.GetOrigin() != nil {
			appRoleData["origin"] = *p.GetOrigin()
		}
		if p.GetValue() != nil {
			appRoleData["value"] = *p.GetValue()
		}
		appRoles = append(appRoles, appRoleData)
	}
	return appRoles
}

func (servicePrincipal *ADServicePrincipalInfo) ServicePrincipalInfo() map[string]interface{} {
	if servicePrincipal.GetInfo() == nil {
		return nil
	}

	return map[string]interface{}{
		"logoUrl":             servicePrincipal.GetInfo().GetLogoUrl(),
		"marketingUrl":        servicePrincipal.GetInfo().GetMarketingUrl(),
		"privacyStatementUrl": servicePrincipal.GetInfo().GetPrivacyStatementUrl(),
		"supportUrl":          servicePrincipal.GetInfo().GetSupportUrl(),
		"termsOfServiceUrl":   servicePrincipal.GetInfo().GetTermsOfServiceUrl(),
	}
}

func (servicePrincipal *ADServicePrincipalInfo) ServicePrincipalKeyCredentials() []map[string]interface{} {
	if servicePrincipal.GetKeyCredentials() == nil {
		return nil
	}

	keyCredentials := []map[string]interface{}{}
	for _, p := range servicePrincipal.GetKeyCredentials() {
		keyCredentialData := map[string]interface{}{}
		if p.GetDisplayName() != nil {
			keyCredentialData["displayName"] = *p.GetDisplayName()
		}
		if p.GetEndDateTime() != nil {
			keyCredentialData["endDateTime"] = *p.GetEndDateTime()
		}
		if p.GetKeyId() != nil {
			keyCredentialData["keyId"] = *p.GetKeyId()
		}
		if p.GetStartDateTime() != nil {
			keyCredentialData["startDateTime"] = *p.GetStartDateTime()
		}
		if p.GetTypeEscaped() != nil {
			keyCredentialData["type"] = *p.GetTypeEscaped()
		}
		if p.GetUsage() != nil {
			keyCredentialData["usage"] = *p.GetUsage()
		}
		if p.GetCustomKeyIdentifier() != nil {
			keyCredentialData["customKeyIdentifier"] = p.GetCustomKeyIdentifier()
		}
		if p.GetKey() != nil {
			keyCredentialData["key"] = p.GetKey()
		}
		keyCredentials = append(keyCredentials, keyCredentialData)
	}
	return keyCredentials
}

func (servicePrincipal *ADServicePrincipalInfo) ServicePrincipalOauth2PermissionScopes() []map[string]interface{} {
	if servicePrincipal.GetOauth2PermissionScopes() == nil {
		return nil
	}

	oauth2PermissionScopes := []map[string]interface{}{}
	for _, p := range servicePrincipal.GetOauth2PermissionScopes() {
		data := map[string]interface{}{}
		if p.GetAdminConsentDescription() != nil {
			data["adminConsentDescription"] = *p.GetAdminConsentDescription()
		}
		if p.GetAdminConsentDisplayName() != nil {
			data["adminConsentDisplayName"] = *p.GetAdminConsentDisplayName()
		}
		if p.GetId() != nil {
			data["id"] = *p.GetId()
		}
		if p.GetIsEnabled() != nil {
			data["isEnabled"] = *p.GetIsEnabled()
		}
		if p.GetTypeEscaped() != nil {
			data["type"] = *p.GetTypeEscaped()
		}
		if p.GetOrigin() != nil {
			data["origin"] = *p.GetOrigin()
		}
		if p.GetUserConsentDescription() != nil {
			data["userConsentDescription"] = p.GetUserConsentDescription()
		}
		if p.GetUserConsentDisplayName() != nil {
			data["userConsentDisplayName"] = p.GetUserConsentDisplayName()
		}
		if p.GetValue() != nil {
			data["value"] = p.GetValue()
		}
		oauth2PermissionScopes = append(oauth2PermissionScopes, data)
	}
	return oauth2PermissionScopes
}

func (servicePrincipal *ADServicePrincipalInfo) ServicePrincipalPasswordCredentials() []map[string]interface{} {
	if servicePrincipal.GetPasswordCredentials() == nil {
		return nil
	}

	passwordCredentials := []map[string]interface{}{}
	for _, p := range servicePrincipal.GetPasswordCredentials() {
		passwordCredentialData := map[string]interface{}{}
		if p.GetDisplayName() != nil {
			passwordCredentialData["displayName"] = *p.GetDisplayName()
		}
		if p.GetHint() != nil {
			passwordCredentialData["hint"] = *p.GetHint()
		}
		if p.GetSecretText() != nil {
			passwordCredentialData["secretText"] = *p.GetSecretText()
		}
		if p.GetKeyId() != nil {
			passwordCredentialData["keyId"] = *p.GetKeyId()
		}
		if p.GetEndDateTime() != nil {
			passwordCredentialData["endDateTime"] = *p.GetEndDateTime()
		}
		if p.GetStartDateTime() != nil {
			passwordCredentialData["startDateTime"] = *p.GetStartDateTime()
		}
		if p.GetCustomKeyIdentifier() != nil {
			passwordCredentialData["customKeyIdentifier"] = p.GetCustomKeyIdentifier()
		}
		passwordCredentials = append(passwordCredentials, passwordCredentialData)
	}
	return passwordCredentials
}

func (signIn *ADSignInReportInfo) SignInAppliedConditionalAccessPolicies() []map[string]interface{} {
	if signIn.GetAppliedConditionalAccessPolicies() == nil {
		return nil
	}

	policies := []map[string]interface{}{}
	for _, p := range signIn.GetAppliedConditionalAccessPolicies() {
		policyData := map[string]interface{}{
			"enforcedGrantControls":   p.GetEnforcedGrantControls(),
			"enforcedSessionControls": p.GetEnforcedSessionControls(),
		}
		if p.GetDisplayName() != nil {
			policyData["displayName"] = *p.GetDisplayName()
		}
		if p.GetId() != nil {
			policyData["id"] = *p.GetId()
		}
		if p.GetResult() != nil {
			policyData["result"] = p.GetResult()
		}
		policies = append(policies, policyData)
	}

	return policies
}

func (signIn *ADSignInReportInfo) SignInDeviceDetail() map[string]interface{} {
	if signIn.GetDeviceDetail() == nil {
		return nil
	}

	deviceDetailInfo := map[string]interface{}{}
	if signIn.GetDeviceDetail().GetBrowser() != nil {
		deviceDetailInfo["browser"] = *signIn.GetDeviceDetail().GetBrowser()
	}
	if signIn.GetDeviceDetail().GetDeviceId() != nil {
		deviceDetailInfo["deviceId"] = *signIn.GetDeviceDetail().GetDeviceId()
	}
	if signIn.GetDeviceDetail().GetDisplayName() != nil {
		deviceDetailInfo["displayName"] = *signIn.GetDeviceDetail().GetDisplayName()
	}
	if signIn.GetDeviceDetail().GetIsCompliant() != nil {
		deviceDetailInfo["isCompliant"] = *signIn.GetDeviceDetail().GetIsCompliant()
	}
	if signIn.GetDeviceDetail().GetIsManaged() != nil {
		deviceDetailInfo["isManaged"] = *signIn.GetDeviceDetail().GetIsManaged()
	}
	if signIn.GetDeviceDetail().GetOperatingSystem() != nil {
		deviceDetailInfo["operatingSystem"] = *signIn.GetDeviceDetail().GetOperatingSystem()
	}
	if signIn.GetDeviceDetail().GetTrustType() != nil {
		deviceDetailInfo["trustType"] = *signIn.GetDeviceDetail().GetTrustType()
	}
	return deviceDetailInfo
}

func (signIn *ADSignInReportInfo) SignInStatus() map[string]interface{} {
	if signIn.GetStatus() == nil {
		return nil
	}

	statusInfo := map[string]interface{}{}
	if signIn.GetStatus().GetErrorCode() != nil {
		statusInfo["errorCode"] = *signIn.GetStatus().GetErrorCode()
	}
	if signIn.GetStatus().GetFailureReason() != nil {
		statusInfo["failureReason"] = *signIn.GetStatus().GetFailureReason()
	}
	if signIn.GetStatus().GetAdditionalDetails() != nil {
		statusInfo["additionalDetails"] = *signIn.GetStatus().GetAdditionalDetails()
	}
	return statusInfo
}

func (signIn *ADSignInReportInfo) SignInLocation() map[string]interface{} {
	if signIn.GetLocation() == nil {
		return nil
	}

	locationInfo := map[string]interface{}{}
	if signIn.GetLocation().GetCity() != nil {
		locationInfo["city"] = *signIn.GetLocation().GetCity()
	}
	if signIn.GetLocation().GetCountryOrRegion() != nil {
		locationInfo["countryOrRegion"] = *signIn.GetLocation().GetCountryOrRegion()
	}
	if signIn.GetLocation().GetState() != nil {
		locationInfo["state"] = *signIn.GetLocation().GetState()
	}
	if signIn.GetLocation().GetGeoCoordinates() != nil {
		coordinateInfo := map[string]interface{}{}
		if signIn.GetLocation().GetGeoCoordinates().GetAltitude() != nil {
			coordinateInfo["altitude"] = *signIn.GetLocation().GetGeoCoordinates().GetAltitude()
		}
		if signIn.GetLocation().GetGeoCoordinates().GetLatitude() != nil {
			coordinateInfo["latitude"] = *signIn.GetLocation().GetGeoCoordinates().GetLatitude()
		}
		if signIn.GetLocation().GetGeoCoordinates().GetLongitude() != nil {
			coordinateInfo["longitude"] = *signIn.GetLocation().GetGeoCoordinates().GetLongitude()
		}
		locationInfo["geoCoordinates"] = coordinateInfo
	}
	return locationInfo
}

func (user *ADUserInfo) UserMemberOf() []map[string]interface{} {
	if user.GetMemberOf() == nil {
		return nil
	}

	members := []map[string]interface{}{}
	for _, i := range user.GetMemberOf() {
		member := map[string]interface{}{
			"@odata.type": i.GetOdataType(),
			"id":          i.GetId(),
		}
		members = append(members, member)
	}
	return members
}

func (user *ADUserInfo) UserPasswordProfile() map[string]interface{} {
	if user.GetPasswordProfile() == nil {
		return nil
	}

	passwordProfileData := map[string]interface{}{}
	if user.GetPasswordProfile().GetForceChangePasswordNextSignIn() != nil {
		passwordProfileData["forceChangePasswordNextSignIn"] = *user.GetPasswordProfile().GetForceChangePasswordNextSignIn()
	}
	if user.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa() != nil {
		passwordProfileData["forceChangePasswordNextSignInWithMfa"] = *user.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa()
	}
	if user.GetPasswordProfile().GetPassword() != nil {
		passwordProfileData["password"] = *user.GetPasswordProfile().GetPassword()
	}

	return passwordProfileData
}

func (user *ADUserInfo) SignInActivity() map[string]interface{} {
	actiity := user.GetSignInActivity()
	if actiity == nil {
		return nil
	}

	return map[string]interface{}{
		"LastSignInDateTime":                actiity.GetLastSignInDateTime(),
		"LastSignInRequestId":               actiity.GetLastSignInRequestId(),
		"LastNonInteractiveSignInDateTime":  actiity.GetLastNonInteractiveSignInDateTime(),
		"LastNonInteractiveSignInRequestId": actiity.GetLastNonInteractiveSignInRequestId(),
	}
}

// Mailbox settings transform functions
func (user *ADUserInfo) GetAutomaticRepliesSetting(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	return mailboxSettings.GetAutomaticRepliesSetting(), nil
}

func (user *ADUserInfo) GetDateFormat(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	if dateFormat := mailboxSettings.GetDateFormat(); dateFormat != nil {
		return *dateFormat, nil
	}
	return nil, nil
}

func (user *ADUserInfo) GetDelegateMeetingMessageDeliveryOptions(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	if options := mailboxSettings.GetDelegateMeetingMessageDeliveryOptions(); options != nil {
		return options.String(), nil
	}
	return nil, nil
}

func (user *ADUserInfo) GetLanguage(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	return mailboxSettings.GetLanguage(), nil
}

func (user *ADUserInfo) GetTimeFormat(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	if timeFormat := mailboxSettings.GetTimeFormat(); timeFormat != nil {
		return *timeFormat, nil
	}
	return nil, nil
}

func (user *ADUserInfo) GetTimeZone(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	if timeZone := mailboxSettings.GetTimeZone(); timeZone != nil {
		return *timeZone, nil
	}
	return nil, nil
}

func (user *ADUserInfo) GetUserPurpose(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	if purpose := mailboxSettings.GetUserPurpose(); purpose != nil {
		return purpose.String(), nil
	}
	return nil, nil
}

func (user *ADUserInfo) GetWorkingHours(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	mailboxSettings := d.HydrateItem.(models.MailboxSettingsable)
	return mailboxSettings.GetWorkingHours(), nil
}
