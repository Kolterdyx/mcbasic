package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

// Accessor is one step in an L-value chain: either [expr] or .field
type Accessor interface {
	ToString() string
}

// IndexAccessor holds a parsed index expression (foo[expr]).
type IndexAccessor struct {
	Index expressions.Expr
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
	Stmt

	Name      tokens.Token
	Accessors []Accessor
	Value     expressions.Expr
}

func (v VariableAssignmentStmt) Accept(visitor StmtVisitor) string {
	return visitor.VisitVariableAssignment(v)
}

func (v VariableAssignmentStmt) StmtType() StmtType {
	return VariableAssignmentStmtType
}
