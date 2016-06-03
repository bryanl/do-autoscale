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
	}{
		{resourceCount: 5, value: 0.5, expected: 5},
		{resourceCount: 5, value: 0.1, expected: 3},
		{resourceCount: 5, value: 0.8, expected: 8},
		{resourceCount: 1, value: 0.1, expected: 1},
		{resourceCount: 9, value: 0.8, expected: 10},
	}

	vp, err := NewValuePolicy(ValuePolicyScale(
		1, 10, 0.8, 3, 0.2, 2,
	))
	require.NoError(t, err)

	for _, c := range cases {
		v := vp.CalculateSize(c.resourceCount, c.value)
		assert.Equal(t, c.expected, v, fmt.Sprintf("case: %#v\n", c))
	}
}
