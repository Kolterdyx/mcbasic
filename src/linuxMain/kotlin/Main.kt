import com.github.ajalt.clikt.core.subcommands
import me.kolterdyx.MCBasic
import me.kolterdyx.command.BuildCommand
import me.kolterdyx.command.InitCommand

fun main(args: Array<String>) = MCBasic()
    .subcommands(InitCommand(), BuildCommand())
    .main(args)