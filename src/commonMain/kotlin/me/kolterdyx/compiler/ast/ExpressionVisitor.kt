package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.*

interface ExpressionVisitor<R> {

    fun visitBinary(expr: BinaryExpression): R

    fun visitUnary(expr: UnaryExpression): R

    fun visitLiteral(expr: LiteralExpression): R

    fun visitGrouping(expr: GroupingExpression): R

    fun visitVariable(expr: VariableExpression): R

    fun visitAssignment(expr: AssignmentExpression): R
}
