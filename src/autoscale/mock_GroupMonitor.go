package autoscale

import "github.com/stretchr/testify/mock"

type MockGroupMonitor struct {
	mock.Mock
}

// Start provides a mock function with given fields: fn
func (_m *MockGroupMonitor) Start(fn AfterMonitorFn) {
	_m.Called(fn)
}

// Stop provides a mock function with given fields:
func (_m *MockGroupMonitor) Stop() {
	_m.Called()
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
