package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arafato/cf-nuke/types"
)

func PrintStatusWithContext(wg *sync.WaitGroup, ctx context.Context, resources types.Resources) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			PrettyPrintStatus(resources)
		case <-ctx.Done():
			PrettyPrintStatus(resources)
			return
		}
	}
}

func PrettyPrintStatus(resources types.Resources) {
	fmt.Print("\n====================================\n")
	for _, resource := range resources {
		fmt.Printf("%s - \033[32m%s\033[0m - %s - %s\n", resource.ProductName, resource.ResourceName, resource.ResourceID, resource.State)
	}
	fmt.Printf("\nStatus: %d resources in total.Removed %d, In-Progress %d, Filtered %d\n", len(resources), resources.NumOf(types.Deleted), resources.NumOf(types.Removing), resources.NumOf(types.Filtered))
}
