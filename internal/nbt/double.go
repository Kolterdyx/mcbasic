package nbt

import "fmt"

type Double struct {
	Value float64
}

func NewDouble(value float64) *Double {
	return &Double{
		Value: value,
	}
}

func (d *Double) ToString() string {
	return fmt.Sprintf("%fd", d.Value)
}
