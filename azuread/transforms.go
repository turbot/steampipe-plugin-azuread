package azuread

import (
	"github.com/microsoftgraph/msgraph-sdk-go/models"
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

type ADServicePrincipalInfo struct {
	models.ServicePrincipalable
}

type ADSignInReportInfo struct {
	models.SignInable
}

type ADUserInfo struct {
	models.Userable
	RefreshTokensValidFromDateTime interface{}
}

type ADUserAppRoleAssignmentInfo struct {
	models.AppRoleAssignmentable
	UserId *string
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

func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicyConditionsSignInRiskLevels() []models.RiskLevel {
	if conditionalAccessPolicy.GetConditions() == nil {
		return nil
	}
	return conditionalAccessPolicy.GetConditions().GetSignInRiskLevels()
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


func (conditionalAccessPolicy *ADConditionalAccessPolicyInfo) ConditionalAccessPolicySessionControlsDisableResilienceDefaults() *bool{ //map[string]interface{} {
	sessionControls := conditionalAccessPolicy.GetSessionControls()
	if sessionControls == nil {
		return nil
	}
	if sessionControls.GetDisableResilienceDefaults() == nil {
		return nil
	}
	
	return sessionControls.GetDisableResilienceDefaults()
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
		"LastSignInDateTime": actiity.GetLastSignInDateTime(),
		"LastSignInRequestId": actiity.GetLastSignInRequestId(),
		"LastNonInteractiveSignInDateTime": actiity.GetLastNonInteractiveSignInDateTime(),
		"LastNonInteractiveSignInRequestId": actiity.GetLastNonInteractiveSignInRequestId(),
	}
}

