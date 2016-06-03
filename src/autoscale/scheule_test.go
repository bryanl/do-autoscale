package autoscale

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestSchedule(t *testing.T) {
	ctx := context.Background()

	actionRan := false

	expectedID := "id"
	as := &ActionStatus{
		Done: make(chan bool, 1),
	}

	fn := func(ctx context.Context, groupID string) *ActionStatus {
		if groupID == expectedID {
			actionRan = true
		}

		as.Done <- true

		return as
	}

	s := NewScheduler(ctx, fn)
	status := s.Status()
	go s.Start()

	status.EnableGroup <- expectedID
	status.Schedule <- expectedID

	activity := <-status.Activity

	require.True(t, actionRan)
	require.Equal(t, expectedID, activity.ID)
	require.NoError(t, activity.Err)
}

func TestSchedule_Disabled(t *testing.T) {
	ctx := context.Background()

	actionRan := false

	expectedID := "id"
	as := &ActionStatus{
		Done: make(chan bool, 1),
	}

	fn := func(ctx context.Context, groupID string) *ActionStatus {
		if groupID == expectedID {
			actionRan = true
		}

		as.Done <- true

		return as
	}

	s := NewScheduler(ctx, fn)
	status := s.Status()
	go s.Start()

	status.DisableGroup <- expectedID
	status.Schedule <- expectedID

	activity := <-status.Activity

	require.False(t, actionRan)
	require.Equal(t, expectedID, activity.ID)
	require.Error(t, activity.Err)
}
