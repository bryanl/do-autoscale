package autoscale

import (
	"pkg/ctxutil"
	"sync"

	"github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
)

// RunList is a list containing groups which are currently
// being autoscaled.
type RunList interface {
	Add(groupName string) error
	Remove(groupName string) error
	IsRunning(groupName string) (bool, error)
	Reset() error
}

// NewRunList creates an instance of RunList.
func NewRunList(ctx context.Context) RunList {
	return &memoryRunList{
		dict: map[string]bool{},
		ctx:  ctx,
	}
}

type memoryRunList struct {
	dict map[string]bool
	mu   sync.Mutex
	ctx  context.Context
}

var _ RunList = (*memoryRunList)(nil)

func (rl *memoryRunList) log() *logrus.Entry {
	return ctxutil.LogFromContext(rl.ctx).WithField("action", "run-list")
}

func (rl *memoryRunList) Add(groupName string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.log().WithField("group", groupName).Info("adding group to run list")

	rl.dict[groupName] = true
	return nil
}

func (rl *memoryRunList) Remove(groupName string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.log().WithField("group", groupName).Info("removing group from run list")

	delete(rl.dict, groupName)
	return nil
}

func (rl *memoryRunList) IsRunning(groupName string) (bool, error) {
	_, ok := rl.dict[groupName]
	return ok, nil
}

func (rl *memoryRunList) Reset() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.dict = map[string]bool{}
	return nil
}
