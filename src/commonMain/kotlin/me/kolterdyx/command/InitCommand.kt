package me.kolterdyx.command

import com.github.ajalt.clikt.core.CliktCommand

class InitCommand : CliktCommand() {
    override fun run() {
        println("Initializing project...")
    }
}