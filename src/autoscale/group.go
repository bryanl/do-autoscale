package autoscale

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
)

// ResourceManagerFactoryFn is a function that returns ResourceManagerFactory.
type ResourceManagerFactoryFn func(g *Group) (ResourceManager, error)

// UpdateGroupRequest is a group update request.
type UpdateGroupRequest struct {
	BaseSize int `json:"base_size"`
}

// Group is an autoscale group
type Group struct {
	ID         string          `json:"id" db:"id"`
	Name       string          `json:"name" db:"name"`
	BaseName   string          `json:"baseName" db:"base_name"`
	TemplateID string          `json:"templateID" db:"template_id"`
	MetricType string          `json:"metricType" db:"metric_type"`
	Metric     Metrics         `json:"metric"`
	RawMetric  json.RawMessage `json:"rawMetric" db:"metric"`
	PolicyType string          `json:"policyType" db:"policy_type"`
	Policy     Policy          `json:"policy" `
	RawPolicy  json.RawMessage `json:"rawPolicy" db:"policy"`
}

var _ json.Marshaler = (*Group)(nil)
var _ json.Unmarshaler = (*Group)(nil)

type groupToJSON struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	BaseName   string          `json:"baseName"`
	TemplateID string          `json:"templateID"`
	MetricType string          `json:"metricType"`
	Metric     json.RawMessage `json:"metric"`
	PolicyType string          `json:"policyType"`
	Policy     json.RawMessage `json:"policy"`
}

type jsonToGroup struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	BaseName   string          `json:"baseName"`
	TemplateID string          `json:"templateID"`
	MetricType string          `json:"metricType"`
	Metric     json.RawMessage `json:"metric"`
	PolicyType string          `json:"policyType"`
	Policy     json.RawMessage `json:"policy"`
}

// MarshalJSON marshals a Group into json.
func (g *Group) MarshalJSON() ([]byte, error) {

	tmp := groupToJSON{
		ID:         g.ID,
		Name:       g.Name,
		BaseName:   g.BaseName,
		TemplateID: g.TemplateID,
		MetricType: g.MetricType,
		PolicyType: g.PolicyType,
	}

	if g.Metric != nil {
		m, err := json.Marshal(g.Metric)
		if err != nil {
			logrus.WithError(err).Error("could not encode metric")
			return nil, err
		}

		tmp.Metric = m
	}

	if g.Policy != nil {
		p, err := json.Marshal(g.Policy)
		if err != nil {
			logrus.WithError(err).Error("could not encode policy")
			return nil, err
		}

		tmp.Policy = p
	}

	return json.Marshal(&tmp)
}

// UnmarshalJSON converts json into a Group.
func (g *Group) UnmarshalJSON(b []byte) error {
	tmp := jsonToGroup{}

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	g.ID = tmp.ID
	g.Name = tmp.Name
	g.BaseName = tmp.BaseName
	g.TemplateID = tmp.TemplateID
	g.MetricType = tmp.MetricType
	g.PolicyType = tmp.PolicyType
	g.RawMetric = tmp.Metric
	g.RawPolicy = tmp.Policy

	switch g.MetricType {
	case "load":
		fl, err := NewFileLoad(FileLoadFromJSON(tmp.Metric))
		if err != nil {
			return err
		}

		g.Metric = fl

	default:
		return fmt.Errorf("unknown metric type: %q", g.MetricType)
	}

	switch g.PolicyType {
	case "value":
		vp, err := NewValuePolicy(ValuePolicyFromJSON(tmp.Policy))
		if err != nil {
			return err
		}

		g.Policy = vp

	default:
		return fmt.Errorf("unknown policy type: %q", g.PolicyType)
	}

	return nil
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

// LoadConfig is the configuration settings for a load based metric.
type LoadConfig struct {
	Utilization float64 `json:"utilization"`
}
