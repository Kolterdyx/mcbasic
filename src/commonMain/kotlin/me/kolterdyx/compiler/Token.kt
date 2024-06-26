package me.kolterdyx.compiler

class Token(
    val type: TokenType,
    val lexeme: String,
    val literal: Any?,
    val pos: Pair<Int, Int>
) {
    override fun toString(): String {
        return "$type $lexeme $literal"
    }
}