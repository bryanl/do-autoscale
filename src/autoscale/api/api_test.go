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

	ts := httptest.NewServer(api.r)
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

	ts := httptest.NewServer(api.r)
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

	ts := httptest.NewServer(api.r)
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
		Region:     "dev0",
		Size:       "512mb",
		Image:      "ubuntu-14-04-x64",
		RawSSHKeys: "123,456,789",
		UserData:   "#userdata",
	}
	repo.On("SaveTemplate", expectedTmpl).Return(1, nil)

	api := New(repo)

	ts := httptest.NewServer(api.r)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	u.Path = "/templates"

	req := []byte(`{
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
