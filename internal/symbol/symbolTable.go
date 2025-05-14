package symbol

import "fmt"

type Table struct {
	symbols    map[string]Symbol
	parent     *Table
	children   map[string]*Table
	scopeName  string
	originFile string
}

func NewTable(parent *Table, scopeName string, originFile string) *Table {
	return &Table{
		symbols:    make(map[string]Symbol),
		parent:     parent,
		children:   make(map[string]*Table),
		scopeName:  scopeName,
		originFile: originFile,
	}
}

func (s *Table) Define(symbol Symbol) error {
	if _, exists := s.symbols[symbol.Name()]; exists {
		return fmt.Errorf("symbol already defined: %s", symbol.Name())
	}
	s.symbols[symbol.Name()] = symbol
	return nil
}

func (s *Table) Lookup(name string) (Symbol, bool) {
	if symbol, exists := s.symbols[name]; exists {
		return symbol, true
	}
	if s.parent != nil {
		return s.parent.Lookup(name)
	}
	return Symbol{}, false
}

func (s *Table) ScopeName() string {
	return s.scopeName
}

func (s *Table) OriginFile() string {
	return s.originFile
}

func (s *Table) AddChild(table *Table) {
	s.children[table.ScopeName()] = table
}

func (s *Table) GetChild(scope string) (*Table, bool) {
	if child, exists := s.children[scope]; exists {
		return child, true
	}
	return nil, false
}

func (s *Table) GetParent() *Table {
	return s.parent
}
