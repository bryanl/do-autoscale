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

// ResourceAllocation is information about an allocated resource.
type ResourceAllocation struct {
	Name    string
	Address string
}

// Available are the available metrics.
func Available() map[string]Gen {
	return map[string]Gen{
		"load": func(*Config) Metrics { return &Load{} },
	}
}

// Gen is a generator for metrics types.
type Gen func(*Config) Metrics

// Config is instance configuration for metrics.
type Config map[string]interface{}

// Metrics pull metrics for a autoscaler.
type Metrics interface {
	Value(groupName string) (float64, error)
	Update(groupName string, resourceAllocations []ResourceAllocation) error
}

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

// Load based on promeetheus metrics.
type Load struct {
}

var _ Metrics = (*Load)(nil)

// Value returns the average load for an entire group.
func (l *Load) Value(groupName string) (float64, error) {
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
func (l *Load) Update(groupName string, resourceAllocations []ResourceAllocation) error {
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

func (l *Load) queryURL() string {
	return "http://localhost:9090/api/v1/query"
}

type targetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
