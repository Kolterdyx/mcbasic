package me.kolterdyx.compiler

object Constants {
    val Keywords = mapOf(
        "and" to TokenType.AND,
        "or" to TokenType.OR,
        "if" to TokenType.IF,
        "else" to TokenType.ELSE,
        "while" to TokenType.WHILE,
        "for" to TokenType.FOR,
        "func" to TokenType.FUNC,
        "break" to TokenType.BREAK,
        "return" to TokenType.RETURN,
        "true" to TokenType.TRUE,
        "false" to TokenType.FALSE,
        "null" to TokenType.NULL,
        "string" to TokenType.VALUETYPE,
        "int" to TokenType.VALUETYPE,
        "float" to TokenType.VALUETYPE,
        "boolean" to TokenType.VALUETYPE,
    )
}