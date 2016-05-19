package autoscale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"pkg/ctxutil"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

// PrometheusLoad based on promeetheus metrics.
type PrometheusLoad struct {
	log       *logrus.Entry
	configDir string
}

// NewPrometheusLoad creates an instance of PrometheusLoad.
func NewPrometheusLoad(ctx context.Context) (*PrometheusLoad, error) {
	log := ctxutil.LogFromContext(ctx)
	if log == nil {
		logger := logrus.New()
		log = logrus.NewEntry(logger)
	}

	configDir := ctx.Value("prometheusConfigDir").(string)
	if configDir == "" {
		var err error
		if configDir, err = ioutil.TempDir("", "promConfig"); err != nil {
			return nil, err
		}
	}

	log.WithFields(logrus.Fields{
		"configDir": configDir,
	}).Info("setting config dir for prometheus")

	return &PrometheusLoad{
		log:       log,
		configDir: configDir,
	}, nil
}

var _ Metrics = (*PrometheusLoad)(nil)

// Value returns the average load for an entire group.
func (l *PrometheusLoad) Value(groupName string) (float64, error) {
	q := fmt.Sprintf(`avg(node_load1{group="%s"})`, groupName)

	// TODO this should be a different func
	config := prometheus.Config{
		Address: "http://localhost:9090",
	}
	pClient, err := prometheus.New(config)
	if err != nil {
		return 0, err
	}

	qapi := prometheus.NewQueryAPI(pClient)
	ctx := context.Background()
	value, err := qapi.Query(ctx, q, time.Now())
	if err != nil {
		return 0, err
	}

	switch t := value.(type) {
	case model.Vector:
		sample := value.(model.Vector)[0]
		return float64(sample.Value), nil
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

	b, err := json.Marshal(&targetGroups)
	if err != nil {
		return err
	}

	l.log.WithField("path", path).Info("writing prometheus target file")

	return ioutil.WriteFile(path, b, 0644)
}

func (l *PrometheusLoad) queryURL() string {
	return "http://localhost:9090/api/v1/query"
}

type targetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
