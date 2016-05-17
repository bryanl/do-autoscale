package autoscale

import (
	"autoscale/metrics"
	"fmt"
	"math/rand"
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

	id := rand.Int()

	log = log.WithField("resource-type", "local").WithField("id", id)

	return &LocalResource{
		log: log,
	}
}

// Actual is the actual number of resources available.
func (r *LocalResource) Actual() (int, error) {

	return r.count, nil
}

// ScaleUp scales resources up.
func (r *LocalResource) ScaleUp(g Group, byN int, repo Repository) error {
	r.log.WithField("by-n", byN).Info("scaling up")

	r.count = r.count + byN
	return nil
}

// ScaleDown scales resources down.
func (r *LocalResource) ScaleDown(g Group, byN int, repo Repository) error {
	r.log.WithField("by-n", byN).Info("scaling down")

	r.count = r.count - byN
	return nil
}

// Allocated returns a slice of ResourceAllocation for this resource.
func (r *LocalResource) Allocated() ([]metrics.ResourceAllocation, error) {
	allocations := []metrics.ResourceAllocation{}
	for i := 0; i < r.count; i++ {
		allocation := metrics.ResourceAllocation{
			Name: fmt.Sprintf("instance-%d", i+1),
		}
		allocations = append(allocations, allocation)
	}

	return allocations, nil
}
