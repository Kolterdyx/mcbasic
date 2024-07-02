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

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is Token) return false
        if (type != other.type) return false
        if (lexeme != other.lexeme) return false
        if (literal != other.literal) return false
        if (pos != other.pos) return false
        return true
    }

    override fun hashCode(): Int {
        var result = type.hashCode()
        result = 31 * result + lexeme.hashCode()
        result = 31 * result + (literal?.hashCode() ?: 0)
        result = 31 * result + pos.hashCode()
        return result
    }
}