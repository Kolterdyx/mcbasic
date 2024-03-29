package ops

import "fmt"

// Eq compares two registers and stores the result in the RX register.
func (o *Op) Eq(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute if score %s %s = %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	cmd += fmt.Sprintf("execute unless score %s %s = %s %s run scoreboard players set %s %s 0\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	return cmd
}

// Neq compares two registers and stores the result in the RX register.
func (o *Op) Neq(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute unless score %s %s = %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	cmd += fmt.Sprintf("execute if score %s %s = %s %s run scoreboard players set %s %s 0\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	return cmd
}

// Gt compares two registers and stores the result in the RX register.
func (o *Op) Gt(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute if score %s %s > %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	cmd += fmt.Sprintf("execute unless score %s %s > %s %s run scoreboard players set %s %s 0\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	return cmd
}

// Gte compares two registers and stores the result in the RX register.
func (o *Op) Gte(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute if score %s %s >= %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	cmd += fmt.Sprintf("execute unless score %s %s >= %s %s run scoreboard players set %s %s 0\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	return cmd
}

// Lt compares two registers and stores the result in the RX register.
func (o *Op) Lt(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute if score %s %s < %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	cmd += fmt.Sprintf("execute unless score %s %s < %s %s run scoreboard players set %s %s 0\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	return cmd
}

// Lte compares two registers and stores the result in the RX register.
func (o *Op) Lte(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute if score %s %s <= %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	cmd += fmt.Sprintf("execute unless score %s %s <= %s %s run scoreboard players set %s %s 0\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
	return cmd
}
