package statements

import "github.com/Kolterdyx/mcbasic/internal/tokens"

type FunctionDeclarationStmt struct {
	Stmt

	Name       tokens.Token
	Parameters []tokens.Token
	Body       []Stmt
}

func (f FunctionDeclarationStmt) Accept(visitor StmtVisitor) {
	visitor.VisitFunctionDeclaration(f)
}

func (f FunctionDeclarationStmt) Type() StatementType {
	return FUNCTION_DECLARATION
}
