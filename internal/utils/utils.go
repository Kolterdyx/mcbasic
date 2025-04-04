package utils

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	log "github.com/sirupsen/logrus"
)

func GetHeaderFuncDefs(headers []interfaces.DatapackHeader) map[string]interfaces.FuncDef {
	funcDefs := make(map[string]interfaces.FuncDef)
	for _, header := range headers {
		log.Debugf("Loading header: %s. Functions: %v", header.Namespace, len(header.Definitions.Functions))
		for _, function := range header.Definitions.Functions {
			funcName := fmt.Sprintf("%s:%s", header.Namespace, function.Name)
			f := interfaces.FuncDef{
				Name:       funcName,
				Args:       make([]interfaces.FuncArg, 0),
				ReturnType: function.ReturnType,
			}
			for _, parameter := range function.Args {
				f.Args = append(f.Args, interfaces.FuncArg{
					Name: parameter.Name,
					Type: parameter.Type,
				})
			}
			funcDefs[funcName] = f
		}
	}
	return funcDefs
}
