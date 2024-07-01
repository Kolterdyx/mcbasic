package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExprVisitor

class LiteralExpression(
    val value: Any?, valueType: ValueType
) : Expression(valueType) {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitLiteralExpr(this)
    }
}