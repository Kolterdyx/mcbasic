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
	Cmp                    = "cmp"    // `cmp <score> <condition> <score>`
	If                     = "if"     // `if <score> <instruction>`
	Unless                 = "unless" // `unless <score> <instruction>`
	Ret                    = "ret"    // `ret`
	Func                   = "func"   // `func <name>\n\t<code>`
	Call                   = "call"   // `call <name>`
)

type Instruction struct {
	Type      InstructionType
	Namespace string
	Storage   string
	Args      []string
}

func (i Instruction) ToString() string {
	return fmt.Sprintf("%s %s", i.Type, strings.Join(i.Args, " "))
}

func (i Instruction) ToMCCommand() string {
	switch i.Type {
	case Set:
		return fmt.Sprintf("data modify storage %s %s set value %s", i.Args[0], i.Args[1], i.Args[2])
	case Copy:
		return fmt.Sprintf("data modify storage %s %s set from storage %s %s", i.Args[2], i.Args[3], i.Args[0], i.Args[1])
	case Math:
		return fmt.Sprintf("execute store result storage %s %s int 1 run scoreboard players operation %s %s %s %s %s", i.Storage, RX, RX, i.Namespace, i.Args[0], RA, i.Namespace)
	case Load:
		return fmt.Sprintf("execute store result score %s %s run data get storage %s %s", i.Args[1], i.Namespace, i.Storage, i.Args[0])
	case Store:
		return fmt.Sprintf("execute store result storage %s %s int 1 run scoreboard players get %s %s", i.Storage, i.Args[0], i.Args[1], i.Namespace)
	case Score:
		return fmt.Sprintf("scoreboard players set %s %s %s", i.Args[0], i.Namespace, i.Args[1])
	case Append:
		panic("not implemented")
	case Size:
		panic("not implemented")
	case Cmp:
		return fmt.Sprintf("execute if score %s %s %s %s %s run data modify storage %s %s set value 1", i.Args[0], i.Namespace, i.Args[1], i.Args[2], i.Namespace, i.Storage, i.Args[3])
	case If:
		return fmt.Sprintf("execute if score %s matches 1.. run %s", i.Args[0], parseInstruction(i.Args[1:], i.Namespace, i.Storage).ToMCCommand())
	case Unless:
		return fmt.Sprintf("execute unless score %s matches 1.. run %s", i.Args[0], parseInstruction(i.Args[1:], i.Namespace, i.Storage).ToMCCommand())
	case Ret:
		return "return 0"
	case Func:
		return ""
	default:
		return ""
	}
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

func (f Function) ToMCFunction() string {
	body := ""
	for _, instr := range f.Instructions {
		body += fmt.Sprintf("%s\n", instr.ToMCCommand())
	}
	return fmt.Sprintf("# Function: %s\n%s", f.Name, body)
}

func ParseIL(source, namespace, storage string) []Function {

	var funcs []Function
	var current *Function

	for _, line := range strings.Split(source, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := split(line)

		if parts[0] == Func && len(parts) > 1 {
			if current != nil {
				funcs = append(funcs, *current)
			}
			current = &Function{Name: parts[1]}
		} else if current != nil {
			instr := parseInstruction(parts, namespace, storage)
			current.Instructions = append(current.Instructions, instr)
		}
	}

	if current != nil {
		funcs = append(funcs, *current)
	}

	return funcs
}

func parseInstruction(parts []string, namespace, storage string) Instruction {
	return Instruction{
		Type:      InstructionType(parts[0]),
		Args:      parts[1:],
		Namespace: namespace,
		Storage:   storage,
	}
}

func split(line string) []string {
	// Split the line by spaces, but keep quoted strings together and snbt data together
	parts := make([]string, 0)
	quote := false
	quoteChar := ' '
	brace := 0
	bracket := 0
	start := 0
	for i, r := range line {
		if r == '"' || r == '\'' {
			if quote && r == quoteChar {
				// End of quoted string
				quote = false
			} else if !quote && line[i-1] != '\\' {
				// Start of quoted string
				quote = true
				quoteChar = r
			}
		} else if r == '{' {
			brace++
		} else if r == '}' {
			brace--
		} else if r == '[' {
			bracket++
		} else if r == ']' {
			bracket--
		} else if r == ' ' && !quote && brace == 0 && bracket == 0 {
			if start < i {
				parts = append(parts, line[start:i])
			}
			start = i + 1
		}
	}
	if start < len(line) {
		parts = append(parts, line[start:])
	}
	return parts
}
