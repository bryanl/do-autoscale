package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValuePolicy(t *testing.T) {
	cases := []struct {
		resourceCount int
		value         float64
		expected      int
		shouldNotify  bool
	}{
		{resourceCount: 5, value: 0.5, expected: 5},
		{resourceCount: 5, value: 0.1, expected: 3, shouldNotify: true},
		{resourceCount: 5, value: 0.8, expected: 8, shouldNotify: true},
		{resourceCount: 1, value: 0.1, expected: 0, shouldNotify: true},
	}

	vp, err := NewValuePolicy(0.8, 3, 0.2, 2)
	assert.NoError(t, err)

	for _, c := range cases {
		mn := &MockMetricNotifier{}

		if c.shouldNotify {
			mn.On("MetricNotify").Return(nil)
		}
		v := vp.Scale(mn, c.resourceCount, c.value)

		assert.Equal(t, c.expected, v)

		assert.True(t, mn.AssertExpectations(t))
	}
}
