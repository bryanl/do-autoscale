package autoscale

import (
	"encoding/json"
	"strings"
)

// Group is an autoscale group
type Group struct {
	ID         string          `json:"ID"`
	BaseName   string          `json:"base_name"`
	BaseSize   int             `json:"base_size"`
	MetricType string          `json:"metric_type"`
	LoadConfig json.RawMessage `json:"load_config"`
}

// Template is a template that will be autoscaled.
type Template struct {
	ID         int    `json:"id" db:"id"`
	Region     string `json:"string" db:"region"`
	Size       string `json:"size" db:"size"`
	Image      string `json:"image" db:"image"`
	RawSSHKeys string `json:"ssh_keys" db:"ssh_keys"`
	UserData   string `json:"user_data" db:"user_data"`
}

// SSHKeys returns ssh keys as a string.
func (t *Template) SSHKeys() []string {
	return strings.Split(t.RawSSHKeys, ",")

}

// LoadConfig is the configuration settings for a load based metric.
type LoadConfig struct {
	Utilization float64 `json:"utilization"`
}
