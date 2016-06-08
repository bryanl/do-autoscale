package autoscale

import (
	"pkg/ctxutil"

	"github.com/satori/go.uuid"

	"golang.org/x/net/context"
)

// Notification is a notification message from the scheduler.
type Notification struct {
	ID      string `json:"id"`
	GroupID string `json:"groupID"`
	Name    string `json:"name"`
	Action  string `json:"action"`
	Delta   int    `json:"delta"`
	Count   int    `json:"count"`
	Message string `json:"message"`
	IsError bool   `json:"isError"`
}

// Notify listens to the scheduler to generate Notification.
type Notify struct {
	ctx                  context.Context
	repo                 Repository
	ActivityListener     chan SchedulerActivity
	NotificationListener chan Notification
}

// NewNotify creates an instance of Notify.
func NewNotify(ctx context.Context, repo Repository) *Notify {
	return &Notify{
		ctx:                  ctx,
		repo:                 repo,
		ActivityListener:     make(chan SchedulerActivity),
		NotificationListener: make(chan Notification, 100),
	}
}

// Start starts the listener.
func (n *Notify) Start() {
	for msg := range n.ActivityListener {
		if msg.Delta == 0 {
			continue
		}

		log := ctxutil.LogFromContext(n.ctx).WithField("action", "notify")
		log.Info("sending notification to websocket clients")

		g, err := n.repo.GetGroup(n.ctx, msg.ID)
		if err != nil {
			log.WithError(err).Error("unable to load group")
			continue
		}

		notif := Notification{
			ID:      uuid.NewV4().String(),
			GroupID: msg.ID,
			Name:    g.Name,
		}
		if msg.Err != nil {
			notif.Message = msg.Err.Error()
			notif.IsError = true
		} else {
			notif.Delta = msg.Delta
			notif.Count = msg.Count
		}

		n.NotificationListener <- notif
	}
}
