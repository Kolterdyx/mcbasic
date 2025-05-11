package nbt

import "fmt"

type String struct {
	Value

	value string
}

func NewString(value string) *String {
	return &String{
		value: escapeSingleQuotes(value),
	}
}

func (s *String) ToString() string {
	return fmt.Sprintf("'%s'", s.value)
}

func escapeSingleQuotes(value string) string {
	escaped := ""
	for _, char := range value {
		if char == '\'' {
			escaped += "\\'"
		} else {
			escaped += string(char)
		}
	}
	return escaped
}
