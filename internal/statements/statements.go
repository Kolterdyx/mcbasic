package statements

type StmtType int

const (
	_ StmtType = iota
	ExpressionStmtType
	PrintStmtType
	VariableDeclarationStmtType
	FunctionDeclarationStmtType
	VariableAssignmentStmtType
	BlockStmtType
)

type StmtVisitor interface {
	VisitExpression(ExpressionStmt) interface{}
	VisitPrint(PrintStmt) interface{}
	VisitVariableDeclaration(VariableDeclarationStmt) interface{}
	VisitFunctionDeclaration(FunctionDeclarationStmt) interface{}
	VisitVariableAssignment(VariableAssignmentStmt) interface{}
	VisitBlock(BlockStmt) interface{}
}

type Stmt interface {
	Accept(StmtVisitor) interface{}
	Type() StmtType
}
