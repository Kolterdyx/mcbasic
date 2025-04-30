package nbt

import "fmt"

type String struct {
	Value

	value string
}

func NewString(value string) *String {
	return &String{
		value: value,
	}
}

func (s *String) ToString() string {
	return fmt.Sprintf("'%s'", s.value)
}
