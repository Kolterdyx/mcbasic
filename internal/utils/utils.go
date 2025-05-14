package utils

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

func CmpOperator(operator interfaces.TokenType) string {
	switch operator {
	case tokens.Greater:
		return ">"
	case tokens.GreaterEqual:
		return ">="
	case tokens.Less:
		return "<"
	case tokens.LessEqual:
		return "<="
	case tokens.EqualEqual:
		return "="
	case tokens.BangEqual:
		return "!="
	default:
	}
	log.Fatalln("unknown operator")
	return ""
}

func SplitFunctionName(lexeme, namespace string) (string, string) {
	parts := strings.Split(lexeme, ":")
	if len(parts) == 1 {
		return namespace, parts[0]
	}
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	panic(fmt.Sprintf("Invalid function name: %s", lexeme))
}

func IsListType(valueType types.ValueType) bool {
	switch valueType.(type) {
	case types.ListTypeStruct:
		return true
	}
	return false
}

func IsStructType(valueType types.ValueType) bool {
	switch valueType.(type) {
	case types.StructTypeStruct:
		return true
	}
	return false
}
