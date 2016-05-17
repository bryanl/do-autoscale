package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Sirupsen/logrus"
)

type prometheusResponse struct {
	Status string         `json:"status"`
	Data   prometheusData `json:"data"`
}

func (pr *prometheusResponse) Value(index int) (float64, error) {
	if len(pr.Data.Results) < index+1 {
		return 0, nil
	}

	return pr.Data.Results[index].Value()
}

type prometheusData struct {
	ResultType string             `json:"data"`
	Results    []prometheusResult `json:"result"`
}

type prometheusResult struct {
	RawValue []interface{} `json:"value"`
}

func (pr *prometheusResult) Value() (float64, error) {
	return strconv.ParseFloat(pr.RawValue[1].(string), 64)
}

// PrometheusLoad based on promeetheus metrics.
type PrometheusLoad struct {
}

// NewPrometheusLoad creates an instance of PrometheusLoad.
func NewPrometheusLoad() (*PrometheusLoad, error) {
	return &PrometheusLoad{}, nil
}

var _ Metrics = (*PrometheusLoad)(nil)

// Value returns the average load for an entire group.
func (l *PrometheusLoad) Value(groupName string) (float64, error) {
	u, err := url.Parse(l.queryURL())
	if err != nil {
		return 0, err
	}

	values := u.Query()
	q := fmt.Sprintf(`avg(node_load1{group="%s"})`, groupName)
	values.Set("query", q)
	u.RawQuery = values.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var pr prometheusResponse
	if err := json.NewDecoder(res.Body).Decode(&pr); err != nil {
		return 0, err
	}

	// we only asked for one thing, so it must be the first result
	return pr.Value(0)
}

// Update updates the prometheus config for a group.
func (l *PrometheusLoad) Update(groupName string, resourceAllocations []ResourceAllocation) error {
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

	path := fmt.Sprintf("/Users/bryan/Development/do/projects/prometheus/configs/%s.json",
		groupName)

	b, err := json.Marshal(&targetGroups)
	if err != nil {
		return err
	}

	logrus.WithField("path", path).Info("writing prometheus target file")

	return ioutil.WriteFile(path, b, 0644)
}

func (l *PrometheusLoad) queryURL() string {
	return "http://localhost:9090/api/v1/query"
}

type targetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
