package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.ast.ExprVisitor

class GroupingExpr(
    val expression: Expr
) : Expr() {
    override fun <R> accept(visitor: ExprVisitor<R>): R {
        return visitor.visitGroupingExpr(this)
    }
}