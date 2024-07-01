package me.kolterdyx.compiler.parser

interface Parser<R, T> {
    fun parse(data: R): T
}