package autoscale

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

func TestCheck(t *testing.T) {
	ogFactory := ResourceManagerFactory
	ogDefaultConfig := DefaultConfig
	defer func() {
		ResourceManagerFactory = ogFactory
		DefaultConfig = ogDefaultConfig
	}()

	tmpPath, err := ioutil.TempDir("", "autoscaler")
	require.NoError(t, err)
	defer os.RemoveAll(tmpPath)

	DefaultConfig[OptionFileLoadPath] = tmpPath

	ctx := context.Background()
	RegisterOfflineMetrics(ctx)

	ResourceManagerFactory = func(g *Group) (ResourceManager, error) {
		r := NewLocalResource(ctx)
		r.(*LocalResource).count = 3
		return r, nil
	}

	policy, err := NewValuePolicy(ValuePolicyScale(
		1, 10, 0.8, 2, 0.2, 1,
	))
	require.NoError(t, err)

	repo := &MockRepository{}

	cases := []struct {
		currentLoad string
		delta       int
	}{
		{currentLoad: "0.5", delta: 0},
		{currentLoad: "0.1", delta: -1},
		{currentLoad: "0.8", delta: 2},
	}

	for _, c := range cases {
		group := &Group{
			ID:         "id",
			Name:       "test-group",
			MetricType: "load",
			PolicyType: "value",
			Policy:     policy,
		}
		repo.On("GetGroup", mock.Anything, "id").Return(group, nil)

		metricPath := filepath.Join(tmpPath, group.Name)
		err = ioutil.WriteFile(metricPath, []byte(c.currentLoad), 0600)
		require.NoError(t, err)

		check := NewCheck(repo)

		as := check.Perform(ctx, "id")

		assert.NoError(t, as.Err)
		assert.Equal(t, c.delta, as.Delta, fmt.Sprintf("delta did not match"))
	}
}
