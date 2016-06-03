package autoscale

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitor_AddGroup(t *testing.T) {
	ctx := context.Background()

	g1 := Group{ID: "1"}
	g2 := Group{ID: "2"}
	g3 := Group{ID: "3"}
	groups := []Group{g1, g2, g3}

	repo := &MockRepository{}
	repo.On("ListGroups", ctx).Return(groups, nil)

	schedulerStatus := &SchedulerStatus{
		EnableGroup:  make(chan string, 10),
		DisableGroup: make(chan string, 10),
		Schedule:     make(chan string, 10),
	}

	m, err := NewMonitor(ctx, repo,
		MonitorGroupCheckTimeout(time.Millisecond))
	require.NoError(t, err)

	go m.Start(schedulerStatus)

	for i := range groups {
		id := <-schedulerStatus.EnableGroup
		assert.Equal(t, groups[i].ID, id)
	}

	m.Stop()

	assert.True(t, repo.AssertExpectations(t))
}

func TestMonitor_RemoveGroup(t *testing.T) {
	ctx := context.Background()

	g1 := Group{ID: "1"}
	g2 := Group{ID: "2"}
	g3 := Group{ID: "3"}

	groups := []Group{g1, g2}
	groupsInRunList := []Group{g1, g2, g3}

	repo := &MockRepository{}
	repo.On("ListGroups", ctx).Return(groups, nil)

	runList := NewRunList(ctx)
	for _, g := range groupsInRunList {
		runList.Add(g.ID)
	}

	schedulerStatus := &SchedulerStatus{
		EnableGroup:  make(chan string, 10),
		DisableGroup: make(chan string, 10),
		Schedule:     make(chan string, 10),
	}

	go func() {
		<-schedulerStatus.EnableGroup
	}()

	m, err := NewMonitor(ctx, repo,
		MonitorRunList(runList),
		MonitorGroupCheckTimeout(100*time.Millisecond))
	require.NoError(t, err)

	go m.Start(schedulerStatus)

	groupID := <-schedulerStatus.DisableGroup
	require.Equal(t, g3.ID, groupID)

	m.Stop()

	assert.True(t, repo.AssertExpectations(t))
}
