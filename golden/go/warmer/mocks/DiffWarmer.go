// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	digest_counter "go.skia.org/infra/golden/go/digest_counter"
	digesttools "go.skia.org/infra/golden/go/digesttools"

	mock "github.com/stretchr/testify/mock"

	summary "go.skia.org/infra/golden/go/summary"

	types "go.skia.org/infra/golden/go/types"
)

// DiffWarmer is an autogenerated mock type for the DiffWarmer type
type DiffWarmer struct {
	mock.Mock
}

// PrecomputeDiffs provides a mock function with given fields: ctx, summaries, testNames, dCounter, diffFinder
func (_m *DiffWarmer) PrecomputeDiffs(ctx context.Context, summaries []*summary.TriageStatus, testNames types.TestNameSet, dCounter digest_counter.DigestCounter, diffFinder digesttools.ClosestDiffFinder) error {
	ret := _m.Called(ctx, summaries, testNames, dCounter, diffFinder)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*summary.TriageStatus, types.TestNameSet, digest_counter.DigestCounter, digesttools.ClosestDiffFinder) error); ok {
		r0 = rf(ctx, summaries, testNames, dCounter, diffFinder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
