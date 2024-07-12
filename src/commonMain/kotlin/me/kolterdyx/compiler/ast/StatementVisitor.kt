package me.kolterdyx.compiler.ast

import me.kolterdyx.compiler.statement.BlockStatement
import me.kolterdyx.compiler.statement.ExpressionStatement
import me.kolterdyx.compiler.statement.VariableDeclarationStatement

interface StatementVisitor<T> {
    fun visitExpression(statement: ExpressionStatement): T
    fun visitVariableDeclaration(statement: VariableDeclarationStatement): T
    fun visitBlock(statement: BlockStatement): T
}