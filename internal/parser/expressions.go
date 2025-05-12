package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"strconv"
)

func (p *Parser) expression() (expressions.Expr, error) {
	return p.or()
}

func (p *Parser) or() (expressions.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = expressions.LogicalExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) and() (expressions.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = expressions.LogicalExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) equality() (expressions.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.BangEqual, tokens.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) comparison() (expressions.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) term() (expressions.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Minus, tokens.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) factor() (expressions.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Slash, tokens.Star, tokens.Percent) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) unary() (expressions.Expr, error) {
	if p.match(tokens.Bang, tokens.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expressions.UnaryExpr{Operator: operator, Expression: right, SourceLocation: p.location()}, nil
	}

	return p.value()
}

func (p *Parser) value() (expressions.Expr, error) {
	expr, err := p.baseValue()
	if err != nil {
		return nil, err
	}
	return p.postfix(expr)
}

func (p *Parser) baseValue() (expressions.Expr, error) {

	if p.match(tokens.Identifier) {
		identifier := p.previous()

		// If this is a function call, delegate to functionCall (handles namespace)
		if p.match(tokens.ParenOpen) {
			//if sym, ok := p.symbols.Lookup(identifier.Lexeme); ok && sym.Type() == symbol.StructSymbol {
			//	return p.structLiteral(identifier, sym.ValueType().(types.StructTypeStruct))
			//} else {
			return p.functionCall(identifier)
		}

		return expressions.VariableExpr{
			Name:           identifier,
			SourceLocation: p.location(),
		}, nil
	}

	// Fallback to literals, grouping, list literals, etc.
	return p.primary()
}

// postfix wraps an Expr in as many [index] or .field as you find:
func (p *Parser) postfix(expr expressions.Expr) (expressions.Expr, error) {
	if p.match(tokens.BracketOpen) {
		return p.bracketPostfix(expr)
	}
	if p.match(tokens.Dot) {
		return p.fieldPostfix(expr)
	}
	return expr, nil
}

func (p *Parser) fieldPostfix(expr expressions.Expr) (expressions.Expr, error) {
	fieldTok, err := p.consume(tokens.Identifier, "Expected field name after '.'")
	if err != nil {
		return nil, err
	}
	expr = expressions.FieldAccessExpr{
		Source:         expr,
		Field:          fieldTok,
		SourceLocation: p.location(),
	}
	return p.postfix(expr)
}

func (p *Parser) bracketPostfix(expr expressions.Expr) (expressions.Expr, error) {
	start, err := p.expression()
	if err != nil {
		return nil, err
	}
	var end expressions.Expr
	if p.match(tokens.Colon) {
		end, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(tokens.BracketClose, "Expected ']'"); err != nil {
		return nil, err
	}
	expr = expressions.SliceExpr{
		TargetExpr:     expr,
		StartIndex:     start,
		EndIndex:       end,
		SourceLocation: p.location(),
	}
	return p.postfix(expr)
}

func (p *Parser) functionCall(name tokens.Token) (expressions.Expr, error) {
	location := p.location()
	args := make([]expressions.Expr, 0)
	if !p.check(tokens.ParenClose) {
		for {
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, exp)
			if len(args) >= 255 {
				return nil, p.error(p.peek(), "Cannot have more than 255 arguments.")
			}
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	_, err := p.consume(tokens.ParenClose, "Expected ')' after arguments.")
	if err != nil {
		return nil, err
	}
	// Find the function in the current scope

	return expressions.FunctionCallExpr{Name: tokens.Token{
		Type:           tokens.Identifier,
		Lexeme:         name.Lexeme,
		Literal:        name.Literal,
		SourceLocation: name.SourceLocation,
	}, Arguments: args, SourceLocation: location}, nil
}

func (p *Parser) slice(expr expressions.Expr) (expressions.Expr, error) {
	start, err := p.expression()
	if err != nil {
		return nil, err
	}
	var end expressions.Expr
	if p.match(tokens.Colon) {
		end, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.BracketClose, "Expected ']' at the end of the slice.")
	if err != nil {
		return nil, err
	}
	return expressions.SliceExpr{StartIndex: start, EndIndex: end, TargetExpr: expr, SourceLocation: p.location()}, nil
}

func (p *Parser) primary() (expressions.Expr, error) {
	if p.match(tokens.False) {
		return expressions.LiteralExpr{Value: nbt.NewInt(0), ValueType: types.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.True) {
		return expressions.LiteralExpr{Value: nbt.NewInt(1), ValueType: types.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.Int) {
		i, err := strconv.ParseInt(p.previous().Literal, 10, 64)
		if err != nil {
			return nil, p.error(p.previous(), "Invalid integer literal.")
		}
		return expressions.LiteralExpr{Value: nbt.NewInt(i), ValueType: types.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.Double) {
		d, err := strconv.ParseFloat(p.previous().Literal, 64)
		if err != nil {
			return nil, p.error(p.previous(), "Invalid integer literal.")
		}
		return expressions.LiteralExpr{Value: nbt.NewDouble(d), ValueType: types.DoubleType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.String) {
		if p.match(tokens.BracketOpen) {
			return p.slice(expressions.LiteralExpr{Value: nbt.NewString(p.peekCount(-2).Literal), ValueType: types.StringType, SourceLocation: p.location()})
		} else {
			return expressions.LiteralExpr{Value: nbt.NewString(p.previous().Literal), ValueType: types.StringType, SourceLocation: p.location()}, nil
		}
	}
	if p.match(tokens.ParenOpen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokens.ParenClose, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}
		return expressions.GroupingExpr{Expression: expr, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.BracketOpen) {
		var elems []expressions.Expr
		var contentType types.ValueType = types.VoidType
		for !p.check(tokens.BracketClose) {
			if p.match(tokens.Comma) {
				continue
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			elems = append(elems, expr)
			if p.check(tokens.BracketClose) {
				break
			}
			_, err = p.consume(tokens.Comma, "Expected ',' after expression.")
			if err != nil {
				return nil, err
			}
		}
		_, err := p.consume(tokens.BracketClose, "Expected ']' after expression.")
		if err != nil {
			return nil, err
		}
		return expressions.ListExpr{
			Elements:       elems,
			SourceLocation: p.location(),
			ValueType:      types.NewListType(contentType),
		}, nil
	}
	errorToken := p.peek()
	p.synchronize()
	return nil, p.error(errorToken, "Expected expression.")
}
