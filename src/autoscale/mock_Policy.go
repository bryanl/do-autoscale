package autoscale

import "github.com/stretchr/testify/mock"

type MockPolicy struct {
	mock.Mock
}

func (_m *MockPolicy) Scale(mn MetricNotifier, resourceCount int, value float64) int {
	ret := _m.Called(mn, resourceCount, value)

	var r0 int
	if rf, ok := ret.Get(0).(func(MetricNotifier, int, float64) int); ok {
		r0 = rf(mn, resourceCount, value)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
