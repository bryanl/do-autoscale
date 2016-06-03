package autoscale

import (
	"pkg/ctxutil"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type GroupActionFn func(ctx context.Context, groupID string) *ActionStatus

type SchedulerActivity struct {
	ID    string
	Err   error
	Delta int
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
	Done  chan bool
	Err   error
	Delta int
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
				actionStatus := s.actionFn(s.ctx, id)
				err := handleActionStatus(s.ctx, actionStatus)
				if err != nil {
					s.log().WithError(err).Error("action did not run with success")
					s.disableGroup(id)
				}

				s.activityChan <- SchedulerActivity{
					ID:    id,
					Err:   err,
					Delta: actionStatus.Delta,
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