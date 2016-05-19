package autoscale

// ResourceManager is a watched resource interface.
type ResourceManager interface {
	Count() (int, error)
	Scale(g Group, byN int, repo Repository) error
	Allocated() ([]ResourceAllocation, error)
}
