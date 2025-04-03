package interfaces

type ValueType string

type ProjectConfig struct {
	Project struct {
		Name        string
		Namespace   string
		Authors     []string
		Entrypoint  string
		Version     string
		Description string
	}
	OutputDir    string
	EnableTraces bool
}

type SourceLocation struct {
	Row int
	Col int
}

type FuncArg struct {
	Name string
	Type ValueType
}

type FuncDef struct {
	Name       string
	Parameters []FuncArg
	ReturnType ValueType
}
