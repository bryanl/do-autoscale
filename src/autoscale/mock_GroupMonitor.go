package autoscale

import "github.com/stretchr/testify/mock"

type MockGroupMonitor struct {
	mock.Mock
}

// Start provides a mock function with given fields: newGroupFn
func (_m *MockGroupMonitor) Start(newGroupFn NewGroupFn) error {
	ret := _m.Called(newGroupFn)

	var r0 error
	if rf, ok := ret.Get(0).(func(NewGroupFn) error); ok {
		r0 = rf(newGroupFn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *MockGroupMonitor) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InRunList provides a mock function with given fields: groupName
func (_m *MockGroupMonitor) InRunList(groupName string) bool {
	ret := _m.Called(groupName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(groupName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
