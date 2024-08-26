package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.TokenType
import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor
import me.kolterdyx.compiler.exception.ParseException


class BinaryExpression(
    val left: Expression,
    val operator: Token,
    val right: Expression
) : Expression(checkCompatibility(left.valueType, right.valueType, operator)) {


    companion object {
        private val ValidIntOperators = mapOf(
            ValueType.INT to setOf(
                TokenType.OP_PLUS,
                TokenType.OP_MINUS,
                TokenType.OP_STAR,
                TokenType.OP_SLASH,
                TokenType.OP_PERCENT,
                TokenType.OP_EQUAL_EQUAL,
                TokenType.OP_BANG_EQUAL,
                TokenType.OP_GREATER,
                TokenType.OP_GREATER_EQUAL,
                TokenType.OP_LESS,
                TokenType.OP_LESS_EQUAL,
            ),
            ValueType.FLOAT to setOf(
                TokenType.OP_PLUS,
                TokenType.OP_MINUS,
                TokenType.OP_STAR,
                TokenType.OP_SLASH,
            ),
        )
        private val ValidIntRightHandTypes = ValidIntOperators.keys

        private val ValidFloatOperators = mapOf(
            ValueType.INT to setOf(
                TokenType.OP_PLUS,
                TokenType.OP_MINUS,
                TokenType.OP_STAR,
                TokenType.OP_SLASH,
            ),
            ValueType.FLOAT to setOf(
                TokenType.OP_PLUS,
                TokenType.OP_MINUS,
                TokenType.OP_STAR,
                TokenType.OP_SLASH,
                TokenType.OP_EQUAL_EQUAL,
                TokenType.OP_BANG_EQUAL,
                TokenType.OP_GREATER,
                TokenType.OP_GREATER_EQUAL,
                TokenType.OP_LESS,
                TokenType.OP_LESS_EQUAL,
            ),
        )
        private val ValidFloatRightHandTypes = ValidFloatOperators.keys

        private val ValidStringOperators = mapOf(
            ValueType.STRING to setOf(
                TokenType.OP_PLUS,
                TokenType.OP_EQUAL_EQUAL,
                TokenType.OP_BANG_EQUAL
            ),
            ValueType.INT to setOf(
                TokenType.OP_PLUS,
            ),
            ValueType.FLOAT to setOf(
                TokenType.OP_PLUS,
            ),
            ValueType.BOOLEAN to setOf(
                TokenType.OP_PLUS,
            )
        )
        private val ValidStringRightHandTypes = ValidStringOperators.keys

        private val ValidBooleanOperators = mapOf(
            ValueType.BOOLEAN to setOf(
                TokenType.KW_AND,
                TokenType.KW_OR,
                TokenType.OP_EQUAL_EQUAL,
                TokenType.OP_BANG_EQUAL
            )
        )
        private val ValidBooleanRightHandTypes = ValidBooleanOperators.keys


        private fun checkCompatibility(right: ValueType, operator: Token, validTypes: Set<ValueType>, validOperators: Map<ValueType, Set<TokenType>>): Boolean {
            return validTypes.contains(right) && validOperators[right]!!.contains(operator.type)
        }

        fun checkCompatibility(left: ValueType, right: ValueType, operator: Token): ValueType {
            val isCompatible = when (left) {
                ValueType.INT -> checkCompatibility(right, operator, ValidIntRightHandTypes, ValidIntOperators)
                ValueType.FLOAT -> checkCompatibility(right, operator, ValidFloatRightHandTypes, ValidFloatOperators)
                ValueType.STRING -> checkCompatibility(right, operator, ValidStringRightHandTypes, ValidStringOperators)
                ValueType.BOOLEAN -> checkCompatibility(right, operator, ValidBooleanRightHandTypes, ValidBooleanOperators)
            }
            if (!isCompatible) {
                throw ParseException(operator, "Incompatible types: $left and $right")
            }
            return left
        }
    }


    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitBinary(this)
    }

    override fun toString(): String {
        return "BinaryExpression($left, $operator, $right)"
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is BinaryExpression) return false
        return this.left == other.left && this.operator == other.operator && this.right == other.right
    }

    override fun hashCode(): Int {
        var result = super.hashCode()
        result = 31 * result + left.hashCode()
        result = 31 * result + operator.hashCode()
        result = 31 * result + right.hashCode()
        return result
    }
}