package interfaces

import "github.com/Kolterdyx/mcbasic/internal/types"

type FunctionDefinition struct {
	Name       string
	Parameters []TypedIdentifier
	ReturnType types.ValueType
}
