package autoscale

import "github.com/stretchr/testify/mock"

type MockRepository struct {
	mock.Mock
}

func (_m *MockRepository) CreateTemplate(tcr CreateTemplateRequest) (Template, error) {
	ret := _m.Called(tcr)

	var r0 Template
	if rf, ok := ret.Get(0).(func(CreateTemplateRequest) Template); ok {
		r0 = rf(tcr)
	} else {
		r0 = ret.Get(0).(Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(CreateTemplateRequest) error); ok {
		r1 = rf(tcr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) GetTemplate(name string) (Template, error) {
	ret := _m.Called(name)

	var r0 Template
	if rf, ok := ret.Get(0).(func(string) Template); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) ListTemplates() ([]Template, error) {
	ret := _m.Called()

	var r0 []Template
	if rf, ok := ret.Get(0).(func() []Template); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Template)
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
func (_m *MockRepository) DeleteTemplate(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRepository) CreateGroup(gcr CreateGroupRequest) (Group, error) {
	ret := _m.Called(gcr)

	var r0 Group
	if rf, ok := ret.Get(0).(func(CreateGroupRequest) Group); ok {
		r0 = rf(gcr)
	} else {
		r0 = ret.Get(0).(Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(CreateGroupRequest) error); ok {
		r1 = rf(gcr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) GetGroup(name string) (Group, error) {
	ret := _m.Called(name)

	var r0 Group
	if rf, ok := ret.Get(0).(func(string) Group); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) ListGroups() ([]Group, error) {
	ret := _m.Called()

	var r0 []Group
	if rf, ok := ret.Get(0).(func() []Group); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Group)
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
func (_m *MockRepository) DeleteGroup(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRepository) SaveGroup(group Group) error {
	ret := _m.Called(group)

	var r0 error
	if rf, ok := ret.Get(0).(func(Group) error); ok {
		r0 = rf(group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
