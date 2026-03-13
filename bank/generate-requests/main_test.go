package main

import (
	"bufio"
	"bytes"
	"io"
	"math"
	"testing"
)

// newTestRequestsWriter creates a RequestsWriter that writes to the given
// io.Writer instead of a real file, useful for unit tests.
func newTestRequestsWriter(w io.Writer) *RequestsWriter {
	return &RequestsWriter{
		writer: bufio.NewWriter(w),
	}
}

// ---------------------------------------------------------------------------
// buildFailureSchedule
// ---------------------------------------------------------------------------

func TestBuildFailureSchedule_CountMatchesRequested(t *testing.T) {
	tests := []struct {
		total, failing uint64
	}{
		{100, 0},
		{100, 1},
		{100, 50},
		{100, 100},
		{1, 1},
		{10, 3},
	}
	for _, tc := range tests {
		schedule := buildFailureSchedule(tc.total, tc.failing)
		if uint64(len(schedule)) != tc.total {
			t.Errorf("len(schedule) = %d, want %d", len(schedule), tc.total)
		}
		count := uint64(0)
		for _, v := range schedule {
			if v {
				count++
			}
		}
		if count != tc.failing {
			t.Errorf("failure count = %d, want %d (total=%d)", count, tc.failing, tc.total)
		}
	}
}

func TestBuildFailureSchedule_ZeroFailures(t *testing.T) {
	schedule := buildFailureSchedule(50, 0)
	for i, v := range schedule {
		if v {
			t.Fatalf("expected no failures, but index %d is true", i)
		}
	}
}

func TestBuildFailureSchedule_AllFailures(t *testing.T) {
	schedule := buildFailureSchedule(20, 20)
	for i, v := range schedule {
		if !v {
			t.Fatalf("expected all failures, but index %d is false", i)
		}
	}
}

// ---------------------------------------------------------------------------
// applyFailure
// ---------------------------------------------------------------------------

func TestApplyFailure_SameAccount(t *testing.T) {
	nextFailureKind = int(SameAccount)
	from, to, amt := applyFailure(1, 2, 500, 100)
	if from != to {
		t.Errorf("SameAccount: expected from == to, got from=%d to=%d", from, to)
	}
	if amt != 500 {
		t.Errorf("SameAccount: the amount should be unchanged, got %d", amt)
	}
}

func TestApplyFailure_FromNonExistent(t *testing.T) {
	nextFailureKind = int(FromNonExistent)
	accounts := 100
	from, to, amt := applyFailure(1, 2, 500, accounts)
	if from != accounts {
		t.Errorf("FromNonExistent: expected from=%d, got %d", accounts, from)
	}
	if to != 2 {
		t.Errorf("FromNonExistent: to should be unchanged, got %d", to)
	}
	if amt != 500 {
		t.Errorf("FromNonExistent: the amount should be unchanged, got %d", amt)
	}
}

func TestApplyFailure_ToNonExistent(t *testing.T) {
	nextFailureKind = int(ToNonExistent)
	accounts := 100
	from, to, amt := applyFailure(1, 2, 500, accounts)
	if from != 1 {
		t.Errorf("ToNonExistent: from should be unchanged, got %d", from)
	}
	if to != accounts {
		t.Errorf("ToNonExistent: expected to=%d, got %d", accounts, to)
	}
	if amt != 500 {
		t.Errorf("ToNonExistent: the amount should be unchanged, got %d", amt)
	}
}

func TestApplyFailure_InsufficientFunds(t *testing.T) {
	nextFailureKind = int(InsufficientFunds)
	from, to, amt := applyFailure(1, 2, 500, 100)
	if from != 1 || to != 2 {
		t.Errorf("InsufficientFunds: IDs should be unchanged, got from=%d to=%d", from, to)
	}
	if amt != math.MaxInt {
		t.Errorf("InsufficientFunds: expected amount=MaxInt, got %d", amt)
	}
}

func TestApplyFailure_CyclesThroughKinds(t *testing.T) {
	nextFailureKind = 0
	for i := 0; i < numFailureKinds*3; i++ {
		expectedNext := (i + 1) % numFailureKinds
		applyFailure(1, 2, 100, 50)
		if nextFailureKind != expectedNext {
			t.Fatalf("iteration %d: nextFailureKind = %d, want %d", i, nextFailureKind, expectedNext)
		}
	}
}

// ---------------------------------------------------------------------------
// writeLine
// ---------------------------------------------------------------------------

func TestWriteLine_SuccessfulRequest(t *testing.T) {
	var buf bytes.Buffer
	rw := newTestRequestsWriter(&buf)

	rw.WriteLine(42, true, 10, 20, 999)
	rw.writer.Flush()

	got := buf.String()
	want := "42 1 10 20 999\n"
	if got != want {
		t.Errorf("writeLine (success) = %q, want %q", got, want)
	}
}

func TestWriteLine_FailedRequest(t *testing.T) {
	var buf bytes.Buffer
	rw := newTestRequestsWriter(&buf)

	rw.WriteLine(7, false, 3, 5, 100)
	rw.writer.Flush()

	got := buf.String()
	want := "7 0 3 5 100\n"
	if got != want {
		t.Errorf("writeLine (failure) = %q, want %q", got, want)
	}
}

func TestWriteLine_ZeroValues(t *testing.T) {
	var buf bytes.Buffer
	rw := newTestRequestsWriter(&buf)

	rw.WriteLine(0, true, 0, 0, 0)
	rw.writer.Flush()

	got := buf.String()
	want := "0 1 0 0 0\n"
	if got != want {
		t.Errorf("writeLine (zeros) = %q, want %q", got, want)
	}
}

func TestWriteLine_LargeValues(t *testing.T) {
	var buf bytes.Buffer
	rw := newTestRequestsWriter(&buf)

	rw.WriteLine(2999999, true, 999999, 500000, math.MaxInt)
	rw.writer.Flush()

	got := buf.String()
	expected := "2999999 1 999999 500000 9223372036854775807\n"
	if got != expected {
		t.Errorf("writeLine (large) = %q, want %q", got, expected)
	}
}

func TestWriteLine_MultipleCallsAppendCorrectly(t *testing.T) {
	var buf bytes.Buffer
	rw := newTestRequestsWriter(&buf)

	rw.WriteLine(0, true, 1, 2, 10)
	rw.WriteLine(1, false, 3, 4, 20)
	rw.WriteLine(2, true, 5, 6, 30)
	rw.writer.Flush()

	got := buf.String()
	want := "0 1 1 2 10\n1 0 3 4 20\n2 1 5 6 30\n"
	if got != want {
		t.Errorf("writeLine (multi) = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// config defaults sanity
// ---------------------------------------------------------------------------

func TestDefaultConstants(t *testing.T) {
	if defaultAlpha <= 1.0 {
		t.Error("defaultAlpha must be > 1.0")
	}
	if defaultTotalRequests <= 0 {
		t.Error("defaultTotalRequests must be > 0")
	}
	if defaultAccounts <= 0 {
		t.Error("defaultAccounts must be > 0")
	}
	if defaultFailingRequests > defaultTotalRequests {
		t.Error("defaultFailingRequests must not exceed defaultTotalRequests")
	}
}

func TestNumFailureKindsMatchesEnum(t *testing.T) {
	// Ensure the constant stays in sync with the enum values.
	kinds := []FailureKind{SameAccount, FromNonExistent, ToNonExistent, InsufficientFunds}
	if numFailureKinds != len(kinds) {
		t.Errorf("numFailureKinds = %d, but there are %d FailureKind values", numFailureKinds, len(kinds))
	}
}
