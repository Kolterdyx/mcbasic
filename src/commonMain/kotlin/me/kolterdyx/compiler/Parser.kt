package me.kolterdyx.compiler

class Parser(
    private val tokens: List<Token>
) {
    fun parse() {
    }

    fun error(token: Token, message: String) {
        println("Error at ${token.pos.first}:${token.pos.second}: $message")
    }
}