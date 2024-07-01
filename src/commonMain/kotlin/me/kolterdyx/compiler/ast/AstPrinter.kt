package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.BinaryExpression
import me.kolterdyx.compiler.expression.GroupingExpression
import me.kolterdyx.compiler.expression.LiteralExpression
import me.kolterdyx.compiler.expression.UnaryExpression
import me.kolterdyx.compiler.statement.ExpressionStatement

class AstPrinter :
    ExpressionVisitor<String>,
    StatementVisitor<String>
{
    // Expression Visitor

    override fun visitBinary(expr: BinaryExpression): String {
        val left = expr.left.accept(this)
        val right = expr.right.accept(this)
        return "(${expr.operator.lexeme} $left $right)"
    }

    override fun visitUnary(expr: UnaryExpression): String {
        val right = expr.right.accept(this)
        return "(${expr.operator.lexeme} $right)"
    }

    override fun visitLiteral(expr: LiteralExpression): String {
        return expr.value.toString()
    }

    override fun visitGrouping(expr: GroupingExpression): String {
        return "(${expr.expression.accept(this)})"
    }

    // End Expression Visitor

    // Statement Visitor

    override fun visitExpression(statement: ExpressionStatement): String {
        return statement.expression.accept(this)
    }

    // End Statement Visitor
}