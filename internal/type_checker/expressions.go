package type_checker

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

func (t *TypeChecker) VisitBinary(expr ast.BinaryExpr) any {
	// TODO: check type compatibility based on operator
	rtype := ast.AcceptExpr[types.ValueType](expr.Right, t)
	ltype := ast.AcceptExpr[types.ValueType](expr.Left, t)
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
	sym, ok := t.table.Lookup(expr.Name.Lexeme)
	if !ok {
		t.error(expr, fmt.Sprintf("variable %s not defined", expr.Name.Lexeme))
	}
	return sym.ValueType()
}

func (t *TypeChecker) VisitFieldAccess(expr ast.FieldAccessExpr) any {
	vtype, ok := ast.AcceptExpr[types.ValueType](expr.Expr, t).GetField(expr.Field.Lexeme)
	if !ok {
		t.error(expr, fmt.Sprintf("field %s not found in %s", expr.Field.Lexeme, expr.Expr.ToString()))
	}
	return vtype
}

func (t *TypeChecker) VisitFunctionCall(expr ast.FunctionCallExpr) any {
	// TODO: check parameter types
	sym, ok := t.table.Lookup(expr.Name.Lexeme)
	if !ok {
		t.error(expr, fmt.Sprintf("function %s not defined", expr.Name.Lexeme))
	}
	return sym.ValueType()
}

func (t *TypeChecker) VisitLogical(expr ast.LogicalExpr) any {
	// TODO: check type compatibility based on operator
	rvalue := ast.AcceptExpr[types.ValueType](expr.Right, t)
	lvalue := ast.AcceptExpr[types.ValueType](expr.Left, t)
	return rvalue
}

func (t *TypeChecker) VisitSlice(expr ast.SliceExpr) any {
	// TODO: fix for strings or lists
	err := ast.AcceptExpr[types.ValueType](expr.TargetExpr, t)
	if err != nil {
		return err
	}
	err = ast.AcceptExpr[types.ValueType](expr.StartIndex, t)
	if err != nil {
		return err
	}
	return ast.AcceptExpr[types.ValueType](expr.EndIndex, t)
}

func (t *TypeChecker) VisitList(expr ast.ListExpr) any {
	// TODO: fix
	for _, item := range expr.Elements {
		err := ast.AcceptExpr[types.ValueType](item, t)
		if err != nil {
			return err
		}
	}
	return nil
}
