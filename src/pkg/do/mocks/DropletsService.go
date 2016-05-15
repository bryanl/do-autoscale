/*
Copyright 2016 The Doctl Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mocks

import (
	"pkg/do"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/mock"
)

type DropletsService struct {
	mock.Mock
}

// List provides a mock function with given fields:
func (_m *DropletsService) List() (do.Droplets, error) {
	ret := _m.Called()

	var r0 do.Droplets
	if rf, ok := ret.Get(0).(func() do.Droplets); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(do.Droplets)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0
func (_m *DropletsService) Get(_a0 int) (*do.Droplet, error) {
	ret := _m.Called(_a0)

	var r0 *do.Droplet
	if rf, ok := ret.Get(0).(func(int) *do.Droplet); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*do.Droplet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *DropletsService) Create(_a0 *godo.DropletCreateRequest, _a1 bool) (*do.Droplet, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *do.Droplet
	if rf, ok := ret.Get(0).(func(*godo.DropletCreateRequest, bool) *do.Droplet); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*do.Droplet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*godo.DropletCreateRequest, bool) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateMultiple provides a mock function with given fields: _a0
func (_m *DropletsService) CreateMultiple(_a0 *godo.DropletMultiCreateRequest) (do.Droplets, error) {
	ret := _m.Called(_a0)

	var r0 do.Droplets
	if rf, ok := ret.Get(0).(func(*godo.DropletMultiCreateRequest) do.Droplets); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(do.Droplets)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*godo.DropletMultiCreateRequest) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0
func (_m *DropletsService) Delete(_a0 int) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Kernels provides a mock function with given fields: _a0
func (_m *DropletsService) Kernels(_a0 int) (do.Kernels, error) {
	ret := _m.Called(_a0)

	var r0 do.Kernels
	if rf, ok := ret.Get(0).(func(int) do.Kernels); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(do.Kernels)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Snapshots provides a mock function with given fields: _a0
func (_m *DropletsService) Snapshots(_a0 int) (do.Images, error) {
	ret := _m.Called(_a0)

	var r0 do.Images
	if rf, ok := ret.Get(0).(func(int) do.Images); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(do.Images)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Backups provides a mock function with given fields: _a0
func (_m *DropletsService) Backups(_a0 int) (do.Images, error) {
	ret := _m.Called(_a0)

	var r0 do.Images
	if rf, ok := ret.Get(0).(func(int) do.Images); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(do.Images)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Actions provides a mock function with given fields: _a0
func (_m *DropletsService) Actions(_a0 int) (do.Actions, error) {
	ret := _m.Called(_a0)

	var r0 do.Actions
	if rf, ok := ret.Get(0).(func(int) do.Actions); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(do.Actions)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Neighbors provides a mock function with given fields: _a0
func (_m *DropletsService) Neighbors(_a0 int) (do.Droplets, error) {
	ret := _m.Called(_a0)

	var r0 do.Droplets
	if rf, ok := ret.Get(0).(func(int) do.Droplets); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(do.Droplets)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
