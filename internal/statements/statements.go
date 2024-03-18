package statements

type StatementType int

const (
	UNDEFINED StatementType = iota
	EXPRESSION
	PRINT
	VARIABLE_DECLARATION
	FUNCTION_DECLARATION
)

type StmtVisitor interface {
	VisitExpression(ExpressionStmt)
	VisitPrint(PrintStmt)
	VisitVariableDeclaration(VariableDeclarationStmt)
	VisitFunctionDeclaration(FunctionDeclarationStmt)
}

type Stmt interface {
	Accept(StmtVisitor)
	Type() StatementType
}
