### Connection

```terraform
connection "azuread" {
  plugin = "azuread"

  # tenant_id     = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  # client_id     = "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"
  # client_secret = "ZZZZZZZZZZZZZZZZZZZZZZZZ"
}
```

#### Optional Qual notes

```go
  List: &plugin.ListConfig{
    Hydrate: listAdUsers,
    KeyColumns: plugin.KeyColumnSlice{
      {Name: "id", Require: plugin.Optional},
      {Name: "user_principal_name", Require: plugin.Optional}, // 'userPrincipalName eq 'lalit@yyyyyy.onmicrosoft.com'
      {Name: "filter", Require: plugin.Optional},              // where filter = 'displayName eq ''Luis'''

      // event fields
      {Name: "user_type", Require: plugin.Optional},                                       // filter=userType eq 'Guest'
      {Name: "account_enabled", Require: plugin.Optional, Operators: []string{"<>", "="}}, // accountEnabled eq true = for true and <> for false
      {Name: "display_name", Require: plugin.Optional},                                    // displayName eq 'Luis'
      {Name: "surname", Require: plugin.Optional},                                         // surname eq 'Luis'
    },
  }
```

````bash
select * from azuread.azuread_user where account_enabled
Column: account_enabled, Operator: '=', Value: 'true'

select * from azuread.azuread_user where not account_enabled
Column: account_enabled, Operator: '<>', Value: 'true'```
````

#### Odata Query based on the columns requested

```sql
select
  user_principal_name,
  jsonb_pretty(data),
  jsonb_pretty(member_of) as member_of
from
  azuread.azuread_user
where
  user_principal_name = 'lalit@turbotad.onmicrosoft.com'
```

```go
  // https://docs.microsoft.com/en-us/graph/api/user-list-memberof?view=graph-rest-1.0&tabs=http

  // https://graph.microsoft.com/v1.0/users?$filter=id eq '5137eeee-aaaa-bbbb-cccc-ddddc06b6636'&&$expand=memberOf($levels=max;$select=id,displayName)

  //  select user_principal_name, jsonb_pretty(data), jsonb_pretty(member_of) as member_of from azuread.azuread_user where user_principal_name = 'lalit@aaaabbbb.onmicrosoft.com'
  input := odata.Query{}
  if helpers.StringSliceContains(d.QueryContext.Columns, "member_of") {
    input.Expand = odata.Expand{
      Relationship: "memberOf",
      Select:       []string{"id", "displayName"},
    }
  }
```

```json
// https://graph.microsoft.com/v1.0/users/5137eeee-aaaa-bbbb-cccc-ddddc06b6636?$expand=memberOf
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users(memberOf())/$entity",
  "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.User",
  "businessPhones": [],
  "displayName": "Abhinash Khuntia",
  "givenName": null,
  "jobTitle": null,
  "mail": null,
  "mobilePhone": null,
  "officeLocation": null,
  "preferredLanguage": null,
  "surname": null,
  "userPrincipalName": "abhinash@turbotad.onmicrosoft.com",
  "id": "5137eeee-aaaa-bbbb-cccc-ddddc06b6636",
  "memberOf": [
    {
      "@odata.type": "#microsoft.graph.directoryRole",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.DirectoryRole",
      "id": "5137eeee-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "description": null,
      "displayName": null,
      "roleTemplateId": null
    },
    {
      "@odata.type": "#microsoft.graph.directoryRole",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.DirectoryRole",
      "id": "5137eeee-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "description": null,
      "displayName": null,
      "roleTemplateId": null
    },
    {
      "@odata.type": "#microsoft.graph.group",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.Group",
      "id": "5137eeee-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "classification": null,
      "createdDateTime": null,
      "creationOptions": [],
      "description": null,
      "displayName": null,
      "expirationDateTime": null,
      "groupTypes": [],
      "isAssignableToRole": null,
      "mail": null,
      "mailEnabled": null,
      "mailNickname": null,
      "membershipRule": null,
      "membershipRuleProcessingState": null,
      "onPremisesDomainName": null,
      "onPremisesLastSyncDateTime": null,
      "onPremisesNetBiosName": null,
      "onPremisesSamAccountName": null,
      "onPremisesSecurityIdentifier": null,
      "onPremisesSyncEnabled": null,
      "preferredDataLocation": null,
      "preferredLanguage": null,
      "proxyAddresses": [],
      "renewedDateTime": null,
      "resourceBehaviorOptions": [],
      "resourceProvisioningOptions": [],
      "securityEnabled": null,
      "securityIdentifier": null,
      "theme": null,
      "visibility": null,
      "onPremisesProvisioningErrors": []
    },
    {
      "@odata.type": "#microsoft.graph.group",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.Group",
      "id": "a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "classification": null,
      "createdDateTime": null,
      "creationOptions": [],
      "description": null,
      "displayName": null,
      "expirationDateTime": null,
      "groupTypes": [],
      "isAssignableToRole": null,
      "mail": null,
      "mailEnabled": null,
      "mailNickname": null,
      "membershipRule": null,
      "membershipRuleProcessingState": null,
      "onPremisesDomainName": null,
      "onPremisesLastSyncDateTime": null,
      "onPremisesNetBiosName": null,
      "onPremisesSamAccountName": null,
      "onPremisesSecurityIdentifier": null,
      "onPremisesSyncEnabled": null,
      "preferredDataLocation": null,
      "preferredLanguage": null,
      "proxyAddresses": [],
      "renewedDateTime": null,
      "resourceBehaviorOptions": [],
      "resourceProvisioningOptions": [],
      "securityEnabled": null,
      "securityIdentifier": null,
      "theme": null,
      "visibility": null,
      "onPremisesProvisioningErrors": []
    },
    {
      "@odata.type": "#microsoft.graph.group",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.Group",
      "id": "a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "classification": null,
      "createdDateTime": null,
      "creationOptions": [],
      "description": null,
      "displayName": null,
      "expirationDateTime": null,
      "groupTypes": [],
      "isAssignableToRole": null,
      "mail": null,
      "mailEnabled": null,
      "mailNickname": null,
      "membershipRule": null,
      "membershipRuleProcessingState": null,
      "onPremisesDomainName": null,
      "onPremisesLastSyncDateTime": null,
      "onPremisesNetBiosName": null,
      "onPremisesSamAccountName": null,
      "onPremisesSecurityIdentifier": null,
      "onPremisesSyncEnabled": null,
      "preferredDataLocation": null,
      "preferredLanguage": null,
      "proxyAddresses": [],
      "renewedDateTime": null,
      "resourceBehaviorOptions": [],
      "resourceProvisioningOptions": [],
      "securityEnabled": null,
      "securityIdentifier": null,
      "theme": null,
      "visibility": null,
      "onPremisesProvisioningErrors": []
    },
    {
      "@odata.type": "#microsoft.graph.group",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.Group",
      "id": "a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "classification": null,
      "createdDateTime": null,
      "creationOptions": [],
      "description": null,
      "displayName": null,
      "expirationDateTime": null,
      "groupTypes": [],
      "isAssignableToRole": null,
      "mail": null,
      "mailEnabled": null,
      "mailNickname": null,
      "membershipRule": null,
      "membershipRuleProcessingState": null,
      "onPremisesDomainName": null,
      "onPremisesLastSyncDateTime": null,
      "onPremisesNetBiosName": null,
      "onPremisesSamAccountName": null,
      "onPremisesSecurityIdentifier": null,
      "onPremisesSyncEnabled": null,
      "preferredDataLocation": null,
      "preferredLanguage": null,
      "proxyAddresses": [],
      "renewedDateTime": null,
      "resourceBehaviorOptions": [],
      "resourceProvisioningOptions": [],
      "securityEnabled": null,
      "securityIdentifier": null,
      "theme": null,
      "visibility": null,
      "onPremisesProvisioningErrors": []
    },
    {
      "@odata.type": "#microsoft.graph.group",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.Group",
      "id": "a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "classification": null,
      "createdDateTime": null,
      "creationOptions": [],
      "description": null,
      "displayName": null,
      "expirationDateTime": null,
      "groupTypes": [],
      "isAssignableToRole": null,
      "mail": null,
      "mailEnabled": null,
      "mailNickname": null,
      "membershipRule": null,
      "membershipRuleProcessingState": null,
      "onPremisesDomainName": null,
      "onPremisesLastSyncDateTime": null,
      "onPremisesNetBiosName": null,
      "onPremisesSamAccountName": null,
      "onPremisesSecurityIdentifier": null,
      "onPremisesSyncEnabled": null,
      "preferredDataLocation": null,
      "preferredLanguage": null,
      "proxyAddresses": [],
      "renewedDateTime": null,
      "resourceBehaviorOptions": [],
      "resourceProvisioningOptions": [],
      "securityEnabled": null,
      "securityIdentifier": null,
      "theme": null,
      "visibility": null,
      "onPremisesProvisioningErrors": []
    },
    {
      "@odata.type": "#microsoft.graph.group",
      "@odata.id": "https://graph.microsoft.com/v2/5137eeee-aaaa-bbbb-cccc-ddddc06b6636/directoryObjects/a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636/Microsoft.DirectoryServices.Group",
      "id": "a0b1c3d4-aaaa-bbbb-cccc-ddddc06b6636",
      "deletedDateTime": null,
      "classification": null,
      "createdDateTime": null,
      "creationOptions": [],
      "description": null,
      "displayName": null,
      "expirationDateTime": null,
      "groupTypes": [],
      "isAssignableToRole": null,
      "mail": null,
      "mailEnabled": null,
      "mailNickname": null,
      "membershipRule": null,
      "membershipRuleProcessingState": null,
      "onPremisesDomainName": null,
      "onPremisesLastSyncDateTime": null,
      "onPremisesNetBiosName": null,
      "onPremisesSamAccountName": null,
      "onPremisesSecurityIdentifier": null,
      "onPremisesSyncEnabled": null,
      "preferredDataLocation": null,
      "preferredLanguage": null,
      "proxyAddresses": [],
      "renewedDateTime": null,
      "resourceBehaviorOptions": [],
      "resourceProvisioningOptions": [],
      "securityEnabled": null,
      "securityIdentifier": null,
      "theme": null,
      "visibility": null,
      "onPremisesProvisioningErrors": []
    }
  ]
}
```
