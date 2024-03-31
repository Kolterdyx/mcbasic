package ops

import "fmt"

func (o *Op) cond(regLeft string, regRight string, operator string) string {
	cmd := fmt.Sprintf("execute if score %s %s %s %s %s run scoreboard players set %s %s 1\n", cs(regLeft), o.Namespace, operator, cs(regRight), o.Namespace, cs(RX), o.Namespace)
	cmd += fmt.Sprintf("execute unless score %s %s %s %s %s run scoreboard players set %s %s 0\n", cs(regLeft), o.Namespace, operator, cs(regRight), o.Namespace, cs(RX), o.Namespace)
	return cmd
}

// Eq compares two registers and stores the result in the RX register.
func (o *Op) Eq(regLeft string, regRight string) string {
	return o.cond(regLeft, regRight, "=")
}

// Neq compares two registers and stores the result in the RX register.
func (o *Op) Neq(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute unless score %s %s = %s %s run scoreboard players set %s %s 1\n", cs(regLeft), o.Namespace, cs(regRight), o.Namespace, cs(RX), o.Namespace)
	cmd += fmt.Sprintf("execute if score %s %s = %s %s run scoreboard players set %s %s 0\n", cs(regLeft), o.Namespace, cs(regRight), o.Namespace, cs(RX), o.Namespace)
	return cmd
}

// Gt compares two registers and stores the result in the RX register.
func (o *Op) Gt(regLeft string, regRight string) string {
	return o.cond(regLeft, regRight, ">")
}

// Gte compares two registers and stores the result in the RX register.
func (o *Op) Gte(regLeft string, regRight string) string {
	return o.cond(regLeft, regRight, ">=")
}

// Lt compares two registers and stores the result in the RX register.
func (o *Op) Lt(regLeft string, regRight string) string {
	return o.cond(regLeft, regRight, "<")
}

// Lte compares two registers and stores the result in the RX register.
func (o *Op) Lte(regLeft string, regRight string) string {
	return o.cond(regLeft, regRight, "<=")
}
