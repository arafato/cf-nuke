package infrastructure

import (
	"context"

	"github.com/arafato/cf-nuke/types"
	"golang.org/x/sync/errgroup"
)

func RemoveCollection(ctx context.Context, resources types.Resources) error {
	g, _ := errgroup.WithContext(ctx)

	for _, resource := range resources {
		resource := resource
		if resource.State() == types.Filtered || resource.State() == types.Hidden {
			continue
		}
		g.Go(func() error {
			return resource.Remove()
		})
	}

	return g.Wait()
}
