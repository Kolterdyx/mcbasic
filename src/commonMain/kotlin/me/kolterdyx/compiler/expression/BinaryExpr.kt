package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExprVisitor

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

class BinaryExpr(
    val left: Expr,
    val operator: Token,
    val right: Expr
) : Expr(checkCompatibility(left.valueType, right.valueType)) {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitBinaryExpr(this)
    }
}