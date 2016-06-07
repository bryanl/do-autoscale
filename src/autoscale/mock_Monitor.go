package autoscale

import "github.com/stretchr/testify/mock"

type MockMonitor struct {
	mock.Mock
}

func (_m *MockMonitor) Start(_a0 *SchedulerStatus) {
	_m.Called(_a0)
}
func (_m *MockMonitor) Stop() {
	_m.Called()
}
