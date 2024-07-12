package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ast.ExpressionVisitor

class AssignmentExpression(
    val name: Token,
    private val value: Expression
) : Expression(value.valueType) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitAssignment(this)
    }

}