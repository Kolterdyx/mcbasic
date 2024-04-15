package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VarDef struct {
	Name string
	Type tokens.TokenType
}

type VariableDeclarationStmt struct {
	Stmt

	Name        tokens.Token
	Type        tokens.TokenType
	Initializer expressions.Expr
}

func (v VariableDeclarationStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVariableDeclaration(v)
}

func (v VariableDeclarationStmt) TType() StmtType {
	return VariableDeclarationStmtType
}
