package api

import (
	"autoscale"
	"net/http"

	"golang.org/x/net/context"
)

// Resource is an autoscale resource.
type Resource interface {
	FindOne(c context.Context, id string) (Response, error)
	Create(c context.Context, obj interface{}) (Response, error)
	Delete(c context.Context, id string) (Response, error)
	Update(c context.Context, obj interface{}) (Response, error)
	FindAll(c context.Context) (Response, error)
}

// Response is a response from a call to an autoscale resource.
type Response interface {
	Result() interface{}
	StatusCode() int
}

type response struct {
	obj        interface{}
	statusCode int
}

var _ Response = (*response)(nil)

func newResponse(obj interface{}, statusCode int) *response {
	return &response{
		obj:        obj,
		statusCode: statusCode,
	}
}

func (r *response) Result() interface{} {
	return r.obj
}

func (r *response) StatusCode() int {
	return r.statusCode
}

type templateWrapper struct {
	Template autoscale.Template `json:"template"`
}

type templatesWrapper struct {
	Templates []autoscale.Template `json:"templates"`
}

type groupWrapper struct {
	Group autoscale.Group `json:"group"`
}

type groupsWrapper struct {
	Groups []autoscale.Group `json:"groups"`
}

type userConfigWrapper struct {
	UserConfig autoscale.UserConfig `json:"userConfigs"`
}

type groupConfigWrapper struct {
	GroupConfig autoscale.GroupConfig `json:"groupConfigs"`
}

type templateResource struct {
	repo autoscale.Repository
}

var _ Resource = (*templateResource)(nil)

func (r *templateResource) FindOne(c context.Context, id string) (Response, error) {
	template, err := r.repo.GetTemplate(c, id)
	if err != nil {
		if err == autoscale.ObjectMissingErr {
			return newResponse(nil, http.StatusNotFound), nil
		}

		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(templateWrapper{Template: *template}, http.StatusOK), nil
}

func (r *templateResource) Create(c context.Context, obj interface{}) (Response, error) {
	in, ok := obj.(autoscale.Template)
	if !ok {
		return newResponse(nil, http.StatusBadRequest), nil
	}

	template, err := r.repo.CreateTemplate(c, in)
	if err != nil {
		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(templateWrapper{Template: *template}, http.StatusCreated), nil
}

func (r *templateResource) Delete(c context.Context, id string) (Response, error) {
	if err := r.repo.DeleteTemplate(c, id); err != nil {
		return newResponse(nil, http.StatusNotFound), nil
	}

	return newResponse(nil, http.StatusNoContent), nil
}

func (r *templateResource) Update(c context.Context, obj interface{}) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *templateResource) FindAll(c context.Context) (Response, error) {
	templates, err := r.repo.ListTemplates(c)
	if err != nil {
		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(templatesWrapper{Templates: templates}, http.StatusOK), nil
}

type groupResource struct {
	repo autoscale.Repository
}

var _ Resource = (*groupResource)(nil)

func (r *groupResource) FindOne(c context.Context, id string) (Response, error) {
	group, err := r.repo.GetGroup(c, id)
	if err != nil {
		if err == autoscale.ObjectMissingErr {
			return newResponse(nil, http.StatusNotFound), nil
		}

		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(groupWrapper{Group: *group}, http.StatusOK), nil
}

func (r *groupResource) Create(c context.Context, obj interface{}) (Response, error) {
	in, ok := obj.(autoscale.Group)
	if !ok {
		return newResponse(nil, http.StatusBadRequest), nil
	}

	group, err := r.repo.CreateGroup(c, in)
	if err != nil {
		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(groupWrapper{Group: *group}, http.StatusCreated), nil
}

func (r *groupResource) Delete(c context.Context, id string) (Response, error) {
	if err := r.repo.DeleteGroup(c, id); err != nil {
		return newResponse(nil, http.StatusNotFound), nil
	}

	return newResponse(nil, http.StatusNoContent), nil
}

func (r *groupResource) Update(c context.Context, obj interface{}) (Response, error) {
	in, ok := obj.(autoscale.Group)
	if !ok {
		return newResponse(nil, http.StatusBadRequest), nil
	}

	err := r.repo.SaveGroup(c, in)
	if err != nil {
		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(groupWrapper{Group: in}, http.StatusOK), nil
}

func (r *groupResource) FindAll(c context.Context) (Response, error) {
	groups, err := r.repo.ListGroups(c)
	if err != nil {
		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(groupsWrapper{Groups: groups}, http.StatusOK), nil
}

type userConfigResource struct {
}

var _ Resource = (*userConfigResource)(nil)

func (r *userConfigResource) FindOne(c context.Context, id string) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *userConfigResource) Create(c context.Context, obj interface{}) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *userConfigResource) Delete(c context.Context, id string) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *userConfigResource) Update(c context.Context, obj interface{}) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *userConfigResource) FindAll(c context.Context) (Response, error) {
	client := autoscale.DOClientFactory()
	uc, err := autoscale.NewUserConfig(c, client)
	if err != nil {
		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(userConfigWrapper{UserConfig: *uc}, http.StatusOK), nil
}

type groupConfigResource struct {
	repo autoscale.Repository
}

var _ Resource = (*groupConfigResource)(nil)

func (r *groupConfigResource) FindOne(c context.Context, id string) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *groupConfigResource) Create(c context.Context, obj interface{}) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *groupConfigResource) Delete(c context.Context, id string) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *groupConfigResource) Update(c context.Context, obj interface{}) (Response, error) {
	return newResponse(nil, http.StatusNotImplemented), nil
}

func (r *groupConfigResource) FindAll(c context.Context) (Response, error) {
	client := autoscale.DOClientFactory()
	gc, err := autoscale.NewGroupConfig(c, client, r.repo)
	if err != nil {
		return newResponse(nil, http.StatusInternalServerError), nil
	}

	return newResponse(groupConfigWrapper{GroupConfig: *gc}, http.StatusOK), nil
}
