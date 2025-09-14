package infrastructure

import (
	"fmt"
	"sync"

	"github.com/arafato/cf-nuke/types"
)

var listers = make(map[string]types.ResourceLister)

func Register(name string, lister types.ResourceLister) error {
	var mu = sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	if _, exists := listers[name]; exists {
		return fmt.Errorf("handler %s already registered", name)
	}
	listers[name] = lister
	return nil
}
