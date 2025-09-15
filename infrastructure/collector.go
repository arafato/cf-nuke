package infrastructure

import (
	"fmt"
	"sync"

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
	var wg sync.WaitGroup
	for _, collector := range collectors {
		wg.Add(1)
		go func(c types.ResourceCollector) {
			defer wg.Done()
			c(creds)
		}(collector)
	}

	go func() {
		wg.Wait()
		close(resourceCollectionChan)
	}()

	for resource := range resourceCollectionChan {
		allResources = append(allResources, resource)
	}

	return allResources
}

func CollectResource(resource *types.Resource) {
	resourceCollectionChan <- resource
}
