package resolver

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
)

type Resolver struct {
	ast.StatementVisitor
	ast.ExpressionVisitor

	table  *symbol.Table
	source []ast.Statement
}

func NewResolver(source []ast.Statement, table *symbol.Table) *Resolver {
	return &Resolver{
		table:  table,
		source: source,
	}
}

func (r *Resolver) Resolve() []error {
	var errors []error

	for _, stmt := range r.source {
		if err := ast.AcceptStmt[error](stmt, r); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
