package types

type Removable interface {
	Remove() error
}

type Resource struct {
	Removable   Removable
	ID          string
	ProductName string
}

type ResourceCollector func(*Credentials) error

type Resources []*Resource
