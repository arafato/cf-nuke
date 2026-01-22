package utils

import "strings"

// IsPermissionError checks if an error is a permission/access error (401, 403)
func IsPermissionError(err error) bool {
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

// IsNotFoundError checks if an error is a 404/not found error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "404") ||
		strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "Not Found")
}

// IsFeatureNotAvailable checks if an error indicates a feature is not available for the account
func IsFeatureNotAvailable(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "not enabled") ||
		strings.Contains(errStr, "not available") ||
		strings.Contains(errStr, "feature flag") ||
		strings.Contains(errStr, "subscription")
}

// IsSkippableError returns true if the error should be skipped (logged as warning) rather than fatal
func IsSkippableError(err error) bool {
	return IsPermissionError(err) || IsNotFoundError(err) || IsFeatureNotAvailable(err)
}
