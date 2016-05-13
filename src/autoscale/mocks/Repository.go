package mocks

import (
	"autoscale"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (_m *Repository) CreateTemplate(tcr autoscale.CreateTemplateRequest) (autoscale.Template, error) {
	ret := _m.Called(tcr)

	var r0 autoscale.Template
	if rf, ok := ret.Get(0).(func(autoscale.CreateTemplateRequest) autoscale.Template); ok {
		r0 = rf(tcr)
	} else {
		r0 = ret.Get(0).(autoscale.Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(autoscale.CreateTemplateRequest) error); ok {
		r1 = rf(tcr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) GetTemplate(name string) (autoscale.Template, error) {
	ret := _m.Called(name)

	var r0 autoscale.Template
	if rf, ok := ret.Get(0).(func(string) autoscale.Template); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(autoscale.Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) ListTemplates() ([]autoscale.Template, error) {
	ret := _m.Called()

	var r0 []autoscale.Template
	if rf, ok := ret.Get(0).(func() []autoscale.Template); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]autoscale.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) DeleteTemplate(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Repository) CreateGroup(gcr autoscale.CreateGroupRequest) (autoscale.Group, error) {
	ret := _m.Called(gcr)

	var r0 autoscale.Group
	if rf, ok := ret.Get(0).(func(autoscale.CreateGroupRequest) autoscale.Group); ok {
		r0 = rf(gcr)
	} else {
		r0 = ret.Get(0).(autoscale.Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(autoscale.CreateGroupRequest) error); ok {
		r1 = rf(gcr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) GetGroup(name string) (autoscale.Group, error) {
	ret := _m.Called(name)

	var r0 autoscale.Group
	if rf, ok := ret.Get(0).(func(string) autoscale.Group); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(autoscale.Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) ListGroups() ([]autoscale.Group, error) {
	ret := _m.Called()

	var r0 []autoscale.Group
	if rf, ok := ret.Get(0).(func() []autoscale.Group); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]autoscale.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) DeleteGroup(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Repository) SaveGroup(group autoscale.Group) error {
	ret := _m.Called(group)

	var r0 error
	if rf, ok := ret.Get(0).(func(autoscale.Group) error); ok {
		r0 = rf(group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
