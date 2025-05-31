package nbt

import (
	"fmt"
)

type StringColor string

const (
	Black       StringColor = "black"
	DarkBlue    StringColor = "dark_blue"
	DarkGreen   StringColor = "dark_green"
	DarkAqua    StringColor = "dark_aqua"
	DarkRed     StringColor = "dark_red"
	DarkPurple  StringColor = "dark_purple"
	Gold        StringColor = "gold"
	Gray        StringColor = "gray"
	DarkGray    StringColor = "dark_gray"
	Blue        StringColor = "blue"
	Green       StringColor = "green"
	Aqua        StringColor = "aqua"
	Red         StringColor = "red"
	LightPurple StringColor = "light_purple"
	Yellow      StringColor = "yellow"
	White       StringColor = "white"
)

func HexColor(hex string) StringColor {
	return StringColor(fmt.Sprintf("#%s", hex))
}

type StringFormat struct {
	Color         StringColor
	Bold          bool
	Italic        bool
	Underlined    bool
	Strikethrough bool
	Obfuscated    bool
}

func (s StringFormat) ToString() string {
	return fmt.Sprintf("StringFormat{Color: %s, Bold: %t, Italic: %t, Underlined: %t, Strikethrough: %t, Obfuscated: %t}", s.Color, s.Bold, s.Italic, s.Underlined, s.Strikethrough, s.Obfuscated)
}

func (s StringFormat) ToMap() map[string]string {
	return map[string]string{
		"color":         string(s.Color),
		"bold":          fmt.Sprintf("%t", s.Bold),
		"italic":        fmt.Sprintf("%t", s.Italic),
		"underlined":    fmt.Sprintf("%t", s.Underlined),
		"strikethrough": fmt.Sprintf("%t", s.Strikethrough),
		"obfuscated":    fmt.Sprintf("%t", s.Obfuscated),
	}
}

type String struct {
	Value

	value string
}

func NewString(value string, fmtParams ...any) *String {
	return &String{
		value: escapeSingleQuotes(fmt.Sprintf(value, fmtParams...)),
	}
}

func NewFormattedString(value string, format StringFormat, fmtParams ...any) Compound {
	compound := NewCompound()
	compound.Set("text", NewString(fmt.Sprintf(value, fmtParams...)))
	for k, v := range format.ToMap() {
		compound.Set(k, NewAny(v))
	}
	return compound
}

func NewErrorString(value string, fmtParams ...any) Compound {
	return NewFormattedString(value, StringFormat{Color: Red, Italic: true}, fmtParams...)
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
