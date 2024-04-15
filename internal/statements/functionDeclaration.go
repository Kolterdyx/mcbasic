package statements

import "github.com/Kolterdyx/mcbasic/internal/tokens"

type FuncArg struct {
	Name string
	Type tokens.TokenType
}

type FuncDef struct {
	Name       string
	Parameters []FuncArg
	ReturnType tokens.TokenType
}

type FunctionDeclarationStmt struct {
	Stmt

	Name       tokens.Token
	Parameters []FuncArg
	ReturnType tokens.TokenType
	Body       BlockStmt
}

func (f FunctionDeclarationStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitFunctionDeclaration(f)
}

func (f FunctionDeclarationStmt) TType() StmtType {
	return FunctionDeclarationStmtType
}

func (f FunctionDeclarationStmt) HasArg(arg string) bool {
	for _, p := range f.Parameters {
		if p.Name == arg {
			return true
		}
	}
	return false
}
