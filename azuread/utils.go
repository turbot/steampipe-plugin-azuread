package azuread

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Constants for Standard Column Descriptions
const (
	ColumnDescriptionTenant = "The Azure Tenant ID where the resource is located."
	ColumnDescriptionTags   = "A map of tags for the resource."
	ColumnDescriptionTitle  = "Title of the resource."
)

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

func getTenant(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getTenant")
	var tenantID string
	var err error

	// Read tenantID from config, or environment variables
	azureADConfig := GetConfig(d.Connection)
	if azureADConfig.TenantID != nil {
		tenantID = *azureADConfig.TenantID
	} else {
		tenantID = os.Getenv("AZURE_TENANT_ID")
	}

	// If not set in config, get tenantID from CLI
	if tenantID == "" {
		tenantID, err = getTenantFromCLI()
		if err != nil {
			return nil, err
		}
	}

	return tenantID, nil
}

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}
