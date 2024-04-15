package statements

type StmtType string

const (
	ExpressionStmtType          StmtType = "Expression"
	VariableDeclarationStmtType StmtType = "VariableDeclaration"
	FunctionDeclarationStmtType StmtType = "FunctionDeclaration"
	VariableAssignmentStmtType  StmtType = "VariableAssignment"
	BlockStmtType               StmtType = "Block"
	WhileStmtType               StmtType = "While"
	IfStmtType                  StmtType = "If"
	ReturnStmtType              StmtType = "Return"
)

type StmtVisitor interface {
	VisitExpression(ExpressionStmt) interface{}
	VisitVariableDeclaration(VariableDeclarationStmt) interface{}
	VisitFunctionDeclaration(FunctionDeclarationStmt) interface{}
	VisitVariableAssignment(VariableAssignmentStmt) interface{}
	VisitBlock(BlockStmt) interface{}
	VisitWhile(WhileStmt) interface{}
	VisitIf(IfStmt) interface{}
	VisitReturn(ReturnStmt) interface{}
}

type Stmt interface {
	Accept(StmtVisitor) interface{}
	TType() StmtType
}
