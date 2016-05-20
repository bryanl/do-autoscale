package autoscale

import "golang.org/x/net/context"

// ResourceManager is a watched resource interface.
type ResourceManager interface {
	Count() (int, error)
	Scale(ctx context.Context, g Group, byN int, repo Repository) error
	Allocated() ([]ResourceAllocation, error)
}
