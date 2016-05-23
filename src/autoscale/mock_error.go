package autoscale

import "github.com/stretchr/testify/mock"

import _ "github.com/lib/pq"

// mockError is an autogenerated mock type for the error type
type mockError struct {
	mock.Mock
}

// Error provides a mock function with given fields:
func (_m *mockError) Error() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}