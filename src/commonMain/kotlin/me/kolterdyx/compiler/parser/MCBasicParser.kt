package me.kolterdyx.compiler.parser

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.ast.AstPrinter
import me.kolterdyx.compiler.exception.ParseException
import me.kolterdyx.compiler.expression.Expression
import me.kolterdyx.compiler.statement.Statement

class MCBasicParser (
    private val expressionParser: Parser<List<Token>, List<Expression>>,
    private val statementParser: Parser<List<Expression>, Statement>
) : Parser<List<Token>, Statement> {

    override fun parse(data: List<Token>): Statement {
        try {
            val expressions = expressionParser.parse(data)
            val astPrinter = AstPrinter()
            for (expression in expressions) {
                println(expression.accept(astPrinter))
            }
            return statementParser.parse(expressions)
        } catch (e: ParseException) {
            println(e)
            return Statement.Empty()
        }
    }
}