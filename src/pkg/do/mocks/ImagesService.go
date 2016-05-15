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

type ImagesService struct {
	mock.Mock
}

// List provides a mock function with given fields: public
func (_m *ImagesService) List(public bool) (do.Images, error) {
	ret := _m.Called(public)

	var r0 do.Images
	if rf, ok := ret.Get(0).(func(bool) do.Images); ok {
		r0 = rf(public)
	} else {
		r0 = ret.Get(0).(do.Images)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(public)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDistribution provides a mock function with given fields: public
func (_m *ImagesService) ListDistribution(public bool) (do.Images, error) {
	ret := _m.Called(public)

	var r0 do.Images
	if rf, ok := ret.Get(0).(func(bool) do.Images); ok {
		r0 = rf(public)
	} else {
		r0 = ret.Get(0).(do.Images)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(public)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListApplication provides a mock function with given fields: public
func (_m *ImagesService) ListApplication(public bool) (do.Images, error) {
	ret := _m.Called(public)

	var r0 do.Images
	if rf, ok := ret.Get(0).(func(bool) do.Images); ok {
		r0 = rf(public)
	} else {
		r0 = ret.Get(0).(do.Images)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(public)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUser provides a mock function with given fields: public
func (_m *ImagesService) ListUser(public bool) (do.Images, error) {
	ret := _m.Called(public)

	var r0 do.Images
	if rf, ok := ret.Get(0).(func(bool) do.Images); ok {
		r0 = rf(public)
	} else {
		r0 = ret.Get(0).(do.Images)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(public)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *ImagesService) GetByID(id int) (*do.Image, error) {
	ret := _m.Called(id)

	var r0 *do.Image
	if rf, ok := ret.Get(0).(func(int) *do.Image); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*do.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBySlug provides a mock function with given fields: slug
func (_m *ImagesService) GetBySlug(slug string) (*do.Image, error) {
	ret := _m.Called(slug)

	var r0 *do.Image
	if rf, ok := ret.Get(0).(func(string) *do.Image); ok {
		r0 = rf(slug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*do.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(slug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, iur
func (_m *ImagesService) Update(id int, iur *godo.ImageUpdateRequest) (*do.Image, error) {
	ret := _m.Called(id, iur)

	var r0 *do.Image
	if rf, ok := ret.Get(0).(func(int, *godo.ImageUpdateRequest) *do.Image); ok {
		r0 = rf(id, iur)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*do.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, *godo.ImageUpdateRequest) error); ok {
		r1 = rf(id, iur)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *ImagesService) Delete(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
