package interfaces

import "fmt"

type SourceLocation struct {
	File string
	Row  int
	Col  int
}

func (l SourceLocation) ToString() string {
	return fmt.Sprintf("%s:%d:%d", l.File, l.Row+1, l.Col+1)
}
