package types

import (
	"context"
	"errors"
	"strings"

	"github.com/cenkalti/backoff/v5"
)

type Removable interface {
	Remove(accountID string, resourceID string) error
}

type Resource struct {
	Removable
	AccountID    string
	ResourceID   string
	ResourceName string
	ProductName  string
	State        ResourceState
}

type ResourceCollector func(*Credentials) (Resources, error)

type Resources []*Resource

//go:generate stringer -type=ResourceState
type ResourceState int

const (
	Removing ResourceState = iota
	Ready
	Deleted
	Failed
	Filtered
)

func (r *Resource) Remove() error {
	operation := func() (struct{}, error) {
		r.State = Removing
		err := r.Removable.Remove(r.AccountID, r.ResourceID)
		if err != nil {
			if strings.Contains(err.Error(), "401 Unauthorized") {
				return struct{}{}, backoff.Permanent(errors.New("Unauthorized Request"))
			}
			return struct{}{}, err
		}
		return struct{}{}, nil
	}

	_, err := backoff.Retry(context.TODO(), operation, backoff.WithBackOff(backoff.NewExponentialBackOff()), backoff.WithMaxTries(3))
	if err != nil {
		r.State = Failed
		return err
	}

	r.State = Deleted
	return nil
}

func (r Resources) NumOf(state ResourceState) int {
	count := 0
	for _, resource := range r {
		if resource.State == state {
			count++
		}
	}
	return count
}
