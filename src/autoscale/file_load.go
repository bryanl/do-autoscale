package autoscale

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
)

var (
	// DefaultStatsDir is the default location for the file based stats metric.
	DefaultStatsDir = "/tmp"
)

// FileLoadOption is a functional option for configuring FileLoad.
type FileLoadOption func(*FileLoad) error

// FileLoad returns hardcoded metrics from files. This is useful if you
// are on a plane and want to test the autoscaler.
type FileLoad struct {
	StatsDir string `json:"stats_dir"`
}

var _ Metrics = (*FileLoad)(nil)

// NewFileLoad creates a new FileLoad instance.
func NewFileLoad(options ...FileLoadOption) (*FileLoad, error) {
	fl := &FileLoad{}

	for _, opt := range options {
		if err := opt(fl); err != nil {
			return nil, err
		}
	}

	if fl.StatsDir == "" {
		fl.StatsDir = DefaultStatsDir
	}

	return fl, nil
}

// FileLoadPath sets the director for FileLoad.
func FileLoadPath(dir string) FileLoadOption {
	return func(fl *FileLoad) error {
		fi, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("unable to stat %q: %v", dir, err)
		}

		if !fi.IsDir() {
			return fmt.Errorf("%q is not a directory", dir)
		}

		fl.StatsDir = dir
		return nil
	}
}

// FileLoadFromJSON configures a FileLoad from JSON.
func FileLoadFromJSON(in json.RawMessage) FileLoadOption {
	return func(fl *FileLoad) error {
		var c map[string]interface{}
		if err := json.Unmarshal(in, &c); err != nil {
			fl.StatsDir = DefaultStatsDir
		} else {
			if dir, ok := c["stats_dir"].(string); ok {
				fl.StatsDir = dir
			}
		}

		return nil
	}
}

// Value converts a FileLoad to JSON to be stored in the databases.
func (l *FileLoad) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan converts a DB value back into a FileLoad.
func (l *FileLoad) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b := []byte(src.([]uint8))
	return json.Unmarshal(b, l)
}

// Measure returns the current value fro a group.
func (l *FileLoad) Measure(ctx context.Context, groupName string) (float64, error) {
	p := filepath.Join(l.StatsDir, groupName)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return 0, err
	}

	str := string(b)
	return strconv.ParseFloat(strings.TrimSpace(str), 64)
}

// Update updates the metric configuration for a group using the resource allocations.
func (l *FileLoad) Update(groupName string, resourceAllocations []ResourceAllocation) error {
	// currently a no-op as the metrics are hard coded.

	return nil
}

// Config returns the configuration for this instances of FileLoad.
func (l *FileLoad) Config() MetricConfig {
	return MetricConfig{
		"statsDir": l.StatsDir,
	}
}

// Values for values
func (l *FileLoad) Values(ctx context.Context, groupName string, tr TimeRange) ([]TimeSeries, error) {
	return []TimeSeries{
		{Timestamp: time.Now(), Value: 9},
		{Timestamp: time.Now().Add(-60 * time.Hour), Value: 9},
	}, nil
}

// InstanceValues for values
func (l *FileLoad) InstanceValues(ctx context.Context, groupName, instanceID string, tr TimeRange) ([]TimeSeries, error) {
	return []TimeSeries{
		{Timestamp: time.Now(), Value: 8},
		{Timestamp: time.Now().Add(-60 * time.Hour), Value: 8},
	}, nil
}

// Remove removes the configuration for a group.
func (l *FileLoad) Remove(ctx context.Context, groupID string) error {
	// no op as values are hard coded.
	return nil
}
