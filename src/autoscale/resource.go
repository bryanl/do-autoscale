package autoscale

import "autoscale/metrics"

// ResourceManager is a watched resource interface.
type ResourceManager interface {
	Count() (int, error)
	Scale(g Group, byN int, repo Repository) error
	Allocated() ([]metrics.ResourceAllocation, error)
}
