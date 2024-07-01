package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ast.ExpressionVisitor

class GroupingExpression(
    val expression: Expression
) : Expression(expression.valueType) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitGrouping(this)
    }
}