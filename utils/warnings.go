package utils

import (
	"fmt"
	"sync"
)

// Warning represents a non-fatal issue encountered during resource collection
type Warning struct {
	ResourceType string
	Context      string // e.g., zone name or account feature
	Message      string
}

var (
	warnings   []Warning
	warningsMu sync.Mutex
)

// AddWarning records a warning for later display and logs it to the scan logger
func AddWarning(resourceType, context, message string) {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	warnings = append(warnings, Warning{
		ResourceType: resourceType,
		Context:      context,
		Message:      message,
	})

	// Also log to the scan logger for file output
	logger := GetLogger()
	if context != "" {
		logger.LogWarning("%s (%s): %s", resourceType, context, message)
	} else {
		logger.LogWarning("%s: %s", resourceType, message)
	}
}

// ClearWarnings resets the warning list
func ClearWarnings() {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	warnings = nil
}

// PrintWarnings outputs all collected warnings to stdout
func PrintWarnings() {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	if len(warnings) == 0 {
		return
	}
	fmt.Printf("\n[WARNINGS] %d issue(s) encountered during collection:\n", len(warnings))
	for _, w := range warnings {
		if w.Context != "" {
			fmt.Printf("  - %s (%s): %s\n", w.ResourceType, w.Context, w.Message)
		} else {
			fmt.Printf("  - %s: %s\n", w.ResourceType, w.Message)
		}
	}
	fmt.Println()
}

// WarningCount returns the number of warnings collected
func WarningCount() int {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	return len(warnings)
}
