package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.BinaryExpression
import me.kolterdyx.compiler.expression.GroupingExpression
import me.kolterdyx.compiler.expression.LiteralExpression
import me.kolterdyx.compiler.expression.UnaryExpression

interface ExprVisitor<R> {

    fun visitBinaryExpr(expr: BinaryExpression): R

    fun visitUnaryExpr(expr: UnaryExpression): R

    fun visitLiteralExpr(expr: LiteralExpression): R

    fun visitGroupingExpr(expr: GroupingExpression): R
}