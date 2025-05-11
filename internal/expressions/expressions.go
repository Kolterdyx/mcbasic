package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type ExprVisitor interface {
	VisitBinary(b BinaryExpr) interfaces.IRCode
	VisitGrouping(g GroupingExpr) interfaces.IRCode
	VisitLiteral(l LiteralExpr) interfaces.IRCode
	VisitUnary(u UnaryExpr) interfaces.IRCode
	VisitVariable(v VariableExpr) interfaces.IRCode
	VisitFieldAccess(v FieldAccessExpr) interfaces.IRCode
	VisitFunctionCall(f FunctionCallExpr) interfaces.IRCode
	VisitLogical(l LogicalExpr) interfaces.IRCode
	VisitSlice(s SliceExpr) interfaces.IRCode
	VisitList(s ListExpr) interfaces.IRCode
	VisitStruct(s StructExpr) interfaces.IRCode
}

type Expr interface {
	ast.Node
	Accept(v ExprVisitor) interfaces.IRCode
	ReturnType() types.ValueType
	ToString() string
	Validate() error
}
