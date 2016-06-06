package autoscale

import (
	"pkg/ctxutil"
	"time"

	"github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
)

// Check checks the statucs of a group.
type Check struct {
	repo Repository
}

// NewCheck creates an instance of Check.
func NewCheck(repo Repository) *Check {
	return &Check{
		repo: repo,
	}
}

// Scale cales the group identified by gruopID.
func (c *Check) Scale(ctx context.Context, groupID string) *ActionStatus {
	log := ctxutil.LogFromContext(ctx).WithField("group-id", groupID)

	as := &ActionStatus{
		Done: make(chan bool, 1),
	}

	defer func() {
		as.Done <- true
	}()

	group, err := c.repo.GetGroup(ctx, groupID)
	if err != nil {
		as.Err = err
		return as
	}

	resource, err := group.Resource()
	if err != nil {
		as.Err = err
		return as
	}

	value, err := group.MetricsValue(ctx)
	if err != nil {
		as.Err = err
		return as
	}

	policy := group.Policy
	count, err := resource.Count()
	if err != nil {
		as.Err = err
		return as
	}

	newCount := policy.CalculateSize(count, value)

	delta := newCount - count

	changed, err := resource.Scale(ctx, *group, delta, c.repo)
	if err != nil {
		as.Err = err
		return as
	}

	if changed {
		log.WithFields(logrus.Fields{
			"metric":       group.MetricType,
			"metric-value": value,
			"new-count":    newCount,
			"delta":        delta,
		}).Info("group change status")

		if err := group.MetricNotify(); err != nil {
			log.WithError(err).Error("notifying metric of current config")
			as.Err = err
			return as
		}

		wup := policy.WarmUpPeriod()
		log.WithField("warm-up-duration", wup).Info("waiting for new service to warm up")
		time.Sleep(wup)
		log.Info("new service has warmed up")
	}

	as.Delta = delta
	as.Count = newCount
	return as
}

// Disable the group identified by groupID.
func (c *Check) Disable(ctx context.Context, groupID string) *ActionStatus {
	log := ctxutil.LogFromContext(ctx).WithField("group-id", groupID)

	as := &ActionStatus{
		Done: make(chan bool, 1),
	}

	defer func() {
		as.Done <- true
	}()

	group, err := c.repo.GetGroup(ctx, groupID)
	if err != nil {
		as.Err = err
		return as
	}

	resource, err := group.Resource()
	if err != nil {
		as.Err = err
		return as
	}

	count, err := resource.Count()
	if err != nil {
		as.Err = err
		return as
	}

	log.Info("disabling group by removing all resources")

	if err := group.Disable(ctx); err != nil {
		as.Err = err
		return as
	}

	toRemove := 0 - count
	_, err = resource.Scale(ctx, *group, toRemove, c.repo)
	if err != nil {
		as.Err = err
		return as
	}

	as.Delta = toRemove
	as.Count = 0

	return as
}
