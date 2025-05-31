package type_checker

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
)

func (t *TypeChecker) VisitBinary(expr *ast.BinaryExpr) any {
	rtype := ast.AcceptExpr[types.ValueType](expr.Right, t)
	ltype := ast.AcceptExpr[types.ValueType](expr.Left, t)
	expr.SetResolvedType(types.VoidType)
	switch expr.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		if ltype != rtype {
			t.error(expr, fmt.Sprintf("cannot compare %s and %s", ltype.ToString(), rtype.ToString()))
		}
		return rtype
	case tokens.Plus:
		switch ltype {
		case types.StringType:
			expr.SetResolvedType(types.StringType)
		}
		fallthrough
	case tokens.Minus, tokens.Slash, tokens.Star, tokens.Percent:
		switch rtype {
		case types.IntType:
			if rtype != types.IntType {
				t.error(expr, fmt.Sprintf("cannot add %s and %s", ltype.ToString(), rtype.ToString()))
			}
			expr.SetResolvedType(types.IntType)
		case types.DoubleType:
			if ltype != types.DoubleType {
				t.error(expr, fmt.Sprintf("cannot add %s and %s", ltype.ToString(), rtype.ToString()))
			}
			expr.SetResolvedType(types.DoubleType)
		}
	default:
		t.error(expr, fmt.Sprintf("unhandled operator %s", expr.Operator.Lexeme))
	}
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitGrouping(expr *ast.GroupingExpr) any {
	expr.SetResolvedType(ast.AcceptExpr[types.ValueType](expr.Expression, t))
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitLiteral(expr *ast.LiteralExpr) any {
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitUnary(expr *ast.UnaryExpr) any {
	expr.SetResolvedType(ast.AcceptExpr[types.ValueType](expr.Expression, t))
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitVariable(expr *ast.VariableExpr) any {
	sym, _ := t.table.Lookup(expr.Name.Lexeme)
	expr.SetResolvedType(sym.ValueType())
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitDotAccess(expr *ast.DotAccessExpr) any {
	sourceType := ast.AcceptExpr[types.ValueType](expr.Source, t)
	expr.SetResolvedType(sourceType)
	vtype, ok := sourceType.GetFieldType(expr.Name.Lexeme)
	if !ok {
		t.error(expr, fmt.Sprintf("field %s not defined in type %s", expr.Name.Lexeme, sourceType.ToString()))
	} else {
		expr.SetResolvedType(vtype)
	}
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitCall(expr *ast.CallExpr) any {
	sym, _ := t.table.Lookup(expr.GetResolvedName())
	expr.SetResolvedType(sym.ValueType())
	switch declarationNode := sym.DeclarationNode().(type) {
	case ast.FunctionDeclarationStmt:
		if len(expr.Arguments) != len(declarationNode.Parameters) {
			t.error(expr, fmt.Sprintf("function %s expects %d arguments, got %d", expr.GetResolvedName(), len(declarationNode.Parameters), len(expr.Arguments)))
			break
		}
		for i, arg := range expr.Arguments {
			ptype := ast.AcceptExpr[types.ValueType](arg, t)
			if ptype != declarationNode.Parameters[i].ValueType && declarationNode.Parameters[i].ValueType != types.VoidType {
				t.error(arg, fmt.Sprintf("argument %d of function %s has invalid type: expected %s, got %s", i, expr.GetResolvedName(), declarationNode.Parameters[i].ValueType.ToString(), ptype.ToString()))
			}
		}
	case ast.StructDeclarationStmt:
		if len(expr.Arguments) != len(declarationNode.StructType.GetFieldNames()) {
			t.error(expr, fmt.Sprintf("struct %s expects %d arguments, got %d", expr.GetResolvedName(), len(declarationNode.StructType.GetFieldNames()), len(expr.Arguments)))
			break
		}
		for i, arg := range expr.Arguments {

			ptype := ast.AcceptExpr[types.ValueType](arg, t)
			fieldName := declarationNode.StructType.GetFieldNames()[i]
			fieldType, _ := declarationNode.StructType.GetFieldType(fieldName)
			if !ptype.Equals(fieldType) {
				t.error(arg, fmt.Sprintf("argument %d of struct %s has invalid type: expected %s, got %s", i, expr.GetResolvedName(), declarationNode.StructType.GetFieldNames()[i], ptype.ToString()))
			}
		}

	}
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitLogical(expr *ast.LogicalExpr) any {
	ast.AcceptExpr[types.ValueType](expr.Left, t)
	ast.AcceptExpr[types.ValueType](expr.Right, t)
	expr.SetResolvedType(types.IntType)
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitSlice(expr *ast.SliceExpr) any {
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
		expr.SetResolvedType(types.StringType)
	} else if utils.IsListType(targetType) {
		expr.SetResolvedType(targetType.(types.ListTypeStruct).ContentType)
	} else {
		t.error(expr.TargetExpr, fmt.Sprintf("target must be string or list, got %s", targetType.ToString()))
		expr.SetResolvedType(types.VoidType)
	}
	return expr.GetResolvedType()
}

func (t *TypeChecker) VisitList(expr *ast.ListExpr) any {
	var itype types.ValueType = types.VoidType
	for _, item := range expr.Elements {
		currType := ast.AcceptExpr[types.ValueType](item, t)
		if itype == types.VoidType {
			itype = currType
		}
		if itype != currType {
			t.error(item, fmt.Sprintf("list elements must be of the same type, got %s and %s", itype.ToString(), currType.ToString()))
		}
	}
	expr.SetResolvedType(types.NewListType(itype))
	return expr.GetResolvedType()
}
