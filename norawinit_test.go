package norawinit_test

import (
	"testing"

	"github.com/joehattori/norawinit"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, norawinit.Analyzer, "b")
}
