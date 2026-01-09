package math_test

import (
	stdmath "math"
	"testing"

	"github.com/OffchainLabs/prysm/v7/math"
	"github.com/OffchainLabs/prysm/v7/testing/require"
)

// Custom type for testing generics
type testUint64 uint64

// assertPanic checks that the provided function panics with the expected message.
func assertPanic(t *testing.T, panicMessage string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic not thrown")
		} else if r != panicMessage {
			t.Errorf("Unexpected panic thrown, want: %#v, got: %#v", panicMessage, r)
		}
	}()
	f()
}

func TestSafeAdd(t *testing.T) {
	tests := []struct {
		a      testUint64
		b      testUint64
		result testUint64
		err    bool
	}{
		{a: 0, b: 0, result: 0},
		{a: 1, b: 2, result: 3},
		{a: 100, b: 200, result: 300},
		{a: testUint64(stdmath.MaxUint64 - 1), b: 1, result: testUint64(stdmath.MaxUint64)},
		{a: testUint64(stdmath.MaxUint64), b: 1, err: true},
		{a: testUint64(stdmath.MaxUint64), b: testUint64(stdmath.MaxUint64), err: true},
	}

	for _, tt := range tests {
		result, err := math.SafeAdd(tt.a, tt.b)
		if tt.err {
			require.ErrorContains(t, "overflows", err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.result, result)
		}
	}
}

func TestMustAdd(t *testing.T) {
	// Test successful addition
	result := math.MustAdd(testUint64(1), testUint64(2))
	require.Equal(t, testUint64(3), result)

	// Test panic on overflow
	assertPanic(t, "addition overflows", func() {
		math.MustAdd(testUint64(stdmath.MaxUint64), testUint64(1))
	})
}

func TestSafeSub(t *testing.T) {
	tests := []struct {
		a      testUint64
		b      testUint64
		result testUint64
		err    bool
	}{
		{a: 0, b: 0, result: 0},
		{a: 3, b: 2, result: 1},
		{a: 300, b: 200, result: 100},
		{a: testUint64(stdmath.MaxUint64), b: testUint64(stdmath.MaxUint64), result: 0},
		{a: 0, b: 1, err: true},
		{a: 5, b: 10, err: true},
	}

	for _, tt := range tests {
		result, err := math.SafeSub(tt.a, tt.b)
		if tt.err {
			require.ErrorContains(t, "underflow", err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.result, result)
		}
	}
}

func TestMustSub(t *testing.T) {
	// Test successful subtraction
	result := math.MustSub(testUint64(5), testUint64(3))
	require.Equal(t, testUint64(2), result)

	// Test panic on underflow
	assertPanic(t, "subtraction underflow", func() {
		math.MustSub(testUint64(0), testUint64(1))
	})
}

func TestFlooredSub(t *testing.T) {
	tests := []struct {
		a      testUint64
		b      testUint64
		result testUint64
	}{
		{a: 0, b: 0, result: 0},
		{a: 5, b: 3, result: 2},
		{a: 100, b: 100, result: 0},
		// Saturating behavior: returns 0 instead of underflow
		{a: 0, b: 1, result: 0},
		{a: 5, b: 10, result: 0},
		{a: 0, b: testUint64(stdmath.MaxUint64), result: 0},
	}

	for _, tt := range tests {
		result := math.FlooredSub(tt.a, tt.b)
		require.Equal(t, tt.result, result)
	}
}

func TestSafeMul(t *testing.T) {
	tests := []struct {
		a      testUint64
		b      testUint64
		result testUint64
		err    bool
	}{
		{a: 0, b: 0, result: 0},
		{a: 2, b: 3, result: 6},
		{a: 100, b: 100, result: 10000},
		{a: testUint64(stdmath.MaxUint64), b: 1, result: testUint64(stdmath.MaxUint64)},
		{a: testUint64(stdmath.MaxUint64), b: 2, err: true},
		{a: testUint64(1 << 32), b: testUint64(1 << 32), err: true},
	}

	for _, tt := range tests {
		result, err := math.SafeMul(tt.a, tt.b)
		if tt.err {
			require.ErrorContains(t, "overflows", err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.result, result)
		}
	}
}

func TestMustMul(t *testing.T) {
	// Test successful multiplication
	result := math.MustMul(testUint64(3), testUint64(4))
	require.Equal(t, testUint64(12), result)

	// Test panic on overflow
	assertPanic(t, "multiplication overflows", func() {
		math.MustMul(testUint64(stdmath.MaxUint64), testUint64(2))
	})
}

func TestSafeDiv(t *testing.T) {
	tests := []struct {
		a      testUint64
		b      testUint64
		result testUint64
		err    bool
	}{
		{a: 0, b: 1, result: 0},
		{a: 6, b: 2, result: 3},
		{a: 100, b: 10, result: 10},
		{a: testUint64(stdmath.MaxUint64), b: 1, result: testUint64(stdmath.MaxUint64)},
		{a: 5, b: 0, err: true},
		{a: 0, b: 0, err: true},
	}

	for _, tt := range tests {
		result, err := math.SafeDiv(tt.a, tt.b)
		if tt.err {
			require.NotNil(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.result, result)
		}
	}
}

func TestMustDiv(t *testing.T) {
	// Test successful division
	result := math.MustDiv(testUint64(12), testUint64(4))
	require.Equal(t, testUint64(3), result)

	// Test panic on div by zero
	assertPanic(t, "integer divide by zero", func() {
		math.MustDiv(testUint64(5), testUint64(0))
	})
}

func TestSafeMod(t *testing.T) {
	tests := []struct {
		a      testUint64
		b      testUint64
		result testUint64
		err    bool
	}{
		{a: 0, b: 1, result: 0},
		{a: 7, b: 3, result: 1},
		{a: 100, b: 10, result: 0},
		{a: 5, b: 0, err: true},
	}

	for _, tt := range tests {
		result, err := math.SafeMod(tt.a, tt.b)
		if tt.err {
			require.NotNil(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.result, result)
		}
	}
}

func TestMustMod(t *testing.T) {
	// Test successful modulo
	result := math.MustMod(testUint64(7), testUint64(3))
	require.Equal(t, testUint64(1), result)

	// Test panic on div by zero
	assertPanic(t, "integer divide by zero", func() {
		math.MustMod(testUint64(5), testUint64(0))
	})
}

// TestTypeInference verifies that Go's type inference works with untyped constants
func TestTypeInference(t *testing.T) {
	// This demonstrates that math.SafeAdd(slot, 5) works via type inference
	var a testUint64 = 10
	result, err := math.SafeAdd(a, 5) // 5 is inferred as testUint64
	require.NoError(t, err)
	require.Equal(t, testUint64(15), result)
}
