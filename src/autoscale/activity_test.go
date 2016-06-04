package autoscale

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActivityManager(t *testing.T) {
	activityChan := make(chan SchedulerActivity, 1)

	am := NewActivityManager(activityChan)

	l1 := make(chan SchedulerActivity, 1)
	am.RegisterListener(l1)
	l2 := make(chan SchedulerActivity, 1)
	am.RegisterListener(l2)

	go am.Start()

	in := SchedulerActivity{ID: "id"}
	activityChan <- in

	var out SchedulerActivity
	out = <-l1
	require.Equal(t, in, out)

	out = <-l2
	require.Equal(t, in, out)

	close(l1)

	activityChan <- in
	out = <-l2
	require.Equal(t, in, out)
}
