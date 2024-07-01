package me.kolterdyx.compiler.exception

import me.kolterdyx.compiler.Token

class ParseException(
    token: Token,
    message: String,
) : Exception("Error at [${token.pos.first}:${token.pos.second}]: $message") {
}