package me.kolterdyx.compiler.statement

import me.kolterdyx.compiler.ast.StatementVisitor

abstract class Statement {

    class Empty : Statement() {
        override fun <R> accept(visitor: StatementVisitor<R>): R {
            throw UnsupportedOperationException("Empty statement should not be visited")
        }
    }

    abstract fun <R> accept(visitor: StatementVisitor<R>): R
}