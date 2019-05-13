package gitinfo

import (
	"context"
	"strings"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
	"go.skia.org/infra/go/testutils/unittest"
	"go.skia.org/infra/go/util"
	vcstu "go.skia.org/infra/go/vcsinfo/testutils"
)

func TestVCSSuite(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	assert.NoError(t, err)
	vcstu.TestDisplay(t, r)
}

func TestFrom(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	r, err := NewGitInfo(context.TODO(), repoDir, false, false)
	assert.NoError(t, err)
	vcstu.TestFrom(t, r)
}

func TestLastN(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		n      int
		values []string
	}{
		{
			n:      0,
			values: []string{},
		},
		{
			n:      1,
			values: []string{"8652a6df7dc8a7e6addee49f6ed3c2308e36bd18"},
		},
		{
			n:      2,
			values: []string{"7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18"},
		},
		{
			n:      5,
			values: []string{"7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18"},
		},
	}
	for _, tc := range testCases {
		if got, want := r.LastN(ctx, tc.n), tc.values; !util.SSliceEqual(got, want) {
			t.Errorf("For N: %d Hashes returned is wrong: Got %#v Want %#v", tc.n, got, want)
		}
	}
}

func TestByIndex(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	r, err := NewGitInfo(context.TODO(), repoDir, false, false)
	assert.NoError(t, err)
	vcstu.TestByIndex(t, r)
}

func TestLastNIndex(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	r, err := NewGitInfo(context.TODO(), repoDir, false, false)
	assert.NoError(t, err)
	vcstu.TestLastNIndex(t, r)
}

func TestIndexOf(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := r.IndexOf(ctx, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f")
	assert.NoError(t, err)
	assert.Equal(t, 0, idx)
	idx, err = r.IndexOf(ctx, "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18")
	assert.NoError(t, err)
	assert.Equal(t, 1, idx)
	idx, err = r.IndexOf(ctx, "foo")
	assert.Error(t, err)

	assert.Equal(t, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", r.firstCommit)
}

func TestRange(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	r, err := NewGitInfo(context.TODO(), repoDir, false, false)
	assert.NoError(t, err)
	vcstu.TestRange(t, r)
}
func TestLog(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.Log(ctx, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", "")
	if err != nil {
		t.Fatal(err)
	}
	want := `commit 7a669cfa3f4cd3482a4fd03989f75efcc7595f7f
Author: Joe Gregorio <jcgregorio@google.com>
Date:   Wed Jul 30 08:00:42 2014 -0400

    First "checkin"
    ` + "\n" + `    With quotes.

README.txt
`

	if got != want {
		t.Errorf("Log failed: \nGot  %q \nWant %q", got, want)
	}

	got, err = r.Log(ctx, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18")
	if err != nil {
		t.Fatal(err)
	}
	want = `commit 8652a6df7dc8a7e6addee49f6ed3c2308e36bd18
Author: Joe Gregorio <jcgregorio@google.com>
Date:   Wed Jul 30 08:01:55 2014 -0400

    Added code. No body.

README.txt
hello.go
`

	if got != want {
		t.Errorf("Log failed: \nGot  %q \nWant %q", got, want)
	}

}

func TestLogFine(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.LogFine(ctx, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", "", "--format=format:%H")
	if err != nil {
		t.Fatal(err)
	}
	want := `7a669cfa3f4cd3482a4fd03989f75efcc7595f7f`

	if got != want {
		t.Errorf("Log failed: \nGot  %q \nWant %q", got, want)
	}

	got, err = r.LogFine(ctx, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18", "--format=format:%H")
	if err != nil {
		t.Fatal(err)
	}
	want = `8652a6df7dc8a7e6addee49f6ed3c2308e36bd18`

	if got != want {
		t.Errorf("Log failed: \nGot  %q \nWant %q", got, want)
	}
}

func TestLogArgs(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}
	got, err := r.LogArgs(ctx, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f..8652a6df7dc8a7e6addee49f6ed3c2308e36bd18", "--format=format:%H")
	if err != nil {
		t.Fatal(err)
	}
	want := `8652a6df7dc8a7e6addee49f6ed3c2308e36bd18`
	if got != want {
		t.Errorf("Log failed: \nGot  %q \nWant %q", got, want)
	}

}

func TestShortList(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}

	l, err := r.ShortList(ctx, "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18", "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(l.Commits), 0; got != want {
		t.Fatalf("Wrong number of zero results: Got %v Want %v", got, want)
	}

	c, err := r.ShortList(ctx, "7a669cfa3f4cd3482a4fd03989f75efcc7595f7f", "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18")
	if err != nil {
		t.Fatal(err)
	}

	expected := []struct {
		Hash    string
		Author  string
		Subject string
	}{
		{
			Hash:    "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18",
			Author:  "Joe Gregorio",
			Subject: "Added code. No body.",
		},
	}
	if got, want := len(c.Commits), len(expected); got != want {
		t.Fatalf("Wrong number of results: Got %v Want %v", got, want)
	}
	for i, o := range c.Commits {
		if got, want := o.Hash, expected[i].Hash; got != want {
			t.Errorf("Wrong hash: Got %v Want %v", got, want)
		}
		if got, want := o.Author, expected[i].Author; got != want {
			t.Errorf("Wrong author: Got %v Want %v", got, want)
		}
		if got, want := o.Subject, expected[i].Subject; got != want {
			t.Errorf("Wrong subject: Got %v Want %v", got, want)
		}
	}
}

func TestTileAddressFromHash(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}

	// The two commits in the repo have timestamps of:
	// 1406721715 and 1406721642.
	testCases := []struct {
		hash   string
		start  time.Time
		num    int
		offset int
	}{
		{
			hash:   "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18",
			start:  time.Unix(1406721642, 0),
			num:    0,
			offset: 1,
		},
		{
			hash:   "8652a6df7dc8a7e6addee49f6ed3c2308e36bd18",
			start:  time.Unix(1406721643, 0),
			num:    0,
			offset: 0,
		},
	}

	for _, tc := range testCases {
		n, o, err := r.TileAddressFromHash(tc.hash, tc.start)
		if err != nil {
			t.Fatal(err)
		}
		if n != tc.num || o != tc.offset {
			t.Errorf("Address unexpected: got (%d, %d), want (%d, %d).", tc.num, tc.offset, n, o)
		}
	}
}

func TestNumCommits(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, false)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := r.NumCommits(), 2; got != want {
		t.Errorf("NumCommit wrong number: Got %v Want %v", got, want)
	}
}

func TestRevList(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, true)
	if err != nil {
		t.Fatal(err)
	}

	revs := []string{
		"8652a6df7dc8a7e6addee49f6ed3c2308e36bd18",
		"7a669cfa3f4cd3482a4fd03989f75efcc7595f7f",
	}
	testCases := []struct {
		Input    []string
		Expected []string
	}{
		{
			Input:    []string{"master"},
			Expected: revs,
		},
		{
			Input:    []string{"HEAD"},
			Expected: revs,
		},
		{
			Input:    []string{"7a669cf..8652a6d"},
			Expected: revs[1:],
		},
		{
			Input:    []string{"8652a6d", "^7a669cf"},
			Expected: revs[1:],
		},
		{
			Input:    []string{"8652a6d..7a669cf"},
			Expected: []string{},
		},
	}
	for _, tc := range testCases {
		actual, err := r.RevList(ctx, tc.Input...)
		if err != nil {
			t.Fatal(err)
		}
		if !util.SSliceEqual(actual, tc.Expected) {
			t.Fatalf("Failed test for: git rev-list %s\nGot:  %v\nWant: %v", strings.Join(tc.Input, " "), actual, tc.Expected)
		}
	}
}

func TestBranchInfo(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, true)
	if err != nil {
		t.Fatal(err)
	}

	allBranches, err := r.GetBranches(ctx)
	assert.NoError(t, err)
	branches := []string{}
	for _, b := range allBranches {
		branches = append(branches, b.Name)
	}
	vcstu.TestBranchInfo(t, r, branches)
}

func TestSetBranch(t *testing.T) {
	unittest.MediumTest(t)
	repoDir, cleanup := vcstu.InitTempRepo()
	defer cleanup()

	ctx := context.Background()
	r, err := NewGitInfo(ctx, repoDir, false, true)
	if err != nil {
		t.Fatal(err)
	}

	branches, err := r.GetBranches(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(branches))

	err = r.Checkout(ctx, "test-branch-1")
	assert.NoError(t, err)

	commits := r.LastN(ctx, 10)
	assert.Equal(t, 3, len(commits))
	assert.Equal(t, "3f5a807d432ac232a952bbf223bc6952e4b49b2c", commits[2])
}
