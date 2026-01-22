package types

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"

	"github.com/cenkalti/backoff/v5"
)

type Removable interface {
	// Some resources require the ID others the name for a delete operation
	Remove(accountID string, resourceID string, resourceName string) error
}

type Resource struct {
	Removable
	AccountID    string
	ResourceID   string
	ResourceName string
	ProductName  string
	state        atomic.Int32 // use State() and SetState() for thread-safe access
	Error        atomic.Value // stores error message string for failed resources
}

type ResourceCollector func(*Credentials) (Resources, error)

type Resources []*Resource

//go:generate stringer -type=ResourceState
type ResourceState int32

const (
	// Ready is the default state (zero value) for new resources
	Ready ResourceState = iota
	Removing
	Deleted
	Failed
	Filtered
	Hidden
)

// State returns the current state of the resource (thread-safe)
func (r *Resource) State() ResourceState {
	return ResourceState(r.state.Load())
}

// SetState sets the state of the resource (thread-safe)
func (r *Resource) SetState(s ResourceState) {
	r.state.Store(int32(s))
}

// SetError stores an error message for the resource (thread-safe)
func (r *Resource) SetError(err error) {
	if err != nil {
		r.Error.Store(err.Error())
	}
}

// GetError returns the error message for the resource (thread-safe)
func (r *Resource) GetError() string {
	if v := r.Error.Load(); v != nil {
		return v.(string)
	}
	return ""
}

func (r *Resource) Remove(ctx context.Context) error {
	operation := func() (struct{}, error) {
		r.SetState(Removing)
		err := r.Removable.Remove(r.AccountID, r.ResourceID, r.ResourceName)
		if err != nil {
			if strings.Contains(err.Error(), "401 Unauthorized") {
				return struct{}{}, backoff.Permanent(errors.New("Unauthorized Request"))
			}
			return struct{}{}, err
		}
		return struct{}{}, nil
	}

	_, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()), backoff.WithMaxTries(3))
	if err != nil {
		r.SetState(Failed)
		r.SetError(err)
		return err
	}

	r.SetState(Deleted)
	return nil
}

func (r Resources) NumOf(state ResourceState) int {
	count := 0
	for _, resource := range r {
		if resource.State() == state {
			count++
		}
	}
	return count
}

func (r Resources) VisibleCount() int {
	return len(r) - r.NumOf(Hidden)
}
