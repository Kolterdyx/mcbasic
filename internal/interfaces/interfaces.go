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
	Dependencies struct {
		Headers []string
	}
	OutputDir    string
	EnableTraces bool
}

type DatapackHeader struct {
	Namespace string `json:"namespace"`
	Meta      struct {
		Name string
	} `json:"meta"`
	Definitions struct {
		Functions []struct {
			Name string `json:"name"`
			Args []struct {
				Name string    `json:"name"`
				Type ValueType `json:"type"`
			}
			Returns struct {
				Type    ValueType `json:"type"`
				Storage string    `json:"storage"`
				Path    string    `json:"path"`
			}
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
