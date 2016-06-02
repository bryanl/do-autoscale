package autoscale

import "github.com/stretchr/testify/mock"

type MockRunList struct {
	mock.Mock
}

func (_m *MockRunList) Add(groupID string) error {
	ret := _m.Called(groupID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(groupID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRunList) Remove(groupID string) error {
	ret := _m.Called(groupID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(groupID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRunList) IsRunning(groupID string) bool {
	ret := _m.Called(groupID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(groupID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
func (_m *MockRunList) List() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
func (_m *MockRunList) Reset() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
