package autoscale

import (
	"encoding/json"
	"regexp"
	"strings"
)

var (
	nameRe = regexp.MustCompile(`^\w[A-Za-z0-9\-]*$`)
)

// Group is an autoscale group
type Group struct {
	ID         string `json:"ID" db:"id"`
	BaseName   string `json:"base_name" db:"base_name"`
	BaseSize   int    `json:"base_size" db:"base_size"`
	MetricType string `json:"metric_type" db:"metric_type"`
	TemplateID int    `json:"template_id" db:"template_id"`
}

// Template is a template that will be autoscaled.
type Template struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Region     string `json:"string" db:"region"`
	Size       string `json:"size" db:"size"`
	Image      string `json:"image" db:"image"`
	RawSSHKeys string `json:"ssh_keys" db:"ssh_keys"`
	UserData   string `json:"user_data" db:"user_data"`
}

// IsValid returns if the template is valid or not.
func (t *Template) IsValid() bool {
	if !nameRe.MatchString(t.Name) {
		return false
	}

	return true
}

// SSHKeys returns ssh keys as a string.
func (t *Template) SSHKeys() []string {
	return strings.Split(t.RawSSHKeys, ",")
}

// MarshalJSON is a custom json marshaler for template.
func (t *Template) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID       int      `json:"id" db:"id"`
		Name     string   `json:"name" db:"name"`
		Region   string   `json:"string" db:"region"`
		Size     string   `json:"size" db:"size"`
		Image    string   `json:"image" db:"image"`
		SSHKeys  []string `json:"ssh_keys" db:"ssh_keys"`
		UserData string   `json:"user_data" db:"user_data"`
	}{
		ID:       t.ID,
		Name:     t.Name,
		Region:   t.Region,
		Size:     t.Size,
		SSHKeys:  t.SSHKeys(),
		UserData: t.UserData,
	})
}

// UnmarshalJSON is a custom json unmarshaler for template.
func (t *Template) UnmarshalJSON(data []byte) error {
	type Alias Template
	aux := &struct {
		Keys []string `json:"ssh_keys"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.RawSSHKeys = strings.Join(aux.Keys, ",")
	return nil
}

// LoadConfig is the configuration settings for a load based metric.
type LoadConfig struct {
	Utilization float64 `json:"utilization"`
}
