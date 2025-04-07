package analyzer

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"path/filepath"
	"strings"
)

// Config holds the configuration for the analyzer.
type Config struct {
	IgnoreDirs []string
}

// Global variable to hold the configuration
var config *Config

// SetConfig sets the configuration for the analyzer.
func SetConfig(cfg *Config) {
	config = cfg
}

// Analyzer defines the import shadow analyzer.
var Analyzer = &analysis.Analyzer{
	Name:             "importshadow",
	Doc:              "Detects variable declarations that shadow imported package names.",
	Run:              runAnalyzer,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
}

func NewAnalyzer(config *Config) *analysis.Analyzer {
	SetConfig(config)
	return Analyzer
}

// visitor holds the import specifications and assignment statements.
type visitor struct {
	importSpec map[string]*ast.Node
	assignStmt map[string]*ast.Ident
}

// Visit implements the ast.Visitor interface.
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch stmt := node.(type) {
	case *ast.ImportSpec:
		if stmt.Name != nil {
			v.importSpec[stmt.Name.Name] = &node
		} else {
			imp := strings.Split(strings.Trim(stmt.Path.Value, "\""), "/")
			v.importSpec[imp[len(imp)-1]] = &node
		}
	case *ast.AssignStmt:
		for _, lhs := range stmt.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				v.assignStmt[ident.Name] = ident
			}
		}
	}
	return v
}

// runAnalyzer runs the import shadow analyzer.
func runAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		// Skip files in ignored directories
		if shouldIgnoreFile(pass.Fset.File(f.Pos()).Name(), config.IgnoreDirs) {
			continue
		}

		v := &visitor{
			importSpec: make(map[string]*ast.Node),
			assignStmt: make(map[string]*ast.Ident),
		}
		ast.Walk(v, f)

		for name, ident := range v.assignStmt {
			if v.importSpec[name] != nil {
				pass.Reportf(ident.Pos(), "Variable '%s' collides with imported package name", name)
			}
		}
	}
	return nil, nil
}

// shouldIgnoreFile checks if a file should be ignored based on the directory.
func shouldIgnoreFile(filePath string, ignoreDirs []string) bool {
	for _, dir := range ignoreDirs {
		if strings.Contains(filepath.ToSlash(filePath), fmt.Sprintf("/%s/", dir)) {
			return true
		}
	}
	return false
}
