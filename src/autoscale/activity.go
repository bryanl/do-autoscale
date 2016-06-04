package autoscale

import "fmt"

// ActivityListener is an entity interested in listening for SchedulerActivity messages.
type ActivityListener chan<- SchedulerActivity

// ActivityManager collects events from Schedule, and fans them out to multiple listerns.
type ActivityManager struct {
	activityChan chan SchedulerActivity
	listeners    []ActivityListener
}

// NewActivityManager creates an instance of ActivityManager.
func NewActivityManager(activityChan chan SchedulerActivity) *ActivityManager {
	return &ActivityManager{
		activityChan: activityChan,
		listeners:    []ActivityListener{},
	}
}

// RegisterListener registers ActivityListeners.
func (a *ActivityManager) RegisterListener(fn ActivityListener) {
	a.listeners = append(a.listeners, fn)
}

// Start starts the fanout process.
func (a *ActivityManager) Start() {
	for {
		select {
		case msg := <-a.activityChan:
			newListeners := []ActivityListener{}

			for _, ch := range a.listeners {
				if err := a.send(msg, ch); err == nil {
					newListeners = append(newListeners, ch)
				}
			}

			a.listeners = newListeners
		}
	}
}

func (a *ActivityManager) send(msg SchedulerActivity, l ActivityListener) error {
	var err error

	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("send to closed listener: %v", x)
		}
	}()

	l <- msg
	return err
}
