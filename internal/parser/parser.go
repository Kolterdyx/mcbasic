package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	log "github.com/sirupsen/logrus"
)

type Parser struct {
	HadError bool
	Tokens   []tokens.Token
	Headers  []interfaces.DatapackHeader

	currentScope string
	current      int

	variables map[string][]statements.VarDef
	functions []interfaces.FuncDef
}

func (p *Parser) Parse() Program {
	var functions []statements.FunctionDeclarationStmt
	funcDefMap := utils.GetHeaderFuncDefs(p.Headers)
	p.functions = make([]interfaces.FuncDef, 0)
	for _, funcDef := range funcDefMap {
		p.functions = append(p.functions, funcDef)
	}
	p.variables = make(map[string][]statements.VarDef)
	p.functions = append(p.functions,
		interfaces.FuncDef{Name: "mcb:print", ReturnType: expressions.VoidType, Args: []interfaces.FuncArg{{Name: "text", Type: expressions.VoidType}}},
		interfaces.FuncDef{Name: "mcb:log", ReturnType: expressions.VoidType, Args: []interfaces.FuncArg{{Name: "text", Type: expressions.VoidType}}},
		interfaces.FuncDef{Name: "mcb:exec", ReturnType: expressions.VoidType, Args: []interfaces.FuncArg{{Name: "command", Type: expressions.StringType}}},
		interfaces.FuncDef{Name: "mcb:len", ReturnType: expressions.IntType, Args: []interfaces.FuncArg{{Name: "string", Type: expressions.StringType}}},

		interfaces.FuncDef{Name: "math:sqrt", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},

		interfaces.FuncDef{Name: "math:sin", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
		interfaces.FuncDef{Name: "math:cos", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
		interfaces.FuncDef{Name: "math:tan", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
		interfaces.FuncDef{Name: "math:asin", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
		interfaces.FuncDef{Name: "math:acos", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
		interfaces.FuncDef{Name: "math:atan", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},

		interfaces.FuncDef{Name: "math:floor", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
		interfaces.FuncDef{Name: "math:ceil", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
		interfaces.FuncDef{Name: "math:round", ReturnType: expressions.DoubleType, Args: []interfaces.FuncArg{{Name: "x", Type: expressions.DoubleType}}},
	)
	for !p.IsAtEnd() {
		statement := p.statement()
		if statement == nil {
			continue
		}
		if statement.TType() != statements.FunctionDeclarationStmtType {
			log.Errorf("Only function declarations are allowed at the top level. Found: %s\n", statement.TType())
			continue
		}
		funcStmt := statement.(statements.FunctionDeclarationStmt)
		functions = append(functions, funcStmt)
		p.functions = append(p.functions, interfaces.FuncDef{Name: funcStmt.Name.Lexeme, ReturnType: funcStmt.ReturnType, Args: funcStmt.Parameters})
	}
	return Program{Functions: functions}
}
