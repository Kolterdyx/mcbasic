package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor

fun checkCompatibility(left: ValueType, right: ValueType): ValueType {
    val isCompatible = when (left) {
        ValueType.INT -> right == ValueType.INT
        ValueType.STRING -> right == ValueType.STRING || right == ValueType.INT || right == ValueType.FLOAT
        ValueType.FLOAT -> right == ValueType.FLOAT
    }
    if (!isCompatible) {
        throw Exception("Incompatible types: $left and $right")
    }
    return left
}

class BinaryExpression(
    val left: Expression,
    val operator: Token,
    val right: Expression
) : Expression(checkCompatibility(left.valueType, right.valueType)) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitBinary(this)
    }
}