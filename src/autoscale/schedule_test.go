package autoscale

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

type testCheck struct {
	ScaleFn   GroupActionFn
	DisableFn GroupActionFn
}

func (tc *testCheck) Scale(ctx context.Context, groupID string) *ActionStatus {
	return tc.ScaleFn(ctx, groupID)
}

func (tc *testCheck) Disable(ctx context.Context, groupID string) *ActionStatus {
	return tc.DisableFn(ctx, groupID)
}

func TestSchedule(t *testing.T) {
	ctx := context.Background()

	actionRan := false

	expectedID := "id"
	as := &ActionStatus{
		Done: make(chan bool, 1),
	}

	tc := &testCheck{
		ScaleFn: func(ctx context.Context, groupID string) *ActionStatus {
			if groupID == expectedID {
				actionRan = true
			}

			as.Done <- true

			return as
		},
	}

	s := NewScheduler(ctx, tc)
	status := s.Status()
	go s.Start()

	status.EnableGroup <- expectedID
	status.Schedule <- expectedID

	activity := <-status.Activity

	require.True(t, actionRan)
	require.Equal(t, expectedID, activity.ID)
	require.NoError(t, activity.Err)
}
