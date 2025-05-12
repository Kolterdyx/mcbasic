package type_checker

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
)

type TypeChecker struct {
	ast.StatementVisitor
	ast.ExpressionVisitor

	table  *symbol.Table
	source []ast.Statement
	errors []error
}

func NewTypeChecker(source []ast.Statement, table *symbol.Table) *TypeChecker {
	return &TypeChecker{
		table:  table,
		source: source,
		errors: make([]error, 0),
	}
}

func (t *TypeChecker) Resolve() []error {
	for _, stmt := range t.source {
		ast.AcceptStmt[any](stmt, t)
	}

	return t.errors
}

func (t *TypeChecker) error(expr ast.Node, message string) any {
	err := fmt.Errorf("error at %s: %s", expr.GetSourceLocation().ToString(), message)
	t.errors = append(t.errors, err)
	return err
}
