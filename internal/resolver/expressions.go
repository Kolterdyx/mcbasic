package resolver

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
)

func (r *Resolver) VisitBinary(expr ast.BinaryExpr) any {
	res := ast.AcceptExpr[Result](expr.Left, r)
	if !res.Ok {
		return res
	}
	return ast.AcceptExpr[Result](expr.Right, r)
}

func (r *Resolver) VisitGrouping(expr ast.GroupingExpr) any {
	return ast.AcceptExpr[Result](expr.Expression, r)
}

func (r *Resolver) VisitLiteral(expr ast.LiteralExpr) any {
	return Result{
		Ok:     true,
		Symbol: symbol.NewSymbol("__literal__", symbol.LiteralSymbol, expr, expr.ValueType),
	}
}

func (r *Resolver) VisitUnary(expr ast.UnaryExpr) any {
	return ast.AcceptExpr[Result](expr.Expression, r)
}

func (r *Resolver) VisitVariable(expr ast.VariableExpr) any {
	vtype, ok := r.table.Lookup(expr.Name.Lexeme)
	if !ok {
		return r.error(expr, fmt.Sprintf("variable %s not defined", expr.Name.Lexeme))
	}
	return Result{
		Ok:     true,
		Symbol: vtype,
	}
}

func (r *Resolver) VisitFieldAccess(expr ast.FieldAccessExpr) any {
	res := ast.AcceptExpr[Result](expr.Source, r)
	if !res.Ok {
		return r.error(expr, fmt.Sprintf("could not resolve source %s", expr.Source.ToString()))
	}
	valueType := res.Symbol.ValueType()
	fieldType, ok := valueType.GetFieldType(expr.Field.Lexeme)
	if !ok {
		return r.error(expr, fmt.Sprintf("field %s not defined", expr.Field.Lexeme))
	}
	return Result{
		Ok:     true,
		Symbol: symbol.NewSymbol(expr.Field.Lexeme, symbol.LiteralSymbol, expr, fieldType),
	}
}

func (r *Resolver) VisitFunctionCall(expr ast.FunctionCallExpr) any {
	sym, ok := r.table.Lookup(expr.Name.Lexeme)
	if !ok {
		return r.error(expr, fmt.Sprintf("function %s not defined", expr.Name.Lexeme))
	}
	return Result{
		Ok:     true,
		Symbol: sym,
	}
}

func (r *Resolver) VisitLogical(expr ast.LogicalExpr) any {
	res := ast.AcceptExpr[Result](expr.Left, r)
	if !res.Ok {
		return res
	}
	return ast.AcceptExpr[Result](expr.Right, r)
}

func (r *Resolver) VisitSlice(expr ast.SliceExpr) any {
	targetRes := ast.AcceptExpr[Result](expr.TargetExpr, r)
	if !targetRes.Ok {
		return targetRes
	}
	res := ast.AcceptExpr[Result](expr.StartIndex, r)
	if !res.Ok {
		return res
	}
	if expr.EndIndex != nil {
		res = ast.AcceptExpr[Result](expr.EndIndex, r)
		if !res.Ok {
			return res
		}
	}
	targetType := targetRes.Symbol.ValueType()
	if targetType.Equals(types.StringType) {
		return Result{
			Ok:     true,
			Symbol: symbol.NewSymbol("__literal__", symbol.LiteralSymbol, expr, types.StringType),
		}
	} else if utils.IsListType(targetType) {
		return Result{
			Ok:     true,
			Symbol: symbol.NewSymbol("__literal__", symbol.LiteralSymbol, expr, targetType.(types.ListTypeStruct).ContentType),
		}
	}
	return r.error(expr, fmt.Sprintf("slice not supported for type %s", targetType.ToString()))
}

func (r *Resolver) VisitList(expr ast.ListExpr) any {
	for _, item := range expr.Elements {
		res := ast.AcceptExpr[Result](item, r)
		if !res.Ok {
			return res
		}
	}
	return Result{
		Ok:     true,
		Symbol: symbol.NewSymbol("__literal__", symbol.LiteralSymbol, expr, expr.ValueType),
	}
}
