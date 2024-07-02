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

    override fun toString(): String {
        return "UnaryExpression($operator, $right)"
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is UnaryExpression) return false
        if (!super.equals(other)) return false
        if (operator != other.operator) return false
        if (right != other.right) return false
        return true
    }

    override fun hashCode(): Int {
        var result = super.hashCode()
        result = 31 * result + operator.hashCode()
        result = 31 * result + right.hashCode()
        return result
    }
}