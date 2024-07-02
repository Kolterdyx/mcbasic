package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.TokenType
import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor
import me.kolterdyx.compiler.exception.ParseException

object Constants {
    val ValidIntRightHandTypes = setOf(ValueType.INT, ValueType.FLOAT)
    val ValidIntOperators = mapOf(
        ValueType.INT to setOf(
            TokenType.PLUS,
            TokenType.MINUS,
            TokenType.STAR,
            TokenType.SLASH,
            TokenType.PERCENT,
        ),
        ValueType.FLOAT to setOf(
            TokenType.PLUS,
            TokenType.MINUS,
            TokenType.STAR,
            TokenType.SLASH,
        ),
    )

    val ValidFloatRightHandTypes = setOf(ValueType.INT, ValueType.FLOAT)
    val ValidFloatOperators = mapOf(
        ValueType.INT to setOf(
            TokenType.PLUS,
            TokenType.MINUS,
            TokenType.STAR,
            TokenType.SLASH,
        ),
        ValueType.FLOAT to setOf(
            TokenType.PLUS,
            TokenType.MINUS,
            TokenType.STAR,
            TokenType.SLASH,
        ),
    )

    val ValidStringRightHandTypes = setOf(
        ValueType.STRING,
        ValueType.INT,
        ValueType.FLOAT,
        ValueType.BOOLEAN
    )
    val ValidStringOperators = mapOf(
        ValueType.STRING to setOf(
            TokenType.PLUS,
            TokenType.EQUAL_EQUAL,
            TokenType.BANG_EQUAL
        ),
        ValueType.INT to setOf(
            TokenType.PLUS,
        ),
        ValueType.FLOAT to setOf(
            TokenType.PLUS,
        ),
        ValueType.BOOLEAN to setOf(
            TokenType.PLUS,
        )
    )

    val ValidBooleanRightHandTypes = setOf(ValueType.BOOLEAN)
    val ValidBooleanOperators = mapOf(
        ValueType.BOOLEAN to setOf(
            TokenType.AND,
            TokenType.OR,
            TokenType.EQUAL_EQUAL,
            TokenType.BANG_EQUAL
        )
    )
}

fun checkCompatibility(right: ValueType, operator: Token, validTypes: Set<ValueType>, validOperators: Map<ValueType, Set<TokenType>>): Boolean {
    return validTypes.contains(right) && validOperators[right]!!.contains(operator.type)
}

fun checkCompatibility(left: ValueType, right: ValueType, operator: Token): ValueType {
    val isCompatible = when (left) {
        ValueType.INT -> checkCompatibility(right, operator, Constants.ValidIntRightHandTypes, Constants.ValidIntOperators)
        ValueType.FLOAT -> checkCompatibility(right, operator, Constants.ValidFloatRightHandTypes, Constants.ValidFloatOperators)
        ValueType.STRING -> checkCompatibility(right, operator, Constants.ValidStringRightHandTypes, Constants.ValidStringOperators)
        ValueType.BOOLEAN -> checkCompatibility(right, operator, Constants.ValidBooleanRightHandTypes, Constants.ValidBooleanOperators)
    }
    if (!isCompatible) {
        throw ParseException(operator, "Incompatible types: $left and $right")
    }
    return left
}

class BinaryExpression(
    val left: Expression,
    val operator: Token,
    val right: Expression
) : Expression(checkCompatibility(left.valueType, right.valueType, operator)) {
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