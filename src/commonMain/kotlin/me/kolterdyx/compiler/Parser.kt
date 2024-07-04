package me.kolterdyx.compiler

import me.kolterdyx.compiler.exception.ParseException
import me.kolterdyx.compiler.expression.*
import me.kolterdyx.compiler.statement.Statement

class Parser(
    private val tokens: List<Token>,
    private var current: Int = 0
) {

    fun parse(): Statement {
        val expressions = mutableListOf<Expression>()
        while (!isAtEnd()) {
            expressions.add(expression())
        }
        return Statement.Empty()
    }

    fun expression(): Expression {
        return equality()
    }

    private fun equality(): Expression {
        var expr = comparison()
        while (match(TokenType.BANG_EQUAL, TokenType.EQUAL_EQUAL)) {
            val operator = previous()
            val right = comparison()
            expr = BinaryExpression(expr, operator, right)
        }
        return expr
    }

    private fun comparison(): Expression {
        var expr = addition()
        while (match(TokenType.GREATER, TokenType.GREATER_EQUAL, TokenType.LESS, TokenType.LESS_EQUAL)) {
            val operator = previous()
            val right = addition()
            expr = BinaryExpression(expr, operator, right)
        }
        return expr
    }

    private fun addition(): Expression {
        var expr = multiplication()
        while (match(TokenType.MINUS, TokenType.PLUS)) {
            val operator = previous()
            val right = multiplication()
            expr = BinaryExpression(expr, operator, right)
        }
        return expr
    }

    private fun multiplication(): Expression {
        var expr = unary()
        while (match(TokenType.SLASH, TokenType.STAR, TokenType.PERCENT)) {
            val operator = previous()
            val right = unary()
            expr = BinaryExpression(expr, operator, right)
        }
        return expr
    }

    private fun unary(): Expression {
        if (match(*UnaryExpression.ValidOperators)) {
            val operator = previous()
            val right = unary()
            return UnaryExpression(operator, right)
        }
        return primary()
    }

    private fun primary(): Expression {
        if (match(TokenType.INT, TokenType.FLOAT, TokenType.STRING, TokenType.BOOLEAN)) return LiteralExpression(previous())
        if (match(TokenType.LEFT_PAREN)) {
            val expr = expression()
            if (!match(TokenType.RIGHT_PAREN)) error("Expected ')' after expression")
            return GroupingExpression(expr)
        }
        error("Expected expression")
    }

    private fun error(message: String): Nothing {
        throw ParseException(tokens[current], message)
    }

    private fun previous(): Token {
        if (current == 0) throw IllegalStateException("No previous token")
        return tokens[current - 1]
    }

    private fun match(vararg tokenTypes: TokenType): Boolean {
        for (type in tokenTypes) {
            if (check(type)) {
                advance()
                return true
            }
        }
        return false
    }

    private fun advance() {
        if (!isAtEnd()) current++
    }

    private fun check(type: TokenType): Boolean {
        if (isAtEnd()) return false
        return peek().type == type
    }

    private fun peek(): Token {
        return tokens[current]
    }

    private fun isAtEnd(): Boolean {
        return peek().type == TokenType.EOF
    }
}