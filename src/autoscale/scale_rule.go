package autoscale

import "errors"

var (
	// ErrOverlap is an error for overlapping rules.
	ErrOverlap = errors.New("rule overlaps with existing rule")
)

// ScaleGroup is a group of ScaleRules.
type ScaleGroup struct {
	rules []ScaleRule
}

// FindAction finds an action for a scenario.
func (sg *ScaleGroup) FindAction(itemCount int, metricValue float64) int {
	for _, rule := range sg.rules {
		if rule.IsMatch(itemCount, metricValue) {
			return rule.Step
		}
	}

	return 0
}

// AddRule adds a rule to a scale group.
func (sg *ScaleGroup) AddRule(rbl, rbu, step int, mbl, mbu float64) error {
	if sg.isOverlap(rbl, rbu, mbl, mbu) {
		return ErrOverlap
	}

	ruleBounds := IntBounds{Lower: rbl, Upper: rbu}
	sr := ScaleRule{
		Bounds: ruleBounds,
		Step:   step,
	}
	sr.SetMetric(mbl, mbu)

	sg.rules = append(sg.rules, sr)

	return nil
}

// isOverlap returns true if the specification overlaps with a current rule. Rules
// overlap when for any item in the rule bounds can have more than one metric.
func (sg *ScaleGroup) isOverlap(rbl, rbu int, mbl, mbu float64) bool {
	for _, rule := range sg.rules {
		isBoundMatch := rule.Bounds.Lower <= rbl && rule.Bounds.Upper >= rbu
		isMetricMatch := rule.Metric.Bounds.Lower <= mbl && rule.Metric.Bounds.Upper >= mbu
		if isBoundMatch && isMetricMatch {
			return true
		}
	}

	return false
}

// ScaleRule handles lower and upper bounds for items which are scaleable.
type ScaleRule struct {
	Bounds IntBounds
	Step   int
	Metric ScaleMetric
}

// IsMatch matches a Metric rule against an itemCount and a metricValue.
func (sr *ScaleRule) IsMatch(itemCount int, metricValue float64) bool {
	return sr.Bounds.Includes(itemCount) && sr.Metric.Bounds.Includes(metricValue)
}

// SetMetric sets the metrics for a ScaleRule.
func (sr *ScaleRule) SetMetric(lower, upper float64) {
	sr.Metric = ScaleMetric{
		Bounds: FloatBounds{
			Lower: lower,
			Upper: upper,
		},
	}
}

// ScaleMetric handles lower and upper bounds for metrics.
type ScaleMetric struct {
	Bounds FloatBounds
}

// FloatBounds are an upper and lower threshold.
type FloatBounds struct {
	Lower float64
	Upper float64
}

// IsValid returns if a FloatBounds is valid or not. A Bounds is valid if
// the lower bound is less than or equal to the upper bound. Also both
// the lower and upper bounds must be greater than or equal to 0.
func (b *FloatBounds) IsValid() bool {
	return b.Lower <= b.Upper && (b.Lower >= 0 && b.Upper >= 0)
}

// Includes returns if the bounds includes an item.
func (b *FloatBounds) Includes(item float64) bool {
	return item >= b.Lower && item <= b.Upper
}

// IntBounds are an upper and lower threshold.
type IntBounds struct {
	Lower int
	Upper int
}

// IsValid returns if an IntBounds is valid or not. A Bounds is valid if
// the lower bound is less than or equal to the upper bound. Also both
// the lower and upper bounds must be greater than or equal to 0.
func (b *IntBounds) IsValid() bool {
	return b.Lower <= b.Upper && (b.Lower >= 0 && b.Upper >= 0)
}

// Includes returns if the bounds includes an item.
func (b *IntBounds) Includes(item int) bool {
	return item >= b.Lower && item <= b.Upper
}
