package infrastructure

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/arafato/cf-nuke/types"
)

var collectors = make(map[string]types.ResourceCollector)

var resourceCollectionChan = make(chan *types.Resource, 100)

func RegisterCollector(name string, collector types.ResourceCollector) error {
	var mu = sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	if _, exists := collectors[name]; exists {
		return fmt.Errorf("handler %s already registered", name)
	}
	collectors[name] = collector
	return nil
}

func ProcessCollection(creds *types.Credentials) types.Resources {
	var allResources types.Resources
	g := new(errgroup.Group)

	for _, collector := range collectors {
		c := collector
		g.Go(func() error {
			resources, err := c(creds)
			if err != nil {
				return err
			}
			for _, resource := range resources {
				resourceCollectionChan <- resource
			}
			return nil
		})
	}

	go func() {
		_ = g.Wait()
		close(resourceCollectionChan)
	}()

	if err := g.Wait(); err != nil {
		fmt.Println("Error during collection, aborting:\n", err)
		os.Exit(1)
	}

	for resource := range resourceCollectionChan {
		allResources = append(allResources, resource)
	}

	return allResources
}
