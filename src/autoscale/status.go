package autoscale

import (
	"pkg/ctxutil"
	"time"

	"golang.org/x/net/context"
)

// Status manages status updates for groups.
type Status struct {
	ctx              context.Context
	repo             Repository
	ActivityListener chan SchedulerActivity
}

// NewStatus creates an instance of Status.
func NewStatus(ctx context.Context, repo Repository) *Status {
	return &Status{
		ctx:              ctx,
		repo:             repo,
		ActivityListener: make(chan SchedulerActivity),
	}
}

func (s *Status) Start() {
	for msg := range s.ActivityListener {
		if msg.Err == nil && msg.Delta != 0 {
			gs := GroupStatus{
				GroupID:   msg.ID,
				Delta:     msg.Delta,
				Total:     msg.Count,
				CreatedAt: time.Now(),
			}

			if err := s.repo.AddGroupStatus(s.ctx, gs); err != nil {
				log := ctxutil.LogFromContext(s.ctx)
				log.WithError(err).Error("unable to add group status")
			}
		}
	}
}
