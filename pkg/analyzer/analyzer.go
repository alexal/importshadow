package analyzer

import (
	"fmt"
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
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch stmt := node.(type) {
	case *ast.ImportSpec:
		fmt.Println(strings.Trim(stmt.Path.Value, "\""))
	case *ast.AssignStmt:
		fmt.Println(stmt.Lhs[0])
	}
	return v
}

func runAnalyzer(pass *analysis.Pass) (interface{}, error) {
	/*ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.ImportSpec)(nil),
		(*ast.AssignStmt)(nil),
	}
	ins.Preorder(nodeFilter, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.ImportSpec:
			fmt.Println(n.Path)
		case *ast.AssignStmt:
			fmt.Println(n.Lhs[0])
		}
	})*/
	v := &visitor{}
	for _, f := range pass.Files {
		fmt.Println(f.Name)
		ast.Walk(v, f)
	}
	return nil, nil
}
