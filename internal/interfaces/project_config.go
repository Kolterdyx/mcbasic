package interfaces

type Project struct {
	Name        string   `toml:"name"`
	Namespace   string   `toml:"namespace"`
	Authors     []string `toml:"authors"`
	Entrypoint  string   `toml:"entrypoint"`
	Version     string   `toml:"version"`
	Description string   `toml:"description"`
}

type ProjectConfig struct {
	CleanBeforeInit bool    `toml:"cleanBeforeInit"`
	Debug           bool    `toml:"debug"`
	Project         Project `toml:"Project"`
	Dependencies    struct {
		Headers []string
	} `toml:"-"`
	OutputDir string `toml:"-"`
}
