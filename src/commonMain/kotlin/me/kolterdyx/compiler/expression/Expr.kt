package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ast.ExprVisitor

abstract class Expr() {
    abstract fun <R> accept(visitor: ExprVisitor<R>): R
}