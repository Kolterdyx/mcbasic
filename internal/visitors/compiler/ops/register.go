package ops

import "fmt"

func (o *Op) Move(from string, to string) string {
	return fmt.Sprintf("data modify storage %s:%s.%s set from storage %s:%s.%s\n", o.Namespace, VarPath, to, o.Namespace, VarPath, from)
}

func (o *Op) MoveConst(value string, to string) string {
	return fmt.Sprintf("data modify storage %s:%s.%s set value %s\n", o.Namespace, VarPath, to, value)
}

func (o *Op) MoveScore(from string, to string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:%s.%s\n", o.Namespace, to, o.Namespace, VarPath, from)
}
