package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"strconv"
)

func (p *Parser) expression() (ast.Expr, error) {
	return p.or()
}

func (p *Parser) or() (ast.Expr, error) {
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
		expr = &ast.LogicalExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
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
		expr = &ast.LogicalExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
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
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
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
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
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
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
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
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(tokens.Bang, tokens.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Operator: operator, Expression: right}, nil
	}

	return p.value()
}

func (p *Parser) value() (ast.Expr, error) {
	expr, err := p.baseValue()
	if err != nil {
		return nil, err
	}
	return p.postfix(expr)
}

func (p *Parser) baseValue() (ast.Expr, error) {

	if p.match(tokens.Identifier) {
		identifier := p.previous()
		return &ast.VariableExpr{
			Name: identifier,
		}, nil
	}

	// Fallback to literals, grouping, list literals, etc.
	return p.primary()
}

// postfix wraps an Expr in as many [index] or .field as you find:
func (p *Parser) postfix(expr ast.Expr) (ast.Expr, error) {
	if p.match(tokens.BracketOpen) {
		return p.bracketPostfix(expr)
	}
	if p.match(tokens.Dot) {
		return p.fieldPostfix(expr)
	}
	if p.match(tokens.ParenOpen) {
		return p.callPostfix(expr)
	}
	return expr, nil
}

func (p *Parser) fieldPostfix(expr ast.Expr) (ast.Expr, error) {
	fieldTok, err := p.consume(tokens.Identifier, "Expected field name after '.'")
	if err != nil {
		return nil, err
	}
	expr = &ast.DotAccessExpr{
		Source:         expr,
		Name:           fieldTok,
		SourceLocation: p.location(),
	}
	return p.postfix(expr)
}

func (p *Parser) bracketPostfix(expr ast.Expr) (ast.Expr, error) {
	start, err := p.expression()
	if err != nil {
		return nil, err
	}
	var end ast.Expr
	if p.match(tokens.Colon) {
		end, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(tokens.BracketClose, "Expected ']'"); err != nil {
		return nil, err
	}
	expr = &ast.SliceExpr{
		TargetExpr:     expr,
		StartIndex:     start,
		EndIndex:       end,
		SourceLocation: p.location(),
	}
	return p.postfix(expr)
}

func (p *Parser) callPostfix(expr ast.Expr) (ast.Expr, error) {
	location := p.location()
	args := make([]ast.Expr, 0)
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
			if p.match(tokens.Comma) {
				if p.check(tokens.ParenClose) {
					break
				}
			} else {
				break
			}
		}
	}
	_, err := p.consume(tokens.ParenClose, "Expected ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.CallExpr{
		Source:         expr,
		Arguments:      args,
		SourceLocation: location,
	}, nil
}

func (p *Parser) slice(expr ast.Expr) (ast.Expr, error) {
	start, err := p.expression()
	if err != nil {
		return nil, err
	}
	var end ast.Expr
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
	return &ast.SliceExpr{StartIndex: start, EndIndex: end, TargetExpr: expr, SourceLocation: p.location()}, nil
}

func (p *Parser) functionDeclarationExpression() (ast.Expr, error) {
	/*
		(arg type, arg type) type {

		}
	*/
	_, err := p.consume(tokens.ParenOpen, "Expected '(' after 'func'.")
	if err != nil {
		return nil, err
	}
	parameters := make([]ast.VariableDeclarationStmt, 0)
	parameterTypes := make([]types.ValueType, 0)
	if !p.check(tokens.ParenClose) {
		for {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Cannot have more than 255 parameters.")
			}
			argName, err := p.consume(tokens.Identifier, "Expected parameter name.")
			if err != nil {
				return nil, err
			}
			argType, err := p.ParseType()
			if err != nil {
				return nil, err
			}
			parameterTypes = append(parameterTypes, argType)
			parameters = append(parameters, ast.VariableDeclarationStmt{
				Name:      argName,
				ValueType: argType,
			})
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	_, err = p.consume(tokens.ParenClose, "Expected ')' after parameters.")
	var returnType types.ValueType = types.VoidType
	if !p.check(tokens.BraceOpen) {
		returnType, err = p.ParseType()
	}
	if returnType == nil {
		return nil, p.error(p.peek(), "Expected return type.")
	}
	if err != nil && !p.check(tokens.BraceOpen) {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &ast.FunctionDeclarationExpr{
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(tokens.False) {
		return ast.NewLiteralExpr(nbt.NewInt(0), types.IntType, p.location()), nil
	}
	if p.match(tokens.True) {
		return ast.NewLiteralExpr(nbt.NewInt(1), types.IntType, p.location()), nil
	}
	if p.match(tokens.Int) {
		i, err := strconv.ParseInt(p.previous().Literal, 10, 64)
		if err != nil {
			return nil, p.error(p.previous(), "Invalid integer literal.")
		}
		return ast.NewLiteralExpr(nbt.NewInt(i), types.IntType, p.location()), nil
	}
	if p.match(tokens.Double) {
		d, err := strconv.ParseFloat(p.previous().Literal, 64)
		if err != nil {
			return nil, p.error(p.previous(), "Invalid integer literal.")
		}
		return ast.NewLiteralExpr(nbt.NewDouble(d), types.DoubleType, p.location()), nil
	}
	if p.match(tokens.String) {
		if p.match(tokens.BracketOpen) {
			return p.slice(ast.NewLiteralExpr(nbt.NewString(p.peekCount(-2).Literal), types.StringType, p.location()))
		} else {
			return ast.NewLiteralExpr(nbt.NewString(p.previous().Literal), types.StringType, p.location()), nil
		}
	}
	if p.match(tokens.Func) {
		return p.functionDeclarationExpression()
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
		return &ast.GroupingExpr{Expression: expr}, nil
	}
	if p.match(tokens.BracketOpen) {
		var elems []ast.Expr
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
		return &ast.ListExpr{
			Elements:       elems,
			SourceLocation: p.location(),
			ValueType:      types.NewListType(contentType),
		}, nil
	}
	errorToken := p.peek()
	p.synchronize()
	return nil, p.error(errorToken, "Expected expression.")
}
