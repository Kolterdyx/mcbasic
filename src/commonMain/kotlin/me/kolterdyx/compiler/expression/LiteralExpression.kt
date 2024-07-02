package me.kolterdyx.compiler.expression

import me.kolterdyx.compiler.Token
import me.kolterdyx.compiler.TokenType
import me.kolterdyx.compiler.ValueType
import me.kolterdyx.compiler.ast.ExpressionVisitor
import me.kolterdyx.compiler.exception.ParseException

fun getLiteralValueType(token: Token): ValueType {
    return when {
        token.type == TokenType.INT -> ValueType.INT
        token.type == TokenType.FLOAT -> ValueType.FLOAT
        token.type == TokenType.STRING -> ValueType.STRING
        token.type == TokenType.BOOLEAN -> ValueType.BOOLEAN
        else -> throw ParseException(token, "Invalid literal type: ${token.type}")
    }
}

class LiteralExpression(
    val value: Token
) : Expression(getLiteralValueType(value)) {
    override fun <R> accept(visitor: ExpressionVisitor<R>): R {
        return visitor.visitLiteral(this)
    }

    override fun toString(): String {
        return "LiteralExpression($value)"
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is LiteralExpression) return false
        if (!super.equals(other)) return false
        if (value != other.value) return false
        return true
    }

    override fun hashCode(): Int {
        var result = super.hashCode()
        result = 31 * result + (value?.hashCode() ?: 0)
        return result
    }
}