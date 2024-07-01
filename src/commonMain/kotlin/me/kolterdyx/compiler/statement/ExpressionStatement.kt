package me.kolterdyx.compiler.statement

import me.kolterdyx.compiler.ast.StatementVisitor
import me.kolterdyx.compiler.expression.Expression

class ExpressionStatement(
    val expression: Expression
) : Statement() {
    override fun <R> accept(visitor: StatementVisitor<R>): R {
        return visitor.visitExpression(this)
    }
}