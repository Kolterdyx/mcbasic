package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ast.ExprVisitor

class UnaryExpression(
    val operator: Token,
    val right: Expression
) : Expression(right.valueType) {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitUnaryExpr(this)
    }
}