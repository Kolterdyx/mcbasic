package interfaces

type PrimitiveType string

type ValueType interface {
	Primitive() ValueType
	ToString() string
	IsType(other ValueType) bool
}

type Project struct {
	Name        string   `toml:"name"`
	Namespace   string   `toml:"namespace"`
	Authors     []string `toml:"authors"`
	Entrypoint  string   `toml:"entrypoint"`
	Version     string   `toml:"version"`
	Description string   `toml:"description"`
}

type ProjectConfig struct {
	Project      Project `toml:"Project"`
	Dependencies struct {
		Headers []string
	} `toml:"-"`
	OutputDir string `toml:"-"`
}

type DatapackHeader struct {
	Namespace   string `json:"namespace"`
	Definitions struct {
		Functions []struct {
			Name string `json:"name"`
			Args []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"args"`
			ReturnType string `json:"returnType"`
		} `json:"functions"`
	} `json:"definitions"`
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
	Args       []FuncArg
	ReturnType ValueType
}

type Author struct {
	Name  string
	Email string
}

func (a Author) String() string {
	return a.Name + " <" + a.Email + ">"
}

// StructField represents a field in a struct declaration.
type StructField struct {
	Name string
	Type ValueType
}
