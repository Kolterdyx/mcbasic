package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VarDef struct {
	Name string
	Type interfaces.ValueType
}

type VariableDeclarationStmt struct {
	Stmt

	Name        tokens.Token
	Type        interfaces.ValueType
	Initializer expressions.Expr
}

func (v VariableDeclarationStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVariableDeclaration(v)
}

func (v VariableDeclarationStmt) TType() StmtType {
	return VariableDeclarationStmtType
}
