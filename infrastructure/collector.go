package infrastructure

import (
	"fmt"
	"sync"

	"github.com/arafato/cf-nuke/types"
)

var collectors = make(map[string]types.ResourceCollector)

var resourceCollectionChan = make(chan *types.Resource, 100)

func RegisterCollector(name string, lister types.ResourceCollector) error {
	var mu = sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	if _, exists := collectors[name]; exists {
		return fmt.Errorf("handler %s already registered", name)
	}
	collectors[name] = lister
	return nil
}

func ProcessCollection() types.Resources {
	return nil
}

func CollectResource(resource *types.Resource) {
	resourceCollectionChan <- resource
}
