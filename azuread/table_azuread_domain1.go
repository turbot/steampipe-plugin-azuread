package azuread

import (
	"context"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/domains"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureAdDomainTest() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_domain_test",
		Description: "Represents an Azure Active Directory (Azure AD) domain",
		Get: &plugin.GetConfig{
			Hydrate:    getAdDomainTest,
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listAdDomainsTest,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The fully qualified name of the domain.", Transform: transform.FromMethod("GetId")},
			{Name: "authentication_type", Type: proto.ColumnType_STRING, Description: "Indicates the configured authentication type for the domain. The value is either Managed or Federated. Managed indicates a cloud managed domain where Azure AD performs user authentication. Federated indicates authentication is federated with an identity provider such as the tenant's on-premises Active Directory via Active Directory Federation Services.", Transform: transform.FromMethod("GetAuthenticationType")},

			// Other fields
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "true if this is the default domain that is used for user creation. There is only one default domain per company.", Transform: transform.FromMethod("GetIsDefault")},
			{Name: "is_admin_managed", Type: proto.ColumnType_BOOL, Description: "The value of the property is false if the DNS record management of the domain has been delegated to Microsoft 365. Otherwise, the value is true.", Transform: transform.FromMethod("GetIsAdminManaged")},
			{Name: "is_initial", Type: proto.ColumnType_BOOL, Description: "true if this is the initial domain created by Microsoft Online Services (companyname.onmicrosoft.com). There is only one initial domain per company.", Transform: transform.FromMethod("GetIsInitial")},
			{Name: "is_root", Type: proto.ColumnType_BOOL, Description: "true if the domain is a verified root domain. Otherwise, false if the domain is a subdomain or unverified.", Transform: transform.FromMethod("GetIsRoot")},
			{Name: "is_verified", Type: proto.ColumnType_BOOL, Description: "true if the domain has completed domain ownership verification.", Transform: transform.FromMethod("GetIsVerified")},

			// Json fields
			{Name: "supported_services", Type: proto.ColumnType_JSON, Description: "The capabilities assigned to the domain. Can include 0, 1 or more of following values: Email, Sharepoint, EmailInternalRelayOnly, OfficeCommunicationsOnline, SharePointDefaultDomain, FullRedelegation, SharePointPublic, OrgIdAuthentication, Yammer, Intune. The values which you can add/remove using Graph API include: Email, OfficeCommunicationsOnline, Yammer.", Transform: transform.FromMethod("GetSupportedServices")},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromMethod("GetId")},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTenant, Hydrate: plugin.HydrateFunc(getTenant).WithCache(), Transform: transform.FromValue()},
		},
	}
}

type ADDomainInfo struct {
	models.Domainable
}

//// LIST FUNCTION

func listAdDomainsTest(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	client, adapter, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	// List operations
	input := &domains.DomainsRequestBuilderGetQueryParameters{}

	// Restrict the limit value to be passed in the query parameter which is not between 1 and 999, otherwise API will throw an error as follow
	// unexpected status 400 with OData error: Request_UnsupportedQuery: Invalid page size specified: '1000'. Must be between 1 and 999 inclusive.
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit > 0 && *limit <= 999 {
			l := int32(*limit)
			input.Top = &l
		}
	}

	options := &domains.DomainsRequestBuilderGetRequestConfiguration{
		QueryParameters: input,
	}

	result, err := client.Domains().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		errObj := getErrorObject(err)
		return nil, errors.New(fmt.Sprintf("failed to list domains. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	pageIterator, err := msgraphcore.NewPageIterator(result, adapter, models.CreateDomainCollectionResponseFromDiscriminatorValue)

	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		domain := pageItem.(models.Domainable)

		d.StreamListItem(ctx, &ADDomainInfo{domain})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return false
		}

		return true
	})

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAdDomainTest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	domainId := d.KeyColumnQuals["id"].GetStringValue()
	if domainId == "" {
		return nil, nil
	}

	// Create client
	client, _, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating client: %v", err))
	}

	domain, err := client.DomainsById(domainId).Get()
	if err != nil {
		errObj := getErrorObject(err)
		if isResourceNotFound(errObj) {
			return nil, nil
		}
		return nil, errors.New(fmt.Sprintf("failed to get domain. Code: %s Message: %s", errObj.Code, errObj.Message))
	}

	return &ADDomainInfo{domain}, nil
}
