package analyzer

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"path/filepath"
	"testing"
)

func TestImportShadowAnalyzer(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")

	analyser := NewAnalyzer(&Config{
		IgnoreDirs: []string{"mocks"},
	})
	analysistest.Run(t, testdata, analyser, "p")
}
