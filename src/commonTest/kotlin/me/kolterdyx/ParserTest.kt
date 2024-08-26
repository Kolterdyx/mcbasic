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
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_PLUS),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_MINUS),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_STAR),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_SLASH),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_PERCENT),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_EQUAL_EQUAL),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_BANG_EQUAL),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_GREATER),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_GREATER_EQUAL),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_LESS),
            row(TokenType.LIT_INT, TokenType.LIT_INT, TokenType.OP_LESS_EQUAL),

            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_PLUS),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_MINUS),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_STAR),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_SLASH),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_EQUAL_EQUAL),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_BANG_EQUAL),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_GREATER),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_GREATER_EQUAL),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_LESS),
            row(TokenType.LIT_FLOAT, TokenType.LIT_FLOAT, TokenType.OP_LESS_EQUAL),

            row(TokenType.LIT_INT, TokenType.LIT_FLOAT, TokenType.OP_PLUS),
            row(TokenType.LIT_INT, TokenType.LIT_FLOAT, TokenType.OP_MINUS),
            row(TokenType.LIT_INT, TokenType.LIT_FLOAT, TokenType.OP_STAR),
            row(TokenType.LIT_INT, TokenType.LIT_FLOAT, TokenType.OP_SLASH),

            row(TokenType.LIT_FLOAT, TokenType.LIT_INT, TokenType.OP_PLUS),
            row(TokenType.LIT_FLOAT, TokenType.LIT_INT, TokenType.OP_MINUS),
            row(TokenType.LIT_FLOAT, TokenType.LIT_INT, TokenType.OP_STAR),
            row(TokenType.LIT_FLOAT, TokenType.LIT_INT, TokenType.OP_SLASH),

            row(TokenType.LIT_STRING, TokenType.LIT_STRING, TokenType.OP_PLUS),
            row(TokenType.LIT_STRING, TokenType.LIT_STRING, TokenType.OP_EQUAL_EQUAL),
            row(TokenType.LIT_STRING, TokenType.LIT_STRING, TokenType.OP_BANG_EQUAL),
            row(TokenType.LIT_STRING, TokenType.LIT_INT, TokenType.OP_PLUS),
            row(TokenType.LIT_STRING, TokenType.LIT_FLOAT, TokenType.OP_PLUS),
            row(TokenType.LIT_STRING, TokenType.LIT_BOOLEAN, TokenType.OP_PLUS),

            row(TokenType.LIT_BOOLEAN, TokenType.LIT_BOOLEAN, TokenType.OP_EQUAL_EQUAL),
            row(TokenType.LIT_BOOLEAN, TokenType.LIT_BOOLEAN, TokenType.OP_BANG_EQUAL),
        ) { left, right, operator ->

            val tokens = mutableListOf(
                Token(left, "1", 1, Pair(1, 1)),
                Token(operator, "", null, Pair(1, 2)),
                Token(right, "1", 1, Pair(1, 3)),
                Token(TokenType.P_EOF, "", null, Pair(1, 4)),
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
            row(TokenType.LIT_INT, TokenType.LIT_STRING, TokenType.OP_PLUS),
            row(TokenType.LIT_INT, TokenType.LIT_BOOLEAN, TokenType.OP_PLUS),
            row(TokenType.LIT_INT, TokenType.LIT_STRING, TokenType.OP_MINUS),
            row(TokenType.LIT_INT, TokenType.LIT_FLOAT, TokenType.OP_EQUAL_EQUAL),
        ) { left, right, operator ->

            val tokens = mutableListOf(
                Token(left, "1", 1, Pair(1, 1)),
                Token(operator, "", null, Pair(1, 2)),
                Token(right, "1", 1, Pair(1, 3)),
                Token(TokenType.P_EOF, "", null, Pair(1, 4)),
            )
            val parser = Parser(tokens)
            shouldThrow<ParseException> {
                parser.expression()
            }
        }
    }

    test("Unary expressions") {
        forAll(
            row(TokenType.OP_PLUS, TokenType.LIT_INT),
            row(TokenType.OP_MINUS, TokenType.LIT_INT),
            row(TokenType.OP_PLUS, TokenType.LIT_FLOAT),
            row(TokenType.OP_MINUS, TokenType.LIT_FLOAT),
            row(TokenType.OP_BANG, TokenType.LIT_BOOLEAN),
        ) { operator, right ->

            val tokens = mutableListOf(
                Token(operator, "", null, Pair(1, 1)),
                Token(right, "1", 1, Pair(1, 2)),
                Token(TokenType.P_EOF, "", null, Pair(1, 3)),
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
            row(TokenType.OP_PLUS, TokenType.LIT_BOOLEAN),
            row(TokenType.OP_MINUS, TokenType.LIT_BOOLEAN),
            row(TokenType.OP_BANG, TokenType.LIT_INT),
            row(TokenType.OP_BANG, TokenType.LIT_FLOAT),
            row(TokenType.OP_BANG, TokenType.LIT_STRING),
        ) { operator, right ->

            val tokens = mutableListOf(
                Token(operator, "", null, Pair(1, 1)),
                Token(right, "1", 1, Pair(1, 2)),
                Token(TokenType.P_EOF, "", null, Pair(1, 3)),
            )
            val parser = Parser(tokens)
            shouldThrow<ParseException> {
                parser.expression()
            }
        }
    }

    test("Literal expression") {
        forAll(
            row(TokenType.LIT_INT, 1),
            row(TokenType.LIT_FLOAT, 1.0),
            row(TokenType.LIT_STRING, "\"string\""),
            row(TokenType.LIT_BOOLEAN, true),
            row(TokenType.LIT_BOOLEAN, false),
        ) { valueType, literal ->

            val tokens = mutableListOf(
                Token(valueType, literal.toString(), literal, Pair(1, 1)),
                Token(TokenType.P_EOF, "", null, Pair(1, 2)),
            )
            val expected = LiteralExpression(Token(valueType, literal.toString(), literal, Pair(1, 1)))
            val parser = Parser(tokens)
            val expression = parser.expression()
            expression shouldBe expected
        }
    }

    test("Grouping expression") {
        forAll(
            row(TokenType.LIT_INT, 1),
            row(TokenType.LIT_FLOAT, 1.0),
            row(TokenType.LIT_STRING, "\"string\""),
            row(TokenType.LIT_BOOLEAN, true),
            row(TokenType.LIT_BOOLEAN, false),
        ) { valueType, literal ->

            val tokens = mutableListOf(
                Token(TokenType.P_LEFT_PAREN, "(", null, Pair(1, 1)),
                Token(valueType, literal.toString(), literal, Pair(1, 2)),
                Token(TokenType.P_RIGHT_PAREN, ")", null, Pair(1, 3)),
                Token(TokenType.P_EOF, "", null, Pair(1, 4)),
            )
            val expected = GroupingExpression(LiteralExpression(Token(valueType, literal.toString(), literal, Pair(1, 2))))
            val parser = Parser(tokens)
            val expression = parser.expression()
            expression shouldBe expected
        }
    }
})