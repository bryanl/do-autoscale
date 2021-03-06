package autoscale

import (
	"fmt"
	"pkg/ctxutil"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
)

// LocalResource is a local resource. Useful for testing on planes.
type LocalResource struct {
	count int
	log   *logrus.Entry
}

var _ ResourceManager = (*LocalResource)(nil)

// NewLocalResource builds an instance of LocalResource.
func NewLocalResource(ctx context.Context) ResourceManager {
	log := ctxutil.LogFromContext(ctx)
	if log == nil {
		l := logrus.New()
		log = logrus.NewEntry(l)
	}

	log = log.WithField("resource-type", "local")

	return &LocalResource{
		log: log,
	}
}

// Count is the actual number of resources available.
func (r *LocalResource) Count() (int, error) {
	return r.count, nil
}

// Scale scales in memory resources byN.
func (r *LocalResource) Scale(ctx context.Context, g Group, byN int, repo Repository) (bool, error) {
	if byN > 0 {
		return true, r.scaleUp(ctx, g, byN, repo)
	} else if byN < 0 {
		return false, r.scaleDown(ctx, g, 0-byN, repo)
	} else {
		return false, nil
	}
}

// ScaleUp scales resources up.
func (r *LocalResource) scaleUp(ctx context.Context, g Group, byN int, repo Repository) error {
	r.log.WithField("by-n", byN).Info("scaling up")

	r.count = r.count + byN
	return nil
}

// ScaleDown scales resources down.
func (r *LocalResource) scaleDown(ctx context.Context, g Group, byN int, repo Repository) error {
	r.log.WithField("by-n", byN).Info("scaling down")

	r.count = r.count - byN
	return nil
}

// Allocated returns a slice of ResourceAllocation for this resource.
func (r *LocalResource) Allocated() ([]ResourceAllocation, error) {
	allocations := []ResourceAllocation{}
	for i := 0; i < r.count; i++ {
		allocation := ResourceAllocation{
			Name: fmt.Sprintf("instance-%d", i+1),
		}
		allocations = append(allocations, allocation)
	}

	return allocations, nil
}
