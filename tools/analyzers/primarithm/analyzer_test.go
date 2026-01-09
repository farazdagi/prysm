package primarithm_test

import (
	"testing"

	"github.com/OffchainLabs/prysm/v7/build/bazel"
	"github.com/OffchainLabs/prysm/v7/tools/analyzers/primarithm"
	"golang.org/x/tools/go/analysis/analysistest"
)

func init() {
	if bazel.BuiltWithBazel() {
		bazel.SetGoEnv()
	}
}

func TestAnalyzer(t *testing.T) {
	testdata := bazel.TestDataPath(t)
	analysistest.Run(t, testdata, primarithm.Analyzer, "testpkg")
}
