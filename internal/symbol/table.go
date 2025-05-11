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

func (t *Table) Define(symbol Symbol) error {
	if _, exists := t.symbols[symbol.Name()]; exists {
		return fmt.Errorf("symbol already defined: %s", symbol.Name())
	}
	t.symbols[symbol.Name()] = symbol
	return nil
}

func (t *Table) Lookup(name string) (Symbol, bool) {
	if symbol, exists := t.symbols[name]; exists {
		return symbol, true
	}
	if t.parent != nil {
		return t.parent.Lookup(name)
	}
	return Symbol{}, false
}
