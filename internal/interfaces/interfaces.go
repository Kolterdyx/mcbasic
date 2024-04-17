package interfaces

type ProjectConfig struct {
	Project struct {
		Name       string
		Namespace  string
		Authors    []string
		Entrypoint string
	}
	OutputDir string
}

type SourceLocation struct {
	Line   int
	Column int
}
