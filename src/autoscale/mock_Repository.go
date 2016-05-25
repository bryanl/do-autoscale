package autoscale

import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"

type MockRepository struct {
	mock.Mock
}

// CreateTemplate provides a mock function with given fields: ctx, tcr
func (_m *MockRepository) CreateTemplate(ctx context.Context, tcr CreateTemplateRequest) (Template, error) {
	ret := _m.Called(ctx, tcr)

	var r0 Template
	if rf, ok := ret.Get(0).(func(context.Context, CreateTemplateRequest) Template); ok {
		r0 = rf(ctx, tcr)
	} else {
		r0 = ret.Get(0).(Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, CreateTemplateRequest) error); ok {
		r1 = rf(ctx, tcr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTemplate provides a mock function with given fields: ctx, name
func (_m *MockRepository) GetTemplate(ctx context.Context, name string) (Template, error) {
	ret := _m.Called(ctx, name)

	var r0 Template
	if rf, ok := ret.Get(0).(func(context.Context, string) Template); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTemplates provides a mock function with given fields: ctx
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

// DeleteTemplate provides a mock function with given fields: ctx, name
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

// CreateGroup provides a mock function with given fields: ctx, gcr
func (_m *MockRepository) CreateGroup(ctx context.Context, gcr CreateGroupRequest) (Group, error) {
	ret := _m.Called(ctx, gcr)

	var r0 Group
	if rf, ok := ret.Get(0).(func(context.Context, CreateGroupRequest) Group); ok {
		r0 = rf(ctx, gcr)
	} else {
		r0 = ret.Get(0).(Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, CreateGroupRequest) error); ok {
		r1 = rf(ctx, gcr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGroup provides a mock function with given fields: ctx, name
func (_m *MockRepository) GetGroup(ctx context.Context, name string) (Group, error) {
	ret := _m.Called(ctx, name)

	var r0 Group
	if rf, ok := ret.Get(0).(func(context.Context, string) Group); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGroups provides a mock function with given fields: ctx
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

// DeleteGroup provides a mock function with given fields: ctx, name
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

// SaveGroup provides a mock function with given fields: ctx, group
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

// Close provides a mock function with given fields:
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
