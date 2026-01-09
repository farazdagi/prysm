// Package primarithm implements a static analyzer that detects raw arithmetic
// operations on primitive types (Slot, Epoch, ValidatorIndex, etc.) and suggests
// using the type's safe arithmetic methods instead.
//
// Raw arithmetic on these uint64-based types can cause integer overflow/underflow
// that wraps around silently, leading to security issues. For example:
//
//	count := uint64(end - begin)  // If end < begin, this wraps to a huge number
//
// Instead, use the safe methods:
//
//	count, err := end.SafeSubSlot(begin)  // Returns error on underflow
//	count := end.FlooredSubSlot(begin)    // Returns 0 on underflow
package primarithm

import (
	"errors"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/comment"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// smallAdditionThreshold is the maximum constant value for which addition is considered safe.
// For `slot + N` to overflow, slot must be >= MaxUint64 - N, which is unrealistic for small N.
const smallAdditionThreshold = 1000

// Doc explaining the tool.
const Doc = "Detects raw arithmetic (+, -, *, /, %) on primitive types (Slot, Epoch, ValidatorIndex, Gwei, etc.) " +
	"which can cause silent integer overflow/underflow. Use the type's safe arithmetic methods instead. " +
	"This check can be suppressed with the `lint:ignore primarithm` comment with proper justification."

// Analyzer runs static analysis.
var Analyzer = &analysis.Analyzer{
	Name:     "primarithm",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// typeInfo contains information about a target primitive type.
type typeInfo struct {
	// suggestions maps arithmetic operators to suggested safe method(s).
	suggestions map[token.Token]string
}

// primitivesPackage is the import path for the primitives package.
const primitivesPackage = "github.com/OffchainLabs/prysm/v7/consensus-types/primitives"

// targetTypes maps fully qualified type names to their safe method suggestions.
// Types without safe methods will still be flagged (suggesting generic math.Safe* functions).
var targetTypes = map[string]typeInfo{
	"Slot": {
		suggestions: map[token.Token]string{
			token.ADD: "SafeAdd/Add/SafeAddSlot/AddSlot",
			token.SUB: "SafeSub/Sub/SafeSubSlot/SubSlot/FlooredSubSlot",
			token.MUL: "SafeMul/Mul/SafeMulSlot/MulSlot",
			token.QUO: "SafeDiv/Div/SafeDivSlot/DivSlot",
			token.REM: "SafeMod/Mod/SafeModSlot/ModSlot",
		},
	},
	"Epoch": {
		suggestions: map[token.Token]string{
			token.ADD: "SafeAdd/Add/SafeAddEpoch/AddEpoch",
			token.SUB: "SafeSub/Sub/FlooredSubEpoch",
			token.MUL: "SafeMul/Mul",
			token.QUO: "SafeDiv/Div",
			token.REM: "SafeMod/Mod",
		},
	},
	"ValidatorIndex": {
		suggestions: map[token.Token]string{
			token.ADD: "SafeAdd/Add",
			token.SUB: "SafeSub/Sub/FlooredSub",
			token.MUL: "math.SafeMul/math.MustMul",
			token.QUO: "SafeDiv/Div",
			token.REM: "SafeMod/Mod",
		},
	},
	"Gwei": {
		suggestions: map[token.Token]string{
			token.ADD: "SafeAdd/Add",
			token.SUB: "SafeSub/Sub/FlooredSub",
			token.MUL: "SafeMul/Mul",
			token.QUO: "SafeDiv/Div",
			token.REM: "math.SafeMod/math.MustMod",
		},
	},
	"CommitteeIndex": {
		suggestions: map[token.Token]string{
			token.ADD: "math.SafeAdd/math.MustAdd",
			token.SUB: "math.SafeSub/math.MustSub/math.FlooredSub",
			token.MUL: "math.SafeMul/math.MustMul",
			token.QUO: "math.SafeDiv/math.MustDiv",
			token.REM: "math.SafeMod/math.MustMod",
		},
	},
	"BuilderIndex": {
		suggestions: map[token.Token]string{
			token.ADD: "math.SafeAdd/math.MustAdd",
			token.SUB: "math.SafeSub/math.MustSub/math.FlooredSub",
			token.MUL: "math.SafeMul/math.MustMul",
			token.QUO: "math.SafeDiv/math.MustDiv",
			token.REM: "math.SafeMod/math.MustMod",
		},
	},
	"SSZUint64": {
		suggestions: map[token.Token]string{
			token.ADD: "math.SafeAdd/math.MustAdd",
			token.SUB: "math.SafeSub/math.MustSub/math.FlooredSub",
			token.MUL: "math.SafeMul/math.MustMul",
			token.QUO: "math.SafeDiv/math.MustDiv",
			token.REM: "math.SafeMod/math.MustMod",
		},
	},
	"BP": {
		suggestions: map[token.Token]string{
			token.ADD: "math.SafeAdd/math.MustAdd",
			token.SUB: "math.SafeSub/math.MustSub/math.FlooredSub",
			token.MUL: "math.SafeMul/math.MustMul",
			token.QUO: "math.SafeDiv/math.MustDiv",
			token.REM: "math.SafeMod/math.MustMod",
		},
	},
}

func run(pass *analysis.Pass) (any, error) {
	// Skip the primitives package itself - safe methods use raw arithmetic internally.
	if strings.HasSuffix(pass.Pkg.Path(), "consensus-types/primitives") {
		return nil, nil
	}

	// Skip the math package - generic safe helpers use raw arithmetic internally.
	if strings.HasSuffix(pass.Pkg.Path(), "/math") {
		return nil, nil
	}

	inspection, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	nodeFilter := []ast.Node{
		(*ast.BinaryExpr)(nil),
	}

	commentMap := comment.New(pass.Fset, pass.Files)

	inspection.Preorder(nodeFilter, func(node ast.Node) {
		binExpr := node.(*ast.BinaryExpr)

		// Only flag arithmetic operators.
		if !isArithmeticOp(binExpr.Op) {
			return
		}

		// Check for lint:ignore comment.
		cg := commentMap.CommentsByPosLine(pass.Fset, node.Pos())
		for _, c := range cg {
			if strings.Contains(c.Text(), "lint:ignore primarithm") {
				return
			}
		}

		// Check if either operand is a target primitive type.
		leftType := pass.TypesInfo.TypeOf(binExpr.X)
		rightType := pass.TypesInfo.TypeOf(binExpr.Y)

		typeName, info, found := getTargetTypeInfo(leftType, rightType)
		if !found {
			return
		}

		// Skip safe patterns (e.g., adding small constants, dividing by non-zero constants).
		if shouldSkipCheck(binExpr, pass.TypesInfo) {
			return
		}

		suggestion := info.suggestions[binExpr.Op]
		if suggestion == "" {
			suggestion = "a safe arithmetic method"
		}

		pass.Reportf(binExpr.Pos(),
			"unsafe arithmetic on primitives.%s; use %s instead to prevent overflow/underflow",
			typeName, suggestion)
	})

	return nil, nil
}

// isArithmeticOp returns true if the operator is an arithmetic operator.
func isArithmeticOp(op token.Token) bool {
	switch op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		return true
	}
	return false
}

// getConstantUint64 returns the uint64 value of an expression if it's a compile-time constant.
// This handles numeric literals, named constants, and constant expressions.
func getConstantUint64(expr ast.Expr, info *types.Info) (uint64, bool) {
	tv, ok := info.Types[expr]
	if !ok || tv.Value == nil {
		return 0, false // Not a constant
	}
	if tv.Value.Kind() != constant.Int {
		return 0, false // Not an integer constant
	}
	val, exact := constant.Uint64Val(tv.Value)
	if !exact {
		return 0, false // Value doesn't fit in uint64 or is negative
	}
	return val, true
}

// shouldSkipCheck returns true if the arithmetic operation is safe and should not be flagged.
// This reduces false positives for patterns that cannot realistically cause overflow/underflow.
func shouldSkipCheck(binExpr *ast.BinaryExpr, info *types.Info) bool {
	switch binExpr.Op {
	case token.ADD:
		// Addition with small constant is safe - overflow needs other operand ≈ MaxUint64.
		if v, ok := getConstantUint64(binExpr.X, info); ok && v <= smallAdditionThreshold {
			return true
		}
		if v, ok := getConstantUint64(binExpr.Y, info); ok && v <= smallAdditionThreshold {
			return true
		}

	case token.SUB:
		// Only skip subtracting literal 0 (no-op).
		// Note: We still flag `slot - 1` because slot could be 0 (genesis), causing underflow.
		if v, ok := getConstantUint64(binExpr.Y, info); ok && v == 0 {
			return true
		}

	case token.MUL:
		// Skip multiplication by 0 (always 0) or 1 (identity).
		// Note: We still flag `* 2` and higher because even small multipliers can overflow
		// with large values (e.g., gwei * 1e9 for wei conversion).
		if v, ok := getConstantUint64(binExpr.X, info); ok && v <= 1 {
			return true
		}
		if v, ok := getConstantUint64(binExpr.Y, info); ok && v <= 1 {
			return true
		}

	case token.QUO, token.REM:
		// Division/modulo by non-zero constant is safe (no div-by-zero risk).
		if v, ok := getConstantUint64(binExpr.Y, info); ok && v != 0 {
			return true
		}
	}

	return false
}

// getTargetTypeInfo checks if either type is a target primitive type.
// Returns the type name, type info, and whether a match was found.
func getTargetTypeInfo(left, right types.Type) (string, typeInfo, bool) {
	for _, t := range []types.Type{left, right} {
		if t == nil {
			continue
		}

		// Handle named types (e.g., primitives.Slot).
		if named, ok := t.(*types.Named); ok {
			obj := named.Obj()
			if obj == nil || obj.Pkg() == nil {
				continue
			}

			// Check if this is from the primitives package.
			pkgPath := obj.Pkg().Path()
			if pkgPath != primitivesPackage {
				continue
			}

			typeName := obj.Name()
			if info, ok := targetTypes[typeName]; ok {
				return typeName, info, true
			}
		}
	}
	return "", typeInfo{}, false
}
