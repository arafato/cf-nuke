package infrastructure

import (
	"fmt"
	"sync"

	"github.com/arafato/cf-nuke/resources"
)

var listers = make(map[string]resources.Resource)

func Register(name string, resource resources.Resource) error {
	var mu = sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	if _, exists := listers[name]; exists {
		return fmt.Errorf("handler %s already registered", name)
	}
	listers[name] = resource
	return nil
}
