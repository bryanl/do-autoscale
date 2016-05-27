package api

import (
	"autoscale"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

type apiTestFn func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL)

func withAPITest(t *testing.T, fn apiTestFn) {
	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	fn(ctx, repo, u)

	repo.AssertExpectations(t)
}

func TestListTemplates(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		ogTmpls := []autoscale.Template{
			{ID: "1"},
			{ID: "2"},
		}

		repo.On("ListTemplates", ctx).Return(ogTmpls, nil)

		u.Path = "/api/templates"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)

		var resp autoscale.TemplatesResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)

		require.Len(t, resp.Templates, 2)
	})
}

func TestDeleteTemplate(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		repo.On("DeleteTemplate", ctx, "1").Return(nil)

		u.Path = "/api/templates/1"

		req, err := http.NewRequest("DELETE", u.String(), nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, 204, res.StatusCode)

	})
}

func TestGetTemplate(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		ogTmpl := autoscale.Template{ID: "1"}

		repo.On("GetTemplate", ctx, "1").Return(ogTmpl, nil)

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
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		repo.On("GetTemplate", ctx, "1").Return(autoscale.Template{}, errors.New("boom"))

		u.Path = "/api/templates/1"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 404, res.StatusCode)
	})
}

func TestCreateTemplate(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {

		ctr := autoscale.CreateTemplateRequest{
			Options: autoscale.TemplateOptions{
				Name:     "a-template",
				Region:   "dev0",
				Size:     "512mb",
				Image:    "ubuntu-14-04-x64",
				SSHKeys:  []string{"123", "456", "789"},
				UserData: "#userdata",
			},
		}

		tmpl := autoscale.Template{
			ID:       "1",
			Name:     "a-template",
			Region:   "dev0",
			Size:     "512mb",
			Image:    "ubuntu-14-04-x64",
			SSHKeys:  []string{"123", "456", "789"},
			UserData: "#userdata",
		}

		repo.On("CreateTemplate", ctx, ctr).Return(tmpl, nil)

		u.Path = "/api/templates"

		req := []byte(`{
    "template":{
      "name": "a-template",
      "region": "dev0",
      "size": "512mb",
      "image": "ubuntu-14-04-x64",
      "ssh_keys": ["123", "456", "789"],
      "user_data": "#userdata"
    }
  }`)

		var buf bytes.Buffer
		_, err := buf.Write(req)
		require.NoError(t, err)

		res, err := http.Post(u.String(), "application/json", &buf)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 201, res.StatusCode)

		var newTmpl autoscale.Template
		err = json.NewDecoder(res.Body).Decode(&newTmpl)
		require.NoError(t, err)

		require.Equal(t, tmpl, newTmpl)

	})
}

func TestListGroups(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		ogGroups := []autoscale.Group{
			{ID: "12345", PolicyType: "value", MetricType: "load", Policy: &autoscale.ValuePolicy{}, Metric: &autoscale.FileLoad{}},
			{ID: "6789", PolicyType: "value", MetricType: "load", Policy: &autoscale.ValuePolicy{}, Metric: &autoscale.FileLoad{}},
		}

		repo.On("ListGroups", ctx).Return(ogGroups, nil)

		u.Path = "/api/groups"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 200, res.StatusCode)

	})
}

func TestDeleteGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		repo.On("DeleteGroup", ctx, "abc").Return(nil)

		u.Path = "/api/groups/abc"

		req, err := http.NewRequest("DELETE", u.String(), nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, 204, res.StatusCode)

	})
}

func TestUpdateGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		ogGroup := autoscale.Group{ID: "abc"}
		ogGroupUpdated := autoscale.Group{ID: "abc"}

		repo.On("GetGroup", ctx, "abc").Return(ogGroup, nil)
		repo.On("SaveGroup", ctx, ogGroupUpdated).Return(nil)

		u.Path = "/api/groups/abc"

		j := `{
    "base_size": 6
  }`

		var buf bytes.Buffer
		_, err := buf.WriteString(j)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", u.String(), &buf)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, 200, res.StatusCode)
	})
}

func TestGetGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		ogGroup := autoscale.Group{ID: "abc"}

		repo.On("GetGroup", ctx, "abc").Return(ogGroup, nil)

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
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {
		repo.On("GetGroup", ctx, "1").Return(autoscale.Group{}, errors.New("missing"))

		u.Path = "/api/groups/1"

		res, err := http.Get(u.String())
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, 404, res.StatusCode)

	})
}

func TestCreateGroup(t *testing.T) {
	withAPITest(t, func(ctx context.Context, repo *autoscale.MockRepository, u *url.URL) {

		cgr := autoscale.CreateGroupRequest{
			Name:         "group",
			BaseName:     "as",
			MetricType:   "load",
			PolicyType:   "value",
			TemplateName: "a-template",
		}

		group := autoscale.Group{
			ID:           "1",
			Name:         "group",
			BaseName:     "as",
			MetricType:   "load",
			PolicyType:   "value",
			TemplateName: "a-template",
		}

		repo.On("CreateGroup", ctx, cgr).Return(group, nil)

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

		var newGroup autoscale.Group
		err = json.NewDecoder(res.Body).Decode(&newGroup)
		require.NoError(t, err)

		require.Equal(t, newGroup, group)

	})
}
