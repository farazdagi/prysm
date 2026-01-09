// Package primitives provides stub types for testing the primarithm analyzer.
package primitives

// Slot represents a single slot (stub).
type Slot uint64

// Add increases slot by x (stub).
func (s Slot) Add(x uint64) Slot { return s }

// AddSlot increases slot by another slot (stub).
func (s Slot) AddSlot(x Slot) Slot { return s }

// SafeAdd increases slot by x safely (stub).
func (s Slot) SafeAdd(x uint64) (Slot, error) { return s, nil }

// SafeAddSlot increases slot by another slot safely (stub).
func (s Slot) SafeAddSlot(x Slot) (Slot, error) { return s, nil }

// Sub subtracts x from the slot (stub).
func (s Slot) Sub(x uint64) Slot { return s }

// SubSlot subtracts another slot from the slot (stub).
func (s Slot) SubSlot(x Slot) Slot { return s }

// SafeSub subtracts x from the slot safely (stub).
func (s Slot) SafeSub(x uint64) (Slot, error) { return s, nil }

// SafeSubSlot subtracts another slot safely (stub).
func (s Slot) SafeSubSlot(x Slot) (Slot, error) { return s, nil }

// FlooredSubSlot safely subtracts, returning 0 on underflow (stub).
func (s Slot) FlooredSubSlot(x Slot) Slot { return s }

// Mul multiplies slot by x (stub).
func (s Slot) Mul(x uint64) Slot { return s }

// SafeMul multiplies slot by x safely (stub).
func (s Slot) SafeMul(x uint64) (Slot, error) { return s, nil }

// Div divides slot by x (stub).
func (s Slot) Div(x uint64) Slot { return s }

// SafeDiv divides slot by x safely (stub).
func (s Slot) SafeDiv(x uint64) (Slot, error) { return s, nil }

// Mod returns slot mod x (stub).
func (s Slot) Mod(x uint64) Slot { return s }

// SafeMod returns slot mod x safely (stub).
func (s Slot) SafeMod(x uint64) (Slot, error) { return s, nil }

// Epoch represents a single epoch (stub).
type Epoch uint64

// Add increases epoch by x (stub).
func (e Epoch) Add(x uint64) Epoch { return e }

// SafeAdd increases epoch by x safely (stub).
func (e Epoch) SafeAdd(x uint64) (Epoch, error) { return e, nil }

// Sub subtracts x from the epoch (stub).
func (e Epoch) Sub(x uint64) Epoch { return e }

// SafeSub subtracts x from the epoch safely (stub).
func (e Epoch) SafeSub(x uint64) (Epoch, error) { return e, nil }

// FlooredSubEpoch safely subtracts, returning 0 on underflow (stub).
func (e Epoch) FlooredSubEpoch(x Epoch) Epoch { return e }

// ValidatorIndex represents a validator index (stub).
type ValidatorIndex uint64

// Add increases the index by x (stub).
func (v ValidatorIndex) Add(x uint64) ValidatorIndex { return v }

// Sub subtracts x from the index (stub).
func (v ValidatorIndex) Sub(x uint64) ValidatorIndex { return v }

// SafeAdd increases the index by x safely (stub).
func (v ValidatorIndex) SafeAdd(x uint64) (ValidatorIndex, error) { return v, nil }

// SafeSub subtracts x from the index safely (stub).
func (v ValidatorIndex) SafeSub(x uint64) (ValidatorIndex, error) { return v, nil }

// FlooredSub safely subtracts, returning 0 on underflow (stub).
func (v ValidatorIndex) FlooredSub(x uint64) ValidatorIndex { return v }

// Gwei is a denomination of Ether (stub).
type Gwei uint64

// Add increases gwei by x (stub).
func (g Gwei) Add(x uint64) Gwei { return g }

// SafeAdd increases gwei by x safely (stub).
func (g Gwei) SafeAdd(x uint64) (Gwei, error) { return g, nil }

// Sub subtracts x from gwei (stub).
func (g Gwei) Sub(x uint64) Gwei { return g }

// SafeSub subtracts x from gwei safely (stub).
func (g Gwei) SafeSub(x uint64) (Gwei, error) { return g, nil }

// FlooredSub safely subtracts, returning 0 on underflow (stub).
func (g Gwei) FlooredSub(x Gwei) Gwei { return g }

// Mul multiplies gwei by x (stub).
func (g Gwei) Mul(x uint64) Gwei { return g }

// SafeMul multiplies gwei by x safely (stub).
func (g Gwei) SafeMul(x uint64) (Gwei, error) { return g, nil }

// Div divides gwei by x (stub).
func (g Gwei) Div(x uint64) Gwei { return g }

// SafeDiv divides gwei by x safely (stub).
func (g Gwei) SafeDiv(x uint64) (Gwei, error) { return g, nil }

// CommitteeIndex represents a committee index (stub).
type CommitteeIndex uint64

// BuilderIndex represents a builder index (stub).
type BuilderIndex uint64

// SSZUint64 is a uint64 with SSZ methods (stub).
type SSZUint64 uint64

// BP represents basis points (stub).
type BP uint64
