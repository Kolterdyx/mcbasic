package type_checker

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
)

func (t *TypeChecker) VisitBinary(expr ast.BinaryExpr) any {
	rtype := ast.AcceptExpr[types.ValueType](expr.Right, t)
	ltype := ast.AcceptExpr[types.ValueType](expr.Left, t)
	switch expr.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		if ltype != rtype {
			t.error(expr, fmt.Sprintf("cannot compare %s and %s", ltype.ToString(), rtype.ToString()))
		}
		return rtype
	case tokens.Plus:
		switch ltype {
		case types.StringType:
			return types.StringType
		}
		fallthrough
	case tokens.Minus, tokens.Slash, tokens.Star, tokens.Percent:
		switch rtype {
		case types.IntType:
			if rtype != types.IntType {
				t.error(expr, fmt.Sprintf("cannot add %s and %s", ltype.ToString(), rtype.ToString()))
			}
			return types.IntType
		case types.DoubleType:
			if ltype != types.DoubleType {
				t.error(expr, fmt.Sprintf("cannot add %s and %s", ltype.ToString(), rtype.ToString()))
			}
			return types.DoubleType
		}
	default:
		t.error(expr, fmt.Sprintf("unhandled operator %s", expr.Operator.Lexeme))
	}
	return types.VoidType
}

func (t *TypeChecker) VisitGrouping(expr ast.GroupingExpr) any {
	return ast.AcceptExpr[types.ValueType](expr.Expression, t)
}

func (t *TypeChecker) VisitLiteral(expr ast.LiteralExpr) any {
	return expr.ValueType
}

func (t *TypeChecker) VisitUnary(expr ast.UnaryExpr) any {
	return ast.AcceptExpr[types.ValueType](expr.Expression, t)
}

func (t *TypeChecker) VisitVariable(expr ast.VariableExpr) any {
	sym, _ := t.table.Lookup(expr.Name.Lexeme)
	return sym.ValueType()
}

func (t *TypeChecker) VisitFieldAccess(expr ast.FieldAccessExpr) any {
	sourceType := ast.AcceptExpr[types.ValueType](expr.Source, t)
	vtype, ok := sourceType.GetFieldType(expr.Field.Lexeme)
	if !ok {
		t.error(expr, fmt.Sprintf("field %s not defined in type %s", expr.Field.Lexeme, sourceType.ToString()))
		return sourceType
	}
	return vtype
}

func (t *TypeChecker) VisitFunctionCall(expr ast.FunctionCallExpr) any {
	sym, _ := t.table.Lookup(expr.Name.Lexeme)
	switch declarationNode := sym.DeclarationNode().(type) {
	case ast.FunctionDeclarationStmt:
		if len(expr.Arguments) != len(declarationNode.Parameters) {
			t.error(expr, fmt.Sprintf("function %s expects %d arguments, got %d", expr.Name.Lexeme, len(declarationNode.Parameters), len(expr.Arguments)))
			return sym.ValueType()
		}
		for i, arg := range expr.Arguments {
			ptype := ast.AcceptExpr[types.ValueType](arg, t)
			if ptype != declarationNode.Parameters[i].ValueType {
				t.error(arg, fmt.Sprintf("argument %d of function %s has invalid type: expected %s, got %s", i, expr.Name.Lexeme, declarationNode.Parameters[i].ValueType.ToString(), ptype.ToString()))
			}
		}
	case ast.StructDeclarationStmt:
		if len(expr.Arguments) != len(declarationNode.StructType.GetFieldNames()) {
			t.error(expr, fmt.Sprintf("struct %s expects %d arguments, got %d", expr.Name.Lexeme, len(declarationNode.StructType.GetFieldNames()), len(expr.Arguments)))
			return sym.ValueType()
		}
		for i, arg := range expr.Arguments {

			ptype := ast.AcceptExpr[types.ValueType](arg, t)
			fieldName := declarationNode.StructType.GetFieldNames()[i]
			fieldType, _ := declarationNode.StructType.GetFieldType(fieldName)
			if !ptype.Equals(fieldType) {
				t.error(arg, fmt.Sprintf("argument %d of struct %s has invalid type: expected %s, got %s", i, expr.Name.Lexeme, declarationNode.StructType.GetFieldNames()[i], ptype.ToString()))
			}
		}

	}
	return sym.ValueType()
}

func (t *TypeChecker) VisitLogical(expr ast.LogicalExpr) any {
	ast.AcceptExpr[types.ValueType](expr.Left, t)
	ast.AcceptExpr[types.ValueType](expr.Right, t)
	return types.IntType
}

func (t *TypeChecker) VisitSlice(expr ast.SliceExpr) any {
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
	} else if utils.IsListType(targetType) {
		return targetType.(types.ListTypeStruct).ContentType
	}
	t.error(expr.TargetExpr, fmt.Sprintf("target must be string or list, got %s", targetType.ToString()))
	return types.VoidType
}

func (t *TypeChecker) VisitList(expr ast.ListExpr) any {
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
	return types.NewListType(itype)
}
