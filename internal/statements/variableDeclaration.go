package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type VarDef struct {
	Name string
	Type types.ValueType
}

type VariableDeclarationStmt struct {
	Stmt

	Name        tokens.Token
	Type        types.ValueType
	Initializer expressions.Expr
}

func (v VariableDeclarationStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitVariableDeclaration(v)
}

func (v VariableDeclarationStmt) StmtType() StmtType {
	return VariableDeclarationStmtType
}
