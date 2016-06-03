package autoscale

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestConvertGroupToJSON(t *testing.T) {
	m, err := NewFileLoad()
	require.NoError(t, err)

	vps := ValuePolicyScale(1, 10, 0.8, 2, 0.2, 1)
	vp, err := NewValuePolicy(vps)
	require.NoError(t, err)

	vp.vpd.WarmUpDuration = defaultWarmUpDuration

	g := Group{
		ID:         "12345",
		Name:       "group",
		BaseName:   "as",
		MetricType: "load",
		Metric:     m,
		PolicyType: "value",
		Policy:     vp,
		TemplateID: "a-template",
	}

	j, err := json.Marshal(&g)
	require.NoError(t, err)

	var newGroup Group
	err = json.Unmarshal(j, &newGroup)
	require.NoError(t, err)

	assert.Equal(t, g.ID, newGroup.ID)
	assert.Equal(t, g.Name, newGroup.Name)
	assert.Equal(t, g.BaseName, newGroup.BaseName)
	assert.Equal(t, g.TemplateID, newGroup.TemplateID)
	assert.Equal(t, g.PolicyType, newGroup.PolicyType)
	assert.Equal(t, g.MetricType, newGroup.MetricType)
	assert.Equal(t, g.Policy, newGroup.Policy)
	assert.Equal(t, g.Metric, newGroup.Metric)
}
