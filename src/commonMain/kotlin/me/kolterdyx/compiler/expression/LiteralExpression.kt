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
}