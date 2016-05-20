package autoscale

import "github.com/stretchr/testify/mock"

type MockMetricNotifier struct {
	mock.Mock
}

func (_m *MockMetricNotifier) MetricNotify() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
