package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"strings"
)

const (
	VarPath    = "vars"
	ArgPath    = "args"
	StructPath = "structs"
)

const (
	RA = "$RA"
	RB = "$RB"

	RX  = "$RX"
	RET = "$RET"

	RCF  = "$RCF"
	CALL = "$CALL"

	MaxCallCounter = 65536
)

type Op struct {
	Namespace string
	Scope     string
	Structs   map[string]statements.StructDeclarationStmt
}

func (o *Op) Macro(argName string) string {
	return fmt.Sprintf("$(%s)", argName)
}

// MacroReplace add $ at the start of each line that uses macros. Macros are found in the pattern $(name)
func (o *Op) MacroReplace(source string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if strings.Contains(line, "$(") && !(line[0:1] == "$") {
			lines[i] = "$" + line
		}
	}
	return strings.Join(lines, "\n")
}

func (o *Op) Return() string {
	return fmt.Sprintf("return 1\n")
}

func Cs(s string) string {
	if s == RCF || s == CALL {
		return s
	}
	suffix := ".$(__call__)"
	if strings.Contains(s, suffix) {
		return s
	}
	return s + suffix
}

func (o *Op) Trace(storage, path string) string {

	text := nbt.NewList(
		nbt.NewCompound().
			Set("text", nbt.NewString(path+": ")).
			Set("color", nbt.NewString("green")),
		nbt.NewCompound().
			Set("type", nbt.NewAny("nbt")).
			Set("source", nbt.NewAny("storage")).
			Set("nbt", nbt.NewString(path)).
			Set("storage", nbt.NewString(storage)).
			Set("color", nbt.NewString("yellow")),
	)

	return fmt.Sprintf("tellraw @a[tag=mcblog] %s\n", text.ToString())
}
