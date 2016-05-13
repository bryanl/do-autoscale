package autoscale

import (
	"database/sql/driver"
	"regexp"
	"strings"
)

var (
	nameRe = regexp.MustCompile(`^\w[A-Za-z0-9\-]*$`)
)

// StringSlice is a slice of strings
type StringSlice []string

// Value converts a string slice to something the driver value can handle. In this case,
// it creates a CSV.
func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}

// Scan converts a DB value back into a StringSlice.
func (s *StringSlice) Scan(src interface{}) error {
	u8 := src.([]uint8)
	ba := make([]byte, 0, len(u8))
	for _, b := range u8 {
		ba = append(ba, byte(b))
	}

	str := string(ba)
	*s = strings.Split(str, ",")
	return nil
}

// CreateTemplateRequest is a template create request.
type CreateTemplateRequest struct {
	Name     string      `json:"name"`
	Region   string      `json:"region"`
	Size     string      `json:"size"`
	Image    string      `json:"image"`
	SSHKeys  StringSlice `json:"ssh_keys"`
	UserData string      `json:"user_data"`
}

// CreateGroupRequest is a group create request.
type CreateGroupRequest struct {
	Name         string `json:"name"`
	BaseName     string `json:"base_name"`
	BaseSize     int    `json:"base_size"`
	MetricType   string `json:"metric_type"`
	TemplateName string `json:"template_name"`
}

// Group is an autoscale group
type Group struct {
	ID           string     `json:"ID" db:"id"`
	Name         string     `json:"name" db:"name"`
	BaseName     string     `json:"base_name" db:"base_name"`
	BaseSize     int        `json:"base_size" db:"base_size"`
	MetricType   string     `json:"metric_type" db:"metric_type"`
	TemplateName string     `json:"template_name" db:"template_name"`
	ScaleGroup   ScaleGroup `json:"scale_group" db:"rules"`
}

// IsValid returns if the template is valid or not.
func (g *Group) IsValid() bool {
	if !nameRe.MatchString(g.Name) {
		return false
	}

	return true
}

// Template is a template that will be autoscaled.
type Template struct {
	ID       string      `json:"id" db:"id"`
	Name     string      `json:"name" db:"name"`
	Region   string      `json:"string" db:"region"`
	Size     string      `json:"size" db:"size"`
	Image    string      `json:"image" db:"image"`
	SSHKeys  StringSlice `json:"ssh_keys" db:"ssh_keys"`
	UserData string      `json:"user_data" db:"user_data"`
}

// IsValid returns if the template is valid or not.
func (t *Template) IsValid() bool {
	if !nameRe.MatchString(t.Name) {
		return false
	}

	return true
}

// LoadConfig is the configuration settings for a load based metric.
type LoadConfig struct {
	Utilization float64 `json:"utilization"`
}
