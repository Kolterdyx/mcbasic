package me.kolterdyx.compiler

class MCBasic {

    var hadError: Boolean = false
        private set

    fun tokenize() {

        val source: String = ""

        val scanner = Scanner(source)
        val tokens = scanner.scanTokens()
    }

    private fun error(position: Pair<Int, Int>, message: String) {
        println("Error at [${position.first}:${position.second}]: $message")
        hadError = true
    }

}