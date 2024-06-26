package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ast.ExprVisitor

class UnaryExpr(
    val operator: Token,
    val right: Expr
) : Expr() {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitUnaryExpr(this)
    }
}