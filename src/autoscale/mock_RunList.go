package autoscale

import "github.com/stretchr/testify/mock"

type MockRunList struct {
	mock.Mock
}

// Add provides a mock function with given fields: groupName
func (_m *MockRunList) Add(groupName string) error {
	ret := _m.Called(groupName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(groupName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: groupName
func (_m *MockRunList) Remove(groupName string) error {
	ret := _m.Called(groupName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(groupName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsRunning provides a mock function with given fields: groupName
func (_m *MockRunList) IsRunning(groupName string) (bool, error) {
	ret := _m.Called(groupName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(groupName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(groupName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reset provides a mock function with given fields:
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
