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

func TagsToMap(tags []string) (*map[string]bool, error) {
	var turbotTagsMap map[string]bool
	if tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]bool{}
	for _, i := range tags {
		turbotTagsMap[i] = true
	}

	return &turbotTagsMap, nil
}