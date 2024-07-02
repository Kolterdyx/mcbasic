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

        addToken(TokenType.EOF)
        return tokens
    }

    private fun scanToken() {
        when (val c = advance()) {
            '(' -> addToken(TokenType.LEFT_PAREN)
            ')' -> addToken(TokenType.RIGHT_PAREN)
            '{' -> addToken(TokenType.LEFT_BRACE)
            '}' -> addToken(TokenType.RIGHT_BRACE)
            ',' -> addToken(TokenType.COMMA)
            '.' -> addToken(TokenType.DOT)
            '-' -> addToken(TokenType.MINUS)
            '+' -> addToken(TokenType.PLUS)
            '%' -> addToken(TokenType.PERCENT)
            ';' -> addToken(TokenType.SEMICOLON)
            '*' -> addToken(TokenType.STAR)
            '!' -> addToken(if (match('=')) TokenType.BANG_EQUAL else TokenType.BANG)
            '=' -> addToken(if (match('=')) TokenType.EQUAL_EQUAL else TokenType.EQUAL)
            '<' -> addToken(if (match('=')) TokenType.LESS_EQUAL else TokenType.LESS)
            '>' -> addToken(if (match('=')) TokenType.GREATER_EQUAL else TokenType.GREATER)
            '/' -> {
                if (match('/')) {
                    // A comment goes until the end of the line.
                    while (peek() != '\n' && !isAtEnd()) advance()
                } else {
                    addToken(TokenType.SLASH)
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
        val type = Constants.Keywords[text] ?: TokenType.IDENTIFIER
        if (type == TokenType.TRUE || type == TokenType.FALSE) {
            addToken(TokenType.BOOLEAN, type == TokenType.TRUE)
        } else addToken(type)
    }

    private fun number() {
        while (peek().isDigit()) advance()

        if (peek() != '.' || !peekNext().isDigit()) {
            addToken(TokenType.INT, source.substring(start, current).toInt())
            return
        }

        // Look for a fractional part.
        if (peek() == '.' && peekNext().isDigit()) {
            // Consume the "."
            advance()

            while (peek().isDigit()) advance()
        }

        addToken(
            TokenType.FLOAT,
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
        addToken(TokenType.STRING, value)
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
        val text = if (type == TokenType.EOF) "" else source.substring(start, current)
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