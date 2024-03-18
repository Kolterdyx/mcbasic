package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VariableDeclarationStmt struct {
	Stmt

	Name        tokens.Token
	Initializer expressions.Expr
}

func (v VariableDeclarationStmt) Accept(visitor StmtVisitor) {
	visitor.VisitVariableDeclaration(v)
}

func (v VariableDeclarationStmt) Type() StatementType {
	return VARIABLE_DECLARATION
}
