package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name:             "govarpkg",
	Doc:              "govarpkg",
	Run:              runAnalyzer,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
}

type visitor struct {
	importSpec map[string]*ast.Node
	assignStmt map[string]*ast.Ident
}

func (v *visitor) walk(n ast.Node) {
	if n != nil {
		ast.Walk(v, n)
	}
}

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
		switch va := stmt.Lhs[0].(type) {
		case *ast.SelectorExpr:
			//TODO?
		case *ast.Ident:
			v.assignStmt[va.Name] = va
		}
	}
	return v
}

func runAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		v := &visitor{
			importSpec: make(map[string]*ast.Node),
			assignStmt: make(map[string]*ast.Ident),
		}
		v.walk(f)

		for k, e := range v.assignStmt {
			if v.importSpec[k] != nil {
				pass.Reportf(e.Pos(), "Variable '%s' collides with imported package name", k)
			}
		}
	}
	return nil, nil
}
