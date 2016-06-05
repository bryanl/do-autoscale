package autoscale

import (
	"pkg/ctxutil"
	"time"

	"github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
)

type Check struct {
	repo Repository
}

func NewCheck(repo Repository) *Check {
	return &Check{
		repo: repo,
	}
}

func (c *Check) Perform(ctx context.Context, groupID string) *ActionStatus {
	log := ctxutil.LogFromContext(ctx)

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

	value, err := group.MetricsValue()
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
	}

	as.Delta = delta
	as.Count = newCount
	return as
}
