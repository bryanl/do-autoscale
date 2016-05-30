package api

import "github.com/stretchr/testify/mock"

import "golang.org/x/net/context"

type MockResource struct {
	mock.Mock
}

func (_m *MockResource) FindOne(c context.Context, id string) (Response, error) {
	ret := _m.Called(c, id)

	var r0 Response
	if rf, ok := ret.Get(0).(func(context.Context, string) Response); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Get(0).(Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockResource) Create(c context.Context, obj interface{}) (Response, error) {
	ret := _m.Called(c, obj)

	var r0 Response
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) Response); ok {
		r0 = rf(c, obj)
	} else {
		r0 = ret.Get(0).(Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(c, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockResource) Delete(c context.Context, id string) (Response, error) {
	ret := _m.Called(c, id)

	var r0 Response
	if rf, ok := ret.Get(0).(func(context.Context, string) Response); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Get(0).(Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockResource) Update(c context.Context, obj interface{}) (Response, error) {
	ret := _m.Called(c, obj)

	var r0 Response
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) Response); ok {
		r0 = rf(c, obj)
	} else {
		r0 = ret.Get(0).(Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(c, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockResource) FindAll(c context.Context) (Response, error) {
	ret := _m.Called(c)

	var r0 Response
	if rf, ok := ret.Get(0).(func(context.Context) Response); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
