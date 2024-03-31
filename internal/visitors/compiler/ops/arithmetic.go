package ops

import (
	"fmt"
)

func (o *Op) arithmetic(var1 string, var2 string, operator string) string {
	return fmt.Sprintf("scoreboard players operation %s %s %s %s %s\n", cs(var1), o.Namespace, operator, cs(var2), o.Namespace)
}

// Add adds two variables together. The result is stored in the first variable.
func (o *Op) Add(var1 string, var2 string) string {
	return o.arithmetic(var1, var2, "+=")
}

// Sub subtracts the second variable from the first variable. The result is stored in the first variable.
func (o *Op) Sub(var1 string, var2 string) string {
	return o.arithmetic(var1, var2, "-=")
}

// Mul multiplies two variables together. The result is stored in the first variable.
func (o *Op) Mul(var1 string, var2 string) string {
	return o.arithmetic(var1, var2, "*=")
}

// Div divides the first variable by the second variable. The result is stored in the first variable.
func (o *Op) Div(var1 string, var2 string) string {
	return o.arithmetic(var1, var2, "/=")
}

// Mod calculates the modulus of the first variable by the second variable. The result is stored in the first variable.
func (o *Op) Mod(var1 string, var2 string) string {
	return o.arithmetic(var1, var2, "%=")
}
