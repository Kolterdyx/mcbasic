package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExprVisitor

abstract class Expr(
    val valueType: ValueType
) {
    abstract fun <R> accept(visitor: ExprVisitor<R>): R
}