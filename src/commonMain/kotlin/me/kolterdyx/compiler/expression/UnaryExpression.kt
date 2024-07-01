package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ast.ExpressionVisitor

class UnaryExpression(
    val operator: Token,
    val right: Expression
) : Expression(right.valueType) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitUnary(this)
    }
}