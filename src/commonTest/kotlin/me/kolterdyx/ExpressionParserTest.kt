package me.kolterdyx

import io.kotest.core.spec.style.FunSpec
import io.kotest.data.forAll
import io.kotest.data.row
import io.kotest.matchers.shouldBe
import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.TokenType
import me.kolterdyx.compiler.expression.BinaryExpression
import me.kolterdyx.compiler.expression.LiteralExpression
import me.kolterdyx.compiler.parser.ExpressionParser

class ExpressionParserTest : FunSpec({
    test("Binary Expressions") {
        forAll(
            row(
                listOf(
                    Token(TokenType.INT, "1", 1, Pair(1, 0)),
                    Token(TokenType.PLUS, "+", null, Pair(1, 1)),
                    Token(TokenType.INT, "2", 2, Pair(1, 2)),
                    Token(TokenType.EOF, "", null, Pair(1, 3)),
                ),
                BinaryExpression(
                    LiteralExpression(1),
                    Token(TokenType.PLUS, "+", null, Pair(1, 1)),
                    LiteralExpression(2)
                )
            ),
        ) { input, expected ->
             val parser = ExpressionParser()
             val expression = parser.parse(input)
             expression[0] shouldBe expected
        }
    }
})