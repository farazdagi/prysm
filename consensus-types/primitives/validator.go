package primitives

import (
	"fmt"

	"github.com/OffchainLabs/prysm/v7/math"
	fssz "github.com/prysmaticlabs/fastssz"
)

var _ fssz.HashRoot = (ValidatorIndex)(0)
var _ fssz.Marshaler = (*ValidatorIndex)(nil)
var _ fssz.Unmarshaler = (*ValidatorIndex)(nil)

// ValidatorIndex in eth2.
type ValidatorIndex uint64

// Div divides validator index by x.
// This method panics if dividing by zero!
func (v ValidatorIndex) Div(x uint64) ValidatorIndex {
	if x == 0 {
		panic("divbyzero") // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return ValidatorIndex(uint64(v) / x)
}

// Add increases validator index by x.
func (v ValidatorIndex) Add(x uint64) ValidatorIndex {
	return ValidatorIndex(uint64(v) + x)
}

// Sub subtracts x from the validator index.
// This method panics if causing an underflow!
func (v ValidatorIndex) Sub(x uint64) ValidatorIndex {
	if uint64(v) < x {
		panic("underflow") // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return ValidatorIndex(uint64(v) - x)
}

// Mod returns result of `validator index % x`.
func (v ValidatorIndex) Mod(x uint64) ValidatorIndex {
	return ValidatorIndex(uint64(v) % x)
}

// SafeAdd increases validator index by x.
// In case of overflow, error is returned.
func (v ValidatorIndex) SafeAdd(x uint64) (ValidatorIndex, error) {
	res, err := math.Add64(uint64(v), x)
	return ValidatorIndex(res), err
}

// SafeSub subtracts x from the validator index.
// In case of underflow, error is returned.
func (v ValidatorIndex) SafeSub(x uint64) (ValidatorIndex, error) {
	res, err := math.Sub64(uint64(v), x)
	return ValidatorIndex(res), err
}

// FlooredSub safely subtracts x from the validator index, returning 0 if the result would underflow.
func (v ValidatorIndex) FlooredSub(x uint64) ValidatorIndex {
	if uint64(v) < x {
		return 0
	}
	return ValidatorIndex(uint64(v) - x)
}

// SafeDiv divides validator index by x.
// In case of division by zero, error is returned.
func (v ValidatorIndex) SafeDiv(x uint64) (ValidatorIndex, error) {
	res, err := math.Div64(uint64(v), x)
	return ValidatorIndex(res), err
}

// SafeMod returns result of `validator index % x`.
// In case of division by zero, error is returned.
func (v ValidatorIndex) SafeMod(x uint64) (ValidatorIndex, error) {
	res, err := math.Mod64(uint64(v), x)
	return ValidatorIndex(res), err
}

// HashTreeRoot --
func (v ValidatorIndex) HashTreeRoot() ([32]byte, error) {
	return fssz.HashWithDefaultHasher(v)
}

// HashTreeRootWith --
func (v ValidatorIndex) HashTreeRootWith(hh *fssz.Hasher) error {
	hh.PutUint64(uint64(v))
	return nil
}

// UnmarshalSSZ --
func (v *ValidatorIndex) UnmarshalSSZ(buf []byte) error {
	if len(buf) != v.SizeSSZ() {
		return fmt.Errorf("expected buffer of length %d received %d", v.SizeSSZ(), len(buf))
	}
	*v = ValidatorIndex(fssz.UnmarshallUint64(buf))
	return nil
}

// MarshalSSZTo --
func (v *ValidatorIndex) MarshalSSZTo(dst []byte) ([]byte, error) {
	marshalled, err := v.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return append(dst, marshalled...), nil
}

// MarshalSSZ --
func (v *ValidatorIndex) MarshalSSZ() ([]byte, error) {
	marshalled := fssz.MarshalUint64([]byte{}, uint64(*v))
	return marshalled, nil
}

// SizeSSZ --
func (v *ValidatorIndex) SizeSSZ() int {
	return 8
}
