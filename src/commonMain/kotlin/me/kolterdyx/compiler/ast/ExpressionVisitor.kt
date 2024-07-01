package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.BinaryExpression
import me.kolterdyx.compiler.expression.GroupingExpression
import me.kolterdyx.compiler.expression.LiteralExpression
import me.kolterdyx.compiler.expression.UnaryExpression

interface ExpressionVisitor<R> {

    fun visitBinary(expr: BinaryExpression): R

    fun visitUnary(expr: UnaryExpression): R

    fun visitLiteral(expr: LiteralExpression): R

    fun visitGrouping(expr: GroupingExpression): R
}