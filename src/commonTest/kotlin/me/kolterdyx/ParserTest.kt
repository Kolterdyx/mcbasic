package me.kolterdyx

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.FunSpec
import io.kotest.data.forAll
import io.kotest.data.row
import io.kotest.matchers.shouldBe
import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.TokenType
import me.kolterdyx.compiler.exception.ParseException
import me.kolterdyx.compiler.expression.BinaryExpression
import me.kolterdyx.compiler.expression.LiteralExpression
import me.kolterdyx.compiler.Parser
import me.kolterdyx.compiler.expression.GroupingExpression
import me.kolterdyx.compiler.expression.UnaryExpression

class ParserTest : FunSpec({
    test("Binary expressions") {
        forAll(
            row(TokenType.INT, TokenType.INT, TokenType.PLUS),
            row(TokenType.INT, TokenType.INT, TokenType.MINUS),
            row(TokenType.INT, TokenType.INT, TokenType.STAR),
            row(TokenType.INT, TokenType.INT, TokenType.SLASH),
            row(TokenType.INT, TokenType.INT, TokenType.PERCENT),
            row(TokenType.INT, TokenType.INT, TokenType.EQUAL_EQUAL),
            row(TokenType.INT, TokenType.INT, TokenType.BANG_EQUAL),
            row(TokenType.INT, TokenType.INT, TokenType.GREATER),
            row(TokenType.INT, TokenType.INT, TokenType.GREATER_EQUAL),
            row(TokenType.INT, TokenType.INT, TokenType.LESS),
            row(TokenType.INT, TokenType.INT, TokenType.LESS_EQUAL),

            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.PLUS),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.MINUS),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.STAR),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.SLASH),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.EQUAL_EQUAL),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.BANG_EQUAL),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.GREATER),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.GREATER_EQUAL),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.LESS),
            row(TokenType.FLOAT, TokenType.FLOAT, TokenType.LESS_EQUAL),

            row(TokenType.INT, TokenType.FLOAT, TokenType.PLUS),
            row(TokenType.INT, TokenType.FLOAT, TokenType.MINUS),
            row(TokenType.INT, TokenType.FLOAT, TokenType.STAR),
            row(TokenType.INT, TokenType.FLOAT, TokenType.SLASH),

            row(TokenType.FLOAT, TokenType.INT, TokenType.PLUS),
            row(TokenType.FLOAT, TokenType.INT, TokenType.MINUS),
            row(TokenType.FLOAT, TokenType.INT, TokenType.STAR),
            row(TokenType.FLOAT, TokenType.INT, TokenType.SLASH),

            row(TokenType.STRING, TokenType.STRING, TokenType.PLUS),
            row(TokenType.STRING, TokenType.STRING, TokenType.EQUAL_EQUAL),
            row(TokenType.STRING, TokenType.STRING, TokenType.BANG_EQUAL),
            row(TokenType.STRING, TokenType.INT, TokenType.PLUS),
            row(TokenType.STRING, TokenType.FLOAT, TokenType.PLUS),
            row(TokenType.STRING, TokenType.BOOLEAN, TokenType.PLUS),

            row(TokenType.BOOLEAN, TokenType.BOOLEAN, TokenType.EQUAL_EQUAL),
            row(TokenType.BOOLEAN, TokenType.BOOLEAN, TokenType.BANG_EQUAL),
        ) { left, right, operator ->

            val tokens = mutableListOf(
                Token(left, "1", 1, Pair(1, 1)),
                Token(operator, "", null, Pair(1, 2)),
                Token(right, "1", 1, Pair(1, 3)),
                Token(TokenType.EOF, "", null, Pair(1, 4)),
            )
            val expected = BinaryExpression(
                LiteralExpression(Token(left, "1", 1, Pair(1, 1))),
                Token(operator, "", null, Pair(1, 2)),
                LiteralExpression(Token(right, "1", 1, Pair(1, 3))),
            )
            val parser = Parser(tokens)
            val expression = parser.expression()
            expression shouldBe expected
        }
    }

    test("Invalid binary expressions") {
        forAll(
            row(TokenType.INT, TokenType.STRING, TokenType.PLUS),
            row(TokenType.INT, TokenType.BOOLEAN, TokenType.PLUS),
            row(TokenType.INT, TokenType.STRING, TokenType.MINUS),
            row(TokenType.INT, TokenType.FLOAT, TokenType.EQUAL_EQUAL),
        ) { left, right, operator ->

            val tokens = mutableListOf(
                Token(left, "1", 1, Pair(1, 1)),
                Token(operator, "", null, Pair(1, 2)),
                Token(right, "1", 1, Pair(1, 3)),
                Token(TokenType.EOF, "", null, Pair(1, 4)),
            )
            val parser = Parser(tokens)
            shouldThrow<ParseException> {
                parser.expression()
            }
        }
    }

    test("Unary expressions") {
        forAll(
            row(TokenType.PLUS, TokenType.INT),
            row(TokenType.MINUS, TokenType.INT),
            row(TokenType.PLUS, TokenType.FLOAT),
            row(TokenType.MINUS, TokenType.FLOAT),
            row(TokenType.BANG, TokenType.BOOLEAN),
        ) { operator, right ->

            val tokens = mutableListOf(
                Token(operator, "", null, Pair(1, 1)),
                Token(right, "1", 1, Pair(1, 2)),
                Token(TokenType.EOF, "", null, Pair(1, 3)),
            )
            val expected = UnaryExpression(
                Token(operator, "", null, Pair(1, 1)),
                LiteralExpression(Token(right, "1", 1, Pair(1, 2))),
            )
            val parser = Parser(tokens)
            val expression = parser.expression()
            expression shouldBe expected
        }
    }

    test("Invalid unary expressions") {
        forAll(
            row(TokenType.PLUS, TokenType.BOOLEAN),
            row(TokenType.MINUS, TokenType.BOOLEAN),
            row(TokenType.BANG, TokenType.INT),
            row(TokenType.BANG, TokenType.FLOAT),
            row(TokenType.BANG, TokenType.STRING),
        ) { operator, right ->

            val tokens = mutableListOf(
                Token(operator, "", null, Pair(1, 1)),
                Token(right, "1", 1, Pair(1, 2)),
                Token(TokenType.EOF, "", null, Pair(1, 3)),
            )
            val parser = Parser(tokens)
            shouldThrow<ParseException> {
                parser.expression()
            }
        }
    }

    test("Literal expression") {
        forAll(
            row(TokenType.INT, 1),
            row(TokenType.FLOAT, 1.0),
            row(TokenType.STRING, "\"string\""),
            row(TokenType.BOOLEAN, true),
            row(TokenType.BOOLEAN, false),
        ) { valueType, literal ->

            val tokens = mutableListOf(
                Token(valueType, literal.toString(), literal, Pair(1, 1)),
                Token(TokenType.EOF, "", null, Pair(1, 2)),
            )
            val expected = LiteralExpression(Token(valueType, literal.toString(), literal, Pair(1, 1)))
            val parser = Parser(tokens)
            val expression = parser.expression()
            expression shouldBe expected
        }
    }

    test("Grouping expression") {
        forAll(
            row(TokenType.INT, 1),
            row(TokenType.FLOAT, 1.0),
            row(TokenType.STRING, "\"string\""),
            row(TokenType.BOOLEAN, true),
            row(TokenType.BOOLEAN, false),
        ) { valueType, literal ->

            val tokens = mutableListOf(
                Token(TokenType.LEFT_PAREN, "(", null, Pair(1, 1)),
                Token(valueType, literal.toString(), literal, Pair(1, 2)),
                Token(TokenType.RIGHT_PAREN, ")", null, Pair(1, 3)),
                Token(TokenType.EOF, "", null, Pair(1, 4)),
            )
            val expected = GroupingExpression(LiteralExpression(Token(valueType, literal.toString(), literal, Pair(1, 2))))
            val parser = Parser(tokens)
            val expression = parser.expression()
            expression shouldBe expected
        }
    }
})