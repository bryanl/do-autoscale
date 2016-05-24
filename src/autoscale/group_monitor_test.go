package autoscale

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

func TestGroupMonitor(t *testing.T) {
	ogTimeout := groupCheckTimeout
	defer func() {
		groupCheckTimeout = ogTimeout
	}()

	groupCheckTimeout = time.Millisecond

	ctx := context.Background()
	repo := &MockRepository{}

	groups := []Group{
		{Name: "g1"},
		{Name: "g2"},
	}

	repo.On("ListGroups", ctx).Return(groups, nil)

	runList := &MockRunList{}
	for _, g := range groups {
		runList.On("Add", g.Name).Return(nil)
		runList.On("IsRunning", g.Name).Return(false, nil).Once()
		runList.On("IsRunning", g.Name).Return(true, nil)
	}
	runList.On("Reset").Return(nil)

	gm, err := NewGroupMonitor(ctx, repo, OptionRunList(runList))
	require.NoError(t, err)

	callCount := 0
	wg := sync.WaitGroup{}
	wg.Add(len(groups))
	done := gm.Start(func(groupName string) {
		callCount++
		wg.Done()
	})

	gm.Stop()

	wg.Wait()

	<-done

	require.Equal(t, 2, callCount)
	require.True(t, repo.AssertExpectations(t))
	require.True(t, runList.AssertExpectations(t))
}
