package me.kolterdyx.compiler

import me.kolterdyx.compiler.exception.ParseException
import me.kolterdyx.compiler.expression.*
import me.kolterdyx.compiler.statement.BlockStatement
import me.kolterdyx.compiler.statement.ExpressionStatement
import me.kolterdyx.compiler.statement.Statement
import me.kolterdyx.compiler.statement.VariableDeclarationStatement

class Parser(
    private val tokens: List<Token>,
    private var current: Int = 0
) {

    fun parse(): List<Statement> {
        val statements = mutableListOf<Statement>()
        while (!isAtEnd()) {
            statements.add(declaration())
        }
        return statements
    }

    private fun statement(): Statement {
        if (match(TokenType.P_LEFT_BRACE)) return block()
        return expressionStatement()
    }

    private fun block(): Statement {
        val statements = mutableListOf<Statement>()
        while (!check(TokenType.P_RIGHT_BRACE) && !isAtEnd()) {
            statements.add(declaration())
        }
        if (!match(TokenType.P_RIGHT_BRACE)) error("Expected '}' after block")
        return BlockStatement(statements)
    }

    private fun declaration(): Statement {
        try {
            if (match(TokenType.KW_VALUETYPE))
                return varDeclaration()
            return statement()
        } catch (e: ParseException) {
            synchronize()
            return Statement.Empty()
        }
    }

    private fun varDeclaration(): Statement {
        val type = previous()
        val name = consume(TokenType.LIT_IDENTIFIER, "Expected variable name")
        val initializer = if (match(TokenType.OP_EQUAL)) expression() else null
        consume(TokenType.P_SEMICOLON, "Expected variable name")
        return VariableDeclarationStatement(name, type, initializer)
    }

    private fun expressionStatement(): Statement {
        val expr = expression()
        if (!match(TokenType.P_SEMICOLON)) error("Expected ';'")
        return ExpressionStatement(expr)
    }

    /* Expressions */

    fun expression(): Expression {
        return assignment()
    }

    private fun assignment(): Expression {
        val expr = equality()
        if (match(TokenType.OP_EQUAL)) {
            val value = assignment()
            if (expr is VariableExpression) {
                val name = expr.name
                return AssignmentExpression(name, value)
            }
            error("Invalid assignment target")
        }
        return expr
    }

    private fun equality(): Expression {
        var expr = comparison()
        while (match(TokenType.OP_BANG_EQUAL, TokenType.OP_EQUAL_EQUAL)) {
            val operator = previous()
            val right = comparison()
            expr = BinaryExpression(expr, operator, right)
        }
        return expr
    }

    private fun comparison(): Expression {
        var expr = addition()
        while (match(TokenType.OP_GREATER, TokenType.OP_GREATER_EQUAL, TokenType.OP_LESS, TokenType.OP_LESS_EQUAL)) {
            val operator = previous()
            val right = addition()
            expr = BinaryExpression(expr, operator, right)
        }
        return expr
    }

    private fun addition(): Expression {
        var expr = multiplication()
        while (match(TokenType.OP_MINUS, TokenType.OP_PLUS)) {
            val operator = previous()
            val right = multiplication()
            expr = BinaryExpression(expr, operator, right)
        }
        return expr
    }

    private fun multiplication(): Expression {
        var expr = unary()
        while (match(TokenType.OP_SLASH, TokenType.OP_STAR, TokenType.OP_PERCENT)) {
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
        if (match(TokenType.LIT_INT, TokenType.LIT_FLOAT, TokenType.LIT_STRING, TokenType.LIT_BOOLEAN)) return LiteralExpression(previous())
        if (match(TokenType.P_LEFT_PAREN)) {
            val expr = expression()
            if (!match(TokenType.P_RIGHT_PAREN)) error("Expected ')' after expression")
            return GroupingExpression(expr)
        }
        error("Expected expression")
    }

    /* Utilities */

    private fun synchronize() {
        advance()
        while (!isAtEnd()) {
            if (previous().type == TokenType.P_SEMICOLON) return
            when (peek().type) {
                TokenType.KW_VALUETYPE -> return
                TokenType.P_LEFT_BRACE -> return
                TokenType.P_EOF -> return
                else -> advance()
            }
        }
    }

    private fun consume(identifier: TokenType, errorMsg: String): Token {
        if (match(identifier)) {
            return previous()
        }
        error(errorMsg)
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
        return peek().type == TokenType.P_EOF
    }
}