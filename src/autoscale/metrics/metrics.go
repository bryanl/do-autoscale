package metrics

var (
	// Available metrics.
	Available = map[string]Gen{
		"neversatisfied": func(*Config) Metrics { return &NeverSatisfied{} },
	}
)

// Gen is a generator for metrics types.
type Gen func(*Config) Metrics

// Config is instance configuration for metrics.
type Config map[string]interface{}

// Metrics pull metrics for a autoscaler.
type Metrics interface {
	IsFull() bool
}

// NeverSatisfied satisfied will always report true for IsFull().
type NeverSatisfied struct {
}

var _ Metrics = (*NeverSatisfied)(nil)

// IsFull always returns true.
func (n *NeverSatisfied) IsFull() bool {
	return true
}
