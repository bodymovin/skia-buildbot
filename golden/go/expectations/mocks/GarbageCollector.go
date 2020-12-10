// Code generated by mockery v2.4.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	expectations "go.skia.org/infra/golden/go/expectations"

	time "time"
)

// GarbageCollector is an autogenerated mock type for the GarbageCollector type
type GarbageCollector struct {
	mock.Mock
}

// GarbageCollect provides a mock function with given fields: _a0
func (_m *GarbageCollector) GarbageCollect(_a0 context.Context) (int, error) {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkUnusedEntriesForGC provides a mock function with given fields: _a0, _a1, _a2
func (_m *GarbageCollector) MarkUnusedEntriesForGC(_a0 context.Context, _a1 expectations.LabelInt, _a2 time.Time) (int, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, expectations.LabelInt, time.Time) int); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, expectations.LabelInt, time.Time) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLastUsed provides a mock function with given fields: _a0, _a1, _a2
func (_m *GarbageCollector) UpdateLastUsed(_a0 context.Context, _a1 []expectations.ID, _a2 time.Time) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []expectations.ID, time.Time) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
