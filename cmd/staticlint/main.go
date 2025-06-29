package main

import (
	"go/ast"

	"github.com/charithe/durationcheck"
	"github.com/gostaticanalysis/sqlrows/passes/sqlrows"
	"github.com/timakin/bodyclose/passes/bodyclose"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/ast/inspector"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Analyzer Анализатор запрещает прямой вызов os.Exit в функции main пакета main
var analyzer = &analysis.Analyzer{
	Name:     "noexit",
	Doc:      "запрет прямого вызова os.Exit в функции main пакета main",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func main() {
	checks := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}

	// SA из staticcheck
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	// S (simple)
	checks = append(checks, simple.Analyzers[0].Analyzer) // S1000
	checks = append(checks, simple.Analyzers[1].Analyzer) // s1001

	// ST (stylecheck)
	checks = append(checks, stylecheck.Analyzers[0].Analyzer) // ST1000
	checks = append(checks, stylecheck.Analyzers[1].Analyzer) // ST1001

	// Публичные анализаторы
	checks = append(checks,
		bodyclose.Analyzer,
		durationcheck.Analyzer,
		sqlrows.Analyzer,
	)

	// noexit
	checks = append(checks, analyzer)

	multichecker.Main(checks...)
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder(nodeFilter, func(n ast.Node) {
		fNode, ok := n.(*ast.FuncDecl)
		if !ok || fNode.Name.Name != "main" {
			return
		}

		ast.Inspect(fNode.Body, func(node ast.Node) bool {
			cNode, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			if isOsExit(cNode) {
				pass.Reportf(cNode.Pos(), "прямой вызов os.Exit в функции main запрещен")
			}

			return true
		})
	})

	return nil, nil
}

// isOsExit это вызов os.Exit?
func isOsExit(call *ast.CallExpr) bool {
	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if selExpr.Sel.Name != "Exit" {
		return false
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "os"
}
