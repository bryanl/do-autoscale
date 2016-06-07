package autoscale

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricsDuration(t *testing.T) {
	cases := []struct {
		tr               TimeRange
		expectedDuration time.Duration
		isError          bool
	}{
		{tr: RangeQuarterDay, expectedDuration: 6 * time.Hour, isError: false},
		{tr: TimeRange("wtf"), isError: true},
	}

	for _, c := range cases {
		d, err := c.tr.Duration()
		if c.isError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.expectedDuration, d)
		}
	}
}
