package autoscale

import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"

type MockGroupAction struct {
	mock.Mock
}

func (_m *MockGroupAction) Scale(ctx context.Context, groupID string) *ActionStatus {
	ret := _m.Called(ctx, groupID)

	var r0 *ActionStatus
	if rf, ok := ret.Get(0).(func(context.Context, string) *ActionStatus); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ActionStatus)
		}
	}

	return r0
}
func (_m *MockGroupAction) Disable(ctx context.Context, groupID string) *ActionStatus {
	ret := _m.Called(ctx, groupID)

	var r0 *ActionStatus
	if rf, ok := ret.Get(0).(func(context.Context, string) *ActionStatus); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ActionStatus)
		}
	}

	return r0
}
