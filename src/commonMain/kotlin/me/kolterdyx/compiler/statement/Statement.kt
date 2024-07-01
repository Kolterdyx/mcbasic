package me.kolterdyx.compiler.statement

import me.kolterdyx.compiler.ast.StatementVisitor

abstract class Statement {
    abstract fun <R> accept(visitor: StatementVisitor<R>): R
}