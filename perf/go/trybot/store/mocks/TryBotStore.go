// Code generated by mockery v2.4.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	store "go.skia.org/infra/perf/go/trybot/store"

	time "time"

	trybot "go.skia.org/infra/perf/go/trybot"

	types "go.skia.org/infra/perf/go/types"
)

// TryBotStore is an autogenerated mock type for the TryBotStore type
type TryBotStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, cl, patch
func (_m *TryBotStore) Get(ctx context.Context, cl types.CL, patch int) ([]store.GetResult, error) {
	ret := _m.Called(ctx, cl, patch)

	var r0 []store.GetResult
	if rf, ok := ret.Get(0).(func(context.Context, types.CL, int) []store.GetResult); ok {
		r0 = rf(ctx, cl, patch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]store.GetResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.CL, int) error); ok {
		r1 = rf(ctx, cl, patch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, since
func (_m *TryBotStore) List(ctx context.Context, since time.Time) ([]store.ListResult, error) {
	ret := _m.Called(ctx, since)

	var r0 []store.ListResult
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) []store.ListResult); ok {
		r0 = rf(ctx, since)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]store.ListResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time) error); ok {
		r1 = rf(ctx, since)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Write provides a mock function with given fields: ctx, tryFile
func (_m *TryBotStore) Write(ctx context.Context, tryFile trybot.TryFile) error {
	ret := _m.Called(ctx, tryFile)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, trybot.TryFile) error); ok {
		r0 = rf(ctx, tryFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
