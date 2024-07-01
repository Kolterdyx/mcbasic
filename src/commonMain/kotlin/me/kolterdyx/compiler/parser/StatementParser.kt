package me.kolterdyx.compiler.parser

import me.kolterdyx.compiler.expression.Expression
import me.kolterdyx.compiler.statement.Statement

class StatementParser : Parser<List<Expression>, Statement> {
    override fun parse(data: List<Expression>): Statement {
        TODO("Not yet implemented")
    }
}