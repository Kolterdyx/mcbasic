package me.kolterdyx.compiler.statement

import me.kolterdyx.compiler.ast.StatementVisitor

class BlockStatement(
    val statements: List<Statement>
) : Statement() {
    override fun <R> accept(visitor: StatementVisitor<R>): R {
        return visitor.visitBlock(this)
    }
}