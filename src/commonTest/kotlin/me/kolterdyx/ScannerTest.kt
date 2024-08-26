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
            row("123", TokenType.LIT_INT, 123),
            row("123.456", TokenType.LIT_FLOAT, 123.456),
            row("true", TokenType.LIT_BOOLEAN, true),
            row("false", TokenType.LIT_BOOLEAN, false),
            row("\"Hello, World!\"", TokenType.LIT_STRING, "Hello, World!"),
            row("null", TokenType.KW_NULL, null),
        ) { input, type, value ->
            val scanner = Scanner(input)
            val token = scanner.scanTokens().first()
            token.type shouldBe type
            token.literal shouldBe value
        }
    }
    test("Keywords") {
        forAll(
            row("and", TokenType.KW_AND, null),
            row("or", TokenType.KW_OR, null),
            row("if", TokenType.KW_IF, null),
            row("else", TokenType.KW_ELSE, null),
            row("break", TokenType.KW_BREAK, null),
            row("while", TokenType.KW_WHILE, null),
            row("for", TokenType.KW_FOR, null),
            row("func", TokenType.KW_FUNC, null),
            row("return", TokenType.KW_RETURN, null),
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
            TokenType.KW_FOR,
            TokenType.P_LEFT_PAREN,
            TokenType.KW_VALUETYPE,
            TokenType.LIT_IDENTIFIER,
            TokenType.OP_EQUAL,
            TokenType.LIT_INT,
            TokenType.P_SEMICOLON,
            TokenType.LIT_IDENTIFIER,
            TokenType.OP_LESS,
            TokenType.LIT_INT,
            TokenType.P_SEMICOLON,
            TokenType.LIT_IDENTIFIER,
            TokenType.OP_EQUAL,
            TokenType.LIT_IDENTIFIER,
            TokenType.OP_PLUS,
            TokenType.LIT_INT,
            TokenType.P_RIGHT_PAREN,
            TokenType.P_LEFT_BRACE,
            TokenType.KW_IF,
            TokenType.P_LEFT_PAREN,
            TokenType.LIT_IDENTIFIER,
            TokenType.OP_PERCENT,
            TokenType.LIT_INT,
            TokenType.OP_EQUAL_EQUAL,
            TokenType.LIT_INT,
            TokenType.P_RIGHT_PAREN,
            TokenType.P_LEFT_BRACE,
            TokenType.LIT_IDENTIFIER,
            TokenType.P_LEFT_PAREN,
            TokenType.LIT_IDENTIFIER,
            TokenType.P_RIGHT_PAREN,
            TokenType.P_SEMICOLON,
            TokenType.P_RIGHT_BRACE,
            TokenType.KW_ELSE,
            TokenType.P_LEFT_BRACE,
            TokenType.KW_WHILE,
            TokenType.P_LEFT_PAREN,
            TokenType.LIT_BOOLEAN,
            TokenType.KW_OR,
            TokenType.LIT_BOOLEAN,
            TokenType.P_RIGHT_PAREN,
            TokenType.P_LEFT_BRACE,
            TokenType.LIT_IDENTIFIER,
            TokenType.P_LEFT_PAREN,
            TokenType.LIT_STRING,
            TokenType.P_RIGHT_PAREN,
            TokenType.P_SEMICOLON,
            TokenType.KW_BREAK,
            TokenType.P_SEMICOLON,
            TokenType.P_RIGHT_BRACE,
            TokenType.P_RIGHT_BRACE,
            TokenType.KW_VALUETYPE,
            TokenType.LIT_IDENTIFIER,
            TokenType.OP_EQUAL,
            TokenType.KW_NULL,
            TokenType.P_SEMICOLON,
            TokenType.P_RIGHT_BRACE,
            TokenType.P_EOF,
        )
    }
})