package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExprVisitor

class LiteralExpr(
    val value: Any?, valueType: ValueType
) : Expr(valueType) {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitLiteralExpr(this)
    }
}