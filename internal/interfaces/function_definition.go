package interfaces

import "github.com/Kolterdyx/mcbasic/internal/types"

type FunctionDefinition struct {
	Name       string
	Args       []TypedIdentifier
	ReturnType types.ValueType
}
