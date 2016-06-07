package autoscale

import (
	"fmt"
	"pkg/ctxutil"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
)

const (
	defaultFileLoadPath = "/tmp"
)

const (
	// OptionFileLoadPath is a file load path option.
	OptionFileLoadPath = iota
)

var (
	metrics = map[string]Metrics{}

	// DefaultConfig is the default configuration for metrics.
	DefaultConfig = Config{}
)

// TimeRange is a time range.
type TimeRange string

const (
	// RangeQuarterDay is six hours.
	RangeQuarterDay TimeRange = "6h"
	// RangeDay is a one day.
	RangeDay TimeRange = "24h"
	// RangeWeek is seven days.
	RangeWeek TimeRange = "168h"
	// RangeMonth is 30 days.
	RangeMonth TimeRange = "720h"
)

// Duration converts a TimeRange to a time.Duration.
func (tr TimeRange) Duration() (time.Duration, error) {
	return time.ParseDuration(string(tr))
}

type TimeSeries struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// MetricConfig is the configuration for a Metric.
type MetricConfig map[string]interface{}

// ResourceAllocation is information about an allocated resource.
type ResourceAllocation struct {
	Name    string
	Address string
}

// Config is configuration for metrics
type Config map[int]interface{}

// Retrieve retrieves the currently configured metrics.
func Retrieve(metricType string) (Metrics, error) {
	m, ok := metrics[metricType]
	if !ok {
		return nil, fmt.Errorf("unknown metric %q", metricType)
	}

	return m, nil
}

// Metrics pull metrics for a autoscaler.
type Metrics interface {
	Measure(ctx context.Context, groupName string) (float64, error)
	Update(groupName string, resourceAllocations []ResourceAllocation) error
	Config() MetricConfig
	Values(ctx context.Context, groupName string, rangeLength TimeRange) ([]TimeSeries, error)
	Remove(ctx context.Context, groupID string) error
}

// RegisterMetric registers metrics.
func RegisterMetric(name string, m Metrics) {
	logrus.WithFields(logrus.Fields{
		"metric-name": name,
		"metric-type": fmt.Sprintf("%T", m)}).Info("registering metric")

	metrics[name] = m
}

// RegisterOfflineMetrics registers an offline set of metrics.
func RegisterOfflineMetrics(ctx context.Context) {
	log := ctxutil.LogFromContext(ctx)
	var path = defaultFileLoadPath
	if p, ok := DefaultConfig[OptionFileLoadPath]; ok {
		path = p.(string)
	}

	m, err := NewFileLoad(FileLoadPath(path))
	if err != nil {
		log.WithError(err).Error("unable to register file based load metric")
		return
	}

	RegisterMetric("load", m)
}

// RegisterDefaultMetrics registers a default set of metrics.
func RegisterDefaultMetrics(ctx context.Context) {
	log := ctxutil.LogFromContext(ctx)
	m, err := NewPrometheusLoad(ctx)
	if err != nil {
		log.WithError(err).Error("unable to register prometheus based load metric")
	}

	RegisterMetric("load", m)
}
