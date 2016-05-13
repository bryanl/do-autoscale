package watcher

import (
	"autoscale/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatcher(t *testing.T) {
	repo := &mocks.Repository{}

	watcher := New(repo)

	done, err := watcher.Watch()
	assert.NoError(t, err)

	_, err = watcher.Watch()
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
