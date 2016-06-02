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
	Add(groupID string) error
	Remove(groupID string) error
	IsRunning(groupID string) bool
	List() []string
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

func (rl *memoryRunList) Add(groupID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.log().WithField("group", groupID).Info("adding group to run list")

	rl.dict[groupID] = true
	return nil
}

func (rl *memoryRunList) Remove(groupID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.log().WithField("group", groupID).Info("removing group from run list")

	delete(rl.dict, groupID)
	return nil
}

func (rl *memoryRunList) List() []string {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	out := []string{}
	for k := range rl.dict {
		out = append(out, k)
	}

	return out
}

func (rl *memoryRunList) IsRunning(groupID string) bool {
	_, ok := rl.dict[groupID]
	return ok
}

func (rl *memoryRunList) Reset() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.dict = map[string]bool{}
	return nil
}
