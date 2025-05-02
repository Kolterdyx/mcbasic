package il

import (
	"fmt"
	"strings"
)

type InstructionType string

const (
	_      InstructionType = ""
	Set                    = "set"    // `set <storage> <path> <value>`
	Copy                   = "copy"   // `copy <storage from> <path from> <storage to> <path to>`
	Math                   = "math"   // `math <operation>`
	Load                   = "load"   // `load <path> <score>`
	Store                  = "store"  // `store <score> <path>`
	Score                  = "score"  // `score <target> <score>`
	Append                 = "append" // `append <listPath> <valuePath>`
	Size                   = "size"   // `size <source> <res>`
	Run                    = "run"    // `run <command>`
	Cmp                    = "cmp"    // `cmp <score> <condition> <score>`
	If                     = "if"     // `if <score> <instruction>`
	Unless                 = "unless" // `unless <score> <instruction>`
	Ret                    = "ret"    // `ret`
	Log                    = "log"    // `log <text>`
	Func                   = "func"   // `func <name>\n\t<code>`
	Call                   = "call"   // `call <name>`
)

type Instruction struct {
	Type InstructionType
	Args []string
}

func (i Instruction) ToString() string {
	return fmt.Sprintf("%s %s", i.Type, strings.Join(i.Args, " "))
}

type Function struct {
	Name         string
	Instructions []Instruction
}

func (f Function) ToString() string {
	body := ""
	for _, instr := range f.Instructions {
		body += fmt.Sprintf("\t%s\n", instr.ToString())
	}
	return fmt.Sprintf("%s %s\n%s", Func, f.Name, body)
}

func ParseIL(source string) []Function {

	var funcs []Function
	var current *Function

	for _, line := range strings.Split(source, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		if parts[0] == "func" && len(parts) > 1 {
			if current != nil {
				funcs = append(funcs, *current)
			}
			current = &Function{Name: parts[1]}
		} else if current != nil {
			instr := Instruction{Type: InstructionType(parts[0]), Args: parts[1:]}
			current.Instructions = append(current.Instructions, instr)
		}
	}

	if current != nil {
		funcs = append(funcs, *current)
	}

	return funcs
}
