package api

import (
	"autoscale"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

type apiTestMocks struct {
	templateResource    *MockResource
	groupResource       *MockResource
	userConfigResource  *MockResource
	groupConfigResource *MockResource
}
type apiTestFn func(ctx context.Context, mocks *apiTestMocks, u *url.URL)

func withAPITest(t *testing.T, fn apiTestFn) {
	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	api := New(ctx, repo)

	mocks := &apiTestMocks{
		templateResource:    &MockResource{},
		groupResource:       &MockResource{},
		userConfigResource:  &MockResource{},
		groupConfigResource: &MockResource{},
	}

	api.templateResourceFactory = func() Resource { return mocks.templateResource }
	api.groupResourceFactory = func() Resource { return mocks.groupResource }
	api.userConfigResourceFactory = func() Resource { return mocks.userConfigResource }
	api.groupConfigResourceFactory = func() Resource { return mocks.groupConfigResource }

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	fn(ctx, mocks, u)

	assert.True(t, repo.AssertExpectations(t))
	assert.True(t, mocks.templateResource.AssertExpectations(t))
	assert.True(t, mocks.groupResource.AssertExpectations(t))
	assert.True(t, mocks.userConfigResource.AssertExpectations(t))
	assert.True(t, mocks.groupConfigResource.AssertExpectations(t))
}

func TestListTemplates(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		ogTmpls := []*autoscale.Template{
			{ID: "1", SSHKeys: autoscale.SSHKeys{{ID: 1}}},
			{ID: "2", SSHKeys: autoscale.SSHKeys{{ID: 1}}},
		}

		resp := newResponse(ogTmpls, 200)
		mocks.templateResource.On("FindAll", mock.Anything).Return(resp, nil)

		u.Path = "/api/templates"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)

		var templates []autoscale.Template
		err = json.NewDecoder(res.Body).Decode(&templates)
		require.NoError(t, err)

		require.Len(t, templates, 2)
	})
}

func TestDeleteTemplate(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		resp := newResponse(nil, 204)
		mocks.templateResource.On("Delete", mock.Anything, "1").Return(resp, nil)

		u.Path = "/api/templates/1"

		req, err := http.NewRequest("DELETE", u.String(), nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, 204, res.StatusCode)

	})
}

func TestGetTemplate(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		ogTmpl := &autoscale.Template{ID: "1"}

		resp := newResponse(ogTmpl, 200)
		mocks.templateResource.On("FindOne", mock.Anything, "1").Return(resp, nil)

		u.Path = "/api/templates/1"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)

		var tmpl autoscale.Template
		err = json.NewDecoder(res.Body).Decode(&tmpl)
		require.NoError(t, err)

		require.Equal(t, "1", tmpl.ID)

	})
}

func TestGetMissingTemplate(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		resp := newResponse(nil, 404)
		mocks.templateResource.On("FindOne", mock.Anything, "1").Return(resp, nil)

		u.Path = "/api/templates/1"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 404, res.StatusCode)
	})
}

func TestCreateTemplate(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {

		newTmpl := autoscale.Template{
			Name:   "a-template",
			Region: "dev0",
			Size:   "512mb",
			Image:  "ubuntu-14-04-x64",
			SSHKeys: autoscale.SSHKeys{
				{ID: 123}, {ID: 456}, {ID: 789},
			},
			UserData: "#userdata",
		}

		tmpl := autoscale.Template{
			ID:     "1",
			Name:   "a-template",
			Region: "dev0",
			Size:   "512mb",
			Image:  "ubuntu-14-04-x64",
			SSHKeys: autoscale.SSHKeys{
				{ID: 123}, {ID: 456}, {ID: 789},
			},
			UserData: "#userdata",
		}

		resp := newResponse(tmpl, 201)
		mocks.templateResource.On("Create", mock.Anything, newTmpl).Return(resp, nil)

		u.Path = "/api/templates"

		req := []byte(`{
    "template":{
      "name": "a-template",
      "region": "dev0",
      "size": "512mb",
      "image": "ubuntu-14-04-x64",
      "sshKeys": [{"id":123}, {"id":456}, {"id":789}],
      "userData": "#userdata"
    }
  }`)

		var buf bytes.Buffer
		_, err := buf.Write(req)
		require.NoError(t, err)

		res, err := http.Post(u.String(), "application/json", &buf)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 201, res.StatusCode)

		err = json.NewDecoder(res.Body).Decode(&newTmpl)
		require.NoError(t, err)

		require.Equal(t, tmpl, newTmpl)

	})
}

func TestListGroups(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		ogGroups := []autoscale.Group{
			{ID: "12345", PolicyType: "value", MetricType: "load", Policy: &autoscale.ValuePolicy{}, Metric: &autoscale.FileLoad{}},
			{ID: "6789", PolicyType: "value", MetricType: "load", Policy: &autoscale.ValuePolicy{}, Metric: &autoscale.FileLoad{}},
		}

		resp := newResponse(ogGroups, 200)
		mocks.groupResource.On("FindAll", mock.Anything).Return(resp, nil)

		u.Path = "/api/groups"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)

	})
}

func TestDeleteGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		resp := newResponse(nil, 204)
		mocks.groupResource.On("Delete", mock.Anything, "abc").Return(resp, nil)

		u.Path = "/api/groups/abc"

		req, err := http.NewRequest("DELETE", u.String(), nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, 204, res.StatusCode)

	})
}

func TestUpdateGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		ogGroup := autoscale.Group{
			ID:         "abc",
			MetricType: "load",
			PolicyType: "value",
		}

		resp := newResponse(&ogGroup, 200)
		mocks.groupResource.On("Update", mock.Anything, mock.AnythingOfType("autoscale.Group")).Return(resp, nil)

		u.Path = "/api/groups/abc"

		j := `
    {
      "group": {
        "policy": {
          "scale_up_value": 6
        },
        "metricType": "load",
        "policyType": "value"
      }
    }`

		var buf bytes.Buffer
		_, err := buf.WriteString(j)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", u.String(), &buf)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, 200, res.StatusCode)
	})
}

func TestGetGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		ogGroup := autoscale.Group{ID: "abc"}

		resp := newResponse(&ogGroup, 200)
		mocks.groupResource.On("FindOne", mock.Anything, "abc").Return(resp, nil)

		u.Path = "/api/groups/abc"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)

		var group autoscale.Group
		err = json.NewDecoder(res.Body).Decode(&group)
		require.NoError(t, err)

		require.Equal(t, "abc", group.ID)
	})
}

func TestGetMissingGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		resp := newResponse(nil, 404)
		mocks.templateResource.On("FindOne", mock.Anything, "1").Return(resp, nil)

		u.Path = "/api/groups/1"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 404, res.StatusCode)

	})
}

func TestCreateGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, mocks *apiTestMocks, u *url.URL) {
		group := autoscale.Group{
			Name:         "group",
			BaseName:     "as",
			MetricType:   "load",
			PolicyType:   "value",
			TemplateName: "a-template",
		}

		newGroup := autoscale.Group{
			ID:           "1",
			Name:         "group",
			BaseName:     "as",
			MetricType:   "load",
			PolicyType:   "value",
			TemplateName: "a-template",
		}

		resp := newResponse(newGroup, 201)
		mocks.groupResource.On("Create", mock.Anything, group).Return(resp, nil)

		u.Path = "/api/groups"

		req := []byte(`{
    "name": "group",
    "base_name": "as",
    "metric_type": "load",
    "policy_type": "value",
    "template_name": "a-template"
  }`)

		var buf bytes.Buffer
		_, err := buf.Write(req)
		require.NoError(t, err)

		res, err := http.Post(u.String(), "application/json", &buf)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 201, res.StatusCode)

		b, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		var g autoscale.Group
		err = jsonapi.Unmarshal(b, &g)
		require.NoError(t, err)

		require.Equal(t, newGroup, g)
	})
}
