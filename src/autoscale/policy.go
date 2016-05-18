package autoscale

import "fmt"

// Policy determine how many resources there should be at the current point in time.
type Policy interface {
	Scale(resourceCount int, value float64) int
}

// ValuePolicy is a Policy that scales up or down based on the current value.
type ValuePolicy struct {
	scaleUpValue   float64
	scaleUpBy      int
	scaleDownValue float64
	scaleDownBy    int
}

var _ Policy = (*ValuePolicy)(nil)

// NewValuePolicy creates an instance of ValuePolicy.
func NewValuePolicy(
	scaleUpValue float64,
	scaleUpBy int,
	scaleDownValue float64,
	scaleDownBy int) (*ValuePolicy, error) {

	if scaleDownValue >= scaleUpValue {
		return nil, fmt.Errorf("scaleDownBalue must be less than scaleUpValue")
	}

	return &ValuePolicy{
		scaleUpValue:   scaleUpValue,
		scaleUpBy:      scaleUpBy,
		scaleDownValue: scaleDownValue,
		scaleDownBy:    scaleDownBy,
	}, nil
}

// Scale returns the amount of items that should exist given a value. If the new value is
// less than 0, then Scale will return 0.
func (p *ValuePolicy) Scale(resourceCount int, value float64) int {
	var newCount = resourceCount
	if value <= p.scaleDownValue {
		newCount = newCount - p.scaleDownBy
	} else if value >= p.scaleUpValue {
		newCount = newCount + p.scaleUpBy
	}

	if newCount < 0 {
		newCount = 0
	}

	return newCount
}
