package me.kolterdyx.compiler

class Scanner(
    private var source: String,
) {
    private var tokens: MutableList<Token> = mutableListOf()
    private var current: Int = 0
    private var start: Int = 0
    private var position: Pair<Int, Int> = Pair(1, 0)

    fun scanTokens(): List<Token> {

        while (!isAtEnd()) {
            // We are at the beginning of the next lexeme.
            start = current
            scanToken()
        }

        addToken(TokenType.P_EOF)
        return tokens
    }

    private fun scanToken() {
        when (val c = advance()) {
            '(' -> addToken(TokenType.P_LEFT_PAREN)
            ')' -> addToken(TokenType.P_RIGHT_PAREN)
            '{' -> addToken(TokenType.P_LEFT_BRACE)
            '}' -> addToken(TokenType.P_RIGHT_BRACE)
            ',' -> addToken(TokenType.P_COMMA)
            '.' -> addToken(TokenType.P_DOT)
            '-' -> addToken(TokenType.OP_MINUS)
            '+' -> addToken(TokenType.OP_PLUS)
            '%' -> addToken(TokenType.OP_PERCENT)
            ';' -> addToken(TokenType.P_SEMICOLON)
            '*' -> addToken(TokenType.OP_STAR)
            '!' -> addToken(if (match('=')) TokenType.OP_BANG_EQUAL else TokenType.OP_BANG)
            '=' -> addToken(if (match('=')) TokenType.OP_EQUAL_EQUAL else TokenType.OP_EQUAL)
            '<' -> addToken(if (match('=')) TokenType.OP_LESS_EQUAL else TokenType.OP_LESS)
            '>' -> addToken(if (match('=')) TokenType.OP_GREATER_EQUAL else TokenType.OP_GREATER)
            '/' -> {
                if (match('/')) {
                    // A comment goes until the end of the line.
                    while (peek() != '\n' && !isAtEnd()) advance()
                } else {
                    addToken(TokenType.OP_SLASH)
                }
            }

            ' ', '\r', '\t' -> {
            }

            '\n' -> {
                position = Pair(position.first + 1, 0)
            }

            '"' -> string()
            else -> {
                when {
                    c.isDigit() -> number()
                    c.isLetter() -> identifier()
                    else -> error("Unexpected character.")
                }
            }
        }
    }

    private fun identifier() {
        while (peek().isLetterOrDigit()) advance()

        // See if the identifier is a reserved word.
        val text = source.substring(start, current)
        val type = Constants.Keywords[text] ?: TokenType.LIT_IDENTIFIER
        if (type == TokenType.KW_TRUE || type == TokenType.KW_FALSE) {
            addToken(TokenType.LIT_BOOLEAN, type == TokenType.KW_TRUE)
        } else addToken(type)
    }

    private fun number() {
        while (peek().isDigit()) advance()

        if (peek() != '.' || !peekNext().isDigit()) {
            addToken(TokenType.LIT_INT, source.substring(start, current).toInt())
            return
        }

        // Look for a fractional part.
        if (peek() == '.' && peekNext().isDigit()) {
            // Consume the "."
            advance()

            while (peek().isDigit()) advance()
        }

        addToken(
            TokenType.LIT_FLOAT,
            source.substring(start, current).toDouble()
        )
    }

    private fun peekNext(): Char {
        return if (current + 1 >= source.length) '\u0000' else source[current + 1]
    }

    private fun string() {
        while (peek() != '"' && !isAtEnd()) {
            if (peek() == '\n') position = Pair(position.first + 1, 0)
            advance()
        }

        // Unterminated string.
        if (isAtEnd()) {
            error("Unterminated string.")
            return
        }

        // The closing "
        advance()

        // Trim the surrounding quotes.
        val value = source.substring(start + 1, current - 1)
        addToken(TokenType.LIT_STRING, value)
    }

    private fun peek(): Char {
        return if (isAtEnd()) '\u0000' else source[current]
    }

    private fun match(c: Char): Boolean {
        if (isAtEnd()) return false
        if (source[current] != c) return false

        current++
        return true
    }

    private fun addToken(type: TokenType) {
        addToken(type, null)
    }

    private fun addToken(type: TokenType, literal: Any?) {
        val text = if (type == TokenType.P_EOF) "" else source.substring(start, current)
        tokens.add(Token(type, text, literal, position))
    }

    private fun advance(): Char {
        current++
        position = Pair(position.first, position.second + 1)
        return source[current - 1]
    }

    private fun isAtEnd(): Boolean {
        return current >= source.length
    }

    private fun error(message: String) {
        println("Error at [${position.first}:${position.second}]: $message")
    }

}