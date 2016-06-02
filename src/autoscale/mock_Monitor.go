package autoscale

import "github.com/stretchr/testify/mock"

type MockMonitor struct {
	mock.Mock
}

func (_m *MockMonitor) Start(s Scheduler) {
	_m.Called(s)
}
func (_m *MockMonitor) Stop() {
	_m.Called()
}
