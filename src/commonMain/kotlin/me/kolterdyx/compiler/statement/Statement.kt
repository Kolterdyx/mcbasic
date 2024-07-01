package me.kolterdyx.compiler.statement

abstract class Statement {
    abstract fun <R> accept(visitor: StatementVisitor<R>): R
}