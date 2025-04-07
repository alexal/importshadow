package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

// Analyzer defines the import shadow analyzer.
var Analyzer = &analysis.Analyzer{
	Name:             "importshadow",
	Doc:              "Detects variable declarations that shadow imported package names.",
	Run:              runAnalyzer,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
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
