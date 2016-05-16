package metrics

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FileLoad returns hardcoded metrics from files. This is useful if you
// are on a plane and want to test the autoscaler.
type FileLoad struct {
	statsDir string
}

var _ Metrics = (*FileLoad)(nil)

// NewFileLoad creates a new FileLoad instance.
func NewFileLoad(statsDir string) (*FileLoad, error) {
	fi, err := os.Stat(statsDir)
	if err != nil {
		return nil, fmt.Errorf("unable to stat %q: %v", statsDir, err)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", statsDir)
	}

	return &FileLoad{
		statsDir: statsDir,
	}, nil
}

// Value returns the current value fro a group.
func (l *FileLoad) Value(groupName string) (float64, error) {
	p := filepath.Join(l.statsDir, groupName)
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
