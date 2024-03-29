package ops

import "fmt"

// And tests if two registers are both truthy and stores the result in the left register.
func (o *Op) And(regLeft string, regRight string) string {
	return fmt.Sprintf("execute if score %s %s matches 1.. if score %s %s matches 1.. run scoreboard players set %s %s 1\n", regLeft, o.Namespace, regRight, o.Namespace, RX, o.Namespace)
}

// Or tests if at least one of two registers is truthy and stores the result in the left register.
func (o *Op) Or(regLeft string, regRight string) string {
	cmd := fmt.Sprintf("execute if score %s %s matches 1.. run scoreboard players set %s %s 1\n", regLeft, o.Namespace, RX, o.Namespace)
	cmd += fmt.Sprintf("execute if score %s %s matches 1.. run scoreboard players set %s %s 1\n", regRight, o.Namespace, RX, o.Namespace)
	return cmd
}
