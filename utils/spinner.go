package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

// SpinWhile runs a function while displaying a spinner with the given message
// Returns the result of the function and the duration it took
func SpinWhile[T any](message string, fn func() T) (T, time.Duration) {
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Start()
	start := time.Now()
	result := fn()
	duration := time.Since(start)
	s.Stop()
	return result, duration
}

// SpinWhileWithError runs a function while displaying a spinner
// Returns the result and error of the function
func SpinWhileWithError[T any](message string, fn func() (T, error)) (T, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Start()
	result, err := fn()
	s.Stop()
	return result, err
}
