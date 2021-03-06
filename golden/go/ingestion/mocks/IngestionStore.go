// Code generated by mockery v2.4.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// IngestionStore is an autogenerated mock type for the IngestionStore type
type IngestionStore struct {
	mock.Mock
}

// SetIngested provides a mock function with given fields: ctx, fileName, md5, ts
func (_m *IngestionStore) SetIngested(ctx context.Context, fileName string, md5 string, ts time.Time) error {
	ret := _m.Called(ctx, fileName, md5, ts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) error); ok {
		r0 = rf(ctx, fileName, md5, ts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WasIngested provides a mock function with given fields: ctx, fileName, md5
func (_m *IngestionStore) WasIngested(ctx context.Context, fileName string, md5 string) (bool, error) {
	ret := _m.Called(ctx, fileName, md5)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, fileName, md5)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, fileName, md5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
