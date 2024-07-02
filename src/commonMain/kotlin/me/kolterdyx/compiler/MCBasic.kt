package me.kolterdyx.compiler

import me.kolterdyx.compiler.parser.ExpressionParser
import me.kolterdyx.compiler.parser.MCBasicParser
import me.kolterdyx.compiler.parser.StatementParser

class MCBasic {

    var hadError: Boolean = false
        private set

    fun tokenize() {

        val source: String = "1 + 2"

        val scanner = Scanner(source)
        val tokens = scanner.scanTokens()
        println(tokens)

        val parser = MCBasicParser(
            ExpressionParser(),
            StatementParser()
        )
        try {
            parser.parse(tokens)
        } catch (e: Exception) {
            println(e)
        }
    }

    private fun error(position: Pair<Int, Int>, message: String) {
        println("Error at [${position.first}:${position.second}]: $message")
        hadError = true
    }

}