package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.BinaryExpr
import me.kolterdyx.compiler.expression.GroupingExpr
import me.kolterdyx.compiler.expression.LiteralExpr
import me.kolterdyx.compiler.expression.UnaryExpr

interface ExprVisitor<R> {

    fun visitBinaryExpr(expr: BinaryExpr): R

    fun visitUnaryExpr(expr: UnaryExpr): R

    fun visitLiteralExpr(expr: LiteralExpr): R

    fun visitGroupingExpr(expr: GroupingExpr): R
}