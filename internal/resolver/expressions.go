package resolver

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
)

func (r *Resolver) VisitBinary(expr *ast.BinaryExpr) any {
	res := ast.AcceptExpr[Result](expr.Left, r)
	if !res.Ok {
		return res
	}
	return ast.AcceptExpr[Result](expr.Right, r)
}

func (r *Resolver) VisitGrouping(expr *ast.GroupingExpr) any {
	return ast.AcceptExpr[Result](expr.Expression, r)
}

func (r *Resolver) VisitLiteral(expr *ast.LiteralExpr) any {
	return Result{
		Ok:     true,
		Symbol: symbol.NewSymbol("__literal__", symbol.LiteralSymbol, expr, expr.GetResolvedType()),
	}
}

func (r *Resolver) VisitUnary(expr *ast.UnaryExpr) any {
	return ast.AcceptExpr[Result](expr.Expression, r)
}

func (r *Resolver) VisitVariable(expr *ast.VariableExpr) any {
	sym, ok := r.table.Lookup(expr.Name.Lexeme)
	if !ok {
		return r.error(expr, fmt.Sprintf("%s is not defined", expr.Name.Lexeme))
	}
	return Result{
		Ok:     true,
		Symbol: sym,
	}
}

func (r *Resolver) VisitDotAccess(expr *ast.DotAccessExpr) any {
	res := ast.AcceptExpr[Result](expr.Source, r)
	if !res.Ok {
		return r.error(expr, fmt.Sprintf("could not resolve '%s'", expr.Source.ToString()))
	}

	valueType := res.Symbol.ValueType()
	fieldType, ok := valueType.GetFieldType(expr.Name.Lexeme)
	if !ok {
		return r.error(expr, fmt.Sprintf("'%s' not defined in '%s'", expr.Name.Lexeme, res.Symbol.Name()))
	}
	return Result{
		Ok:     true,
		Symbol: symbol.NewSymbol(expr.Name.Lexeme, symbol.LiteralSymbol, expr, fieldType),
	}

}

func (r *Resolver) VisitCall(expr *ast.CallExpr) any {
	for _, arg := range expr.Arguments {
		res := ast.AcceptExpr[Result](arg, r)
		if !res.Ok {
			return res
		}
	}
	res := ast.AcceptExpr[Result](expr.Source, r)
	if !res.Ok {
		return r.error(expr, fmt.Sprintf("could not resolve '%s'", expr.Source.ToString()))
	}
	if res.Symbol.Type() != symbol.FunctionSymbol && res.Symbol.Type() != symbol.StructSymbol {
		return r.error(expr, fmt.Sprintf("'%s' is neither a function nor a struct.", expr.Source.ToString()))
	}
	symName := res.Symbol.Name()
	sym, ok := r.table.Lookup(res.Symbol.Alias())
	if !ok {
		return r.error(expr, fmt.Sprintf("%s '%s' not defined", sym.Type(), symName))
	}
	expr.SetResolvedName(res.Symbol.Alias())
	return Result{
		Ok:     true,
		Symbol: sym,
	}
}

func (r *Resolver) VisitLogical(expr *ast.LogicalExpr) any {
	res := ast.AcceptExpr[Result](expr.Left, r)
	if !res.Ok {
		return res
	}
	return ast.AcceptExpr[Result](expr.Right, r)
}

func (r *Resolver) VisitSlice(expr *ast.SliceExpr) any {
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

func (r *Resolver) VisitList(expr *ast.ListExpr) any {
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
