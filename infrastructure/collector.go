package infrastructure

import (
	"fmt"
	"os"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/arafato/cf-nuke/types"
)

var collectors = make(map[string]types.ResourceCollector)

func RegisterCollector(name string, collector types.ResourceCollector) {
	if _, exists := collectors[name]; exists {
		panic(fmt.Errorf("handler %s already registered", name))
	}
	collectors[name] = collector
}

// ListCollector returns an alphabetically sorted list of registered collector names.
func ListCollectors() []string {
	var collectorNames []string
	for name, _ := range collectors {
		collectorNames = append(collectorNames, name)
	}
	slices.Sort(collectorNames)
	return collectorNames
}

func ProcessCollection(creds *types.Credentials) types.Resources {
	var resourceCollectionChan = make(chan *types.Resource, 100)
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

	var collectedErr error
	go func() {
		collectedErr = g.Wait()
		close(resourceCollectionChan)
	}()

	for resource := range resourceCollectionChan {
		allResources = append(allResources, resource)
	}

	if collectedErr != nil {
		fmt.Println("Error during collection, aborting:\n", collectedErr)
		os.Exit(1)
	}

	return allResources
}
