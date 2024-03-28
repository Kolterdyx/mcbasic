package ops

import "fmt"

// Eq compares two registers and stores the result in the left register.
func (o *Op) Eq(regLeft string, regRight string) string {
	return fmt.Sprintf("execute if score %s %s = %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, regLeft, o.Namespace)
}

// Neq compares two registers and stores the result in the left register.
func (o *Op) Neq(regLeft string, regRight string) string {
	return fmt.Sprintf("execute unless score %s %s = %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, regLeft, o.Namespace)
}

// Gt compares two registers and stores the result in the left register.
func (o *Op) Gt(regLeft string, regRight string) string {
	return fmt.Sprintf("execute if score %s %s > %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, regLeft, o.Namespace)
}

// Gte compares two registers and stores the result in the left register.
func (o *Op) Gte(regLeft string, regRight string) string {
	return fmt.Sprintf("execute if score %s %s >= %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, regLeft, o.Namespace)
}

// Lt compares two registers and stores the result in the left register.
func (o *Op) Lt(regLeft string, regRight string) string {
	return fmt.Sprintf("execute if score %s %s < %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, regLeft, o.Namespace)
}

// Lte compares two registers and stores the result in the left register.
func (o *Op) Lte(regLeft string, regRight string) string {
	return fmt.Sprintf("execute if score %s %s <= %s %s run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, regLeft, o.Namespace)
}
