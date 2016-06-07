package autoscale

import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"

type MockMetrics struct {
	mock.Mock
}

func (_m *MockMetrics) Measure(ctx context.Context, groupName string) (float64, error) {
	ret := _m.Called(ctx, groupName)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context, string) float64); ok {
		r0 = rf(ctx, groupName)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockMetrics) Update(groupName string, resourceAllocations []ResourceAllocation) error {
	ret := _m.Called(groupName, resourceAllocations)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []ResourceAllocation) error); ok {
		r0 = rf(groupName, resourceAllocations)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockMetrics) Config() MetricConfig {
	ret := _m.Called()

	var r0 MetricConfig
	if rf, ok := ret.Get(0).(func() MetricConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(MetricConfig)
	}

	return r0
}
func (_m *MockMetrics) Values(ctx context.Context, groupName string, rangeLength TimeRange) ([]TimeSeries, error) {
	ret := _m.Called(ctx, groupName, rangeLength)

	var r0 []TimeSeries
	if rf, ok := ret.Get(0).(func(context.Context, string, TimeRange) []TimeSeries); ok {
		r0 = rf(ctx, groupName, rangeLength)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]TimeSeries)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, TimeRange) error); ok {
		r1 = rf(ctx, groupName, rangeLength)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockMetrics) Remove(ctx context.Context, groupID string) error {
	ret := _m.Called(ctx, groupID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, groupID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
