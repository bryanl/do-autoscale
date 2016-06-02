package autoscale

import (
	"errors"
	"fmt"
	"pkg/ctxutil"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	ErrActionTimedOut = fmt.Errorf("action timed out")
	ErrDisabledGroup  = fmt.Errorf("group is disabled")

	SchedulerActionTimeout = 60 * time.Minute
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
		ctx:  ctx,
		repo: repo,
		quit: make(chan bool),
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

type GroupActionFn func(groupID string) *ActionStatus

type SchedulerActivity struct {
	ID  string
	Err error
}

type SchedulerStatus struct {
	EnableGroup  chan string
	DisableGroup chan string
	Schedule     chan string
	Activity     chan SchedulerActivity
}

type Scheduler struct {
	ctx              context.Context
	enableGroupChan  chan string
	disableGroupChan chan string
	scheduleChan     chan string
	activityChan     chan SchedulerActivity
	actionFn         GroupActionFn
	disabledIDs      map[string]bool
}

func NewScheduler(ctx context.Context, fn GroupActionFn) *Scheduler {
	return &Scheduler{
		ctx:              ctx,
		enableGroupChan:  make(chan string),
		disableGroupChan: make(chan string),
		scheduleChan:     make(chan string),
		activityChan:     make(chan SchedulerActivity, 1000),
		actionFn:         fn,
		disabledIDs:      map[string]bool{},
	}
}

func (s *Scheduler) Status() *SchedulerStatus {
	return &SchedulerStatus{
		EnableGroup:  s.enableGroupChan,
		DisableGroup: s.disableGroupChan,
		Schedule:     s.scheduleChan,
		Activity:     s.activityChan,
	}
}

type ActionStatus struct {
	Done chan bool
	Err  error
}

func (s *Scheduler) Start() {
	for {
		select {
		case id := <-s.scheduleChan:
			s.log().WithField("group-id", id).Info("scheduling group")

			if _, ok := s.disabledIDs[id]; ok {
				s.log().WithField("group-id", id).Warn("will not schedule group as it is disabled")
				s.activityChan <- SchedulerActivity{
					ID:  id,
					Err: ErrDisabledGroup,
				}
				continue
			}
			go func() {
				actionStatus := s.actionFn(id)
				err := handleActionStatus(s.ctx, actionStatus)
				if err != nil {
					s.log().WithError(err).Error("action did not run with success")
					s.disableGroup(id)
				}

				s.activityChan <- SchedulerActivity{
					ID:  id,
					Err: err,
				}
			}()

		case id := <-s.enableGroupChan:
			s.log().WithField("group-id", id).Info("enabling group")
			delete(s.disabledIDs, id)

		case id := <-s.disableGroupChan:
			s.log().WithField("group-id", id).Info("disabling group")
			s.disableGroup(id)
		}
	}
}

func (s *Scheduler) disableGroup(id string) {
	s.disabledIDs[id] = true
}

func (s *Scheduler) log() *logrus.Entry {
	return ctxutil.LogFromContext(s.ctx).WithField("action", "schedule")
}

func handleActionStatus(ctx context.Context, as *ActionStatus) error {
	log := ctxutil.LogFromContext(ctx).WithField("action", "action-handler")
	log.Info("handle action status")

	timer := time.NewTimer(SchedulerActionTimeout)
	select {
	case <-timer.C:
		log.Warn("action timed out")
		return ErrActionTimedOut
	case <-as.Done:
		return as.Err
	}

}
