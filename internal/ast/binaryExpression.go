package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type BinaryExpr struct {
	Expr

	Left     Expr
	Operator tokens.Token
	Right    Expr
}

func (b BinaryExpr) Accept(v ExpressionVisitor) any {
	return v.VisitBinary(b)
}

func (b BinaryExpr) Type() NodeType {
	return BinaryExpression
}

func (b BinaryExpr) ToString() string {
	return "(" + b.Left.ToString() + " " + b.Operator.Lexeme + " " + b.Right.ToString() + ")"
}

func (b BinaryExpr) GetSourceLocation() interfaces.SourceLocation {
	return b.Left.GetSourceLocation()
}

//
//func (b BinaryExpr) Validate() error {
//	if b.Left == nil {
//		return fmt.Errorf("BinaryExpr expression must have a left hand expression")
//	}
//	if b.Right == nil {
//		return fmt.Errorf("BinaryExpr expression must have a right hand expression")
//	}
//	switch b.Left.ReturnType() {
//	case types.IntType:
//		switch b.Right.ReturnType() {
//		case types.IntType:
//		default:
//			return fmt.Errorf("invalid righ hand time for type 'int' and operator '%s': %s", b.Operator.Lexeme, b.Right.ReturnType().ToString())
//		}
//	case types.StringType:
//	case types.DoubleType:
//		switch b.Right.ReturnType() {
//		case types.DoubleType:
//		default:
//			return fmt.Errorf("invalid righ hand time for type 'double' and operator '%s': %s", b.Operator.Lexeme, b.Right.ReturnType().ToString())
//		}
//	}
//	return nil
//}
