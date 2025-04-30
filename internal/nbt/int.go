package nbt

import "fmt"

type Int struct {
	Value

	value int64
}

func NewInt(value int64) *Int {
	return &Int{
		value: value,
	}
}

func (i *Int) ToString() string {
	return fmt.Sprintf("%d", i.value)
}
