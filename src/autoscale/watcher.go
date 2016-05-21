package autoscale

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
)

var (
	groupWatchDuration = 5 * time.Second
)

type watchedJob struct {
	name string
}

// Watcher watches groups.
type Watcher struct {
	repo       Repository
	log        *logrus.Entry
	groupNames []string
	workChan   chan watchedJob
	quitChan   chan int

	wg sync.Mutex
}

func makeJobQueue() chan watchedJob {
	return make(chan watchedJob, 1000)
}

// NewWatcher creates an instance of Watcher.
func NewWatcher(repo Repository) *Watcher {
	return &Watcher{
		repo:     repo,
		log:      logrus.WithField("action", "watcher"),
		workChan: makeJobQueue(),
	}
}

// AddGroup adds a group to the watch list..
func (w *Watcher) AddGroup(name string) error {
	w.wg.Lock()
	defer w.wg.Unlock()

	for _, n := range w.groupNames {
		if n == name {
			return fmt.Errorf("group %s is already being watched", name)
		}
	}

	w.groupNames = append(w.groupNames, name)

	w.log.WithField("group-name", name).Info("adding group")
	w.queueJob(name)

	return nil
}

func (w *Watcher) queueJob(name string) {
	w.log.WithField("name", name).Info("queueing job")
	job := watchedJob{
		name: name,
	}

	w.workChan <- job
}

// RemoveGroup removes a group from the watch list.
func (w *Watcher) RemoveGroup(name string) {
	w.wg.Lock()
	defer w.wg.Unlock()

	for i, groupName := range w.groupNames {
		if name == groupName {
			w.log.WithField("group-name", name).Info("removing group")
			w.groupNames = append(w.groupNames[:i], w.groupNames[i+1:]...)
			break
		}
	}
}

// Groups are the currently watched groups.
func (w *Watcher) Groups() []string {
	return w.groupNames
}

// Watch starts the watching process.
func (w *Watcher) Watch(ctx context.Context) (chan bool, error) {
	w.wg.Lock()
	defer w.wg.Unlock()

	log := w.log

	if w.quitChan != nil {
		log.Warn("watcher is already running")
		return nil, fmt.Errorf("watcher is already running")
	}

	log.Info("starting watcher")

	done := make(chan bool, 1)
	w.quitChan = make(chan int, 1)

	groupWatchTicker := time.NewTicker(groupWatchDuration)

	if w.workChan == nil {
		w.workChan = makeJobQueue()
	}

	go func() {
		for _, name := range w.groupNames {
			w.queueJob(name)
		}

		for {
			select {
			case job := <-w.workChan:
				log := log.WithField("group-name", job.name)
				log.Info("watching group")

				g, err := w.repo.GetGroup(ctx, job.name)
				if err != nil {
					log.WithError(err).Error("retrieve group")
				}

				go w.queueCheck(ctx, g)

			case <-groupWatchTicker.C:
				groups, err := w.repo.ListGroups(ctx)
				if err != nil {
					log.WithError(err).Error("unable to load groups to watch")
				}

				for _, group := range groups {
					w.AddGroup(group.Name)
				}

			case <-w.quitChan:
				log.Info("watcher is shutting down")
				close(w.workChan)
				w.quitChan = nil
				w.workChan = nil
				done <- true
				log.Info("watcher is stopped")
				break
			}
		}
	}()

	return done, nil

}

// Stop stops the watcher.
func (w *Watcher) Stop() {
	w.wg.Lock()
	defer w.wg.Unlock()

	if w.quitChan != nil {
		w.quitChan <- 1
	} else {
		w.log.Info("watcher was not running, so it can't be stopped")
	}

}

// check group to make sure it is at capacity.
func (w *Watcher) queueCheck(ctx context.Context, g Group) {
	checkDelay := 10 * time.Second

	if err := w.check(ctx, g); err != nil {
		checkDelay = 15 * time.Second
		w.log.
			WithError(err).
			WithField("delay", checkDelay).
			Error("check failed and will be tried again")
	} else {
		w.log.
			WithField("delay", checkDelay).
			Info("scheduling future check")
	}

	timer := time.NewTimer(checkDelay)
	<-timer.C

	w.queueJob(g.Name)
}

func (w *Watcher) check(ctx context.Context, g Group) error {
	log := w.log.WithField("group-name", g.Name)

	resource, err := g.Resource()
	if err != nil {
		return err
	}

	value, err := g.MetricsValue()
	if err != nil {
		return err
	}

	policy := g.Policy
	count, err := resource.Count()
	if err != nil {
		return err
	}

	newCount := policy.Scale(&g, count, value)

	delta := newCount - count

	log.WithFields(logrus.Fields{
		"metric":       g.MetricType,
		"metric-value": value,
		"new-count":    newCount,
		"delta":        delta,
	}).Info("current metric value")

	changed, err := resource.Scale(ctx, g, delta, w.repo)
	if err != nil {
		return err
	}

	if changed {
		wup := policy.WarmUpPeriod()
		log.WithField("warm-up-duration", wup).Info("waiting for new service to warm up")
		time.Sleep(wup)
	}

	return nil
}
