package me.kolterdyx.compiler.statement

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ast.StatementVisitor
import me.kolterdyx.compiler.expression.Expression

class VariableDeclarationStatement(
    val name: Token,
    val type: Token,
    val initializer: Expression?
) : Statement() {
    override fun <R> accept(visitor: StatementVisitor<R>): R {
        return visitor.visitVariableDeclaration(this)
    }
}