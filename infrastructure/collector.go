package infrastructure

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/arafato/cf-nuke/types"
	"github.com/arafato/cf-nuke/utils"
)

var (
	accountCollectors = make(map[string]types.AccountCollector)
	zoneCollectors    = make(map[string]types.ZoneCollector)
)

// RegisterAccountCollector registers an account-level resource collector
func RegisterAccountCollector(name string, collector types.AccountCollector) {
	if _, exists := accountCollectors[name]; exists {
		panic(fmt.Errorf("account collector %s already registered", name))
	}
	accountCollectors[name] = collector
}

// RegisterZoneCollector registers a zone-level resource collector
func RegisterZoneCollector(name string, collector types.ZoneCollector) {
	if _, exists := zoneCollectors[name]; exists {
		panic(fmt.Errorf("zone collector %s already registered", name))
	}
	zoneCollectors[name] = collector
}

// isPermissionError checks if an error is a permission/access error (401, 403)
func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "403") ||
		strings.Contains(errStr, "401") ||
		strings.Contains(errStr, "Forbidden") ||
		strings.Contains(errStr, "Unauthorized") ||
		strings.Contains(errStr, "not have permission") ||
		strings.Contains(errStr, "Access denied")
}

// isNotFoundError checks if an error is a 404/not found error
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "404") ||
		strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "Not Found")
}

// isFeatureNotAvailable checks if an error indicates a feature is not available for the account
func isFeatureNotAvailable(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "not enabled") ||
		strings.Contains(errStr, "not available") ||
		strings.Contains(errStr, "feature flag") ||
		strings.Contains(errStr, "subscription")
}

// isSkippableError returns true if the error should be skipped (logged as warning) rather than fatal
func isSkippableError(err error) bool {
	return isPermissionError(err) || isNotFoundError(err) || isFeatureNotAvailable(err)
}

// isTransientError returns true if the error is a transient network error that can be retried
func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return err == io.EOF ||
		strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "timeout")
}

// isThrottlingError returns true if the error is a rate limiting error
func isThrottlingError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "too many requests")
}

// collectTask represents a single collection task (account-level or zone-level)
type collectTask struct {
	name    string                          // collector name
	context string                          // zone name (empty for account-level)
	collect func() (types.Resources, error) // the collection function
}

// formatTaskName returns a formatted name for logging (includes zone if present)
func (t *collectTask) formatTaskName() string {
	if t.context != "" {
		return fmt.Sprintf("%s (%s)", t.name, t.context)
	}
	return t.name
}

// collectWithRetry attempts to collect resources with retries for transient errors
func collectWithRetry(collect func() (types.Resources, error), maxRetries int) (types.Resources, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resources, err := collect()
		if err == nil {
			return resources, nil
		}
		lastErr = err

		// Don't retry non-transient errors
		if !isTransientError(err) {
			return nil, err
		}

		// Wait before retry (exponential backoff: 1s, 2s, 4s)
		if attempt < maxRetries {
			backoff := time.Duration(1<<attempt) * time.Second
			time.Sleep(backoff)
		}
	}
	return nil, lastErr
}

// ProcessCollection collects resources from all registered collectors
func ProcessCollection(creds *types.Credentials, zones []*types.Zone) types.Resources {
	logger := utils.GetLogger()
	resourceCollectionChan := make(chan *types.Resource, 100)
	var allResources types.Resources
	g := new(errgroup.Group)
	// Limit concurrent API calls to avoid rate limiting
	g.SetLimit(20)

	// Build task list
	var tasks []collectTask

	// Account-level tasks
	for name, collector := range accountCollectors {
		c := collector
		tasks = append(tasks, collectTask{
			name:    name,
			context: "",
			collect: func() (types.Resources, error) { return c(creds) },
		})
	}

	// Zone-level tasks (one per collector per zone)
	for name, collector := range zoneCollectors {
		for _, zone := range zones {
			c, z := collector, zone
			tasks = append(tasks, collectTask{
				name:    name,
				context: z.Name,
				collect: func() (types.Resources, error) { return c(creds, z) },
			})
		}
	}

	// Process all tasks with unified error handling
	for _, task := range tasks {
		t := task
		g.Go(func() error {
			resources, err := collectWithRetry(t.collect, 3)
			if err != nil {
				if isThrottlingError(err) {
					logger.LogWarning("%s: rate limited - %v", t.formatTaskName(), err)
				} else {
					logger.LogWarning("%s: %v", t.formatTaskName(), err)
				}
				return nil
			}
			for _, resource := range resources {
				resourceCollectionChan <- resource
			}
			return nil
		})
	}

	go func() {
		g.Wait()
		close(resourceCollectionChan)
	}()

	for resource := range resourceCollectionChan {
		allResources = append(allResources, resource)
	}

	return allResources
}

// ListCollectors returns an alphabetically sorted list of all registered collector names.
func ListCollectors() []string {
	var collectorNames []string
	for name := range accountCollectors {
		collectorNames = append(collectorNames, name)
	}
	for name := range zoneCollectors {
		collectorNames = append(collectorNames, name)
	}
	slices.Sort(collectorNames)
	return collectorNames
}
