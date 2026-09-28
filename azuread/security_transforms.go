package azuread

import (
	"time"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	securitymodels "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

type ADSecurityAlertInfo struct {
	securitymodels.Alertable
}

type ADSecureScoreControlScoreInfo struct {
	models.ControlScoreable
	MaxScore *float64
}

type ADSubscribedSkuInfo struct {
	models.SubscribedSkuable
}

func (a *ADSecurityAlertInfo) Severity() *string {
	if a.GetSeverity() == nil {
		return nil
	}
	return stringPtr(a.GetSeverity().String())
}

func (a *ADSecurityAlertInfo) Status() *string {
	if a.GetStatus() == nil {
		return nil
	}
	return stringPtr(a.GetStatus().String())
}

func (a *ADSecurityAlertInfo) Classification() *string {
	if a.GetClassification() == nil {
		return nil
	}
	return stringPtr(a.GetClassification().String())
}

func (a *ADSecurityAlertInfo) Determination() *string {
	if a.GetDetermination() == nil {
		return nil
	}
	return stringPtr(a.GetDetermination().String())
}

func (a *ADSecurityAlertInfo) ServiceSource() *string {
	if a.GetServiceSource() == nil {
		return nil
	}
	return stringPtr(a.GetServiceSource().String())
}

func (a *ADSecurityAlertInfo) Evidence() []map[string]interface{} {
	evidence := a.GetEvidence()
	if evidence == nil {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(evidence))
	for _, e := range evidence {
		if e == nil {
			continue
		}
		data := map[string]interface{}{}
		if e.GetOdataType() != nil {
			data["@odata.type"] = *e.GetOdataType()
		}
		if e.GetCreatedDateTime() != nil {
			data["createdDateTime"] = *e.GetCreatedDateTime()
		}
		if e.GetVerdict() != nil {
			data["verdict"] = e.GetVerdict().String()
		}
		if e.GetRemediationStatus() != nil {
			data["remediationStatus"] = e.GetRemediationStatus().String()
		}
		if e.GetRemediationStatusDetails() != nil {
			data["remediationStatusDetails"] = *e.GetRemediationStatusDetails()
		}
		if roles := e.GetRoles(); roles != nil {
			roleStrings := make([]string, 0, len(roles))
			for _, r := range roles {
				roleStrings = append(roleStrings, r.String())
			}
			data["roles"] = roleStrings
		}
		if e.GetDetailedRoles() != nil {
			data["detailedRoles"] = e.GetDetailedRoles()
		}
		if e.GetTags() != nil {
			data["tags"] = e.GetTags()
		}
		for k, v := range e.GetAdditionalData() {
			data[k] = v
		}
		out = append(out, data)
	}

	return out
}

func (a *ADSecurityAlertInfo) Comments() []map[string]interface{} {
	comments := a.GetComments()
	if comments == nil {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(comments))
	for _, c := range comments {
		if c == nil {
			continue
		}
		data := map[string]interface{}{}
		if c.GetComment() != nil {
			data["comment"] = *c.GetComment()
		}
		if c.GetCreatedByDisplayName() != nil {
			data["createdByDisplayName"] = *c.GetCreatedByDisplayName()
		}
		if c.GetCreatedDateTime() != nil {
			data["createdDateTime"] = *c.GetCreatedDateTime()
		}
		out = append(out, data)
	}

	return out
}

func (c *ADSecureScoreControlScoreInfo) ImplementationStatus() *string {
	return additionalDataString(c.GetAdditionalData(), "implementationStatus")
}

// LastSyncedDateTime reads the control score's lastSynced value, which Graph returns only in the
// control score's additional data as an ISO-8601 string.
func (c *ADSecureScoreControlScoreInfo) LastSyncedDateTime() *time.Time {
	raw := additionalDataString(c.GetAdditionalData(), "lastSynced")
	if raw == nil {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, *raw)
	if err != nil {
		return nil
	}
	return &parsed
}

func (s *ADSubscribedSkuInfo) SkuId() *string {
	if s.GetSkuId() == nil {
		return nil
	}
	return stringPtr(s.GetSkuId().String())
}

func (s *ADSubscribedSkuInfo) ServicePlans() []map[string]interface{} {
	plans := s.GetServicePlans()
	if plans == nil {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(plans))
	for _, p := range plans {
		if p == nil {
			continue
		}
		data := map[string]interface{}{}
		if p.GetServicePlanId() != nil {
			data["servicePlanId"] = p.GetServicePlanId().String()
		}
		if p.GetServicePlanName() != nil {
			data["servicePlanName"] = *p.GetServicePlanName()
		}
		if p.GetProvisioningStatus() != nil {
			data["provisioningStatus"] = *p.GetProvisioningStatus()
		}
		if p.GetAppliesTo() != nil {
			data["appliesTo"] = *p.GetAppliesTo()
		}
		out = append(out, data)
	}

	return out
}

func stringPtr(s string) *string {
	return &s
}

func additionalDataString(data map[string]any, key string) *string {
	if data == nil {
		return nil
	}
	value, ok := data[key]
	if !ok || value == nil {
		return nil
	}
	switch v := value.(type) {
	case *string:
		return v
	case string:
		return &v
	}
	return nil
}
