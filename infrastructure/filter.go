package infrastructure

import (
	"slices"

	"github.com/arafato/cf-nuke/config"
	"github.com/arafato/cf-nuke/types"
)

func FilterCollection(resources types.Resources, config *config.Config) {
	resourceTypeFilter := config.ResourceTypes.Excludes
	resourceTypeFilterSet := make(map[string]struct{}, len(resourceTypeFilter))
	for _, filter := range resourceTypeFilter {
		resourceTypeFilterSet[filter] = struct{}{}
	}

	resourceIDFilters := config.ResourceIDs.Excludes
	resourceIDLookup := make(map[string]map[string]struct{})
	for _, resourceIDFilter := range resourceIDFilters {
		if _, ok := resourceIDLookup[resourceIDFilter.ResourceType]; !ok {
			resourceIDLookup[resourceIDFilter.ResourceType] = make(map[string]struct{})
		}
		resourceIDLookup[resourceIDFilter.ResourceType][resourceIDFilter.ID] = struct{}{}
	}

	for _, resource := range resources {
		if resource.State() == types.Hidden {
			continue
		}
		if _, ok := resourceTypeFilterSet[resource.ProductName]; ok {
			resource.SetState(types.Filtered)
			continue
		}

		if idSet, ok := resourceIDLookup[resource.ProductName]; ok {
			if _, ok := idSet[resource.ResourceID]; ok {
				resource.SetState(types.Filtered)
				continue
			}
			if _, ok := idSet[resource.ResourceName]; ok {
				resource.SetState(types.Filtered)
				continue
			}
		}

		if slices.Contains(config.Zones.Excludes, resource.ResourceName) {
			resource.SetState(types.Filtered)
			continue
		}

		resource.SetState(types.Ready)
	}
}
