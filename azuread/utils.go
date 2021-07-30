package azuread

import "strings"

// Constants for Standard Column Descriptions
const (
	ColumnDescriptionTenant = "The Azure Tenant ID where the resource is located."
	ColumnDescriptionTags   = "A map of tags for the resource."
	ColumnDescriptionTitle  = "Title of the resource."
)

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "Request_ResourceNotFound")
}

// func isNotFoundError(err error, notFoundErrors []string) bool {
// 	errorString := err.Error()
// 	for _, errorCode := range notFoundErrors {
// 		if strings.Contains(errorString, errorCode) {
// 			return true
// 		}
// 	}
// 	return false
// }
