package math

// Uint64Based is a constraint for types with uint64 as their underlying type.
// This includes primitives.Slot, primitives.Epoch, primitives.ValidatorIndex, etc.
type Uint64Based interface {
	~uint64
}

// --- Addition ---

// SafeAdd adds two values of the same type, returning error on overflow.
func SafeAdd[T Uint64Based](a, b T) (T, error) {
	res, err := Add64(uint64(a), uint64(b))
	return T(res), err
}

// MustAdd adds two values of the same type, panicking on overflow.
func MustAdd[T Uint64Based](a, b T) T {
	res, err := SafeAdd(a, b)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// --- Subtraction ---

// SafeSub subtracts b from a, returning error on underflow.
func SafeSub[T Uint64Based](a, b T) (T, error) {
	res, err := Sub64(uint64(a), uint64(b))
	return T(res), err
}

// MustSub subtracts b from a, panicking on underflow.
func MustSub[T Uint64Based](a, b T) T {
	res, err := SafeSub(a, b)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// FlooredSub subtracts b from a, returning 0 on underflow (saturating subtraction).
func FlooredSub[T Uint64Based](a, b T) T {
	if a < b {
		return 0
	}
	return a - b
}

// --- Multiplication ---

// SafeMul multiplies two values of the same type, returning error on overflow.
func SafeMul[T Uint64Based](a, b T) (T, error) {
	res, err := Mul64(uint64(a), uint64(b))
	return T(res), err
}

// MustMul multiplies two values of the same type, panicking on overflow.
func MustMul[T Uint64Based](a, b T) T {
	res, err := SafeMul(a, b)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// --- Division ---

// SafeDiv divides a by b, returning error on division by zero.
func SafeDiv[T Uint64Based](a, b T) (T, error) {
	res, err := Div64(uint64(a), uint64(b))
	return T(res), err
}

// MustDiv divides a by b, panicking on division by zero.
func MustDiv[T Uint64Based](a, b T) T {
	res, err := SafeDiv(a, b)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}

// --- Modulo ---

// SafeMod returns a % b, error on division by zero.
func SafeMod[T Uint64Based](a, b T) (T, error) {
	res, err := Mod64(uint64(a), uint64(b))
	return T(res), err
}

// MustMod returns a % b, panicking on division by zero.
func MustMod[T Uint64Based](a, b T) T {
	res, err := SafeMod(a, b)
	if err != nil {
		panic(err.Error()) // lint:nopanic -- Panic is communicated in the godoc commentary.
	}
	return res
}
