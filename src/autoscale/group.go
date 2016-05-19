package autoscale

import (
	"crypto/md5"
	"database/sql/driver"
	"fmt"
	"io"
	"pkg/doclient"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
)

var (
	nameRe = regexp.MustCompile(`^\w[A-Za-z0-9\-]*$`)

	// ResourceManagerFactory creates a ResourceManager given a group.
	ResourceManagerFactory ResourceManagerFactoryFn = func(g *Group) (ResourceManager, error) {
		doClient := doclient.New(DOAccessToken())
		tag := fmt.Sprintf("do:as:%s", g.Name)

		h := md5.New()
		io.WriteString(h, tag)
		hash := fmt.Sprintf("%x", h.Sum(nil))

		newTag := hash[0:8]

		log := logrus.WithField("group-name", g.Name)
		return NewDropletResource(doClient, newTag, log)
	}
)

// ResourceManagerFactoryFn is a function that returns ResourceManagerFactory.
type ResourceManagerFactoryFn func(g *Group) (ResourceManager, error)

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

// UpdateGroupRequest is a group update request.
type UpdateGroupRequest struct {
	BaseSize int `json:"base_size"`
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

	policy Policy
}

// IsValid returns if the template is valid or not.
func (g *Group) IsValid() bool {
	if !nameRe.MatchString(g.Name) {
		return false
	}

	return true
}

// Policy is the scaling policy for the group.
func (g *Group) Policy() (Policy, error) {
	if g.policy == nil {
		p, err := NewValuePolicy(0.75, 2, 0.2, 1)
		if err != nil {
			logrus.
				WithError(err).
				WithField("group-name", g.Name).
				Error("unable to create policy")
			return nil, err
		}

		g.policy = p
	}

	return g.policy, nil
}

// Resource is a resource than can be managed for a group.
func (g *Group) Resource() (ResourceManager, error) {
	return ResourceManagerFactory(g)
}

// NotifyMetrics notifies the metrics system that the instance configuration has changed.
func (g *Group) NotifyMetrics() error {
	r, err := g.Resource()
	if err != nil {
		return err
	}

	allocated, err := r.Allocated()
	if err != nil {
		logrus.WithError(err).Error("unable to retrieve allocated resources")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"group-name":  g.Name,
		"metric-type": g.MetricType,
	}).Info("fetching metric for group")

	m, err := Retrieve(g.MetricType)
	if err != nil {
		logrus.WithError(err).WithField("metric-type", g.MetricType).Error("unable to retrieve metric")
		return err
	}

	return m.Update(g.Name, allocated)
}

// MetricsValue retrieves current metric value for group.
func (g *Group) MetricsValue() (float64, error) {
	logrus.WithFields(logrus.Fields{
		"group-name":  g.Name,
		"metric-type": g.MetricType,
	}).Info("fetching metric value for group")

	m, err := Retrieve(g.MetricType)
	if err != nil {
		logrus.WithError(err).WithField("metric-type", g.MetricType).Error("unable to retrieve metric")
		return 0, err
	}

	return m.Value(g.Name)
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
