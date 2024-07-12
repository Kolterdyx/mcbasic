package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor

class VariableExpression(
    val name: Token,
    valueType: ValueType
) : Expression(valueType) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitVariable(this)
    }

}