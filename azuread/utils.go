package azuread

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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

func commonColumns(c []*plugin.Column) []*plugin.Column {
	return append([]*plugin.Column{
		{
			Name:        "tenant_id",
			Type:        proto.ColumnType_STRING,
			Description: ColumnDescriptionTenant,
			Hydrate:     getTenant,
			Transform:   transform.FromValue(),
		},
	}, c...)
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize
// since getTenant is a call, caching should be per connection
var getTenantMemoized = plugin.HydrateFunc(getTenantUncached).Memoize(memoize.WithCacheKeyFunction(getTenantCacheKey))

// Build a cache key for the call to getTenant.
func getTenantCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getTenant"
	return key, nil
}

func getTenant(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	projectId, err := getTenantMemoized(ctx, d, h)
	if err != nil {
		return nil, err
	}

	return projectId, nil
}

func getTenantUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getTenant")
	var tenantID string
	var err error
	cacheKey := "getTenant"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		tenantID = cachedData.(string)
	} else {
		// Read tenant ID from config, or environment variables
		microsoft365Config := GetConfig(d.Connection)
		if microsoft365Config.TenantID != nil {
			tenantID = *microsoft365Config.TenantID
		} else if os.Getenv("AZURE_TENANT_ID") != "" {
			tenantID = os.Getenv("AZURE_TENANT_ID")
		}

		// If not set in config, get tenant ID from CLI
		if tenantID == "" {
			tenantID, err = getTenantFromCLI()
			if err != nil {
				return nil, err
			}
		}
		// save to extension cache
		d.ConnectionManager.Cache.Set(cacheKey, tenantID)
	}

	return tenantID, nil
}

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}
