package statements

import "github.com/Kolterdyx/mcbasic/internal/tokens"

type FunctionDeclarationStmt struct {
	Stmt

	Name       tokens.Token
	Parameters []tokens.Token
	Body       BlockStmt
}

func (f FunctionDeclarationStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitFunctionDeclaration(f)
}

func (f FunctionDeclarationStmt) Type() StmtType {
	return FunctionDeclarationStmtType
}
