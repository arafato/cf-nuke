package types

// Resource represents a cloud resource that can be removed
type Resource interface {
	Remove() error
}

// ResourceLister represents a function that can list resources of a specific type
type ResourceLister func() ([]Resource, error)
