package me.kolterdyx

import io.kotest.core.spec.style.FunSpec
import io.kotest.data.forAll
import io.kotest.data.row
import io.kotest.matchers.shouldBe
import me.kolterdyx.compiler.Scanner
import me.kolterdyx.compiler.TokenType


class ScannerTest : FunSpec({
    test("Literals") {
        forAll(
            row("123", TokenType.INT, 123),
            row("123.456", TokenType.FLOAT, 123.456),
            row("true", TokenType.BOOLEAN, true),
            row("false", TokenType.BOOLEAN, false),
            row("\"Hello, World!\"", TokenType.STRING, "Hello, World!"),
            row("null", TokenType.NULL, null),
        ) { input, type, value ->
            val scanner = Scanner(input)
            val token = scanner.scanTokens().first()
            token.type shouldBe type
            token.literal shouldBe value
        }
    }
    test("Keywords") {
        forAll(
            row("and", TokenType.AND, null),
            row("or", TokenType.OR, null),
            row("if", TokenType.IF, null),
            row("else", TokenType.ELSE, null),
            row("while", TokenType.WHILE, null),
            row("for", TokenType.FOR, null),
            row("func", TokenType.FUNC, null),
            row("return", TokenType.RETURN, null),
        ) {input, type, value ->
            val scanner = Scanner(input)
            val token = scanner.scanTokens().first()
            token.type shouldBe type
            token.literal shouldBe value
        }
    }
    test("Multiple tokens") {
        /// This source code isn't supposed to be valid syntax, it's just a test
        val scanner = Scanner("""
            for (int i = 0; i < 10; i = i + 1) {
                if (i % 2 == 0) {
                    print(i);
                } else {
                    while (true or false) {
                        print("Hello, World!");
                        break;
                    }
                }
                boolean test = null;
            }
        """)
        val tokens = scanner.scanTokens()
        tokens.map { it.type } shouldBe listOf(
            TokenType.FOR,
            TokenType.LEFT_PAREN,
            TokenType.VALUETYPE,
            TokenType.IDENTIFIER,
            TokenType.EQUAL,
            TokenType.INT,
            TokenType.SEMICOLON,
            TokenType.IDENTIFIER,
            TokenType.LESS,
            TokenType.INT,
            TokenType.SEMICOLON,
            TokenType.IDENTIFIER,
            TokenType.EQUAL,
            TokenType.IDENTIFIER,
            TokenType.PLUS,
            TokenType.INT,
            TokenType.RIGHT_PAREN,
            TokenType.LEFT_BRACE,
            TokenType.IF,
            TokenType.LEFT_PAREN,
            TokenType.IDENTIFIER,
            TokenType.PERCENT,
            TokenType.INT,
            TokenType.EQUAL_EQUAL,
            TokenType.INT,
            TokenType.RIGHT_PAREN,
            TokenType.LEFT_BRACE,
            TokenType.IDENTIFIER,
            TokenType.LEFT_PAREN,
            TokenType.IDENTIFIER,
            TokenType.RIGHT_PAREN,
            TokenType.SEMICOLON,
            TokenType.RIGHT_BRACE,
            TokenType.ELSE,
            TokenType.LEFT_BRACE,
            TokenType.WHILE,
            TokenType.LEFT_PAREN,
            TokenType.BOOLEAN,
            TokenType.OR,
            TokenType.BOOLEAN,
            TokenType.RIGHT_PAREN,
            TokenType.LEFT_BRACE,
            TokenType.IDENTIFIER,
            TokenType.LEFT_PAREN,
            TokenType.STRING,
            TokenType.RIGHT_PAREN,
            TokenType.SEMICOLON,
            TokenType.BREAK,
            TokenType.SEMICOLON,
            TokenType.RIGHT_BRACE,
            TokenType.RIGHT_BRACE,
            TokenType.VALUETYPE,
            TokenType.IDENTIFIER,
            TokenType.EQUAL,
            TokenType.NULL,
            TokenType.SEMICOLON,
            TokenType.RIGHT_BRACE,
            TokenType.EOF,
        )
    }
})