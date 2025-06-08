package symbol

import "strings"

type ModuleNode struct {
	name     string
	table    *Table
	children map[string]*ModuleNode
}

type Manager struct {
	root     *ModuleNode
	rootPath string
}

func NewManager(rootPath string) *Manager {
	return &Manager{
		rootPath: rootPath,
		root: &ModuleNode{
			name:     "",
			children: make(map[string]*ModuleNode),
		},
	}
}

func (m *Manager) AddFile(filePath string, table *Table) {
	if table == nil {
		return
	}

	// Strip the root path and split file path into components
	relPath := strings.TrimPrefix(filePath, m.rootPath)
	relPath = strings.Trim(relPath, "/")
	parts := strings.Split(relPath, "/")

	curr := m.root
	for _, part := range parts {
		if curr.children[part] == nil {
			curr.children[part] = &ModuleNode{
				name:     part,
				children: make(map[string]*ModuleNode),
			}
		}
		curr = curr.children[part]
	}
	curr.table = table
}

func (m *Manager) GetFile(filePath string) (*Table, bool) {
	node := m.getNode(filePath)
	if node != nil && node.table != nil {
		return node.table, true
	}
	return nil, false
}

func (m *Manager) GetSymbol(filePath, symbolName string) (Symbol, bool) {
	node := m.getNode(filePath)
	if node != nil && node.table != nil {
		return node.table.Lookup(symbolName)
	}
	return Symbol{}, false
}

func (m *Manager) getNode(filePath string) *ModuleNode {
	relPath := strings.TrimPrefix(filePath, m.rootPath)
	relPath = strings.Trim(relPath, "/")
	parts := strings.Split(relPath, "/")

	curr := m.root
	for _, part := range parts {
		next := curr.children[part]
		if next == nil {
			return nil
		}
		curr = next
	}
	return curr
}
