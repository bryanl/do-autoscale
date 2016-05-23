package autoscale

import (
	"crypto/md5"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"pkg/doclient"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/context"

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

	defaultValuePolicy = valuePolicyData{
		ScaleUpValue:   0.8,
		ScaleUpBy:      2,
		ScaleDownValue: 0.2,
		ScaleDownBy:    1,
		WarmUpDuration: "10s",
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
	Name         string          `json:"name"`
	BaseName     string          `json:"base_name"`
	TemplateName string          `json:"template_name"`
	MetricType   string          `json:"metric_type"`
	Metric       json.RawMessage `json:"metric,omitempty"`
	PolicyType   string          `json:"policy_type"`
	Policy       json.RawMessage `json:"policy,omitempty"`
}

// ConvertToGroup convertes a CreateGroupRequest to a Group
func (cgr *CreateGroupRequest) ConvertToGroup(ctx context.Context) (*Group, error) {
	g := &Group{
		Name:         cgr.Name,
		BaseName:     cgr.BaseName,
		TemplateName: cgr.TemplateName,
		MetricType:   cgr.MetricType,
		PolicyType:   cgr.PolicyType,
	}

	switch g.MetricType {
	case "load":
		fl, err := NewFileLoad(FileLoadFromJSON(cgr.Metric))
		if err != nil {
			return nil, err
		}

		g.Metric = fl

	default:
		return nil, fmt.Errorf("unknown metric type: %q", g.MetricType)
	}

	switch g.PolicyType {
	case "value":
		vp, err := NewValuePolicy(ValuePolicyFromJSON(cgr.Policy))
		if err != nil {
			return nil, err
		}

		g.Policy = vp

	default:
		return nil, fmt.Errorf("unknown policy type: %q", g.PolicyType)
	}

	return g, nil
}

// UpdateGroupRequest is a group update request.
type UpdateGroupRequest struct {
	BaseSize int `json:"base_size"`
}

// Group is an autoscale group
type Group struct {
	ID           string  `json:"ID" db:"id"`
	Name         string  `json:"name" db:"name"`
	BaseName     string  `json:"base_name" db:"base_name"`
	TemplateName string  `json:"template_name" db:"template_name"`
	MetricType   string  `json:"metric_type" db:"metric_type"`
	Metric       Metrics `json:"metric" db:"metric"`
	PolicyType   string  `json:"policy_type" db:"policy_type"`
	Policy       Policy  `json:"policy" db:"policy"`
}

// IsValid returns if the template is valid or not.
func (g *Group) IsValid() bool {
	if !nameRe.MatchString(g.Name) {
		return false
	}

	return true
}

// Resource is a resource than can be managed for a group.
func (g *Group) Resource() (ResourceManager, error) {
	return ResourceManagerFactory(g)
}

// MetricNotify notifies the metrics system that the instance configuration has changed.
func (g *Group) MetricNotify() error {
	logrus.Info("notifying metric configuration has changed")
	r, err := g.Resource()
	if err != nil {
		return err
	}

	allocated, err := r.Allocated()
	if err != nil {
		logrus.WithError(err).Error("unable to retrieve allocated resources")
		return err
	}

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
	}).Debug("fetching metric value for group")

	m, err := Retrieve(g.MetricType)
	if err != nil {
		logrus.WithError(err).WithField("metric-type", g.MetricType).Error("unable to retrieve metric")
		return 0, err
	}

	return m.Measure(g.Name)
}

// LoadPolicy loads policies.
func (g *Group) LoadPolicy(in interface{}) error {
	switch g.PolicyType {
	default:
		return fmt.Errorf("unknown policy type: %v", g.PolicyType)
	case "value":
		vp := ValuePolicy{mu: &sync.Mutex{}}
		if err := vp.Scan(in); err != nil {
			return err
		}

		g.Policy = &vp
	}

	return nil
}

// LoadMetric loads metrics.
func (g *Group) LoadMetric(in interface{}) error {
	switch g.MetricType {
	default:
		return fmt.Errorf("unknown metric type: %v", g.PolicyType)
	case "load":
		fl := FileLoad{}
		if err := fl.Scan(in); err != nil {
			return err
		}

		g.Metric = &fl
	}

	return nil
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
