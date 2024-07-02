package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor

fun getLiteralValueType(value: String): ValueType {
    return when {
        value == "true" || value == "false" -> ValueType.BOOLEAN
        value.toIntOrNull() != null -> ValueType.INT
        value.toDoubleOrNull() != null -> ValueType.FLOAT
        else -> ValueType.STRING
    }
}

class LiteralExpression(
    val value: Any?
) : Expression(getLiteralValueType(value.toString())) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitLiteral(this)
    }

    override fun toString(): String {
        return "LiteralExpression($value)"
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is LiteralExpression) return false
        if (!super.equals(other)) return false
        if (value != other.value) return false
        return true
    }

    override fun hashCode(): Int {
        var result = super.hashCode()
        result = 31 * result + (value?.hashCode() ?: 0)
        return result
    }
}