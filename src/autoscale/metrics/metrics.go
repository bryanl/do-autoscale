package metrics

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

// ResourceAllocation is information about an allocated resource.
type ResourceAllocation struct {
	Name    string
	Address string
}

var (
	metrics = map[string]Metrics{}
)

// Retrieve retrieves the currently configured metrics.
func Retrieve(metricType string) (Metrics, error) {
	m, ok := metrics[metricType]
	if !ok {
		return nil, fmt.Errorf("unknown metric %q", metricType)
	}

	return m, nil
}

// Gen is a generator for metrics types.
type Gen func(*Config) Metrics

// Config is instance configuration for metrics.
type Config map[string]interface{}

// Metrics pull metrics for a autoscaler.
type Metrics interface {
	Value(groupName string) (float64, error)
	Update(groupName string, resourceAllocations []ResourceAllocation) error
}

// RegisterMetric registers metrics.
func RegisterMetric(name string, m Metrics) {
	logrus.WithFields(logrus.Fields{
		"metric-name": name,
		"metric-type": fmt.Sprintf("%T", m)}).Info("registering metric")

	metrics[name] = m
}

// RegisterDefaultMetrics registers a default set of metrics.
func RegisterDefaultMetrics() {
	m, err := NewFileLoad("/tmp")
	if err != nil {
		logrus.WithError(err).Error("unable to register file based load metric")
		return
	}

	RegisterMetric("load", m)
}
