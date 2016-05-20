package autoscale

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplate_IsValid(t *testing.T) {
	cases := []struct {
		Name    string
		IsValid bool
	}{
		{Name: "1234", IsValid: true},
		{Name: "-1234", IsValid: false},
		{Name: "a-template", IsValid: true},
	}

	for _, c := range cases {
		tmpl := Template{
			Name: c.Name,
		}

		assert.Equal(t, c.IsValid, tmpl.IsValid())
	}
}

func TestGroup_IsValid(t *testing.T) {
	cases := []struct {
		Name    string
		IsValid bool
	}{
		{Name: "1234", IsValid: true},
		{Name: "-1234", IsValid: false},
		{Name: "a-template", IsValid: true},
	}

	for _, c := range cases {
		group := Group{
			Name: c.Name,
		}

		assert.Equal(t, c.IsValid, group.IsValid())
	}
}

func TestConvertCreateGroupRequestToGroup(t *testing.T) {
	policyJSON := []byte(`{
    "scale_up_value": 0.2,
    "scale_up_by": 2,
    "scale_down_value": 0.8,
    "scale_down_by": 1
  }`)

	metricJSON := []byte(`{
    "stats_dir": "/tmp"
  }`)

	cgr := CreateGroupRequest{
		Name:         "name",
		BaseName:     "base_name",
		TemplateName: "template-name",
		MetricType:   "load",
		Metric:       metricJSON,
		PolicyType:   "value",
		Policy:       policyJSON,
	}

	ctx := context.Background()

	group, err := cgr.ConvertToGroup(ctx)
	require.NoError(t, err)

	expected := &Group{
		Name:         "name",
		BaseName:     "base_name",
		TemplateName: "template-name",
		MetricType:   "load",
		Metric: &FileLoad{
			StatsDir: "/tmp",
		},
		PolicyType: "value",
		Policy: &ValuePolicy{
			ScaleUpValue:   0.2,
			ScaleUpBy:      2,
			ScaleDownValue: 0.8,
			ScaleDownBy:    1,
		},
	}

	require.Equal(t, expected.Name, group.Name)
	require.Equal(t, expected.BaseName, group.BaseName)
	require.Equal(t, expected.TemplateName, group.TemplateName)
	require.Equal(t, expected.MetricType, group.MetricType)
	require.Equal(t, expected.PolicyType, group.PolicyType)
	require.Equal(t, expected.Metric, group.Metric)
	require.Equal(t, expected.Policy, group.Policy)
}

func TestConvertCreateGroupRequestToGroup_WithDefaults(t *testing.T) {
	cgr := CreateGroupRequest{
		Name:         "name",
		BaseName:     "base_name",
		TemplateName: "template-name",
		MetricType:   "load",
		PolicyType:   "value",
	}

	ctx := context.Background()

	group, err := cgr.ConvertToGroup(ctx)
	require.NoError(t, err)

	expected := &Group{
		Name:         "name",
		BaseName:     "base_name",
		TemplateName: "template-name",
		MetricType:   "load",
		Metric: &FileLoad{
			StatsDir: "/tmp",
		},
		PolicyType: "value",
		Policy: &ValuePolicy{
			ScaleUpValue:   0.2,
			ScaleUpBy:      2,
			ScaleDownValue: 0.8,
			ScaleDownBy:    1,
		},
	}

	require.Equal(t, expected.Name, group.Name)
	require.Equal(t, expected.BaseName, group.BaseName)
	require.Equal(t, expected.TemplateName, group.TemplateName)
	require.Equal(t, expected.MetricType, group.MetricType)
	require.Equal(t, expected.PolicyType, group.PolicyType)

	vp := group.Policy.(*ValuePolicy)
	require.Equal(t, defaultValuePolicy, *vp)

	m := group.Metric.(*FileLoad)
	require.Equal(t, defaultLoadMetric, *m)
}
