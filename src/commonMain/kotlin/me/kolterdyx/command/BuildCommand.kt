package me.kolterdyx.command

import com.github.ajalt.clikt.core.CliktCommand

class BuildCommand : CliktCommand() {
    override fun run() {
        println("Building project...")
    }
}