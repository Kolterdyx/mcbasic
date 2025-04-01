package interfaces

type ProjectConfig struct {
	Project struct {
		Name        string
		Namespace   string
		Authors     []string
		Entrypoint  string
		Version     string
		Description string
	}
	OutputDir           string
	EnableTraces        bool
	FixedPointPrecision int
}

type SourceLocation struct {
	Row int
	Col int
}
