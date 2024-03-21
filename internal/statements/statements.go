package statements

type StmtType string

const (
	ExpressionStmtType          StmtType = "Expression"
	PrintStmtType               StmtType = "Print"
	VariableDeclarationStmtType StmtType = "VariableDeclaration"
	FunctionDeclarationStmtType StmtType = "FunctionDeclaration"
	VariableAssignmentStmtType  StmtType = "VariableAssignment"
	BlockStmtType               StmtType = "Block"
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
