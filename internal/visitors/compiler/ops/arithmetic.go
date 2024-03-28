package ops

import (
	"fmt"
)

// Add adds two variables together. The result is stored in the first variable.
func (o *Op) Add(var1 string, var2 string) string {
	return fmt.Sprintf("scoreboard players operation %s %s += %s %s\n", var1, o.Namespace, var2, o.Namespace)
}

// Sub subtracts the second variable from the first variable. The result is stored in the first variable.
func (o *Op) Sub(var1 string, var2 string) string {
	return fmt.Sprintf("scoreboard players operation %s %s -= %s %s\n", var1, o.Namespace, var2, o.Namespace)
}

// Mul multiplies two variables together. The result is stored in the first variable.
func (o *Op) Mul(var1 string, var2 string) string {
	return fmt.Sprintf("scoreboard players operation %s %s *= %s %s\n", var1, o.Namespace, var2, o.Namespace)
}

// Div divides the first variable by the second variable. The result is stored in the first variable.
func (o *Op) Div(var1 string, var2 string) string {
	return fmt.Sprintf("scoreboard players operation %s %s /= %s %s\n", var1, o.Namespace, var2, o.Namespace)
}

// Mod calculates the modulus of the first variable by the second variable. The result is stored in the first variable.
func (o *Op) Mod(var1 string, var2 string) string {
	return fmt.Sprintf("scoreboard players operation %s %s %%= %s %s\n", var1, o.Namespace, var2, o.Namespace)
}
