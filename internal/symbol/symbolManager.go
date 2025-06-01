package symbol

type Manager struct {
	tables map[string]*Table

	rootPath string
}

func NewManager(rootPath string) *Manager {
	return &Manager{
		tables:   make(map[string]*Table),
		rootPath: rootPath,
	}
}

func (s *Manager) AddFile(filePath string, table *Table) {
	s.tables[filePath] = table
}

func (s *Manager) GetFile(filePath string) (*Table, bool) {
	if table, exists := s.tables[filePath]; exists {
		return table, true
	}
	return nil, false
}

func (s *Manager) GetSymbol(filePath, symbolName string) (Symbol, bool) {
	if table, exists := s.GetFile(filePath); exists {
		return table.Lookup(symbolName)
	}
	return Symbol{}, false
}
