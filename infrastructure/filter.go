package infrastructure

import (
	"github.com/arafato/cf-nuke/config"
	"github.com/arafato/cf-nuke/types"
)

func FilterCollection(resources types.Resources, config *config.Config) {
	resourceFilter := config.ResourceTypes.Excludes
	filterSet := make(map[string]struct{}, len(resourceFilter))
	for _, filter := range resourceFilter {
		filterSet[filter] = struct{}{}
	}
	for _, resource := range resources {
		if _, ok := filterSet[resource.ProductName]; ok {
			resource.State = types.Filtered
			continue
		}
		resource.State = types.Ready
	}
}
