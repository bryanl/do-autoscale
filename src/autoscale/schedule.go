package autoscale

import (
	"pkg/ctxutil"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type GroupActionFn func(ctx context.Context, groupID string) *ActionStatus

type GroupAction interface {
	Scale(ctx context.Context, groupID string) *ActionStatus
	Disable(ctx context.Context, groupID string) *ActionStatus
}

type SchedulerActivity struct {
	ID    string
	Err   error
	Delta int
	Count int
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
	groupAction      GroupAction
	disabledIDs      map[string]bool
}

func NewScheduler(ctx context.Context, ga GroupAction) *Scheduler {
	return &Scheduler{
		ctx:              ctx,
		enableGroupChan:  make(chan string, 1),
		disableGroupChan: make(chan string, 1),
		scheduleChan:     make(chan string, 1),
		activityChan:     make(chan SchedulerActivity, 1),
		groupAction:      ga,
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
	Count int
}

func (s *Scheduler) Start() {
	for {
		select {
		case id := <-s.scheduleChan:
			s.log().WithField("group-id", id).Debug("scheduling group")

			if _, ok := s.disabledIDs[id]; ok {
				s.log().WithField("group-id", id).Warn("will not schedule group as it is disabled")
				s.activityChan <- SchedulerActivity{
					ID:  id,
					Err: ErrDisabledGroup,
				}
				continue
			}

			go func() {
				actionStatus := s.groupAction.Scale(s.ctx, id)
				err := handleActionStatus(s.ctx, actionStatus)
				if err != nil {
					s.log().WithError(err).Error("action did not run with success")
					s.disableGroup(id)
				}

				s.activityChan <- SchedulerActivity{
					ID:    id,
					Err:   err,
					Delta: actionStatus.Delta,
					Count: actionStatus.Count,
				}

				time.AfterFunc(ScheduleReenqueueTimeout, func() {
					s.scheduleChan <- id
				})
			}()

		case id := <-s.enableGroupChan:
			s.log().WithField("group-id", id).Info("enabling group")
			delete(s.disabledIDs, id)

		case id := <-s.disableGroupChan:
			s.log().WithField("group-id", id).Info("disabling group")

			if !s.disabledIDs[id] {
				go func(id string) {
					actionStatus := s.groupAction.Disable(s.ctx, id)
					err := handleActionStatus(s.ctx, actionStatus)
					if err != nil {
						s.log().WithError(err).Error("action did not run with success")
					}

					s.activityChan <- SchedulerActivity{
						ID:    id,
						Err:   err,
						Delta: actionStatus.Delta,
						Count: actionStatus.Count,
					}
				}(id)
			}

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

	timer := time.NewTimer(SchedulerActionTimeout)
	select {
	case <-timer.C:
		log.Warn("action timed out")
		return ErrActionTimedOut
	case <-as.Done:
		return as.Err
	}

}
