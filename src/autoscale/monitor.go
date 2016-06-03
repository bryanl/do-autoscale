package autoscale

import (
	"errors"
	"pkg/ctxutil"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type Monitor interface {
	Start(*SchedulerStatus)
	Stop()
}

type monitor struct {
	ctx               context.Context
	repo              Repository
	runList           RunList
	groupCheckTimeout time.Duration
	quit              chan bool
}

var _ Monitor = (*monitor)(nil)

type MonitorOption func(Monitor) error

func NewMonitor(ctx context.Context, repo Repository, opts ...MonitorOption) (Monitor, error) {
	m := &monitor{
		ctx:               ctx,
		repo:              repo,
		groupCheckTimeout: DefaultGroupCheckTimeout,
		quit:              make(chan bool),
	}

	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}

	if m.runList == nil {
		m.runList = NewRunList(ctx)
	}

	return m, nil
}

func MonitorRunList(runList RunList) MonitorOption {
	return func(m Monitor) error {
		myMonitor, ok := m.(*monitor)
		if !ok {
			return errors.New("could not set run list")
		}

		myMonitor.runList = runList
		return nil
	}
}

func MonitorGroupCheckTimeout(timeout time.Duration) MonitorOption {
	return func(m Monitor) error {
		myMonitor, ok := m.(*monitor)
		if !ok {
			return errors.New("could not set group check timeout")
		}

		myMonitor.groupCheckTimeout = timeout
		return nil
	}
}

func (m *monitor) Start(schedulerStatus *SchedulerStatus) {
	log := m.log()
	log.Debug("starting monitor")

	timer := time.NewTimer(m.groupCheckTimeout)
	for {
		select {
		case <-timer.C:
			groups, err := m.repo.ListGroups(m.ctx)
			if err != nil {
				log.WithError(err).Error("could not retrieve groups")
				continue
			}

			// add groups from the datastore.
			groupIDs := []string{}
			for _, g := range groups {
				groupIDs = append(groupIDs, g.ID)
				if !m.runList.IsRunning(g.ID) {
					schedulerStatus.EnableGroup <- g.ID
					schedulerStatus.Schedule <- g.ID
					m.runList.Add(g.ID)
				}
			}

			// remove groups that no longer exist in the datastore
			for _, groupID := range m.runList.List() {
				if !stringInSlice(groupID, groupIDs) {
					m.runList.Remove(groupID)
					schedulerStatus.DisableGroup <- groupID
				}
			}

			timer.Reset(m.groupCheckTimeout)
		case <-m.quit:
			log.Debug("monitor stopped")
			return
		}
	}
}

func (m *monitor) Stop() {
	m.quit <- true
}

func (m *monitor) log() *logrus.Entry {
	return ctxutil.LogFromContext(m.ctx).WithField("action", "monitor")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
