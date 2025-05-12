package type_checker

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

func (t *TypeChecker) VisitBinary(expr ast.BinaryExpr) any {
	// TODO: check type compatibility based on operator
	rtype := ast.AcceptExpr[types.ValueType](expr.Right, t)
	ast.AcceptExpr[types.ValueType](expr.Left, t)
	return rtype
}

func (t *TypeChecker) VisitGrouping(expr ast.GroupingExpr) any {
	return ast.AcceptExpr[types.ValueType](expr.Expr, t)
}

func (t *TypeChecker) VisitLiteral(expr ast.LiteralExpr) any {
	return expr.ValueType
}

func (t *TypeChecker) VisitUnary(expr ast.UnaryExpr) any {
	return ast.AcceptExpr[types.ValueType](expr.Expr, t)
}

func (t *TypeChecker) VisitVariable(expr ast.VariableExpr) any {
	sym, _ := t.table.Lookup(expr.Name.Lexeme)
	return sym.ValueType()
}

func (t *TypeChecker) VisitFieldAccess(expr ast.FieldAccessExpr) any {
	vtype, _ := ast.AcceptExpr[types.ValueType](expr.Expr, t).GetField(expr.Field.Lexeme)
	return vtype
}

func (t *TypeChecker) VisitFunctionCall(expr ast.FunctionCallExpr) any {
	sym, _ := t.table.Lookup(expr.Name.Lexeme)
	funcStmt := sym.DeclarationNode().(ast.FunctionDeclarationStmt)
	if len(expr.Arguments) != len(funcStmt.Parameters) {
		t.error(expr, fmt.Sprintf("function %s expects %d arguments, got %d", expr.Name.Lexeme, len(funcStmt.Parameters), len(expr.Arguments)))
		return sym.ValueType()
	}
	for i, arg := range expr.Arguments {
		ptype := ast.AcceptExpr[types.ValueType](arg, t)
		if ptype != funcStmt.Parameters[i].ValueType {
			t.error(arg, fmt.Sprintf("argument %d of function %s has invalid type: expected %s, got %s", i, expr.Name.Lexeme, funcStmt.Parameters[i].ValueType.ToString(), ptype.ToString()))
		}
	}
	return sym.ValueType()
}

func (t *TypeChecker) VisitLogical(expr ast.LogicalExpr) any {
	// TODO: check type compatibility based on operator
	rvalue := ast.AcceptExpr[types.ValueType](expr.Right, t)
	ast.AcceptExpr[types.ValueType](expr.Left, t)
	return rvalue
}

func (t *TypeChecker) VisitSlice(expr ast.SliceExpr) any {
	// TODO: fix for strings or lists
	targetType := ast.AcceptExpr[types.ValueType](expr.TargetExpr, t)
	sIndexType := ast.AcceptExpr[types.ValueType](expr.StartIndex, t)
	if sIndexType != types.IntType {
		t.error(expr.StartIndex, fmt.Sprintf("index must be int, got %s", sIndexType.ToString()))
	}
	if expr.EndIndex != nil {
		eIndexType := ast.AcceptExpr[types.ValueType](expr.EndIndex, t)
		if eIndexType != types.IntType {
			t.error(expr.EndIndex, fmt.Sprintf("index must be int, got %s", eIndexType.ToString()))
		}
	}
	if targetType == types.StringType {
		return types.StringType
	} else if listType, ok := targetType.(*types.ListTypeStruct); ok {
		return listType.ContentType
	}
	t.error(expr.TargetExpr, fmt.Sprintf("target must be string or list, got %s", targetType.ToString()))
	return types.VoidType
}

func (t *TypeChecker) VisitList(expr ast.ListExpr) any {
	// TODO: fix
	for _, item := range expr.Elements {
		err := ast.AcceptExpr[types.ValueType](item, t)
		if err != nil {
			return err
		}
	}
	return types.VoidType
}
