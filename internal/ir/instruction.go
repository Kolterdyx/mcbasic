package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
)

const (
	_      interfaces.InstructionType = ""
	Set                               = "set"    // `set <storage> <path> <value>`
	Copy                              = "copy"   // `copy <storage from> <path from> <storage to> <path to>`
	Remove                            = "remove" // `remove <storage> <path>`
	Math                              = "math"   // `math <operation>`
	Load                              = "load"   // `load <path> <score>`
	Store                             = "store"  // `store <score> <path>`
	Score                             = "score"  // `score <target> <score>`
	Append                            = "append" // `append <listPath> <valuePath>`
	Size                              = "size"   // `size <source> <res>`
	Cmp                               = "cmp"    // `cmp <score> <condition> <score>`
	If                                = "if"     // `if <score> <instruction>`
	Unless                            = "unless" // `unless <score> <instruction>`
	Ret                               = "ret"    // `ret`
	Func                              = "func"   // `func <name>\n\t<code>`
	Call                              = "call"   // `call <name>
	Branch                            = "branch" // `branch <name>`
	Raw                               = "raw"    // `raw <command>`
)

type Instruction struct {
	interfaces.Instruction

	Type        interfaces.InstructionType
	DPNamespace string
	Storage     string
	Args        []string
}

func (i Instruction) ToString() string {
	return fmt.Sprintf("%s %s", i.Type, strings.Join(i.Args, " "))
}

func (i Instruction) ToMCCommand() string {
	switch i.Type {
	case Set:
		return fmt.Sprintf("data modify storage %s %s set value %s\n", i.Args[0], i.Args[1], i.Args[2])
	case Copy:
		return fmt.Sprintf("data modify storage %s %s set from storage %s %s\n", i.Args[2], i.Args[3], i.Args[0], i.Args[1])
	case Remove:
		return fmt.Sprintf("data remove storage %s %s\n", i.Args[0], i.Args[1])
	case Math:
		return fmt.Sprintf("execute store result storage %s %s int 1 run scoreboard players operation %s %s %s %s %s\n", i.Storage, RX, RX, i.DPNamespace, i.Args[0], RA, i.DPNamespace)
	case Load:
		return fmt.Sprintf("execute store result score %s %s run data get storage %s %s\n", i.Args[1], i.DPNamespace, i.Storage, i.Args[0])
	case Store:
		return fmt.Sprintf("execute store result storage %s %s int 1 run scoreboard players get %s %s\n", i.Storage, i.Args[1], i.Args[0], i.DPNamespace)
	case Score:
		return fmt.Sprintf("scoreboard players set %s %s %s\n", i.Args[0], i.DPNamespace, i.Args[1])
	case Append:
		log.Fatalln("Not implemented")
	case Size:
		return fmt.Sprintf("execute store result storage %s %s int 1 run data get storage %s %s\n", i.Storage, i.Args[1], i.Storage, i.Args[0])
	case Cmp:
		return fmt.Sprintf("execute if score %s %s %s %s %s run data modify storage %s %s set value 1\n", i.Args[0], i.DPNamespace, i.Args[1], i.Args[2], i.DPNamespace, i.Storage, i.Args[3])
	case If:
		return fmt.Sprintf("execute if score %s %s matches 1.. run %s", i.Args[0], i.DPNamespace, parseInstruction(strings.Join(i.Args[1:], " "), i.DPNamespace, i.Storage).ToMCCommand())
	case Unless:
		return fmt.Sprintf("execute unless score %s %s matches 1.. run %s", i.Args[0], i.DPNamespace, parseInstruction(strings.Join(i.Args[1:], " "), i.DPNamespace, i.Storage).ToMCCommand())
	case Ret:
		return "return 0\n"
	case Call, Branch:
		return fmt.Sprintf("function mcb:%s {function:'%s', function_namespace:'%s', args:'%s', namespace:'%s'}\n", path.Join(paths.Internal, string(i.Type)), i.Args[0], i.Args[1], i.Args[2], i.DPNamespace)
	case Raw:
		return fmt.Sprintf("%s\n", i.Args[0])
	case Func:
		// This is not a valid command, but a placeholder for the function definition
	default:
	}
	log.Debugf("Not implemented: %s", i.ToString())
	return ""
}

func (i Instruction) GetArgs() []string {
	return i.Args
}

func (i Instruction) GetType() interfaces.InstructionType {
	return i.Type
}
