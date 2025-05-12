package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

// Accessor is one step in an L-value chain: either [expr] or .field
type Accessor interface {
	ToString() string
}

// IndexAccessor holds a parsed index expression (foo[expr]).
type IndexAccessor struct {
	Index Expr
}

func (i IndexAccessor) ToString() string {
	return "[x]"
}

// FieldAccessor holds a parsed field name (foo.bar).
type FieldAccessor struct {
	Field tokens.Token
}

func (f FieldAccessor) ToString() string {
	return "." + f.Field.Lexeme
}

type VariableAssignmentStmt struct {
	Statement

	Name      tokens.Token
	Accessors []Accessor
	Value     Expr
}

func (v VariableAssignmentStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitVariableAssignment(v)
}

func (v VariableAssignmentStmt) Type() NodeType {
	return VariableAssignmentStatement
}

func (v VariableAssignmentStmt) ToString() string {
	body := v.Name.Lexeme
	for _, accessor := range v.Accessors {
		body += accessor.ToString()
	}
	body += " = " + v.Value.ToString() + ";"
	return body
}

func (v VariableAssignmentStmt) GetSourceLocation() interfaces.SourceLocation {
	return v.Name.SourceLocation
}
