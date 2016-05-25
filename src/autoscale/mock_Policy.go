package autoscale

import "github.com/stretchr/testify/mock"

import "time"

type MockPolicy struct {
	mock.Mock
}

// CalculateSize provides a mock function with given fields: resourceCount, value
func (_m *MockPolicy) CalculateSize(resourceCount int, value float64) int {
	ret := _m.Called(resourceCount, value)

	var r0 int
	if rf, ok := ret.Get(0).(func(int, float64) int); ok {
		r0 = rf(resourceCount, value)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// WarmUpPeriod provides a mock function with given fields:
func (_m *MockPolicy) WarmUpPeriod() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// Config provides a mock function with given fields:
func (_m *MockPolicy) Config() PolicyConfig {
	ret := _m.Called()

	var r0 PolicyConfig
	if rf, ok := ret.Get(0).(func() PolicyConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(PolicyConfig)
	}

	return r0
}

// MarshalJSON provides a mock function with given fields:
func (_m *MockPolicy) MarshalJSON() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
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
