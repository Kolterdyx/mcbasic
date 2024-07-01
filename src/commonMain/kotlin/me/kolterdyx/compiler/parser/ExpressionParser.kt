package me.kolterdyx.compiler.parser

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ast.ExpressionVisitor
import me.kolterdyx.compiler.expression.Expression

class ExpressionParser : Parser<List<Token>, List<Expression>> {
    override fun parse(data: List<Token>): List<Expression> {
        val expressions = mutableListOf<Expression>()
        val tokens = data.toMutableList()
        while (tokens.isNotEmpty()) {
            expressions.add(parseExpression(tokens))
        }
        return expressions
    }

    private fun parseExpression(tokens: MutableList<Token>): Expression {
        TODO("Not yet implemented")
    }
}