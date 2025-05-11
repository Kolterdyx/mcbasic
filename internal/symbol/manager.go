package symbol

type Manager struct {
	tables map[string]*Table
}

func NewManager() *Manager {
	return &Manager{
		tables: make(map[string]*Table),
	}
}

func (m *Manager) AddFile(filePath string, table *Table) {
	m.tables[filePath] = table
}

func (m *Manager) GetSymbol(filePath, symbolName string) (Symbol, bool) {
	if table, exists := m.tables[filePath]; exists {
		return table.Lookup(symbolName)
	}
	return Symbol{}, false
}
