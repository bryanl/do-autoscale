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

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestListTemplates(t *testing.T) {
	ogTmpls := []autoscale.Template{
		{ID: "1"},
		{ID: "2"},
	}

	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	repo.On("ListTemplates", ctx).Return(ogTmpls, nil)

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/templates"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var tmpls []autoscale.Template
	err = json.NewDecoder(res.Body).Decode(&tmpls)
	assert.NoError(t, err)

	assert.Len(t, tmpls, 2)

	repo.AssertExpectations(t)
}

func TestDeleteTemplate(t *testing.T) {
	ctx := context.Background()

	repo := &autoscale.MockRepository{}
	repo.On("DeleteTemplate", ctx, "1").Return(nil)

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/templates/1"

	req, err := http.NewRequest("DELETE", u.String(), nil)
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, 204, res.StatusCode)

	repo.AssertExpectations(t)
}

func TestGetTemplate(t *testing.T) {
	ogTmpl := autoscale.Template{ID: "1"}

	ctx := context.Background()

	repo := &autoscale.MockRepository{}
	repo.On("GetTemplate", ctx, "1").Return(ogTmpl, nil)

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/templates/1"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var tmpl autoscale.Template
	err = json.NewDecoder(res.Body).Decode(&tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "1", tmpl.ID)

	repo.AssertExpectations(t)
}

func TestGetMissingTemplate(t *testing.T) {
	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	repo.On("GetTemplate", ctx, "1").Return(autoscale.Template{}, errors.New("boom"))

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/templates/1"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 404, res.StatusCode)

	repo.AssertExpectations(t)
}

func TestCreateTemplate(t *testing.T) {
	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	ctr := autoscale.CreateTemplateRequest{
		Name:     "a-template",
		Region:   "dev0",
		Size:     "512mb",
		Image:    "ubuntu-14-04-x64",
		SSHKeys:  []string{"123", "456", "789"},
		UserData: "#userdata",
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

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/templates"

	req := []byte(`{
    "name": "a-template",
    "region": "dev0",
    "size": "512mb",
    "image": "ubuntu-14-04-x64",
    "ssh_keys": ["123", "456", "789"],
    "user_data": "#userdata"
  }`)

	var buf bytes.Buffer
	_, err = buf.Write(req)
	assert.NoError(t, err)

	res, err := http.Post(u.String(), "application/json", &buf)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 201, res.StatusCode)

	var newTmpl autoscale.Template
	err = json.NewDecoder(res.Body).Decode(&newTmpl)
	assert.NoError(t, err)

	assert.Equal(t, tmpl, newTmpl)

	repo.AssertExpectations(t)
}

func TestListGroups(t *testing.T) {
	ogGroups := []autoscale.Group{
		{ID: "12345"},
		{ID: "6789"},
	}

	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	repo.On("ListGroups", ctx).Return(ogGroups, nil)

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var groups []autoscale.Group
	err = json.NewDecoder(res.Body).Decode(&groups)
	assert.NoError(t, err)

	assert.Len(t, groups, 2)

	repo.AssertExpectations(t)
}

func TestDeleteGroup(t *testing.T) {
	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	repo.On("DeleteGroup", ctx, "abc").Return(nil)

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups/abc"

	req, err := http.NewRequest("DELETE", u.String(), nil)
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, 204, res.StatusCode)

	repo.AssertExpectations(t)
}

func TestUpdateGroup(t *testing.T) {
	ogGroup := autoscale.Group{ID: "abc"}
	ogGroupUpdated := autoscale.Group{ID: "abc"}

	ctx := context.Background()
	repo := &autoscale.MockRepository{}
	repo.On("GetGroup", ctx, "abc").Return(ogGroup, nil)
	repo.On("SaveGroup", ctx, ogGroupUpdated).Return(nil)

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups/abc"

	j := `{
    "base_size": 6
  }`

	var buf bytes.Buffer
	_, err = buf.WriteString(j)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", u.String(), &buf)
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, 200, res.StatusCode)

	repo.AssertExpectations(t)
}

func TestGetGroup(t *testing.T) {
	ctx := context.Background()
	ogGroup := autoscale.Group{ID: "abc"}

	repo := &autoscale.MockRepository{}
	repo.On("GetGroup", ctx, "abc").Return(ogGroup, nil)

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups/abc"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var group autoscale.Group
	err = json.NewDecoder(res.Body).Decode(&group)
	assert.NoError(t, err)

	assert.Equal(t, "abc", group.ID)

	repo.AssertExpectations(t)
}

func TestGetMissingGroup(t *testing.T) {
	ctx := context.Background()

	repo := &autoscale.MockRepository{}
	repo.On("GetGroup", ctx, "1").Return(autoscale.Group{}, errors.New("boom"))

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups/1"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 404, res.StatusCode)

	repo.AssertExpectations(t)
}

func TestCreateGroup(t *testing.T) {
	ctx := context.Background()

	repo := &autoscale.MockRepository{}
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

	api := New(ctx, repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups"

	req := []byte(`{
    "name": "group",
    "base_name": "as",
    "metric_type": "load",
    "policy_type": "value",
    "template_name": "a-template"
  }`)

	var buf bytes.Buffer
	_, err = buf.Write(req)
	assert.NoError(t, err)

	res, err := http.Post(u.String(), "application/json", &buf)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 201, res.StatusCode)

	var newGroup autoscale.Group
	err = json.NewDecoder(res.Body).Decode(&newGroup)
	assert.NoError(t, err)

	assert.Equal(t, newGroup, group)

	repo.AssertExpectations(t)
}
