package autoscale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"pkg/ctxutil"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

const (
	// PrometheusURLContextKey is the context key for the prometheus URL.
	PrometheusURLContextKey = "prometheusURL"

	// PrometheusConfigDirContextKey is the context key for the prometheus config dir.
	PrometheusConfigDirContextKey = "prometheusConfigDir"
)

// PrometheusLoad based on promeetheus metrics.
type PrometheusLoad struct {
	log           *logrus.Entry
	configDir     string
	prometheusURL string
}

// NewPrometheusLoad creates an instance of PrometheusLoad.
func NewPrometheusLoad(ctx context.Context) (*PrometheusLoad, error) {
	log := ctxutil.LogFromContext(ctx)
	if log == nil {
		logger := logrus.New()
		log = logrus.NewEntry(logger)
	}

	promURL := ctxutil.StringFromContext(ctx, PrometheusURLContextKey)
	if promURL == "" {
		return nil, fmt.Errorf("prometheus url wasn't supplied in context")
	}

	configDir := ctxutil.StringFromContext(ctx, PrometheusConfigDirContextKey)
	if configDir == "" {
		var err error
		if configDir, err = ioutil.TempDir("", "promConfig"); err != nil {
			return nil, err
		}
	}

	log.WithFields(logrus.Fields{
		"configDir":     configDir,
		"prometheusURL": promURL,
	}).Info("setting config dir for prometheus")

	return &PrometheusLoad{
		log:           log,
		configDir:     configDir,
		prometheusURL: promURL,
	}, nil
}

var _ Metrics = (*PrometheusLoad)(nil)

// Measure returns the average load for an entire group.
func (l *PrometheusLoad) Measure(ctx context.Context, groupName string) (float64, error) {
	q := fmt.Sprintf(`avg(node_load1{group="%s"})`, groupName)

	config := prometheus.Config{
		Address: l.prometheusURL,
	}
	pClient, err := prometheus.New(config)
	if err != nil {
		return 0, err
	}

	qapi := prometheus.NewQueryAPI(pClient)

	value, err := qapi.Query(ctx, q, time.Now())
	if err != nil {
		return 0, err
	}

	switch t := value.(type) {
	case model.Vector:
		var f float64
		v := value.(model.Vector)
		if len(v) > 0 {
			f = float64(v[0].Value)
		}
		return f, nil

	default:
		l.log.WithField("query-value-type", t).Warning("unknown prometheus query response")
		return 0, nil
	}
}

// Update updates the prometheus config for a group.
func (l *PrometheusLoad) Update(groupName string, resourceAllocations []ResourceAllocation) error {
	l.log.WithField("group", groupName).Info("updating prometheus")
	tg := targetGroup{
		Labels: map[string]string{
			"group": groupName,
		},
	}

	for _, allocation := range resourceAllocations {
		target := fmt.Sprintf("%s:%d", allocation.Address, 9100)
		tg.Targets = append(tg.Targets, target)
	}

	targetGroups := []targetGroup{tg}

	path := fmt.Sprintf("%s/%s.json", l.configDir, groupName)
	if err := os.MkdirAll(l.configDir, 0755); err != nil {
		return err
	}

	b, err := json.Marshal(&targetGroups)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("unable to write metrics json: %v", err)
	}

	return nil
}

// Config returns the configuration for this instance of PrometheusLoad.
func (l *PrometheusLoad) Config() MetricConfig {
	return MetricConfig{
		"configDir": l.configDir,
	}
}

func (l *PrometheusLoad) Values(ctx context.Context, groupName string) ([]TimeSeries, error) {
	q := fmt.Sprintf(`avg(node_load1{group="%s"})`, groupName)

	config := prometheus.Config{
		Address: l.prometheusURL,
	}

	pClient, err := prometheus.New(config)
	if err != nil {
		return nil, err
	}

	qapi := prometheus.NewQueryAPI(pClient)

	now := time.Now()
	then := now.Add(-6 * time.Hour)

	r := prometheus.Range{
		Start: then,
		End:   now,
		Step:  30 * time.Second,
	}
	value, err := qapi.QueryRange(ctx, q, r)
	if err != nil {
		return nil, err
	}

	ts := []TimeSeries{}

	switch value.(type) {
	case model.Matrix:
		v := value.(model.Matrix)
		for _, i := range v {
			for _, sp := range i.Values {
				ts = append(ts, TimeSeries{
					Timestamp: sp.Timestamp.Time(),
					Value:     float64(sp.Value),
				})
			}
		}
	}

	return ts, nil
}

type targetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
