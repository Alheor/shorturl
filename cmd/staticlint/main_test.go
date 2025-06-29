package main

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

type tmpExpr struct {
	ast.Expr
}

// TestAnalyzer проверяет работу анализатора noexit.
// Структура тестов:
// - src/a/test.go - пакет main с функцией main (должны быть ошибки)
// - src/b/lib.go - пакет не main (не должно быть ошибок)
// - src/c/util.go - пакет main без функции main (не должно быть ошибок)
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "a", "b", "c")
}

func TestIsOsExitTrue(t *testing.T) {
	exp := ast.CallExpr{Fun: &ast.SelectorExpr{Sel: &ast.Ident{Name: `Exit`}, X: &ast.Ident{Name: `os`}}}
	assert.True(t, isOsExit(&exp))
}

func TestIsOsExitFalse(t *testing.T) {
	exp := ast.CallExpr{Fun: &ast.SelectorExpr{Sel: &ast.Ident{Name: `none`}, X: &ast.Ident{Name: `os`}}}
	assert.False(t, isOsExit(&exp))

	exp = ast.CallExpr{Fun: &ast.SelectorExpr{Sel: &ast.Ident{Name: `Exit`}, X: &ast.Ident{Name: `none`}}}
	assert.False(t, isOsExit(&exp))

	exp = ast.CallExpr{Fun: &ast.SelectorExpr{Sel: &ast.Ident{Name: `Exit`}, X: &tmpExpr{}}}
	assert.False(t, isOsExit(&exp))

	exp = ast.CallExpr{Fun: &tmpExpr{}}
	assert.False(t, isOsExit(&exp))
}
