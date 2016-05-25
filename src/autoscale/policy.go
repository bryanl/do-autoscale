package autoscale

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MetricNotifier notifies that a metric has changed.
type MetricNotifier interface {
	MetricNotify() error
}

// PolicyConfig is the configuration for a Policy.
type PolicyConfig map[string]interface{}

// Policy determine how many resources there should be at the current point in time.
type Policy interface {
	CalculateSize(resourceCount int, value float64) int
	WarmUpPeriod() time.Duration
	Config() PolicyConfig
	MarshalJSON() ([]byte, error)
}

type valuePolicyData struct {
	MinSize        int     `json:"min_size"`
	MaxSize        int     `json:"max_size"`
	ScaleUpValue   float64 `json:"scale_up_value"`
	ScaleUpBy      int     `json:"scale_up_by"`
	ScaleDownValue float64 `json:"scale_down_value"`
	ScaleDownBy    int     `json:"scale_down_by"`
	WarmUpDuration string  `json:"warm_up_duration"`
}

// ValuePolicyOption is a functional option for configuring a ValuePolicy.
type ValuePolicyOption func(*ValuePolicy) error

// ValuePolicy is a Policy that scales up or down based on the current value.
type ValuePolicy struct {
	vpd          valuePolicyData
	warmUpPeriod time.Duration
	mu           *sync.Mutex
}

var _ Policy = (*ValuePolicy)(nil)

// NewValuePolicy creates an instance of ValuePolicy.
func NewValuePolicy(options ...ValuePolicyOption) (*ValuePolicy, error) {
	vp := &ValuePolicy{
		mu: &sync.Mutex{},
	}

	for _, opt := range options {
		if err := opt(vp); err != nil {
			return nil, err
		}
	}

	return vp, nil
}

// ValuePolicyScale sets scale parameters for a ValuePolicy.
func ValuePolicyScale(minSize int, scaleUpValue float64, scaleUpBy int, scaleDownValue float64, scaleDownBy int) ValuePolicyOption {
	return func(vp *ValuePolicy) error {
		vp.vpd.ScaleUpValue = scaleUpValue
		vp.vpd.ScaleUpBy = scaleUpBy
		vp.vpd.ScaleDownValue = scaleDownValue
		vp.vpd.ScaleDownBy = scaleDownBy

		return nil
	}
}

// ValuePolicyFromJSON configures a ValuePolicy from JSON.
func ValuePolicyFromJSON(in json.RawMessage) ValuePolicyOption {
	return func(vp *ValuePolicy) error {
		var vpd valuePolicyData
		if err := json.Unmarshal(in, &vpd); err != nil {
			vpd = defaultValuePolicy
		}

		if vpd.MaxSize < vpd.MinSize {
			return fmt.Errorf("maxSize (%d) must be greater or equal to minSize(%d)", vpd.MaxSize, vpd.MinSize)
		}

		if vpd.ScaleDownValue >= vpd.ScaleUpValue {
			return fmt.Errorf("scaleUpValue (%f) must be more than scaleDownValue (%f)",
				vpd.ScaleUpValue, vpd.ScaleDownValue)
		}

		dur, err := time.ParseDuration(vpd.WarmUpDuration)
		if err != nil {
			return err
		}

		vp.vpd = vpd
		vp.warmUpPeriod = dur

		return nil
	}
}

// Value converts a ValuePolicy to JSON to be stored in the databases.
func (p *ValuePolicy) Value() (driver.Value, error) {
	return json.Marshal(p.vpd)
}

// Scan converts a DB value back into a FileLoad.
func (p *ValuePolicy) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b := json.RawMessage(src.([]uint8))

	return ValuePolicyFromJSON(b)(p)
}

// CalculateSize returns the amount of items that should exist given a value. If the new value is
// less than 0, then Scale will return 0.
func (p *ValuePolicy) CalculateSize(resourceCount int, value float64) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	newCount := resourceCount
	if value <= p.vpd.ScaleDownValue {
		newCount = newCount - p.vpd.ScaleDownBy
	} else if value >= p.vpd.ScaleUpValue {
		newCount = newCount + p.vpd.ScaleUpBy
	}

	if newCount <= p.vpd.MinSize {
		return p.vpd.MinSize
	}

	if newCount > p.vpd.MaxSize {
		return p.vpd.MaxSize
	}

	return newCount
}

// WarmUpPeriod is the time needed for the new service to warm up. No checks should happen in this period.
func (p *ValuePolicy) WarmUpPeriod() time.Duration {
	return p.warmUpPeriod
}

// Config is the current configuration for ValuePolicy.
func (p *ValuePolicy) Config() PolicyConfig {
	return PolicyConfig{
		"scaleUpBy":      p.vpd.ScaleUpBy,
		"scaleUpValue":   p.vpd.ScaleUpValue,
		"scaleDownBy":    p.vpd.ScaleDownBy,
		"scaleDownValue": p.vpd.ScaleDownValue,
		"warmUpPeriod":   p.warmUpPeriod,
	}
}

// MarshalJSON converts policy to JSON.
func (p *ValuePolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(&p.vpd)
}
