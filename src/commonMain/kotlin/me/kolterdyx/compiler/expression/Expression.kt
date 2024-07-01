package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor

abstract class Expression(
    val valueType: ValueType
) {
    abstract fun <R> accept(visitor: ExpressionVisitor<R>): R
}