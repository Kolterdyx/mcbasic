package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ast.ExprVisitor

class LiteralExpr(
    val value: Any?
) : Expr() {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitLiteralExpr(this)
    }
}