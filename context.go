package main

type Context struct {
	name   string
	parent *Context
	entry  *Position
	table  *SymbolTable
}

func NewContext(name string, parent *Context, entry *Position) *Context {
	return &Context{name, parent, entry, nil}
}
