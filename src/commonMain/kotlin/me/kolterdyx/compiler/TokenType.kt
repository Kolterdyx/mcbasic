package me.kolterdyx.compiler

enum class TokenType {
    // Single-character tokens.
    LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE,
    COMMA, DOT, MINUS, PLUS, SEMICOLON, SLASH, STAR,
    PERCENT,

    // One or two character tokens.
    BANG, BANG_EQUAL,
    EQUAL, EQUAL_EQUAL,
    GREATER, GREATER_EQUAL,
    LESS, LESS_EQUAL,

    // Literals.
    IDENTIFIER, STRING, INT, FLOAT, BOOLEAN,

    // Keywords.
    AND, CLASS, ELSE, FALSE, FUNC, FOR, IF, NULL, OR,
    RETURN, TRUE, WHILE, VALUETYPE, BREAK,

    EOF,
}