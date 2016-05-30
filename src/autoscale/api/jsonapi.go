package api

import "github.com/manyminds/api2go"

type jsonAPIError struct {
	ID     string              `json:"-"`
	Links  *api2go.ErrorLinks  `json:"links,omitempty"`
	Status string              `json:"status,omitempty"`
	Code   string              `json:"code,omitempty"`
	Title  string              `json:"title,omitempty"`
	Detail string              `json:"detail,omitempty"`
	Source *api2go.ErrorSource `json:"source,omitempty"`
	Meta   interface{}         `json:"meta,omitempty"`
}

func (j *jsonAPIError) GetID() string {
	return j.ID
}

func (j *jsonAPIError) GetName() string {
	return "errors"
}

type jsonAPIErrors []*jsonAPIError
