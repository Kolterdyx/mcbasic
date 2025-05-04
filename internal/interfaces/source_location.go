package interfaces

import "fmt"

type SourceLocation struct {
	Row int
	Col int
}

func (l SourceLocation) ToString() string {
	return fmt.Sprintf("%d:%d", l.Row+1, l.Col+1)
}
