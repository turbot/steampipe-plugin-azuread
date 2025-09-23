package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	betamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func tableAzureAdUserAuthenticationMethod(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user_authentication_method",
		Description: "Represents authentication methods registered for users in Azure AD, including passwords, phone numbers, and hardware tokens.",
		Get: &plugin.GetConfig{
			Hydrate: getAdUserAuthenticationMethod,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.AllColumns([]string{"user_id", "id"}),
		},
		List: &plugin.ListConfig{
			Hydrate:    listAdUserAuthenticationMethods,
			KeyColumns: plugin.SingleColumn("user_id"),
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the user.", Transform: transform.FromQual("user_id")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for the authentication method.", Transform: transform.FromMethod("GetId")},
			{Name: "method_type", Type: proto.ColumnType_STRING, Description: "The type of authentication method (password, phone, email, etc.).", Transform: transform.FromMethod("GetMethodType")},
			{Name: "created_date_time", Type: proto.ColumnType_STRING, Description: "The date and time when the authentication method was created.", Transform: transform.FromMethod("GetCreatedDateTime")},
			{Name: "last_used_date_time", Type: proto.ColumnType_STRING, Description: "The date and time when the authentication method was last used.", Transform: transform.FromMethod("GetLastUsedDateTime")},
			{Name: "method_details", Type: proto.ColumnType_JSON, Description: "Method-specific details and properties.", Transform: transform.FromMethod("GetMethodDetails")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.From(adUserAuthenticationMethodTitle)},
		}),
	}
}

type ADUserAuthenticationMethodInfo struct {
	betamodels.AuthenticationMethodable
	UserId string
}

func (method *ADUserAuthenticationMethodInfo) GetId() *string {
	return method.AuthenticationMethodable.GetId()
}

func (method *ADUserAuthenticationMethodInfo) GetMethodType() *string {
	// Determine method type based on the concrete type
	switch method.AuthenticationMethodable.(type) {
	case betamodels.PasswordAuthenticationMethodable:
		methodType := "password"
		return &methodType
	case betamodels.PhoneAuthenticationMethodable:
		methodType := "phone"
		return &methodType
	case betamodels.EmailAuthenticationMethodable:
		methodType := "email"
		return &methodType
	case betamodels.Fido2AuthenticationMethodable:
		methodType := "fido2"
		return &methodType
	case betamodels.MicrosoftAuthenticatorAuthenticationMethodable:
		methodType := "microsoftAuthenticator"
		return &methodType
	case betamodels.TemporaryAccessPassAuthenticationMethodable:
		methodType := "temporaryAccessPass"
		return &methodType
	case betamodels.SoftwareOathAuthenticationMethodable:
		methodType := "softwareOath"
		return &methodType
	case betamodels.HardwareOathAuthenticationMethodable:
		methodType := "hardwareOath"
		return &methodType
	case betamodels.WindowsHelloForBusinessAuthenticationMethodable:
		methodType := "windowsHelloForBusiness"
		return &methodType
	case betamodels.PlatformCredentialAuthenticationMethodable:
		methodType := "platformCredential"
		return &methodType
	case betamodels.ExternalAuthenticationMethodable:
		methodType := "external"
		return &methodType
	case betamodels.QrCodePinAuthenticationMethodable:
		methodType := "qrCodePin"
		return &methodType
	case betamodels.PasswordlessMicrosoftAuthenticatorAuthenticationMethodable:
		methodType := "passwordlessMicrosoftAuthenticator"
		return &methodType
	case betamodels.ServiceNowAuthenticationMethodable:
		methodType := "serviceNow"
		return &methodType
	default:
		methodType := "unknown"
		return &methodType
	}
}

func (method *ADUserAuthenticationMethodInfo) GetCreatedDateTime() *string {
	if method.AuthenticationMethodable.GetCreatedDateTime() != nil {
		dateTime := method.AuthenticationMethodable.GetCreatedDateTime().String()
		return &dateTime
	}
	return nil
}

func (method *ADUserAuthenticationMethodInfo) GetLastUsedDateTime() *string {
	if method.AuthenticationMethodable.GetLastUsedDateTime() != nil {
		dateTime := method.AuthenticationMethodable.GetLastUsedDateTime().String()
		return &dateTime
	}
	return nil
}

func (method *ADUserAuthenticationMethodInfo) GetMethodDetails() map[string]interface{} {
	details := make(map[string]interface{})

	// Add common properties
	if method.AuthenticationMethodable.GetId() != nil {
		details["id"] = *method.AuthenticationMethodable.GetId()
	}
	if method.AuthenticationMethodable.GetCreatedDateTime() != nil {
		details["createdDateTime"] = method.AuthenticationMethodable.GetCreatedDateTime().String()
	}
	if method.AuthenticationMethodable.GetLastUsedDateTime() != nil {
		details["lastUsedDateTime"] = method.AuthenticationMethodable.GetLastUsedDateTime().String()
	}

	// Add method-specific properties based on type
	switch methodType := method.AuthenticationMethodable.(type) {
	case betamodels.PasswordAuthenticationMethodable:
		if methodType.GetPassword() != nil {
			details["password"] = *methodType.GetPassword()
		}
	case betamodels.PhoneAuthenticationMethodable:
		if methodType.GetPhoneNumber() != nil {
			details["phoneNumber"] = *methodType.GetPhoneNumber()
		}
		if methodType.GetPhoneType() != nil {
			details["phoneType"] = methodType.GetPhoneType().String()
		}
		if methodType.GetSmsSignInState() != nil {
			details["smsSignInState"] = methodType.GetSmsSignInState().String()
		}
	case betamodels.EmailAuthenticationMethodable:
		if methodType.GetEmailAddress() != nil {
			details["emailAddress"] = *methodType.GetEmailAddress()
		}
	case betamodels.Fido2AuthenticationMethodable:
		if methodType.GetAaGuid() != nil {
			details["aaGuid"] = *methodType.GetAaGuid()
		}
		if methodType.GetAttestationLevel() != nil {
			details["attestationLevel"] = methodType.GetAttestationLevel().String()
		}
		if methodType.GetAttestationCertificates() != nil {
			details["attestationCertificates"] = methodType.GetAttestationCertificates()
		}
		if methodType.GetCreatedDateTime() != nil {
			details["createdDateTime"] = methodType.GetCreatedDateTime().String()
		}
		if methodType.GetDisplayName() != nil {
			details["displayName"] = *methodType.GetDisplayName()
		}
		if methodType.GetModel() != nil {
			details["model"] = *methodType.GetModel()
		}
	case betamodels.MicrosoftAuthenticatorAuthenticationMethodable:
		if methodType.GetCreatedDateTime() != nil {
			details["createdDateTime"] = methodType.GetCreatedDateTime().String()
		}
		if methodType.GetDeviceTag() != nil {
			details["deviceTag"] = *methodType.GetDeviceTag()
		}
		if methodType.GetDisplayName() != nil {
			details["displayName"] = *methodType.GetDisplayName()
		}
		if methodType.GetPhoneAppVersion() != nil {
			details["phoneAppVersion"] = *methodType.GetPhoneAppVersion()
		}
		if methodType.GetDeviceTag() != nil {
			details["deviceTag"] = *methodType.GetDeviceTag()
		}
	case betamodels.TemporaryAccessPassAuthenticationMethodable:
		if methodType.GetCreatedDateTime() != nil {
			details["createdDateTime"] = methodType.GetCreatedDateTime().String()
		}
		if methodType.GetIsUsableOnce() != nil {
			details["isUsableOnce"] = *methodType.GetIsUsableOnce()
		}
		if methodType.GetLifetimeInMinutes() != nil {
			details["lifetimeInMinutes"] = *methodType.GetLifetimeInMinutes()
		}
		if methodType.GetStartDateTime() != nil {
			details["startDateTime"] = methodType.GetStartDateTime().String()
		}
		if methodType.GetTemporaryAccessPass() != nil {
			details["temporaryAccessPass"] = *methodType.GetTemporaryAccessPass()
		}
	case betamodels.SoftwareOathAuthenticationMethodable:
		if methodType.GetSecretKey() != nil {
			details["secretKey"] = *methodType.GetSecretKey()
		}
	case betamodels.HardwareOathAuthenticationMethodable:
		if methodType.GetDevice() != nil {
			device := methodType.GetDevice()
			deviceDetails := make(map[string]interface{})
			if device.GetId() != nil {
				deviceDetails["id"] = *device.GetId()
			}
			if device.GetDisplayName() != nil {
				deviceDetails["displayName"] = *device.GetDisplayName()
			}
			details["device"] = deviceDetails
		}
	case betamodels.WindowsHelloForBusinessAuthenticationMethodable:
		if methodType.GetCreatedDateTime() != nil {
			details["createdDateTime"] = methodType.GetCreatedDateTime().String()
		}
		if methodType.GetDevice() != nil {
			device := methodType.GetDevice()
			deviceDetails := make(map[string]interface{})
			if device.GetId() != nil {
				deviceDetails["id"] = *device.GetId()
			}
			if device.GetDisplayName() != nil {
				deviceDetails["displayName"] = *device.GetDisplayName()
			}
			details["device"] = deviceDetails
		}
		if methodType.GetDisplayName() != nil {
			details["displayName"] = *methodType.GetDisplayName()
		}
		if methodType.GetKeyStrength() != nil {
			details["keyStrength"] = methodType.GetKeyStrength().String()
		}
	case betamodels.PlatformCredentialAuthenticationMethodable:
		if methodType.GetCreatedDateTime() != nil {
			details["createdDateTime"] = methodType.GetCreatedDateTime().String()
		}
		if methodType.GetDevice() != nil {
			device := methodType.GetDevice()
			deviceDetails := make(map[string]interface{})
			if device.GetId() != nil {
				deviceDetails["id"] = *device.GetId()
			}
			if device.GetDisplayName() != nil {
				deviceDetails["displayName"] = *device.GetDisplayName()
			}
			details["device"] = deviceDetails
		}
		if methodType.GetDisplayName() != nil {
			details["displayName"] = *methodType.GetDisplayName()
		}
		if methodType.GetKeyStrength() != nil {
			details["keyStrength"] = methodType.GetKeyStrength().String()
		}
		if methodType.GetPlatform() != nil {
			details["platform"] = methodType.GetPlatform().String()
		}
	}

	return details
}

func adUserAuthenticationMethodTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	method := d.HydrateItem.(*ADUserAuthenticationMethodInfo)
	if method.GetId() != nil {
		return *method.GetId(), nil
	}
	return nil, nil
}

func listAdUserAuthenticationMethods(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access user authentication methods
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_user_authentication_method.listAdUserAuthenticationMethods", "connection_error", err)
		return nil, err
	}

	userId := d.EqualsQuals["user_id"].GetStringValue()
	if userId == "" {
		return nil, nil
	}

	result, err := client.Users().ByUserId(userId).Authentication().Methods().Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("listAdUserAuthenticationMethods", "list_user_authentication_methods_error", errObj)
		return nil, errObj
	}

	if result.GetValue() != nil {
		for _, method := range result.GetValue() {
			d.StreamListItem(ctx, &ADUserAuthenticationMethodInfo{
				AuthenticationMethodable: method,
				UserId:                   userId,
			})
		}
	}

	return nil, nil
}

func getAdUserAuthenticationMethod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create beta client to access user authentication methods
	client, _, err := GetGraphBetaClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuread_user_authentication_method.getAdUserAuthenticationMethod", "connection_error", err)
		return nil, err
	}

	userId := d.EqualsQuals["user_id"].GetStringValue()
	methodId := d.EqualsQuals["id"].GetStringValue()

	if userId == "" || methodId == "" {
		return nil, nil
	}

	result, err := client.Users().ByUserId(userId).Authentication().Methods().ByAuthenticationMethodId(methodId).Get(ctx, nil)
	if err != nil {
		errObj := getErrorObject(err)
		plugin.Logger(ctx).Error("getAdUserAuthenticationMethod", "get_user_authentication_method_error", errObj)
		return nil, errObj
	}

	return &ADUserAuthenticationMethodInfo{
		AuthenticationMethodable: result,
		UserId:                   userId,
	}, nil
}
