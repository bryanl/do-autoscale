package autoscale

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	jsonMessage := []byte(`{
    "scale_up_value": 0.8,
    "scale_up_by": 3,
    "scale_down_value": 0.2,
    "scale_down_by": 2,
    "warm_up_duration": "1m"
  }`)

	vp, err := NewValuePolicy(ValuePolicyFromJSON(jsonMessage))
	require.NoError(t, err)

	for _, c := range cases {
		mn := &MockMetricNotifier{}

		if c.shouldNotify {
			mn.On("MetricNotify").Return(nil)
		}
		v := vp.Scale(mn, c.resourceCount, c.value)

		assert.Equal(t, c.expected, v, fmt.Sprintf("case: %#v\n", c))

		assert.True(t, mn.AssertExpectations(t))
	}
}
