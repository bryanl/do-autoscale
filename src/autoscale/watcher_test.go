package autoscale

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func OffTestWatcher(t *testing.T) {
	repo := &MockRepository{}
	ctx := context.Background()

	watcher := NewWatcher(repo)

	done, err := watcher.Watch(ctx)
	assert.NoError(t, err)

	_, err = watcher.Watch(ctx)
	assert.Error(t, err)

	err = watcher.AddGroup("group-1")
	assert.NoError(t, err)
	assert.Len(t, watcher.Groups(), 1)

	err = watcher.AddGroup("group-1")
	assert.Error(t, err)

	watcher.RemoveGroup("group-1")
	assert.Len(t, watcher.Groups(), 0)

	watcher.Stop()
	<-done

}
