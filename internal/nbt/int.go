package nbt

import "fmt"

type Int struct {
	Value

	value int
}

func NewInt(value int) *Int {
	return &Int{
		value: value,
	}
}

func (i *Int) ToString() string {
	return fmt.Sprintf("%dL", i.value)
}
