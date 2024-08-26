package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.expression.*
import me.kolterdyx.compiler.statement.BlockStatement
import me.kolterdyx.compiler.statement.ExpressionStatement
import me.kolterdyx.compiler.statement.VariableDeclarationStatement

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

    override fun visitVariable(expr: VariableExpression): String {
        TODO("Not yet implemented")
    }

    override fun visitAssignment(expr: AssignmentExpression): String {
        TODO("Not yet implemented")
    }

    // End Expression Visitor

    // Statement Visitor

    override fun visitExpression(statement: ExpressionStatement): String {
        return statement.expression.accept(this)
    }

    override fun visitVariableDeclaration(statement: VariableDeclarationStatement): String {
        TODO("Not yet implemented")
    }

    override fun visitBlock(statement: BlockStatement): String {
        TODO("Not yet implemented")
    }

    // End Statement Visitor
}