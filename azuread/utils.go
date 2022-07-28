package azuread

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// Constants for Standard Column Descriptions
const (
	ColumnDescriptionTenant = "The Azure Tenant ID where the resource is located."
	ColumnDescriptionTags   = "A map of tags for the resource."
	ColumnDescriptionTitle  = "Title of the resource."
)

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "Request_ResourceNotFound")
}

func isNotFoundErrorPredicate(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if err != nil {
			for _, item := range notFoundErrors {
				if strings.Contains(err.Error(), item) {
					return true
				}
			}
		}
		return false
	}
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

type QualsColumn struct {
	ColumnName string
	ColumnType string
	FilterName string
}

func getTenantId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getTenantId")

	session, err := GetNewSession(ctx, d)
	if err != nil {
		return nil, err
	}

	return session.TenantID, nil
}
