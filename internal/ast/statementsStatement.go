package ast

type StatementVisitor interface {
	VisitExpression(ExpressionStmt) any
	VisitVariableDeclaration(VariableDeclarationStmt) any
	VisitFunctionDeclaration(FunctionDeclarationStmt) any
	VisitVariableAssignment(VariableAssignmentStmt) any
	VisitStructDeclaration(StructDeclarationStmt) any
	VisitBlock(BlockStmt) any
	VisitWhile(WhileStmt) any
	VisitIf(IfStmt) any
	VisitReturn(ReturnStmt) any
	VisitSetReturnFlag(SetReturnFlagStmt) any
	VisitImport(ImportStmt) any
}

type Statement interface {
	Node
	Accept(StatementVisitor) any
	ToString() string
}

func AcceptStmt[T any](stmt Statement, v StatementVisitor) T {
	return stmt.Accept(v).(T)
}
