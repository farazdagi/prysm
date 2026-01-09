package primitives

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/OffchainLabs/prysm/v7/math"
	fssz "github.com/prysmaticlabs/fastssz"
)

// ZW returns a non-nil zero value for primitives.Wei
func ZeroWei() Wei {
	return big.NewInt(0)
}

// Wei is the smallest unit of Ether, represented as a pointer to a bigInt.
type Wei *big.Int

// Gwei is a denomination of 1e9 Wei represented as an uint64.
type Gwei uint64

var _ fssz.HashRoot = (Gwei)(0)
var _ fssz.Marshaler = (*Gwei)(nil)
var _ fssz.Unmarshaler = (*Gwei)(nil)

// Add increases Gwei by x.
// In case of overflow, panic is thrown.
func (g Gwei) Add(x uint64) Gwei {
	res, err := g.SafeAdd(x)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// SafeAdd increases Gwei by x.
// In case of overflow, error is returned.
func (g Gwei) SafeAdd(x uint64) (Gwei, error) {
	res, err := math.Add64(uint64(g), x)
	return Gwei(res), err
}

// Sub subtracts x from Gwei.
// In case of underflow, panic is thrown.
func (g Gwei) Sub(x uint64) Gwei {
	res, err := g.SafeSub(x)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// SafeSub subtracts x from Gwei.
// In case of underflow, error is returned.
func (g Gwei) SafeSub(x uint64) (Gwei, error) {
	res, err := math.Sub64(uint64(g), x)
	return Gwei(res), err
}

// FlooredSub safely subtracts x from Gwei, returning 0 if the result would underflow.
func (g Gwei) FlooredSub(x Gwei) Gwei {
	if g < x {
		return 0
	}
	return g - x
}

// Mul multiplies Gwei by x.
// In case of overflow, panic is thrown.
func (g Gwei) Mul(x uint64) Gwei {
	res, err := g.SafeMul(x)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// SafeMul multiplies Gwei by x.
// In case of overflow, error is returned.
func (g Gwei) SafeMul(x uint64) (Gwei, error) {
	res, err := math.Mul64(uint64(g), x)
	return Gwei(res), err
}

// Div divides Gwei by x.
// In case of division by zero, panic is thrown.
func (g Gwei) Div(x uint64) Gwei {
	res, err := g.SafeDiv(x)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// SafeDiv divides Gwei by x.
// In case of division by zero, error is returned.
func (g Gwei) SafeDiv(x uint64) (Gwei, error) {
	res, err := math.Div64(uint64(g), x)
	return Gwei(res), err
}

// HashTreeRoot --
func (g Gwei) HashTreeRoot() ([32]byte, error) {
	return fssz.HashWithDefaultHasher(g)
}

// HashTreeRootWith --
func (g Gwei) HashTreeRootWith(hh *fssz.Hasher) error {
	hh.PutUint64(uint64(g))
	return nil
}

// UnmarshalSSZ --
func (g *Gwei) UnmarshalSSZ(buf []byte) error {
	if len(buf) != g.SizeSSZ() {
		return fmt.Errorf("expected buffer of length %d received %d", g.SizeSSZ(), len(buf))
	}
	*g = Gwei(fssz.UnmarshallUint64(buf))
	return nil
}

// MarshalSSZTo --
func (g *Gwei) MarshalSSZTo(dst []byte) ([]byte, error) {
	marshalled, err := g.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return append(dst, marshalled...), nil
}

// MarshalSSZ --
func (g *Gwei) MarshalSSZ() ([]byte, error) {
	marshalled := fssz.MarshalUint64([]byte{}, uint64(*g))
	return marshalled, nil
}

// SizeSSZ --
func (g *Gwei) SizeSSZ() int {
	return 8
}

// WeiToBigInt is a convenience method to cast a wei back to a big int
func WeiToBigInt(w Wei) *big.Int {
	return w
}

// Uint64ToWei creates a new Wei (aka big.Int) representing the given uint64 value.
func Uint64ToWei(v uint64) Wei {
	return big.NewInt(0).SetUint64(v)
}

// LittleEndianBytesToWei returns a Wei value given a little-endian binary representation.
// The only places we use this representation are in protobuf types that hold either the
// local execution payload bid or the builder bid. Going forward we should avoid that representation
// so this function being used in new places should be considered a code smell.
func LittleEndianBytesToWei(value []byte) Wei {
	if len(value) == 0 {
		return big.NewInt(0)
	}
	v := make([]byte, len(value))
	copy(v, value)
	// SetBytes expects a big-endian representation of the value, so we reverse the byte slice.
	slices.Reverse(v)
	return big.NewInt(0).SetBytes(v)
}

// WeiToGwei converts big int wei to uint64 gwei.
// The input `v` is copied before being modified.
func WeiToGwei(v Wei) Gwei {
	if v == nil {
		return 0
	}
	gweiPerEth := big.NewInt(1e9)
	copied := big.NewInt(0).Set(v)
	copied.Div(copied, gweiPerEth)
	return Gwei(copied.Uint64())
}
