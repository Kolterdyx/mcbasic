package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

type Program struct {
	Functions []statements.FunctionDeclarationStmt
}

func (p Program) Visit(visitor statements.StmtVisitor) []interface{} {
	var res []interface{}
	for _, f := range p.Functions {
		res = append(res, f.Accept(visitor))
	}
	return res
}