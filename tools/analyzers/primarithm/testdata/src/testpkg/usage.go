// Package testpkg provides test cases for the primarithm analyzer.
package testpkg

import (
	"github.com/OffchainLabs/prysm/v7/consensus-types/primitives"
)

// Named constant for testing constant detection.
const slotsPerEpoch = 32

// =============================================================================
// ADDITION TESTS
// =============================================================================

func AdditionFlagged() {
	var a, b primitives.Slot
	_ = a + b // want "unsafe arithmetic on primitives.Slot; use SafeAdd/Add/SafeAddSlot/AddSlot instead to prevent overflow/underflow"

	// Large constant (> threshold) - should flag
	_ = a + 1001 // want "unsafe arithmetic on primitives.Slot; use SafeAdd/Add/SafeAddSlot/AddSlot instead to prevent overflow/underflow"
	_ = 1001 + a // want "unsafe arithmetic on primitives.Slot; use SafeAdd/Add/SafeAddSlot/AddSlot instead to prevent overflow/underflow"
}

func AdditionSkipped() {
	var a primitives.Slot
	// Small constants - safe because overflow needs a ≈ MaxUint64
	_ = a + 1    // OK: small constant
	_ = 1 + a    // OK: small constant on left
	_ = a + 100  // OK: small constant
	_ = a + 1000 // OK: at threshold

	// Named constant (small value)
	_ = a + slotsPerEpoch // OK: named constant, small value
}

// =============================================================================
// SUBTRACTION TESTS
// =============================================================================

func SubtractionFlagged() {
	var a, b primitives.Slot
	// Two variables - THE DANGEROUS PATTERN (PR #16227 bug)
	_ = a - b // want "unsafe arithmetic on primitives.Slot; use SafeSub/Sub/SafeSubSlot/SubSlot/FlooredSubSlot instead to prevent overflow/underflow"

	// Even small constants - underflow at slot=0 is realistic (genesis)
	_ = a - 1 // want "unsafe arithmetic on primitives.Slot; use SafeSub/Sub/SafeSubSlot/SubSlot/FlooredSubSlot instead to prevent overflow/underflow"
	_ = a - slotsPerEpoch // want "unsafe arithmetic on primitives.Slot; use SafeSub/Sub/SafeSubSlot/SubSlot/FlooredSubSlot instead to prevent overflow/underflow"
}

func SubtractionSkipped() {
	var a primitives.Slot
	// Only subtracting 0 is safe (no-op)
	_ = a - 0 // OK: subtracting zero is a no-op
}

// =============================================================================
// MULTIPLICATION TESTS
// =============================================================================

func MultiplicationFlagged() {
	var a, b primitives.Slot
	// Two variables
	_ = a * b // want "unsafe arithmetic on primitives.Slot; use SafeMul/Mul/SafeMulSlot/MulSlot instead to prevent overflow/underflow"

	// Even small multipliers can overflow with large values (e.g., gwei * 1e9)
	_ = a * 2 // want "unsafe arithmetic on primitives.Slot; use SafeMul/Mul/SafeMulSlot/MulSlot instead to prevent overflow/underflow"
	_ = 2 * a // want "unsafe arithmetic on primitives.Slot; use SafeMul/Mul/SafeMulSlot/MulSlot instead to prevent overflow/underflow"
}

func MultiplicationSkipped() {
	var a primitives.Slot
	// Only 0 and 1 are safe
	_ = a * 0 // OK: always 0
	_ = 0 * a // OK: always 0
	_ = a * 1 // OK: identity
	_ = 1 * a // OK: identity
}

// =============================================================================
// DIVISION TESTS
// =============================================================================

func DivisionFlagged() {
	var a, b primitives.Slot
	// Variable divisor could be zero
	_ = a / b // want "unsafe arithmetic on primitives.Slot; use SafeDiv/Div/SafeDivSlot/DivSlot instead to prevent overflow/underflow"
}

func DivisionSkipped() {
	var a primitives.Slot
	// Non-zero constant divisor is always safe
	_ = a / 1  // OK: non-zero constant
	_ = a / 2  // OK: non-zero constant
	_ = a / 32 // OK: non-zero constant
	_ = a / slotsPerEpoch // OK: non-zero named constant
}

// =============================================================================
// MODULO TESTS
// =============================================================================

func ModuloFlagged() {
	var a, b primitives.Slot
	// Variable divisor could be zero
	_ = a % b // want "unsafe arithmetic on primitives.Slot; use SafeMod/Mod/SafeModSlot/ModSlot instead to prevent overflow/underflow"
}

func ModuloSkipped() {
	var a primitives.Slot
	// Non-zero constant divisor is always safe
	_ = a % 2  // OK: non-zero constant
	_ = a % 32 // OK: non-zero constant
	_ = a % slotsPerEpoch // OK: non-zero named constant
}

// =============================================================================
// OTHER TYPES - Ensure they're also handled correctly
// =============================================================================

func EpochArithmetic() {
	var a, b primitives.Epoch
	_ = a - b // want "unsafe arithmetic on primitives.Epoch; use SafeSub/Sub/FlooredSubEpoch instead to prevent overflow/underflow"
	_ = a + 1 // OK: small constant addition
	_ = a / 2 // OK: non-zero constant division
}

func ValidatorIndexArithmetic() {
	var a, b primitives.ValidatorIndex
	_ = a + b // want "unsafe arithmetic on primitives.ValidatorIndex; use SafeAdd/Add instead to prevent overflow/underflow"
	_ = a + 1 // OK: small constant addition
}

func GweiArithmetic() {
	var a, b primitives.Gwei
	// Gwei is especially dangerous for multiplication (wei conversion)
	_ = a * b // want "unsafe arithmetic on primitives.Gwei; use SafeMul/Mul instead to prevent overflow/underflow"
	_ = a * 2 // want "unsafe arithmetic on primitives.Gwei; use SafeMul/Mul instead to prevent overflow/underflow"
	_ = a + 1 // OK: small constant addition
}

func CommitteeIndexArithmetic() {
	var a, b primitives.CommitteeIndex
	_ = a + b // want "unsafe arithmetic on primitives.CommitteeIndex; use math.SafeAdd/math.MustAdd instead to prevent overflow/underflow"
	_ = a + 1 // OK: small constant addition
}

// =============================================================================
// COMPARISON AND BITWISE - Should NOT be flagged
// =============================================================================

func ComparisonOperations() {
	var a, b primitives.Slot
	_ = a < b  // OK: comparison, not arithmetic
	_ = a > b  // OK: comparison, not arithmetic
	_ = a <= b // OK: comparison, not arithmetic
	_ = a >= b // OK: comparison, not arithmetic
	_ = a == b // OK: comparison, not arithmetic
	_ = a != b // OK: comparison, not arithmetic
}

func BitwiseOperations() {
	var a primitives.Slot
	_ = a & 0xFF // OK: bitwise, not arithmetic
	_ = a | 0xFF // OK: bitwise, not arithmetic
	_ = a ^ 0xFF // OK: bitwise, not arithmetic
	_ = a << 2   // OK: shift, not arithmetic
	_ = a >> 2   // OK: shift, not arithmetic
}

// =============================================================================
// LINT IGNORE - Should suppress warning
// =============================================================================

func IgnoredArithmetic() {
	var a, b primitives.Slot
	_ = a - b // lint:ignore primarithm -- intentional for test
}

// =============================================================================
// THE ORIGINAL BUG PATTERN (PR #16227)
// =============================================================================

func OriginalBugPattern() {
	var begin, end primitives.Slot
	// This is the exact pattern that caused the panic - MUST be flagged
	_ = uint64(end - begin) // want "unsafe arithmetic on primitives.Slot; use SafeSub/Sub/SafeSubSlot/SubSlot/FlooredSubSlot instead to prevent overflow/underflow"

	// Safe alternatives (not flagged)
	diff := end.FlooredSubSlot(begin)
	_ = uint64(diff)

	diff2, _ := end.SafeSubSlot(begin)
	_ = uint64(diff2)
}
