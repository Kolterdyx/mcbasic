package expressions

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
)

type BinaryExpr struct {
	Expr
	interfaces.SourceLocation

	Left     Expr
	Operator tokens.Token
	Right    Expr
}

func (b BinaryExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitBinary(b)
}

func (b BinaryExpr) Type() ast.NodeType {
	return ast.BinaryExpression
}

func (b BinaryExpr) ReturnType() types.ValueType {
	switch b.Left.ReturnType() {
	case types.IntType:
		switch b.Right.ReturnType() {
		case types.IntType:
			return types.IntType
		default:
			log.Warnf("Invalid righ hand time for type 'int' and operator '%s': %s", b.Operator.Lexeme, b.Right.ReturnType())
			return nil
		}
	case types.StringType:
		return types.StringType
	case types.DoubleType:
		switch b.Right.ReturnType() {
		case types.DoubleType:
			return types.DoubleType
		default:
			log.Warnf("Invalid righ hand time for type 'double' and operator '%s': %s", b.Operator.Lexeme, b.Right.ReturnType())
			return nil
		}
	}
	return b.Left.ReturnType()
}

func (b BinaryExpr) ToString() string {
	return "(" + b.Left.ToString() + " " + b.Operator.Lexeme + " " + b.Right.ToString() + ")"
}

func (b BinaryExpr) Validate() error {
	if b.Left == nil {
		return fmt.Errorf("BinaryExpr expression must have a left hand expression")
	}
	if b.Right == nil {
		return fmt.Errorf("BinaryExpr expression must have a right hand expression")
	}
	switch b.Left.ReturnType() {
	case types.IntType:
		switch b.Right.ReturnType() {
		case types.IntType:
		default:
			return fmt.Errorf("invalid righ hand time for type 'int' and operator '%s': %s", b.Operator.Lexeme, b.Right.ReturnType().ToString())
		}
	case types.StringType:
	case types.DoubleType:
		switch b.Right.ReturnType() {
		case types.DoubleType:
		default:
			return fmt.Errorf("invalid righ hand time for type 'double' and operator '%s': %s", b.Operator.Lexeme, b.Right.ReturnType().ToString())
		}
	}
	return nil
}
