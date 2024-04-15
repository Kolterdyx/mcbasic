package ops

import (
	"fmt"
	"strconv"
)

func (o *Op) Move(from string, to string) string {
	return fmt.Sprintf("data modify storage %s:%s %s set from storage %s:%s %s\n", o.Namespace, VarPath, cs(to), o.Namespace, VarPath, cs(from))
}

func (o *Op) MoveConst(value string, to string) string {
	if _, err := strconv.Atoi(value); err != nil {
		value = strconv.Quote(value)
	}
	return fmt.Sprintf("data modify storage %s:%s %s set value %s\n", o.Namespace, VarPath, cs(to), value)
}

func (o *Op) MoveScore(from string, to string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:%s %s\n", o.Namespace, cs(to), o.Namespace, VarPath, cs(from))
}
