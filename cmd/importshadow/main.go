package main

import (
	"flag"
	"github.com/alexal/importshadow/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
	"strings"
)

func main() {
	ignoreDirs := flag.String("ignoreDirs", "mocks", "Comma-separated list of directories to ignore")
	flag.Parse()

	// Split the ignoreDirs flag into a slice
	ignoreDirsSlice := strings.Split(*ignoreDirs, ",")

	// Create a Config instance
	config := &analyzer.Config{
		IgnoreDirs: ignoreDirsSlice,
	}

	// Pass the config to the analyzer
	importShadowAnalyzer := analyzer.NewAnalyzer(config)

	singlechecker.Main(importShadowAnalyzer)
}
