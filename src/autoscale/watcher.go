package autoscale

import (
	"fmt"
	"pkg/ctxutil"
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
	repo         Repository
	groupMonitor GroupMonitor
	ctx          context.Context

	groupNames []string
	workChan   chan watchedJob
	quitChan   chan int

	wg sync.Mutex
}

func makeJobQueue() chan watchedJob {
	return make(chan watchedJob, 1000)
}

// NewWatcher creates an instance of Watcher.
func NewWatcher(ctx context.Context, repo Repository) (*Watcher, error) {
	gm, err := NewGroupMonitor(ctx, repo)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		repo:         repo,
		groupMonitor: gm,
		ctx:          ctx,
		workChan:     makeJobQueue(),
	}, nil
}

func (w *Watcher) log() *logrus.Entry {
	return ctxutil.LogFromContext(w.ctx).WithField("action", "watcher")
}

func (w *Watcher) queueJob(name string) {
	w.log().WithField("name", name).Info("queueing job")
	job := watchedJob{
		name: name,
	}

	w.workChan <- job
}

// Groups are the currently watched groups.
func (w *Watcher) Groups() []string {
	return w.groupNames
}

// Watch starts the watching process.
func (w *Watcher) Watch(ctx context.Context) (chan bool, error) {
	w.wg.Lock()
	defer w.wg.Unlock()

	log := w.log()

	if w.quitChan != nil {
		log.Warn("watcher is already running")
		return nil, fmt.Errorf("watcher is already running")
	}

	log.Info("starting watcher")

	done := make(chan bool, 1)
	w.quitChan = make(chan int, 1)

	if w.workChan == nil {
		w.workChan = makeJobQueue()
	}

	w.groupMonitor.Start(func(groupName string) {
		log.WithField("group", groupName).Info("watch group")
		w.workChan <- watchedJob{name: groupName}
	})

	go func() {

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
		w.log().Info("watcher was not running, so it can't be stopped")
	}

}

// check group to make sure it is at capacity.
func (w *Watcher) queueCheck(ctx context.Context, g Group) {
	checkDelay := 10 * time.Second

	if err := w.check(ctx, g); err != nil {
		checkDelay = 15 * time.Second
		w.log().
			WithError(err).
			WithField("delay", checkDelay).
			Error("check failed and will be tried again")
	} else {
		w.log().
			WithField("delay", checkDelay).
			Info("scheduling future check")
	}

	if w.groupMonitor.InRunList(g.Name) {
		timer := time.NewTimer(checkDelay)
		<-timer.C

		w.queueJob(g.Name)
	}
}

func (w *Watcher) check(ctx context.Context, g Group) error {
	log := w.log().WithField("group-name", g.Name)

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

	newCount := policy.CalculateSize(&g, count, value)

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
		if err := g.MetricNotify(); err != nil {
			log.WithError(err).Error("notifying metric of new config")
		}
		wup := policy.WarmUpPeriod()
		log.WithField("warm-up-duration", wup).Info("waiting for new service to warm up")
		time.Sleep(wup)
	}

	return nil
}
