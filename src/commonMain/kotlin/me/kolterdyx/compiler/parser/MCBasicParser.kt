package me.kolterdyx.compiler.parser

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.exception.ParseException
import me.kolterdyx.compiler.expression.Expression
import me.kolterdyx.compiler.statement.Statement

class MCBasicParser (
    private val expressionParser: Parser<List<Token>, List<Expression>>,
    private val statementParser: Parser<List<Expression>, Statement>
) : Parser<List<Token>, Statement> {

    override fun parse(data: List<Token>): Statement {
        try {
            return statementParser.parse(expressionParser.parse(data))
        } catch (e: ParseException) {
            println(e)
            return Statement.Empty()
        }
    }
}