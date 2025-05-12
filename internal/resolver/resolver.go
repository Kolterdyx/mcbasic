package resolver

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
)

type Resolver struct {
	ast.StatementVisitor
	ast.ExpressionVisitor

	table  *symbol.Table
	source []ast.Statement
	errors []error
}

func NewResolver(source []ast.Statement, table *symbol.Table) *Resolver {
	return &Resolver{
		table:  table,
		source: source,
		errors: make([]error, 0),
	}
}

func (r *Resolver) Resolve() []error {
	for _, stmt := range r.source {
		ast.AcceptStmt[any](stmt, r)
	}

	return r.errors
}

func (r *Resolver) error(expr ast.Node, message string) any {
	err := fmt.Errorf("error: %s", message)
	r.errors = append(r.errors, err)
	return err
}
