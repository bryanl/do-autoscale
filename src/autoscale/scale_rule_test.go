package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScaleGroup_AddRules(t *testing.T) {
	sg := ScaleGroup{}

	err := sg.AddRule(0, 10, 5, 10, 20)
	assert.NoError(t, err)

	err = sg.AddRule(0, 10, 2, 15, 20)
	assert.Error(t, err, "metric bounds can't overlap")

	err = sg.AddRule(11, 20, 3, 10, 20)
	assert.NoError(t, err)

}

func TestScaleGroup_Match(t *testing.T) {
	sg := ScaleGroup{}
	err := sg.AddRule(0, 10, 5, 10, 20)
	assert.NoError(t, err)

	err = sg.AddRule(11, 20, 3, 10, 20)
	assert.NoError(t, err)

	i := sg.FindAction(5, 15)
	assert.Equal(t, 5, i)

	i = sg.FindAction(11, 15)
	assert.Equal(t, 3, i)

	i = sg.FindAction(15, 30)
	assert.Equal(t, 0, i)

	i = sg.FindAction(35, 10)
	assert.Equal(t, 0, i)

}

func TestFloatBounds_IsValid(t *testing.T) {
	cases := []struct {
		lower   float64
		upper   float64
		isValid bool
	}{
		{lower: 0, upper: 20, isValid: true},
		{lower: 10, upper: 5, isValid: false},
		{lower: -10, upper: 5, isValid: false},
	}

	for _, c := range cases {
		b := FloatBounds{Lower: c.lower, Upper: c.upper}
		assert.Equal(t, c.isValid, b.IsValid())
	}
}

func TestIntBounds_IsValid(t *testing.T) {
	cases := []struct {
		lower   int
		upper   int
		isValid bool
	}{
		{lower: 0, upper: 20, isValid: true},
		{lower: 10, upper: 5, isValid: false},
		{lower: -10, upper: 5, isValid: false},
	}

	for _, c := range cases {
		b := IntBounds{Lower: c.lower, Upper: c.upper}
		assert.Equal(t, c.isValid, b.IsValid())
	}
}

func TestScaleRule_IsMatch(t *testing.T) {
	rule := ScaleRule{
		Bounds: IntBounds{Lower: 0, Upper: 10},
		Step:   5,
		Metric: ScaleMetric{
			Bounds: FloatBounds{Lower: 10, Upper: 20},
		},
	}

	assert.True(t, rule.IsMatch(5, 15))
}
