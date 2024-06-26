package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.BinaryExpr
import me.kolterdyx.compiler.expression.GroupingExpr
import me.kolterdyx.compiler.expression.LiteralExpr
import me.kolterdyx.compiler.expression.UnaryExpr

class AstPrinter : ExprVisitor<String> {
    override fun visitBinaryExpr(expr: BinaryExpr): String {
        val left = expr.left.accept(this)
        val right = expr.right.accept(this)
        return "(${expr.operator.lexeme} $left $right)"
    }

    override fun visitUnaryExpr(expr: UnaryExpr): String {
        val right = expr.right.accept(this)
        return "(${expr.operator.lexeme} $right)"
    }

    override fun visitLiteralExpr(expr: LiteralExpr): String {
        return expr.value.toString()
    }

    override fun visitGroupingExpr(expr: GroupingExpr): String {
        return "(${expr.expression.accept(this)})"
    }
}