package utils

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"

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
	data := [][]string{{"Product", "ID/Name", "Status"}}
	for _, resource := range resources {
		if resource.State == types.Hidden {
			continue
		}
		data = append(data, []string{resource.ProductName, resource.ResourceName, resource.State.String()})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	table.Bulk(data[1:])
	table.Render()

	visibleCount := resources.VisibleCount()
	fmt.Printf("\nStatus: %d resources in total. Removed %d, In-Progress %d, Filtered %d, Failed %d\n", visibleCount, resources.NumOf(types.Deleted), resources.NumOf(types.Removing), resources.NumOf(types.Filtered), resources.NumOf(types.Failed))
}
