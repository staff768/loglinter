package analyzer


import (
	"testing"
	"path/filepath"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T){
	testdata, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatalf("failed to get absolute path to testdata: %v", err)
	}
	analysistest.Run(t, testdata, Analyzer, "a")
}