package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor

class LiteralExpression(
    val value: Any?, valueType: ValueType
) : Expression(valueType) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitLiteral(this)
    }
}