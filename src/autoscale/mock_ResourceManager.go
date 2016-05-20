package autoscale

import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"

type MockResourceManager struct {
	mock.Mock
}

func (_m *MockResourceManager) Count() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockResourceManager) Scale(ctx context.Context, g Group, byN int, repo Repository) error {
	ret := _m.Called(ctx, g, byN, repo)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Group, int, Repository) error); ok {
		r0 = rf(ctx, g, byN, repo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockResourceManager) Allocated() ([]ResourceAllocation, error) {
	ret := _m.Called()

	var r0 []ResourceAllocation
	if rf, ok := ret.Get(0).(func() []ResourceAllocation); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ResourceAllocation)
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
