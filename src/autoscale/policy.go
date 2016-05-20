package autoscale

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
)

// MetricNotifier notifies that a metric has changed.
type MetricNotifier interface {
	MetricNotify() error
}

// Policy determine how many resources there should be at the current point in time.
type Policy interface {
	Scale(mn MetricNotifier, resourceCount int, value float64) int
}

// ValuePolicy is a Policy that scales up or down based on the current value.
type ValuePolicy struct {
	ScaleUpValue   float64 `json:"scale_up_value"`
	ScaleUpBy      int     `json:"scale_up_by"`
	ScaleDownValue float64 `json:"scale_down_value"`
	ScaleDownBy    int     `json:"scale_down_by"`

	mu *sync.Mutex
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
		ScaleUpValue:   scaleUpValue,
		ScaleUpBy:      scaleUpBy,
		ScaleDownValue: scaleDownValue,
		ScaleDownBy:    scaleDownBy,
		mu:             &sync.Mutex{},
	}, nil
}

// Value converts a ValuePolicy to JSON to be stored in the databases.
func (p *ValuePolicy) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan converts a DB value back into a FileLoad.
func (p *ValuePolicy) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b := []byte(src.([]uint8))
	return json.Unmarshal(b, p)
}

// Scale returns the amount of items that should exist given a value. If the new value is
// less than 0, then Scale will return 0.
func (p *ValuePolicy) Scale(mn MetricNotifier, resourceCount int, value float64) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	ogCount := resourceCount
	newCount := resourceCount
	if value <= p.ScaleDownValue {
		newCount = newCount - p.ScaleDownBy
	} else if value >= p.ScaleUpValue {
		newCount = newCount + p.ScaleUpBy
	}

	if newCount < 0 {
		newCount = 0
	}

	if ogCount != newCount {
		mn.MetricNotify()
	}

	return newCount
}
