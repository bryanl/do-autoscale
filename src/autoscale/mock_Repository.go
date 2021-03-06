package autoscale

import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"

type MockRepository struct {
	mock.Mock
}

func (_m *MockRepository) CreateTemplate(ctx context.Context, t Template) (*Template, error) {
	ret := _m.Called(ctx, t)

	var r0 *Template
	if rf, ok := ret.Get(0).(func(context.Context, Template) *Template); ok {
		r0 = rf(ctx, t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, Template) error); ok {
		r1 = rf(ctx, t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) GetTemplate(ctx context.Context, name string) (*Template, error) {
	ret := _m.Called(ctx, name)

	var r0 *Template
	if rf, ok := ret.Get(0).(func(context.Context, string) *Template); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) ListTemplates(ctx context.Context) ([]Template, error) {
	ret := _m.Called(ctx)

	var r0 []Template
	if rf, ok := ret.Get(0).(func(context.Context) []Template); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) DeleteTemplate(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRepository) CreateGroup(ctx context.Context, g Group) (*Group, error) {
	ret := _m.Called(ctx, g)

	var r0 *Group
	if rf, ok := ret.Get(0).(func(context.Context, Group) *Group); ok {
		r0 = rf(ctx, g)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, Group) error); ok {
		r1 = rf(ctx, g)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) GetGroup(ctx context.Context, name string) (*Group, error) {
	ret := _m.Called(ctx, name)

	var r0 *Group
	if rf, ok := ret.Get(0).(func(context.Context, string) *Group); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) ListGroups(ctx context.Context) ([]Group, error) {
	ret := _m.Called(ctx)

	var r0 []Group
	if rf, ok := ret.Get(0).(func(context.Context) []Group); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) DeleteGroup(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRepository) SaveGroup(ctx context.Context, group Group) error {
	ret := _m.Called(ctx, group)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Group) error); ok {
		r0 = rf(ctx, group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRepository) AddGroupStatus(ctx context.Context, g GroupStatus) error {
	ret := _m.Called(ctx, g)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, GroupStatus) error); ok {
		r0 = rf(ctx, g)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockRepository) ListGroupStatus(ctx context.Context) ([]GroupStatus, error) {
	ret := _m.Called(ctx)

	var r0 []GroupStatus
	if rf, ok := ret.Get(0).(func(context.Context) []GroupStatus); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]GroupStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) GetGroupStatus(ctx context.Context, groupID string) (*GroupStatus, error) {
	ret := _m.Called(ctx, groupID)

	var r0 *GroupStatus
	if rf, ok := ret.Get(0).(func(context.Context, string) *GroupStatus); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) GetGroupHistory(ctx context.Context, groupID string, tr TimeRange) ([]GroupStatus, error) {
	ret := _m.Called(ctx, groupID, tr)

	var r0 []GroupStatus
	if rf, ok := ret.Get(0).(func(context.Context, string, TimeRange) []GroupStatus); ok {
		r0 = rf(ctx, groupID, tr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]GroupStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, TimeRange) error); ok {
		r1 = rf(ctx, groupID, tr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRepository) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
