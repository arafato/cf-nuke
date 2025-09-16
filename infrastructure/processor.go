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
		g.Go(func() error {
			return resource.Remove()
		})
	}

	return g.Wait()
}
