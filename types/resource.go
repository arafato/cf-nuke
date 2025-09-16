package types

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

type ResourceState int

const (
	ResourceStatePending = iota
	ResourceStateReady
	ResourceStateDeleted
	ResourceStateFailed
	ResourceStateFiltered
)
