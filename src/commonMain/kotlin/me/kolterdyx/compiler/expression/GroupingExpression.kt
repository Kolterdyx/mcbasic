package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ast.ExpressionVisitor

class GroupingExpression(
    val expression: Expression
) : Expression(expression.valueType) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitGrouping(this)
    }

    override fun toString(): String {
        return "GroupingExpression($expression)"
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is GroupingExpression) return false
        if (!super.equals(other)) return false
        if (expression != other.expression) return false
        return true
    }

    override fun hashCode(): Int {
        var result = super.hashCode()
        result = 31 * result + expression.hashCode()
        return result
    }
}