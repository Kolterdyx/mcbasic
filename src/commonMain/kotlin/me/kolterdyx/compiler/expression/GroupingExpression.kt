package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ast.ExprVisitor

class GroupingExpression(
    val expression: Expression
) : Expression(expression.valueType) {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitGroupingExpr(this)
    }
}