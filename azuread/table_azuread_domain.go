package azuread

import (
	"context"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureAdDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_domain",
		Description: "Represents an Azure Active Directory (Azure AD) domain",
		Get: &plugin.GetConfig{
			Hydrate:           getAdDomain,
			ShouldIgnoreError: isNotFoundErrorPredicate([]string{"Invalid object identifier"}),
			KeyColumns:        plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDomains,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The fully qualified name of the domain.", Transform: transform.FromGo()},
			{Name: "authentication_type", Type: proto.ColumnType_STRING, Description: "Indicates the configured authentication type for the domain. The value is either Managed or Federated. Managed indicates a cloud managed domain where Azure AD performs user authentication. Federated indicates authentication is federated with an identity provider such as the tenant's on-premises Active Directory via Active Directory Federation Services."},

			// Other fields
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "true if this is the default domain that is used for user creation. There is only one default domain per company."},
			{Name: "is_admin_managed", Type: proto.ColumnType_BOOL, Description: "The value of the property is false if the DNS record management of the domain has been delegated to Microsoft 365. Otherwise, the value is true."},
			{Name: "is_initial", Type: proto.ColumnType_BOOL, Description: "true if this is the initial domain created by Microsoft Online Services (companyname.onmicrosoft.com). There is only one initial domain per company."},
			{Name: "is_root", Type: proto.ColumnType_BOOL, Description: "true if the domain is a verified root domain. Otherwise, false if the domain is a subdomain or unverified."},
			{Name: "is_verified", Type: proto.ColumnType_BOOL, Description: "true if the domain has completed domain ownership verification."},

			// // Json fields
			{Name: "supported_services", Type: proto.ColumnType_JSON, Description: "The capabilities assigned to the domain. Can include 0, 1 or more of following values: Email, Sharepoint, EmailInternalRelayOnly, OfficeCommunicationsOnline, SharePointDefaultDomain, FullRedelegation, SharePointPublic, OrgIdAuthentication, Yammer, Intune. The values which you can add/remove using Graph API include: Email, OfficeCommunicationsOnline, Yammer."},

			// // Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("DisplayName", "ID")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenantId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAdDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewDomainsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer

	input := odata.Query{
		Top: 999,
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < 999 {
			input.Top = int(*limit)
		}
	}

	domains, _, err := client.List(ctx, input)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, domain := range *domains {
		d.StreamListItem(ctx, domain)
	}

	return nil, err
}

//// Hydrate Functions

func getAdDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var domainId string
	if h.Item != nil {
		domainId = *h.Item.(msgraph.ServicePrincipal).ID
	} else {
		domainId = d.KeyColumnQuals["id"].GetStringValue()
	}

	if domainId == "" {
		return nil, nil
	}
	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	client := msgraph.NewDomainsClient(session.TenantID)
	client.BaseClient.Authorizer = session.Authorizer
	client.BaseClient.DisableRetries = true

	domain, _, err := client.Get(ctx, domainId, odata.Query{})
	if err != nil {
		return nil, err
	}
	return *domain, nil
}
