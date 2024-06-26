package me.kolterdyx.command

import com.github.ajalt.clikt.core.CliktCommand
import me.kolterdyx.compiler.MCBasic

class BuildCommand : CliktCommand() {
    override fun run() {
        val mcBasic = MCBasic()
        mcBasic.tokenize()
        if (mcBasic.hadError) {
            return
        }
    }
}