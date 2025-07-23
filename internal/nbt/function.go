package nbt

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Function struct {
	Name    string
	Context string // Path to the data context, e.g., "ns:path/to/data.context"
}

func NewFunction(name, context string) *Function {
	return &Function{
		Name:    name,
		Context: context,
	}
}

func NewAnonymousFunction(context string) *Function {
	return &Function{
		Name:    generateAnonymousFunctionName(),
		Context: context,
	}
}

func (f *Function) ToString() string {
	return fmt.Sprintf("{name:\"%s\", context:\"%s\"}", f.Name, f.Context)
}

func generateAnonymousFunctionName() string {
	uuidValue, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to generate uuid for anonymous function: %s", err.Error())
		return ""
	}
	return uuidValue.String()
}
