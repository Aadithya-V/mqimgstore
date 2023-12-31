// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	product "github.com/Aadithya-V/mqimgstore/product"
	mock "github.com/stretchr/testify/mock"
)

// IProductRepository is an autogenerated mock type for the IProductRepository type
type IProductRepository struct {
	mock.Mock
}

// InsertProduct provides a mock function with given fields: addProduct
func (_m *IProductRepository) InsertProduct(addProduct product.AddableProduct) (int64, error) {
	ret := _m.Called(addProduct)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(product.AddableProduct) (int64, error)); ok {
		return rf(addProduct)
	}
	if rf, ok := ret.Get(0).(func(product.AddableProduct) int64); ok {
		r0 = rf(addProduct)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(product.AddableProduct) error); ok {
		r1 = rf(addProduct)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIProductRepository creates a new instance of IProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIProductRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IProductRepository {
	mock := &IProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
