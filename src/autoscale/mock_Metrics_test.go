package autoscale

import "github.com/stretchr/testify/mock"

type MockMetrics struct {
	mock.Mock
}

func (_m *MockMetrics) Value(groupName string) (float64, error) {
	ret := _m.Called(groupName)

	var r0 float64
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(groupName)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(groupName)
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
