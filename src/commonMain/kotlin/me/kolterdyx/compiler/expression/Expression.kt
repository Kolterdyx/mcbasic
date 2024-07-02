package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor

abstract class Expression(
    val valueType: ValueType
) {
    abstract fun <R> accept(visitor: ExpressionVisitor<R>): R

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is Expression) return false
        return this.valueType == other.valueType
    }

    override fun hashCode(): Int {
        return valueType.hashCode()
    }
}