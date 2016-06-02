package autoscale

import (
	"pkg/ctxutil"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	groupCheckTimeout = 5 * time.Second
)

// AfterMonitorFn is a function called after a group is monitored.
type AfterMonitorFn func(groupName string)

// GroupMonitor monitors groups.
type GroupMonitor interface {
	Start(fn AfterMonitorFn) chan struct{}
	Stop()
	InRunList(groupName string) bool
	SetRunList(runList RunList)
}

// GroupMonitorOption is an option for configuring a GroupMonitor instance.
type GroupMonitorOption func(GroupMonitor) error

// NewGroupMonitor creates a new instance of GroupMonitor.
func NewGroupMonitor(ctx context.Context, repo Repository, opts ...GroupMonitorOption) (GroupMonitor, error) {

	gm := &groupMonitor{
		ctx:  ctx,
		repo: repo,
	}

	for _, opt := range opts {
		if err := opt(gm); err != nil {
			return nil, err
		}
	}

	if gm.runList == nil {
		gm.runList = NewRunList(ctx)
	}

	return gm, nil
}

// OptionRunList sets a custom RunList for the GroupMonitor.
func OptionRunList(runList RunList) GroupMonitorOption {
	return func(gm GroupMonitor) error {
		gm.SetRunList(runList)
		return nil
	}
}

type groupMonitor struct {
	ctx     context.Context
	repo    Repository
	runList RunList

	quit chan struct{}
}

var _ GroupMonitor = (*groupMonitor)(nil)

func (gm *groupMonitor) log() *logrus.Entry {
	return ctxutil.LogFromContext(gm.ctx).WithField("action", "group-monitor")
}

func (gm *groupMonitor) Start(fn AfterMonitorFn) chan struct{} {
	gm.log().Info("starting group monitoring")

	gm.quit = make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	timer := time.NewTimer(groupCheckTimeout)

	go func() {
		for {
			select {
			case <-timer.C:
				groups, err := gm.repo.ListGroups(gm.ctx)
				if err != nil {
					gm.log().WithError(err).Error("could not retrieve groups")
				}

				for _, g := range groups {
					if !gm.InRunList(g.Name) {
						if err := gm.runList.Add(g.Name); err == nil {
							fn(g.Name)
						}
					}
				}

				timer.Reset(groupCheckTimeout)
			case <-gm.quit:
				gm.quit = nil
				gm.runList.Reset()
				done <- struct{}{}
				break
			}
		}
	}()

	return done
}

func (gm *groupMonitor) Stop() {
	gm.log().Info("stopping group monitoring")
	gm.quit <- struct{}{}

}

func (gm *groupMonitor) InRunList(groupName string) bool {
	return gm.runList.IsRunning(groupName)
}

func (gm *groupMonitor) SetRunList(runList RunList) {
	gm.runList = runList
}
