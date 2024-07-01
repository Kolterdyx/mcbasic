package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.statement.ExpressionStatement

interface StatementVisitor<T> {
    fun visitExpression(statement: ExpressionStatement): T
}