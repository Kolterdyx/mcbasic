package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.BinaryExpression
import me.kolterdyx.compiler.expression.GroupingExpression
import me.kolterdyx.compiler.expression.LiteralExpression
import me.kolterdyx.compiler.expression.UnaryExpression

class AstPrinter : ExprVisitor<String> {
    override fun visitBinaryExpr(expr: BinaryExpression): String {
        val left = expr.left.accept(this)
        val right = expr.right.accept(this)
        return "(${expr.operator.lexeme} $left $right)"
    }

    override fun visitUnaryExpr(expr: UnaryExpression): String {
        val right = expr.right.accept(this)
        return "(${expr.operator.lexeme} $right)"
    }

    override fun visitLiteralExpr(expr: LiteralExpression): String {
        return expr.value.toString()
    }

    override fun visitGroupingExpr(expr: GroupingExpression): String {
        return "(${expr.expression.accept(this)})"
    }
}