package api

import (
	"autoscale"
	"autoscale/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTemplates(t *testing.T) {
	ogTmpls := []autoscale.Template{
		{ID: 1},
		{ID: 2},
	}

	repo := &mocks.Repository{}
	repo.On("ListTemplates").Return(ogTmpls, nil)

	api := New(repo)

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
}

func TestGetTemplate(t *testing.T) {
	ogTmpl := autoscale.Template{ID: 1}

	repo := &mocks.Repository{}
	repo.On("GetTemplate", 1).Return(&ogTmpl, nil)

	api := New(repo)

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

	assert.Equal(t, 1, tmpl.ID)
}

func TestGetMissingTemplate(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetTemplate", 1).Return(nil, errors.New("boom"))

	api := New(repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/templates/1"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 404, res.StatusCode)
}

func TestCreateTemplate(t *testing.T) {
	repo := &mocks.Repository{}
	expectedTmpl := &autoscale.Template{
		Name:       "a-template",
		Region:     "dev0",
		Size:       "512mb",
		Image:      "ubuntu-14-04-x64",
		RawSSHKeys: "123,456,789",
		UserData:   "#userdata",
	}
	repo.On("CreateTemplate", expectedTmpl).Return(1, nil)

	api := New(repo)

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

	var tmpl autoscale.Template
	err = json.NewDecoder(res.Body).Decode(&tmpl)
	assert.NoError(t, err)

	assert.Equal(t, 1, tmpl.ID)
}

func TestListGroups(t *testing.T) {
	ogGroups := []autoscale.Group{
		{ID: "12345"},
		{ID: "6789"},
	}

	repo := &mocks.Repository{}
	repo.On("ListGroups").Return(ogGroups, nil)

	api := New(repo)

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
}

func TestGetGroup(t *testing.T) {
	ogGroup := autoscale.Group{ID: "abc"}

	repo := &mocks.Repository{}
	repo.On("GetGroup", "abc").Return(&ogGroup, nil)

	api := New(repo)

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
}

func TestGetMissingGroup(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetGroup", "1").Return(nil, errors.New("boom"))

	api := New(repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups/1"

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 404, res.StatusCode)
}

func TestCreateGroup(t *testing.T) {
	repo := &mocks.Repository{}
	expectedGroup := &autoscale.Group{
		BaseName:   "as",
		BaseSize:   3,
		MetricType: "load",
		TemplateID: 1,
	}
	repo.On("CreateGroup", expectedGroup).Return("abc", nil)

	api := New(repo)

	ts := httptest.NewServer(api.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/groups"

	req := []byte(`{
    "base_name": "as",
    "base_size": 3,
    "metric_type": "load",
    "template_id": 1
  }`)

	var buf bytes.Buffer
	_, err = buf.Write(req)
	assert.NoError(t, err)

	res, err := http.Post(u.String(), "application/json", &buf)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 201, res.StatusCode)

	var group autoscale.Group
	err = json.NewDecoder(res.Body).Decode(&group)
	assert.NoError(t, err)

	assert.Equal(t, "abc", group.ID)
}
