package resolver

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

func (r *Resolver) VisitBinary(expr ast.BinaryExpr) any {
	err := ast.AcceptExpr[error](expr.Left, r)
	if err != nil {
		return err
	}
	return ast.AcceptExpr[error](expr.Right, r)
}

func (r *Resolver) VisitGrouping(expr ast.GroupingExpr) any {
	return ast.AcceptExpr[error](expr.Expr, r)
}

func (r *Resolver) VisitLiteral(expr ast.LiteralExpr) any {
	return nil
}

func (r *Resolver) VisitUnary(expr ast.UnaryExpr) any {
	return ast.AcceptExpr[error](expr.Expr, r)
}

func (r *Resolver) VisitVariable(expr ast.VariableExpr) any {
	_, ok := r.table.Lookup(expr.Name.Lexeme)
	if !ok {
		return r.error(expr, fmt.Sprintf("variable %s not defined", expr.Name.Lexeme))
	}
	return nil
}

func (r *Resolver) VisitFieldAccess(expr ast.FieldAccessExpr) any {
	_, ok := ast.AcceptExpr[types.ValueType](expr.Expr, r).GetField(expr.Field.Lexeme)
	if !ok {
		return r.error(expr, fmt.Sprintf("field %s not defined", expr.Field.Lexeme))
	}
	return nil
}

func (r *Resolver) VisitFunctionCall(expr ast.FunctionCallExpr) any {
	_, ok := r.table.Lookup(expr.Name.Lexeme)
	if !ok {
		return fmt.Errorf("function %s not defined", expr.Name.Lexeme)
	}
	return nil
}

func (r *Resolver) VisitLogical(expr ast.LogicalExpr) any {
	err := ast.AcceptExpr[error](expr.Left, r)
	if err != nil {
		return err
	}
	return ast.AcceptExpr[error](expr.Right, r)
}

func (r *Resolver) VisitSlice(expr ast.SliceExpr) any {
	err := ast.AcceptExpr[error](expr.TargetExpr, r)
	if err != nil {
		return err
	}
	err = ast.AcceptExpr[error](expr.StartIndex, r)
	if err != nil {
		return err
	}
	return ast.AcceptExpr[error](expr.EndIndex, r)
}

func (r *Resolver) VisitList(expr ast.ListExpr) any {
	for _, item := range expr.Elements {
		err := ast.AcceptExpr[error](item, r)
		if err != nil {
			return err
		}
	}
	return nil
}
