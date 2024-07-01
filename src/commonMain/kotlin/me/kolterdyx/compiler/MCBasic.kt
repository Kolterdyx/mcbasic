package me.kolterdyx.compiler

import me.kolterdyx.compiler.parser.ExpressionParser
import me.kolterdyx.compiler.parser.MCBasicParser
import me.kolterdyx.compiler.parser.StatementParser

class MCBasic {

    var hadError: Boolean = false
        private set

    fun tokenize() {

        val source: String = ""

        val scanner = Scanner(source)
        val tokens = scanner.scanTokens()

        val parser = MCBasicParser(
            ExpressionParser(),
            StatementParser()
        )
        parser.parse(tokens)
    }

    private fun error(position: Pair<Int, Int>, message: String) {
        println("Error at [${position.first}:${position.second}]: $message")
        hadError = true
    }

}