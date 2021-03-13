package main

type SymbolTable struct {
	symbols map[string]Value
	parent  *SymbolTable
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{make(map[string]Value), parent}
}

func (self *SymbolTable) Get(name string) Value {
	value := self.symbols[name]
	if value == nil && self.parent != nil {
		return self.parent.Get(name)
	}
	return value
}

func (self *SymbolTable) Set(name string, value Value) {
	self.symbols[name] = value
}

func (self *SymbolTable) Remove(name string) {
	delete(self.symbols, name)
}
