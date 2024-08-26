package me.kolterdyx.compiler

object Constants {
    val Keywords = mapOf(
        "and" to TokenType.KW_AND,
        "or" to TokenType.KW_OR,
        "if" to TokenType.KW_IF,
        "else" to TokenType.KW_ELSE,
        "while" to TokenType.KW_WHILE,
        "for" to TokenType.KW_FOR,
        "func" to TokenType.KW_FUNC,
        "break" to TokenType.KW_BREAK,
        "return" to TokenType.KW_RETURN,
        "true" to TokenType.KW_TRUE,
        "false" to TokenType.KW_FALSE,
        "null" to TokenType.KW_NULL,
        "string" to TokenType.KW_VALUETYPE,
        "int" to TokenType.KW_VALUETYPE,
        "float" to TokenType.KW_VALUETYPE,
        "boolean" to TokenType.KW_VALUETYPE,
    )
}