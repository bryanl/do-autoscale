package mocks

import (
	"autoscale"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (_m *Repository) SaveTemplate(t *autoscale.Template) (int, error) {
	ret := _m.Called(t)

	var r0 int
	if rf, ok := ret.Get(0).(func(*autoscale.Template) int); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*autoscale.Template) error); ok {
		r1 = rf(t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) GetTemplate(id int) (*autoscale.Template, error) {
	ret := _m.Called(id)

	var r0 *autoscale.Template
	if rf, ok := ret.Get(0).(func(int) *autoscale.Template); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*autoscale.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
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
func (_m *Repository) CreateGroup(t *autoscale.Group) (string, error) {
	ret := _m.Called(t)

	var r0 string
	if rf, ok := ret.Get(0).(func(*autoscale.Group) string); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*autoscale.Group) error); ok {
		r1 = rf(t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Repository) GetGroup(id string) (*autoscale.Group, error) {
	ret := _m.Called(id)

	var r0 *autoscale.Group
	if rf, ok := ret.Get(0).(func(string) *autoscale.Group); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*autoscale.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
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
